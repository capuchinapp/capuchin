<script lang="ts">
    import Fa from 'svelte-fa';
    import {
        faEllipsis,
        faPencil,
        faTrash,
        faPlus,
        faSave,
        faCalendarDay,
        faHome,
        faPlay,
        faStop,
        faRubleSign,
        faCircleCheck,
        faHistory,
    } from '@fortawesome/free-solid-svg-icons';
    import {faStar, faCirclePlay, faCircle, faClock} from '@fortawesome/free-regular-svg-icons';
    import dayjs from 'dayjs';
    import AirDatepicker from 'air-datepicker';
    import airDatepickerLocaleRu from 'air-datepicker/locale/ru';
    import TomSelect from 'tom-select';
    import {Modal as BSModal} from 'bootstrap/dist/js/bootstrap.esm';
    import IMask from 'imask';
    import {onMount} from 'svelte';
    import {t} from '../i18n';
    import {timeTotalSeconds, timeTotalClass, settings as settingsStore} from '../stores';
    import PageHeader from '../components/PageHeader.svelte';
    import NoContent from '../components/NoContent.svelte';
    import Modal from '../components/Modal.svelte';
    import {toastWarning} from '../services/Toast';
    import TimelogAPI from '../services/Timelog';
    import ClientAPI from '../services/Clients';
    import ProjectAPI from '../services/Projects';
    import TaskAPI from '../services/Tasks';
    import FavoritesAPI from '../services/Favorites';
    import {moneyInteger, moneyFormatted} from '../services/MoneyHelper';
    import {nowTime, getTime, getDiffSeconds, calculateTime, hoursFromSeconds} from '../services/DateTime';
    import Validation from '../services/Validation';
    import type {Client, Favorite, ProjectOption, TaskOption, Timelog} from '../types';

    type ValidateErrors = {
        projectId: string;
        taskId: string;
        date: string;
        timeStart: string;
    };

    type ValidateFavoriteErrors = {
        projectId: string;
        taskId: string;
    };

    const now = dayjs();
    const timelogModalId = 'js-timelog-form';
    const timelogDeleterModalId = 'js-timelog-deleter-form';
    const favoritesModalId = 'js-favorites-form';
    const historyModalId = 'js-history-form';
    const imaskTime = {
        mask: 'HH:MM',
        lazy: false,
        blocks: {
            HH: {
                mask: IMask.MaskedRange,
                from: 0,
                to: 23,
            },
            MM: {
                mask: IMask.MaskedRange,
                from: 0,
                to: 59,
            },
        },
    };
    const projectTpl: ProjectOption = {
        clientId: '',
        id: '',
        name: '',
        billableRate: '0.00',
    };
    const taskTpl: TaskOption = {
        id: '',
        name: '',
    };
    const timelogTpl = {
        id: '',
        project: {...projectTpl},
        task: {...taskTpl},
        date: '',
        timeStart: '',
        timeEnd: null as string | null,
        billableRate: '0.00',
        comment: '',
        _dateFormatted: '',
        _timeStartFormatted: '',
        _timeEndFormatted: '',
        _duration: '',
        _amount: '',
    };
    const favoriteTpl = {
        id: '',
        name: '',
        project: {...projectTpl},
        task: {...taskTpl},
        billableRate: '0.00',
        comment: '',
    };

    let formCalendar = $state<AirDatepicker<HTMLElement>>();
    let selectedCalendar = $state<AirDatepicker<HTMLElement>>();

    let selectedDate = $state<dayjs.Dayjs>();
    let projectSelect!: TomSelect;
    let taskSelect!: TomSelect;
    let projectFavoriteSelect!: TomSelect;
    let taskFavoriteSelect!: TomSelect;
    let timelogModal!: BSModal;
    let timelogDeleterModal!: BSModal;
    let favoritesModal = $state<BSModal>();
    let historyModal!: BSModal;
    let validation!: Validation;

    let startedTimelogTimer = $state<number>();
    let startedTimelogId = $state<string>();
    let startedTimelogTime = $state<dayjs.Dayjs | undefined>();

    let timelogDeleterTimelogId = $state<string>();

    let validateErrors = $state<ValidateErrors>({
        projectId: '',
        taskId: '',
        date: '',
        timeStart: '',
    });

    let validateFavoriteErrors = $state<ValidateFavoriteErrors>({
        projectId: '',
        taskId: '',
    });

    let title = $state('');
    let timelogList = $state<Timelog[]>([]);
    let clientList: Client[] = [];
    let projectList: ProjectOption[] = [];
    let taskList: TaskOption[] = [];
    let favoriteList = $state<Favorite[]>([]);
    let historyList = $state<Timelog[]>([]);
    let timelog = $state({...timelogTpl});
    let startedTimelogSeconds = $state(0);
    let editingFavorite = $state({...favoriteTpl});

    $effect(() => {
        timelog.billableRate = timelog.project.billableRate;
        editingFavorite.billableRate = editingFavorite.project.billableRate;
    });

    $effect(() => {
        if (selectedDate) {
            selectedCalendar!.selectDate(selectedDate.toDate(), {
                updateTime: true,
                silent: true,
            });

            TimelogAPI.getList(selectedDate!.format('YYYY-MM-DD'), selectedDate!.format('YYYY-MM-DD')).then((result) => {
                timelogList = result;

                let dayName = $t('today');
                if (selectedDate!.format('YYYY-MM-DD') !== dayjs().format('YYYY-MM-DD')) {
                    dayName = selectedDate!.format('ddd');
                }

                title = `${dayName}, ${selectedDate!.format('DD MMM')}`;

                calculateTimeTotalSeconds();
            });
        }
    });

    $effect(() => {
        const idx = timelogList.findIndex((e) => e.timeEnd === null);
        if (idx !== -1) {
            const record = timelogList[idx];

            if (startedTimelogId !== record.id) {
                stopTimer();

                startedTimelogId = record.id;
                startedTimelogTime = dayjs(`${record.date} ${record.timeStart}`);
                calculateTimeTotalSeconds();

                startedTimelogTimer = setInterval(() => {
                    calculateTimeTotalSeconds();
                }, 1000);

                $timeTotalClass = 'text-success';
            }
        }
    });

    onMount(async () => {
        await FavoritesAPI.getList().then((res) => (favoriteList = res));
        await ClientAPI.getList().then((res) => (clientList = res.filter((e) => !e.archivedAt)));
        await ProjectAPI.getList(true).then(
            (res) =>
                (projectList = res
                    .filter((e) => !e.archivedAt)
                    .map((e) => ({
                        clientId: e.clientId,
                        id: e.id,
                        name: e.name,
                        billableRate: moneyFormatted(e.billableRate / 100),
                    })))
        );

        validation = new Validation();

        projectSelect = new TomSelect('#js-project-select', {
            options: projectList,
            optgroups: clientList,
            valueField: 'id',
            labelField: 'name',
            sortField: 'name',
            searchField: ['name'],
            optgroupField: 'clientId',
            optgroupValueField: 'id',
            optgroupLabelField: 'name',
            onChange: (value: string) => {
                timelog.project = findProject(String(value));
                loadTaskList('track', timelog.project.id, null);
            },
            allowEmptyOption: true,
            lockOptgroupOrder: true,
        });

        taskSelect = new TomSelect('#js-task-select', {
            options: taskList,
            valueField: 'id',
            labelField: 'name',
            sortField: 'name',
            searchField: ['name'],
            onChange: (value: string) => (timelog.task = findTask(String(value))),
            allowEmptyOption: true,
        });

        projectFavoriteSelect = new TomSelect('#js-project-favorite-select', {
            options: projectList,
            optgroups: clientList,
            valueField: 'id',
            labelField: 'name',
            sortField: 'name',
            searchField: ['name'],
            optgroupField: 'clientId',
            optgroupValueField: 'id',
            optgroupLabelField: 'name',
            onChange: (value: string) => {
                editingFavorite.project = findProject(String(value));
                loadTaskList('favorite', editingFavorite.project.id, null);
            },
            allowEmptyOption: true,
            lockOptgroupOrder: true,
        });

        taskFavoriteSelect = new TomSelect('#js-task-favorite-select', {
            options: taskList,
            valueField: 'id',
            labelField: 'name',
            sortField: 'name',
            searchField: ['name'],
            onChange: (value: string) => (editingFavorite.task = findTask(String(value))),
            allowEmptyOption: true,
        });

        timelogModal = new BSModal(document.getElementById(timelogModalId)!);
        timelogDeleterModal = new BSModal(document.getElementById(timelogDeleterModalId)!);
        favoritesModal = new BSModal(document.getElementById(favoritesModalId)!);
        historyModal = new BSModal(document.getElementById(historyModalId)!);

        let inputTimeStart = IMask(document.getElementById('js-time-start')!, imaskTime);
        inputTimeStart.on('complete', function () {
            timelog.timeStart = `${inputTimeStart.value}:00`;
            timelog._timeStartFormatted = inputTimeStart.value;
        });

        let inputTimeEnd = IMask(document.getElementById('js-time-end')!, imaskTime);
        inputTimeEnd.on('complete', function () {
            timelog.timeEnd = `${inputTimeEnd.value}:00`;
            timelog._timeEndFormatted = inputTimeEnd.value;
            calculateAmount();
        });

        formCalendar = new AirDatepicker('#js-form-calendar', {
            locale: airDatepickerLocaleRu,
            dateFormat: $settingsStore.dateFormat,
            autoClose: true,
            isMobile: true,
            toggleSelected: false,
            onSelect: ({date}) => {
                let dt = dayjs(date.toString());

                timelog.date = dt.format('YYYY-MM-DD');
                timelog._dateFormatted = dt.format($settingsStore.dateFormat);
            },
        });

        selectedCalendar = new AirDatepicker('#js-selected-calendar', {
            locale: airDatepickerLocaleRu,
            dateFormat: 'MMM yyyy',
            autoClose: true,
            isMobile: true,
            toggleSelected: false,
            onSelect: ({date}) => (selectedDate = dayjs(date.toString())),
        });

        selectedDate = now;
    });

    async function onSubmitForm() {
        validateErrors = {
            projectId: '',
            taskId: '',
            date: '',
            timeStart: '',
        };

        const data = {
            projectId: timelog.project.id,
            taskId: timelog.task.id !== '' ? timelog.task.id : null,
            date: timelog.date,
            timeStart: timelog.timeStart,
            timeEnd: timelog.timeEnd === '' ? null : timelog.timeEnd,
            billableRate: moneyInteger(timelog.billableRate),
            comment: timelog.comment,
        };

        validation.validateTimelog(
            data,
            async () => {
                if (timelog.id) {
                    const res = await TimelogAPI.update(timelog.id, data);
                    if (res !== null) {
                        const idx = timelogList.findIndex((e) => e.id === timelog.id);
                        if (idx !== -1) {
                            if (
                                selectedDate!.format('YYYY-MM-DD') !==
                                dayjs(`${res.date} ${res.timeStart}`).format('YYYY-MM-DD')
                            ) {
                                timelogList.splice(idx, 1);
                                timelogList = timelogList; // Необходимо для срабатывания реактивности
                                timelogModal.hide();

                                return;
                            }

                            timelogList[idx] = res;

                            if (startedTimelogId === timelog.id) {
                                startedTimelogTime = dayjs(`${timelog.date} ${timelog.timeStart}`);
                            }

                            timelogModal.hide();
                        }
                    }
                } else {
                    const res = await TimelogAPI.create(data);
                    if (res !== null) {
                        if (
                            selectedDate!.format('YYYY-MM-DD') ===
                            dayjs(`${res.date} ${res.timeStart}`).format('YYYY-MM-DD')
                        ) {
                            timelogList = [res, ...timelogList];
                        }

                        onAfterCreate(res.id);
                        timelogModal.hide();
                    }
                }

                calculateTimeTotalSeconds();
            },
            (err) => {
                validateErrors = err;
            }
        );
    }

    function onShowForm() {
        const now = dayjs();

        timelog = {
            ...timelogTpl,
            ...{
                date: now.format('YYYY-MM-DD'),
                timeStart: nowTime(true),
                _dateFormatted: now.format($settingsStore.dateFormat),
                _timeStartFormatted: nowTime(false),
            },
        };

        validateErrors = {
            projectId: '',
            taskId: '',
            date: '',
            timeStart: '',
        };

        projectSelect.setValue('', false);

        formCalendar!.selectDate(now.toDate(), {
            updateTime: true,
            silent: true,
        });

        timelogModal.show();
    }

    async function onEditTimelog(timelogId: string) {
        const res = await TimelogAPI.getById(timelogId);
        if (res !== null) {
            const project = findProject(res.projectId);
            if (project.id === '') {
                toastWarning($t('track.warnCannotChangeRecordInArchive'), false);

                return;
            }

            const dt = dayjs(res.date);

            timelog.id = res.id;
            timelog.project = project;
            timelog.date = res.date;
            timelog.timeStart = res.timeStart;
            timelog.timeEnd = res.timeEnd;
            timelog.billableRate = moneyFormatted(res.billableRate / 100);
            timelog.comment = res.comment;

            timelog._dateFormatted = dt.format($settingsStore.dateFormat);
            timelog._timeStartFormatted = res.timeStart.substr(0, 5);

            if (res.timeEnd !== null) {
                timelog._timeEndFormatted = res.timeEnd.substr(0, 5);
                calculateAmount();
            } else {
                timelog._timeEndFormatted = '';
                timelog._duration = '';
                timelog._amount = '';
            }

            projectSelect.setValue(project.id, true);

            await loadTaskList('track', project.id, res.taskId);
            timelog.task = findTask(res.taskId);

            formCalendar!.selectDate(dt.toDate(), {
                updateTime: true,
                silent: true,
            });

            timelogModal.show();
        }
    }

    async function onDeleteTimelogModal(timelogId: string) {
        timelogDeleterTimelogId = timelogId;
        timelogDeleterModal.show();
    }

    async function onDeleteTimelog(timelogId: string) {
        const res = await TimelogAPI.delete(timelogId);
        if (res !== null) {
            timelogList = timelogList.filter((e) => e.id !== timelogId);
        }

        timelogDeleterModal.hide();
    }

    async function onStart(timelogId: string) {
        const idx = timelogList.findIndex((e) => e.id === timelogId);
        if (idx !== -1) {
            const record = timelogList[idx];

            const res = await TimelogAPI.create({
                projectId: record.projectId,
                taskId: record.taskId !== '' ? record.taskId : null,
                date: dayjs().format('YYYY-MM-DD'),
                timeStart: nowTime(true),
                timeEnd: null,
                billableRate: record.billableRate,
                comment: record.comment,
            });
            if (res !== null) {
                timelogList = [res, ...timelogList];
                onAfterCreate(res.id);
            }
        }
    }

    async function onStop(timelogId: string) {
        const now = dayjs();
        const res = await TimelogAPI.stop(timelogId, now.format('YYYY-MM-DD'), now.format('HH:mm:ss'));
        if (res !== null) {
            const idx = timelogList.findIndex((e) => e.id === timelogId);
            if (idx !== -1) {
                timelogList[idx] = res;
                stopTimer();
            }
        }
    }

    async function onAfterCreate(createdTimelogId: string) {
        // Если есть запущенная запись переполучим ее (она уже остановлена)
        const idx = timelogList.findIndex((e) => e.id !== createdTimelogId && e.timeEnd === null);
        if (idx !== -1) {
            const res = await TimelogAPI.getById(timelogList[idx].id);
            if (res !== null) {
                timelogList[idx] = res;
            }
        }

        if (selectedDate!.format('YYYY-MM-DD') !== now.format('YYYY-MM-DD')) {
            selectedDate = now;
        }
    }

    async function onSubmitFavoriteForm() {
        const data = {
            name: editingFavorite.name,
            projectId: editingFavorite.project.id,
            taskId: editingFavorite.task.id !== '' ? editingFavorite.task.id : null,
            billableRate: moneyInteger(editingFavorite.billableRate),
            comment: editingFavorite.comment,
        };

        if (editingFavorite.id) {
            const res = await FavoritesAPI.update(editingFavorite.id, data);
            if (res !== null) {
                const idx = favoriteList.findIndex((e) => e.id === editingFavorite.id);
                if (idx !== -1) {
                    favoriteList[idx] = res;
                    onClearFavoriteForm();
                }
            }
        } else {
            const res = await FavoritesAPI.create(data);
            if (res !== null) {
                favoriteList = [res, ...favoriteList];
                onClearFavoriteForm();
            }
        }
    }

    async function onEditFavorite(favoriteId: string) {
        const res = await FavoritesAPI.getById(favoriteId);
        if (res !== null) {
            const project = findProject(res.projectId);
            if (project.id === '') {
                toastWarning($t('track.warnCannotChangeRecordInArchive'), false);

                return;
            }

            editingFavorite.id = res.id;
            editingFavorite.name = res.name;
            editingFavorite.project = project;
            editingFavorite.billableRate = moneyFormatted(res.billableRate / 100);
            editingFavorite.comment = res.comment;

            projectFavoriteSelect.setValue(project.id, true);

            await loadTaskList('favorite', project.id, res.taskId);
            editingFavorite.task = findTask(res.taskId);
        }
    }

    async function onDeleteFavorite(favoriteId: string) {
        const res = await FavoritesAPI.delete(favoriteId);
        if (res !== null) {
            favoriteList = favoriteList.filter((e) => e.id !== favoriteId);
        }
    }

    async function onStartFavorite(favoriteId: string) {
        const idx = favoriteList.findIndex((e) => e.id === favoriteId);
        if (idx !== -1) {
            const record = favoriteList[idx];

            const res = await TimelogAPI.create({
                projectId: record.projectId,
                taskId: record.taskId !== '' ? record.taskId : null,
                date: dayjs().format('YYYY-MM-DD'),
                timeStart: nowTime(true),
                timeEnd: null,
                billableRate: record.billableRate,
                comment: record.comment,
            });
            if (res !== null) {
                timelogList = [res, ...timelogList];
                onAfterCreate(res.id);
            }
        }
    }

    function onClearFavoriteForm() {
        validateFavoriteErrors = {
            projectId: '',
            taskId: '',
        };
        editingFavorite = {...favoriteTpl};
        projectFavoriteSelect.setValue('', false);

        taskFavoriteSelect.clear(true);
        taskFavoriteSelect.clearOptions();
        taskFavoriteSelect.refreshOptions(false);
    }

    async function onCompleteTask(taskId: string) {
        const res = await TaskAPI.complete(taskId);
        if (res !== null) {
            timelogList = timelogList.map((e) => {
                if (e.taskId === taskId) {
                    return {...e, taskCompletedAt: res.completedAt};
                }
                return e;
            });
        }
    }

    async function onShowHistory(n: number) {
        TimelogAPI.getLastN(n).then((result) => {
            historyList = result;
            historyModal.show();
        });
    }

    async function onStartHistory(historyId: string) {
        const idx = historyList.findIndex((e) => e.id === historyId);
        if (idx !== -1) {
            const record = historyList[idx];

            const res = await TimelogAPI.create({
                projectId: record.projectId,
                taskId: record.taskId !== '' ? record.taskId : null,
                date: dayjs().format('YYYY-MM-DD'),
                timeStart: nowTime(true),
                timeEnd: null,
                billableRate: record.billableRate,
                comment: record.comment,
            });
            if (res !== null) {
                timelogList = [res, ...timelogList];
                onAfterCreate(res.id);
                historyModal.hide();
            }
        }
    }

    async function loadTaskList(type: string, projectId: string, selectedTaskId: string | null) {
        await TaskAPI.getList(projectId).then((res) => {
            taskList = res
                .filter((e) => !e.archivedAt)
                .map((e) => ({
                    id: e.id,
                    name: e.name,
                }));

            switch (type) {
                case 'track':
                    taskSelect.clear(true);
                    taskSelect.clearOptions();
                    taskSelect.addOptions(taskList, false);
                    taskSelect.refreshOptions(false);

                    if (selectedTaskId !== null) {
                        taskSelect.setValue(selectedTaskId, true);
                    }
                    break;
                case 'favorite':
                    taskFavoriteSelect.clear(true);
                    taskFavoriteSelect.clearOptions();
                    taskFavoriteSelect.addOptions(taskList, false);
                    taskFavoriteSelect.refreshOptions(false);

                    if (selectedTaskId !== null) {
                        taskFavoriteSelect.setValue(selectedTaskId, true);
                    }
                    break;
            }
        });
    }

    function findProject(projectId: string): ProjectOption {
        let project: ProjectOption = {...projectTpl};

        const idx = projectList.findIndex((e) => e.id === projectId);
        if (idx !== -1) {
            project = projectList[idx];
        }

        return project;
    }

    function findTask(taskId: string | null): TaskOption {
        let task: TaskOption = {...taskTpl};

        const idx = taskList.findIndex((e) => e.id === taskId);
        if (idx !== -1) {
            task = taskList[idx];
        }

        return task;
    }

    function calculateTimeTotalSeconds() {
        startedTimelogSeconds = 0;
        if (startedTimelogTime !== undefined) {
            startedTimelogSeconds = dayjs().diff(startedTimelogTime, 'seconds');
        }

        $timeTotalSeconds = timelogList.reduce(
            (preVal, curVal) => preVal + curVal.durationSeconds,
            startedTimelogSeconds
        );
    }

    function stopTimer() {
        if (startedTimelogTimer) {
            clearInterval(startedTimelogTimer);
            startedTimelogTime = undefined;
            startedTimelogSeconds = 0;
            $timeTotalClass = '';
        }
    }

    function calculateAmount() {
        if (timelog.timeStart !== '' && timelog.timeEnd !== '' && timelog.timeEnd !== null) {
            const seconds = getDiffSeconds(timelog.timeStart, timelog.timeEnd);
            const hours = hoursFromSeconds(seconds, -1);
            const amount = moneyInteger(timelog.billableRate) * hours;

            timelog._duration = calculateTime(seconds, true);
            timelog._amount = moneyFormatted(amount / 100);
        } else {
            timelog._duration = '';
            timelog._amount = '';
        }
    }
