<script>
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
        faMoneyBillAlt,
    } from '@fortawesome/free-solid-svg-icons';
    import {faClock} from '@fortawesome/free-regular-svg-icons';
    import dayjs from 'dayjs';
    import AirDatepicker from 'air-datepicker';
    import airDatepickerLocaleEn from 'air-datepicker/locale/en';
    import airDatepickerLocaleRu from 'air-datepicker/locale/ru';
    import TomSelect from 'tom-select';
    import {Modal as BSModal} from 'bootstrap/dist/js/bootstrap.esm';
    import IMask from 'imask';
    import {onMount} from 'svelte';
    import {_} from 'svelte-i18n';
    import PageHeader from '../components/PageHeader.svelte';
    import NoContent from '../components/NoContent.svelte';
    import Modal from '../components/Modal.svelte';
    import {toastWarning} from '../services/Toast';
    import TimelogAPI from '../services/Timelog';
    import ClientAPI from '../services/Clients';
    import ProjectAPI from '../services/Projects';
    import {moneyInteger, moneyFormatted} from '../services/MoneyHelper';
    import {nowTime, getTime, getDiffSeconds, calculateTime, hoursFromSeconds} from '../services/DateTime';
    import {currentLocale, timeTotalSeconds, timeTotalClass, settings as settingsStore} from '../stores';
    import Validation from '../services/Validation';

    const now = dayjs();
    const timelogModalID = 'js-timelog-form';
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
    const projectTpl = {
        clientUUID: '',
        uuid: '',
        name: '',
        billableRate: '0.00',
    };
    const timelogTpl = {
        uuid: '',
        project: projectTpl,
        date: '',
        timeStart: '',
        timeEnd: null,
        billableRate: '0.00',
        comment: null,
        _dateFormatted: '',
        _timeStartFormatted: '',
        _timeEndFormatted: '',
        _duration: '',
        _amount: '',
    };

    let /** @type {AirDatepicker<HTMLElement>} */ formCalendar;
    let /** @type {AirDatepicker<HTMLElement>} */ selectedCalendar;

    let /** @type {dayjs.Dayjs} */ selectedDate;
    let /** @type {TomSelect} */ projectSelect;
    let /** @type {BSModal} */ timelogModal;
    let /** @type {Validation} */ validation;

    let /** @type {number} */ startedTimelogTimer;
    let /** @type {string} */ startedTimelogUUID;

    let title = '';
    let startedTimelogSeconds = 0;
    let validateErrors = {};
    let timelogList = [];
    let clientList = [];
    let projectList = [];
    let timelog = timelogTpl;

    $: {
        timelog.billableRate = timelog.project.billableRate;
    }

    $: if (selectedDate) {
        selectedCalendar.selectDate(selectedDate.toDate(), {
            updateTime: true,
            silent: true,
        });

        TimelogAPI.getList(selectedDate.format('YYYY-MM-DD'), selectedDate.format('YYYY-MM-DD')).then((result) => {
            timelogList = result;

            let dayName = $_('today');
            if (selectedDate.format('YYYY-MM-DD') !== dayjs().format('YYYY-MM-DD')) {
                dayName = selectedDate.format('ddd');
            }

            title = `${dayName}, ${selectedDate.format('DD MMM')}`;

            calculateTimeTotalSeconds();
        });
    }

    $: {
        const idx = timelogList.findIndex((e) => e.timeEnd === null);
        if (idx !== -1) {
            const record = timelogList[idx];

            if (startedTimelogUUID !== record.uuid) {
                stopTimer();

                startedTimelogSeconds = dayjs().diff(dayjs(`${record.date} ${record.timeStart}`), 'seconds');
                startedTimelogUUID = record.uuid;
                calculateTimeTotalSeconds();

                startedTimelogTimer = setInterval(() => {
                    startedTimelogSeconds++;
                    $timeTotalSeconds++;
                }, 1000);

                $timeTotalClass = 'text-success';
            }
        }
    }

    onMount(async () => {
        await ClientAPI.getList().then((res) => (clientList = res.filter((e) => !e.archivedAt)));
        await ProjectAPI.getList(true).then(
            (res) =>
                (projectList = res
                    .filter((e) => !e.archivedAt)
                    .map((e) => ({
                        clientUUID: e.clientUUID,
                        uuid: e.uuid,
                        name: e.name,
                        billableRate: moneyFormatted(e.billableRate / 100),
                    })))
        );

        validation = new Validation();

        projectSelect = new TomSelect('#js-project-select', {
            options: projectList,
            optgroups: clientList,
            valueField: 'uuid',
            labelField: 'name',
            sortField: 'name',
            searchField: ['name'],
            optgroupField: 'clientUUID',
            optgroupValueField: 'uuid',
            optgroupLabelField: 'name',
            onChange: (value) => (timelog.project = findProject(String(value))),
            allowEmptyOption: true,
            lockOptgroupOrder: true,
        });

        timelogModal = new BSModal(document.getElementById(timelogModalID));

        let inputTimeStart = IMask(document.getElementById('js-time-start'), imaskTime);
        inputTimeStart.on('complete', function () {
            timelog.timeStart = `${inputTimeStart.value}:00`;
            timelog._timeStartFormatted = inputTimeStart.value;
        });

        let inputTimeEnd = IMask(document.getElementById('js-time-end'), imaskTime);
        inputTimeEnd.on('complete', function () {
            timelog.timeEnd = `${inputTimeEnd.value}:00`;
            timelog._timeEndFormatted = inputTimeEnd.value;
            calculateAmount();
        });

        formCalendar = new AirDatepicker('#js-form-calendar', {
            locale: $currentLocale === 'ru' ? airDatepickerLocaleRu : airDatepickerLocaleEn,
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
            locale: $currentLocale === 'ru' ? airDatepickerLocaleRu : airDatepickerLocaleEn,
            dateFormat: 'MMM yyyy',
            autoClose: true,
            isMobile: true,
            toggleSelected: false,
            onSelect: ({date}) => (selectedDate = dayjs(date.toString())),
        });

        selectedDate = now;
    });

    async function onSubmitForm() {
        validateErrors = {};

        const data = {
            projectUUID: timelog.project.uuid,
            date: timelog.date,
            timeStart: timelog.timeStart,
            timeEnd: timelog.timeEnd === '' ? null : timelog.timeEnd,
            billableRate: moneyInteger(timelog.billableRate),
            comment: timelog.comment,
        };

        validation.validateTimelog(
            data,
            async () => {
                if (timelog.uuid) {
                    const res = await TimelogAPI.update(timelog.uuid, data);
                    if (res !== null) {
                        const idx = timelogList.findIndex((e) => e.uuid === timelog.uuid);
                        if (idx !== -1) {
                            timelogList[idx] = res;

                            if (startedTimelogUUID === timelog.uuid) {
                                startedTimelogSeconds = dayjs().diff(
                                    dayjs(`${timelog.date} ${timelog.timeStart}`),
                                    'seconds'
                                );
                            }

                            timelogModal.hide();
                        }
                    }
                } else {
                    const res = await TimelogAPI.create(data);
                    if (res !== null) {
                        timelogList = [res, ...timelogList];

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
                _timeStartFormatted: nowTime(),
            },
        };

        validateErrors = {};

        projectSelect.setValue('', false);

        formCalendar.selectDate(now.toDate(), {
            updateTime: true,
            silent: true,
        });

        timelogModal.show();
    }

    async function onEditTimelog(/** @type {string} */ timelogUUID) {
        const res = await TimelogAPI.getById(timelogUUID);
        if (res !== null) {
            const project = findProject(res.projectUUID);

            if (project.uuid === '') {
                toastWarning($_('track.warnCannotChangeRecordInArchive'), false);

                return;
            }

            const dt = dayjs(res.date);

            timelog.uuid = res.uuid;
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

            projectSelect.setValue(project.uuid, true);

            formCalendar.selectDate(dt.toDate(), {
                updateTime: true,
                silent: true,
            });

            timelogModal.show();
        }
    }

    async function onDeleteTimelog(/** @type {string} */ timelogUUID) {
        const res = await TimelogAPI.delete(timelogUUID);
        if (res !== null) {
            timelogList = timelogList.filter((e) => e.uuid !== timelogUUID);
        }
    }

    async function onStart(/** @type {string} */ timelogUUID) {
        const idx = timelogList.findIndex((e) => e.uuid === timelogUUID);
        if (idx !== -1) {
            const record = timelogList[idx];

            const res = await TimelogAPI.create({
                projectUUID: record.projectUUID,
                date: dayjs().format('YYYY-MM-DD'),
                timeStart: nowTime(true),
                timeEnd: null,
                billableRate: record.billableRate,
                comment: null,
            });
            if (res !== null) {
                timelogList = [res, ...timelogList];

                stopPrevRecord(res.uuid);

                if (selectedDate.format('YYYY-MM-DD') !== now.format('YYYY-MM-DD')) {
                    selectedDate = now;
                }
            }
        }
    }

    async function onStop(/** @type {string} */ timelogUUID) {
        const now = dayjs();
        const res = await TimelogAPI.stop(timelogUUID, now.format('YYYY-MM-DD'), now.format('HH:mm:ss'));
        if (res !== null) {
            const idx = timelogList.findIndex((e) => e.uuid === timelogUUID);
            if (idx !== -1) {
                timelogList[idx] = res;
                stopTimer();
            }
        }
    }

    /**
     * Если есть запущенная запись переполучим ее (она уже остановлена)
     *
     * @param {string} timelogUUID
     */
    async function stopPrevRecord(timelogUUID) {
        const idx = timelogList.findIndex((e) => e.uuid !== timelogUUID && e.timeEnd === null);
        if (idx !== -1) {
            const res = await TimelogAPI.getById(timelogList[idx].uuid);
            if (res !== null) {
                timelogList[idx] = res;
            }
        }
    }

    /**
     * Поиск проекта в списке
     *
     * @param {string} projectUUID
     */
    function findProject(projectUUID) {
        let project = projectTpl;

        const idx = projectList.findIndex((e) => e.uuid === projectUUID);
        if (idx !== -1) {
            project = projectList[idx];
        }

        return project;
    }

    /**
     * Расчёт суммы времени всех логов
     */
    function calculateTimeTotalSeconds() {
        $timeTotalSeconds = timelogList.reduce(
            (/** @type {number} */ preVal, /** @type {Object} */ curVal) => preVal + curVal.durationSeconds,
            startedTimelogSeconds
        );
    }

    /**
     * Остановка таймера
     */
    function stopTimer() {
        if (startedTimelogTimer) {
            clearInterval(startedTimelogTimer);
            startedTimelogSeconds = 0;
            $timeTotalClass = '';
        }
    }

    /**
     * Расчёт диапазона времени и суммы записи
     */
    function calculateAmount() {
        if (timelog.timeStart !== '' && timelog.timeEnd !== '') {
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

<!-- svelte-ignore a11y-label-has-associated-control -->
<PageHeader {title} firstColClass="col-auto text-capitalize">
    <div class="col-auto pt2">
        <input type="hidden" id="js-selected-calendar" />
        <button
            type="button"
            on:click={() => selectedCalendar.show()}
            title={$_('goToSpecificDate')}
            class="btn btn-sm btn-primary"
        >
            <Fa fw icon={faCalendarDay} />
        </button>
        &nbsp;
        <button type="button" on:click={() => (selectedDate = now)} title={$_('today')} class="btn btn-sm btn-primary">
            <Fa fw icon={faHome} />
        </button>
    </div>
    {#if timelogList.length > 0}
        <div class="col pt2 text-end">
            <button type="button" on:click={onShowForm} class="btn btn-sm btn-primary">
                <Fa icon={faPlus} />
                {$_('track.add')}
            </button>
        </div>
    {/if}
</PageHeader>

{#if timelogList.length > 0}
    <table class="table table-hover mt-3">
        <tbody>
            {#each timelogList as timelog (timelog.uuid)}
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
                                        on:click={() => onEditTimelog(timelog.uuid)}
                                        class="dropdown-item"
                                    >
                                        <Fa fw icon={faPencil} />
                                        {$_('edit')}
                                    </button>
                                </li>
                                <li>
                                    <button
                                        type="button"
                                        on:click={() => onDeleteTimelog(timelog.uuid)}
                                        class="dropdown-item text-danger"
                                    >
                                        <Fa fw icon={faTrash} />
                                        {$_('delete')}
                                    </button>
                                </li>
                            </ul>
                        </div>
                    </td>
                    <td class="align-middle text-start">
                        {timelog.clientName} - {timelog.projectName}
                        {#if timelog.comment}
                            <i>{timelog.comment}</i>
                        {/if}
                    </td>
                    <td class="align-middle text-end w1 text-nowrap">
                        {#if timelog.timeEnd}
                            <button
                                type="button"
                                on:click={() => onStart(timelog.uuid)}
                                class="btn btn-sm btn-outline-success btn-play"
                                title={$_('track.toRun')}
                            >
                                <Fa fw icon={faPlay} />
                            </button>
                        {:else}
                            <button
                                type="button"
                                on:click={() => onStop(timelog.uuid)}
                                class="btn btn-sm btn-outline-danger"
                                title={$_('track.toStop')}
                            >
                                <Fa fw icon={faStop} />
                                {$_('track.stop')}
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

<!-- svelte-ignore a11y-label-has-associated-control -->
<Modal id={timelogModalID} size="lg" title={$_('track.title')}>
    <div slot="body">
        <form on:submit|preventDefault={onSubmitForm}>
            <div class="mb-3 row">
                <div class="col">
                    <label class="form-label">{$_('projects.project')}</label>
                    <select id="js-project-select" class="form-select" class:is-invalid={validateErrors.projectUUID} />
                    {#if validateErrors.projectUUID}
                        <div class="invalid-feedback">{validateErrors.projectUUID}</div>
                    {/if}
                </div>
                <div class="col-auto">
                    <label class="form-label">{$_('rate')}</label>
                    <input type="text" value={timelog.billableRate} class="form-control" disabled />
                </div>
            </div>
            <div class="mb-3 row">
                <div class="col-auto">
                    <label class="form-label">{$_('date')}</label>
                    <input type="hidden" id="js-form-calendar" />
                    <input
                        type="text"
                        value={timelog._dateFormatted}
                        on:click={() => formCalendar.show()}
                        class="form-control"
                        class:is-invalid={validateErrors.date}
                    />
                    {#if validateErrors.date}
                        <div class="invalid-feedback">{validateErrors.date}</div>
                    {/if}
                </div>
                <div class="col">
                    <div class="clearfix">
                        <label class="form-label float-start">{$_('time')}</label>
                        {#if timelog._duration !== ''}
                            <small class="float-end text-success">
                                <Fa icon={faClock} />
                                {timelog._duration}
                                &nbsp;
                                <Fa icon={faMoneyBillAlt} />
                                {timelog._amount}
                            </small>
                        {/if}
                    </div>
                    <div class="input-group">
                        <span class="input-group-text">{$_('since')}</span>
                        <input
                            type="text"
                            id="js-time-start"
                            value={timelog._timeStartFormatted}
                            class="form-control"
                            class:is-invalid={validateErrors.timeStart}
                        />
                        <span class="input-group-text">{$_('till')}</span>
                        <input type="text" id="js-time-end" value={timelog._timeEndFormatted} class="form-control" />
                    </div>
                    {#if validateErrors.timeStart}
                        <div class="invalid-feedback">{validateErrors.timeStart}</div>
                    {/if}
                </div>
            </div>
            <div>
                <label class="form-label">{$_('comment')}</label>
                <textarea bind:value={timelog.comment} rows="3" class="form-control" />
            </div>
        </form>
    </div>
    <div slot="footer">
        {#if timelog.uuid}
            <button type="button" on:click={onSubmitForm} class="btn btn-success">
                <Fa icon={faSave} />
                {$_('save')}
            </button>
        {:else}
            <button type="button" on:click={onSubmitForm} class="btn btn-primary">
                <Fa icon={faPlus} />
                {$_('add')}
            </button>
        {/if}
        &nbsp;
        <button type="button" data-bs-dismiss="modal" class="btn btn-secondary">{$_('close')}</button>
    </div>
</Modal>

<style>
    table .btn-play {
        display: none;
    }

    table tr:hover .btn-play {
        display: inline-block;
    }
</style>