</script>

<!-- svelte-ignore a11y_label_has_associated_control -->
<PageHeader {title} firstColClass="col-auto text-capitalize">
    <div class="col-auto pt2">
        <input type="hidden" id="js-selected-calendar" />
        <button
            type="button"
            onclick={() => selectedCalendar!.show()}
            title={$t('goToSpecificDate')}
            class="btn btn-sm btn-primary"
        >
            <Fa fw icon={faCalendarDay} />
        </button>
        &nbsp;
        <button type="button" onclick={() => (selectedDate = now)} title={$t('today')} class="btn btn-sm btn-primary">
            <Fa fw icon={faHome} />
        </button>
    </div>
    <div class="col pt2 text-end">
        <div class="btn-group" role="group">
            {#if timelogList.length > 0}
                <button type="button" onclick={onShowForm} class="btn btn-sm btn-primary">
                    <Fa icon={faPlus} />
                    <span class="d-none d-sm-inline">{$t('track.add')}</span>
                </button>
            {/if}
            <button type="button" onclick={() => onShowHistory(10)} class="btn btn-sm btn-outline-primary">
                <Fa icon={faHistory} />
            </button>
        </div>
    </div>
</PageHeader>

<div class="d-flex flex-wrap gap-1 mt-3">
    {#each favoriteList as favorite (favorite.id)}
        <button
            type="button"
            onclick={() => onStartFavorite(favorite.id)}
            class="btn btn-sm btn-outline-dark opacity-50"
        >
            {favorite.name}
            <Fa icon={faCirclePlay} />
        </button>
    {/each}
    <button type="button" onclick={() => favoritesModal!.show()} class="btn btn-sm btn-outline-primary">
        <Fa icon={faStar} />
        <span class="d-none d-sm-inline">{$t('favorites.title')}</span>
    </button>
</div>

{#if timelogList.length > 0}
    <table class="table align-middle table-hover mt-3 mb-0">
        <tbody>
            {#each timelogList as timelog (timelog.id)}
                <tr>
                    <td class="align-middle w1">
                        <div class="dropdown">
                            <button
                                class="btn btn-outline-primary btn-sm dropdown-toggle"
                                type="button"
                                data-bs-toggle="dropdown"
                            >
                                <Fa icon={faEllipsis} />
                            </button>
                            <ul class="dropdown-menu">
                                <li>
                                    <button
                                        type="button"
                                        onclick={() => onEditTimelog(timelog.id)}
                                        class="dropdown-item"
                                    >
                                        <Fa fw icon={faPencil} />
                                        {$t('edit')}
                                    </button>
                                </li>
                                <li>
                                    <button
                                        type="button"
                                        onclick={() => onDeleteTimelogModal(timelog.id)}
                                        class="dropdown-item text-danger"
                                    >
                                        <Fa fw icon={faTrash} />
                                        {$t('delete')}
                                    </button>
                                </li>
                            </ul>
                        </div>
                    </td>
                    <td class="align-middle text-start">
                        {timelog.clientName} - {timelog.projectName}
                        {#if timelog.taskId}
                            <span class="badge text-bg-secondary">
                                {#if timelog.taskCompletedAt}
                                    <Fa fw icon={faCircleCheck} class="text-success" />
                                {:else}
                                    <!-- svelte-ignore a11y_invalid_attribute -->
                                    <a
                                        href="javascript:;"
                                        onclick={(e) => {
                                            e.preventDefault();
                                            onCompleteTask(timelog.taskId!);
                                        }}
                                    >
                                        <Fa fw icon={faCircle} title={$t('projects.tasks.toComplete')} />
                                    </a>
                                {/if}
                                {timelog.taskName}
                            </span>
                        {/if}
                        {#if timelog.comment}
                            <br />
                            <i>{timelog.comment}</i>
                        {/if}
                    </td>
                    <td class="align-middle text-end w1 text-nowrap">
                        {#if timelog.timeEnd}
                            <button
                                type="button"
                                onclick={() => onStart(timelog.id)}
                                class="btn btn-sm btn-outline-success btn-play"
                                title={$t('track.toRun')}
                            >
                                <Fa fw icon={faPlay} />
                            </button>
                        {:else}
                            <button
                                type="button"
                                onclick={() => onStop(timelog.id)}
                                class="btn btn-sm btn-outline-danger"
                                title={$t('track.toStop')}
                            >
                                <Fa fw icon={faStop} />
                                {$t('track.stop')}
                            </button>
                        {/if}
                    </td>
                    <td class="align-middle text-end w1 text-nowrap">
                        {#if timelog.timeEnd}
                            <span class="fs-4">{calculateTime(timelog.durationSeconds, true)}</span>
                            <br />
                            <small>
                                {getTime(`${timelog.date} ${timelog.timeStart}`)}
                                -
                                {getTime(`${timelog.date} ${timelog.timeEnd}`)}
                            </small>
                        {:else}
                            <span class="fs-4 text-success">{calculateTime(startedTimelogSeconds, true)}</span>
                        {/if}
                    </td>
                </tr>
            {/each}
        </tbody>
    </table>
{:else}
    <NoContent
        callback={() => {
            onShowForm();
        }}
    />
{/if}

<!-- svelte-ignore a11y_label_has_associated_control -->
<Modal id={timelogModalId} size="lg" title={$t('track.title')}>
    {#snippet body()}
        <div>
            <form
                onsubmit={(e) => {
                    e.preventDefault();
                    onSubmitForm();
                }}
            >
                <div class="mb-3 row">
                    <div class="col">
                        <label class="form-label">{$t('projects.project')}</label>
                        <select
                            id="js-project-select"
                            class="form-select"
                            class:is-invalid={validateErrors.projectId}
                        ></select>
                        {#if validateErrors.projectId}
                            <div class="invalid-feedback">{validateErrors.projectId}</div>
                        {/if}
                    </div>
                    <div class="col-auto">
                        <label class="form-label">{$t('rate')}</label>
                        <input type="text" value={timelog.billableRate} class="form-control" disabled />
                    </div>
                </div>
                <div class="mb-3 row">
                    <div class="col">
                        <label class="form-label">{$t('projects.tasks.task')}</label>
                        <select
                            id="js-task-select"
                            class="form-select"
                            class:is-invalid={validateErrors.taskId}
                        ></select>
                        {#if validateErrors.taskId}
                            <div class="invalid-feedback">{validateErrors.taskId}</div>
                        {/if}
                    </div>
                </div>
                <div class="mb-3 row">
                    <div class="col-auto">
                        <label class="form-label">{$t('date')}</label>
                        <input type="hidden" id="js-form-calendar" />
                        <input
                            type="text"
                            value={timelog._dateFormatted}
                            onclick={() => formCalendar!.show()}
                            class="form-control"
                            class:is-invalid={validateErrors.date}
                        />
                        {#if validateErrors.date}
                            <div class="invalid-feedback">{validateErrors.date}</div>
                        {/if}
                    </div>
                    <div class="col">
                        <div class="clearfix">
                            <label class="form-label float-start">{$t('time')}</label>
                            {#if timelog._duration !== ''}
                                <small class="float-end text-success">
                                    <Fa icon={faClock} />
                                    {timelog._duration}
                                    &nbsp;
                                    <Fa icon={faRubleSign} />
                                    {timelog._amount}
                                </small>
                            {/if}
                        </div>
                        <div class="input-group">
                            <span class="input-group-text">{$t('since')}</span>
                            <input
                                type="text"
                                id="js-time-start"
                                value={timelog._timeStartFormatted}
                                class="form-control"
                                class:is-invalid={validateErrors.timeStart}
                            />
                            <span class="input-group-text">{$t('till')}</span>
                            <input
                                type="text"
                                id="js-time-end"
                                value={timelog._timeEndFormatted}
                                class="form-control"
                            />
                        </div>
                        {#if validateErrors.timeStart}
                            <div class="invalid-feedback">{validateErrors.timeStart}</div>
                        {/if}
                    </div>
                </div>
                <div>
                    <label class="form-label">{$t('comment')}</label>
                    <textarea bind:value={timelog.comment} rows="3" class="form-control"></textarea>
                </div>
            </form>
        </div>
    {/snippet}
    {#snippet footer()}
        <div>
            {#if timelog.id}
                <button type="button" onclick={onSubmitForm} class="btn btn-success">
                    <Fa icon={faSave} />
                    {$t('save')}
                </button>
            {:else}
                <button type="button" onclick={onSubmitForm} class="btn btn-primary">
                    <Fa icon={faPlus} />
                    {$t('add')}
                </button>
            {/if}
            &nbsp;
            <button type="button" data-bs-dismiss="modal" class="btn btn-secondary">{$t('close')}</button>
        </div>
    {/snippet}
</Modal>

<!-- svelte-ignore a11y_label_has_associated_control -->
<Modal id={favoritesModalId} size="lg" title={$t('favorites.title')}>
    {#snippet body()}
        <div>
            <form
                onsubmit={(e) => {
                    e.preventDefault();
                    onSubmitFavoriteForm();
                }}
            >
                <div class="row">
                    <div class="col col-form-label-sm">
                        <div class="mb-2">
                            <label class="form-label">{$t('favorites.name')}</label>
                            <input type="text" bind:value={editingFavorite.name} class="form-control form-control-sm" />
                        </div>
                        <div class="mb-2 row">
                            <div class="col">
                                <label class="form-label">{$t('projects.project')}</label>
                                <select
                                    id="js-project-favorite-select"
                                    class="form-select form-select-sm"
                                    class:is-invalid={validateFavoriteErrors.projectId}
                                ></select>
                                {#if validateFavoriteErrors.projectId}
                                    <div class="invalid-feedback">{validateFavoriteErrors.projectId}</div>
                                {/if}
                            </div>
                            <div class="col-auto">
                                <label class="form-label">{$t('rate')}</label>
                                <input
                                    type="text"
                                    value={editingFavorite.billableRate}
                                    class="form-control form-control-sm"
                                    disabled
                                />
                            </div>
                        </div>
                        <div class="mb-2 row">
                            <div class="col">
                                <label class="form-label">{$t('projects.tasks.task')}</label>
                                <select
                                    id="js-task-favorite-select"
                                    class="form-select form-select-sm"
                                    class:is-invalid={validateFavoriteErrors.taskId}
                                ></select>
                                {#if validateFavoriteErrors.taskId}
                                    <div class="invalid-feedback">{validateFavoriteErrors.taskId}</div>
                                {/if}
                            </div>
                        </div>
                    </div>
                    <div class="col col-form-label-sm">
                        <div class="mb-2">
                            <label class="form-label">{$t('comment')}</label>
                            <textarea
                                bind:value={editingFavorite.comment}
                                rows="6"
                                class="form-control form-control-sm"
                            ></textarea>
                        </div>
                        {#if editingFavorite.id}
                            <button type="submit" onclick={onSubmitFavoriteForm} class="btn btn-sm btn-success">
                                <Fa icon={faSave} />
                                {$t('save')}
                            </button>
                        {:else}
                            <button type="button" onclick={onSubmitFavoriteForm} class="btn btn-sm btn-primary">
                                <Fa icon={faPlus} />
                                {$t('add')}
                            </button>
                        {/if}
                        &nbsp;
                        <button type="button" onclick={onClearFavoriteForm} class="btn btn-sm btn-secondary">
                            {$t('clear')}
                        </button>
                    </div>
                </div>
            </form>
            <hr />
            <table class="table align-middle mt-3">
                <tbody class="border-top">
                    {#each favoriteList as favorite (favorite.id)}
                        <tr>
                            <td>
                                <button type="button" onclick={() => onEditFavorite(favorite.id)} class="btn btn-link">
                                    {favorite.name}
                                </button>
                                <br />
                                {favorite.clientName} - {favorite.projectName}
                                {#if favorite.taskId}
                                    <span class="badge text-bg-secondary">{favorite.taskName}</span>
                                {/if}
                            </td>
                            <td class="d-none d-md-table-cell">
                                {#if favorite.comment}
                                    {favorite.comment}
                                {/if}
                            </td>
                            <td class="w1">
                                <button
                                    type="button"
                                    onclick={() => onDeleteFavorite(favorite.id)}
                                    class="btn btn-link text-danger"
                                >
                                    {$t('delete')}
                                </button>
                            </td>
                        </tr>
                    {/each}
                </tbody>
            </table>
        </div>
    {/snippet}
    {#snippet footer()}
        <div>
            <button type="button" data-bs-dismiss="modal" class="btn btn-secondary">{$t('close')}</button>
        </div>
    {/snippet}
</Modal>

<!-- svelte-ignore a11y_label_has_associated_control -->
<Modal id={historyModalId} size="lg" title={$t('history.title')}>
    {#snippet body()}
        <div>
            <table class="table table-striped table-borderless align-middle m-0">
                <tbody>
                    {#each historyList as history (history.id)}
                        <tr>
                            <td class="align-middle text-end w1 text-nowrap">
                                <button
                                    type="button"
                                    onclick={() => onStartHistory(history.id)}
                                    class="btn btn-sm btn-outline-success"
                                    title={$t('track.toRun')}
                                >
                                    <Fa fw icon={faPlay} />
                                </button>
                            </td>
                            <td>
                                {history.clientName} - {history.projectName}
                                {#if history.taskId}
                                    <span class="badge text-bg-secondary">
                                        {history.taskName}
                                    </span>
                                {/if}
                                {#if history.comment}
                                    <br />
                                    <small><i>{history.comment}</i></small>
                                {:else}
                                    <br />
                                    <small><i>&nbsp;</i></small>
                                {/if}
                            </td>
                        </tr>
                    {/each}
                </tbody>
            </table>
        </div>
    {/snippet}
    {#snippet footer()}
        <div>
            <button type="button" data-bs-dismiss="modal" class="btn btn-secondary">{$t('close')}</button>
        </div>
    {/snippet}
</Modal>

<!-- svelte-ignore a11y_label_has_associated_control -->
<div id={timelogDeleterModalId} class="modal fade" tabindex="-1">
    <div class="modal-dialog modal-dialog-centered">
        <div class="modal-content">
            <div class="modal-header">
                <h5 class="modal-title">{$t('track.deleter.title')}</h5>
                <!-- svelte-ignore a11y_consider_explicit_label -->
                <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
            </div>
            <div class="modal-body">
                {$t('track.deleter.message')}
            </div>
            <div class="modal-footer">
                <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">
                    {$t('track.deleter.no')}
                </button>
                <button type="button" class="btn btn-danger" onclick={() => onDeleteTimelog(timelogDeleterTimelogId!)}>
                    {$t('track.deleter.yes')}
                </button>
            </div>
        </div>
    </div>
</div>