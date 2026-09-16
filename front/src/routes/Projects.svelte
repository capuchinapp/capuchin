<script lang="ts">
    import Fa from 'svelte-fa';
    import {
        faArrowLeft,
        faBoxArchive,
        faBoxOpen,
        faEllipsis,
        faPencil,
        faPlus,
        faSave,
        faListCheck,
        faCircle,
        faCircleCheck,
        faPlay,
        faBusinessTime,
        faRubleSign,
        faCalendarDays,
    } from '@fortawesome/free-solid-svg-icons';
    import {faClock} from '@fortawesome/free-regular-svg-icons';
    import dayjs from 'dayjs';
    import TomSelect from 'tom-select';
    import {Modal as BSModal} from 'bootstrap/dist/js/bootstrap.esm';
    import {onMount} from 'svelte';
    import {push} from 'svelte-spa-router';
    import {t} from '../i18n';
    import PageHeader from '../components/PageHeader.svelte';
    import MoneyValue from '../components/MoneyValue.svelte';
    import NoContent from '../components/NoContent.svelte';
    import NothingFound from '../components/NothingFound.svelte';
    import Modal from '../components/Modal.svelte';
    import ClientAPI from '../services/Clients';
    import ProjectAPI from '../services/Projects';
    import TaskAPI from '../services/Tasks';
    import TimelogAPI from '../services/Timelog';
    import {moneyFormat, moneyInteger, moneyFormatted} from '../services/MoneyHelper';
    import {getDate, getTime, nowTime, hoursFromSeconds, formatDurationSmart} from '../services/DateTime';
    import Validation from '../services/Validation';
    import type {ClientOption, Project, Task, TaskReport} from '../types';

    type ValidateErrors = {
        clientId: string;
        name: string;
        billableRate: string;
    };

    type TaskListItem = Task & {
        hasReport?: boolean;
        durationSeconds?: number;
        billableAmount?: number;
        uniqueDatesCount?: number;
    };

    type TaskForm = {
        id: string;
        projectId: string;
        name: string;
        comment: string;
        archivedAt?: string | null;
        completedAt?: string | null;
    };

    const projectModalId = 'js-project-form';
    const tasksModalId = 'js-tasks-form';
    const clientTpl: ClientOption = {
        id: '',
        name: '',
        billableRate: '0.00',
    };
    const projectTpl = {
        id: '',
        client: {...clientTpl},
        name: '',
        billableRate: '0.00',
        comment: '',
    };

    let clientSelect!: TomSelect;
    let statusSelect!: TomSelect;
    let formClientSelect!: TomSelect;

    let projectModal!: BSModal;
    let tasksModal!: BSModal;
    let validation!: Validation;

    let validateErrors = $state<ValidateErrors>({
        clientId: '',
        name: '',
        billableRate: '0.00',
    });

    let statusSelected = $state('active');
    let clientSelectedId = $state('');
    let projectList = $state<Project[]>([]);
    let projectFilteredList = $state<Project[]>([]);
    let clientList = $state<ClientOption[]>([]);
    let taskList = $state<TaskListItem[]>([]);
    let project = $state({...projectTpl});
    let editingTask = $state<TaskForm | null>(null);
    let projectSelectedForTasks = $state<Project | null>(null);

    $effect(() => {
        projectFilteredList = projectList;

        switch (statusSelected) {
            case 'active':
                projectFilteredList = projectList.filter((c) => !c.archivedAt);
                break;

            case 'archived':
                projectFilteredList = projectList.filter((c) => c.archivedAt);
                break;
        }

        if (clientSelectedId !== '') {
            projectFilteredList = projectFilteredList.filter((c) => c.clientId === clientSelectedId);
        }
    });

    onMount(async () => {
        validation = new Validation();
        projectModal = new BSModal(document.getElementById(projectModalId)!);
        tasksModal = new BSModal(document.getElementById(tasksModalId)!);
        projectList = await ProjectAPI.getList(true);

        await ClientAPI.getList().then(
            (res) =>
                (clientList = res
                    .filter((e) => !e.archivedAt)
                    .map((e) => ({
                        id: e.id,
                        name: e.name,
                        billableRate: moneyFormatted(e.billableRate / 100),
                    })))
        );

        formClientSelect = new TomSelect('#js-form-client-select', {
            options: clientList,
            items: [''],
            valueField: 'id',
            labelField: 'name',
            sortField: 'name',
            searchField: ['name'],
            onChange: (value: string) => (project.client = findClient(String(value))),
            allowEmptyOption: true,
        });

        const _waiter = setInterval(() => {
            if (
                !clientSelect &&
                !statusSelect &&
                document.getElementById('js-client-select') &&
                document.getElementById('js-status-select')
            ) {
                clearInterval(_waiter);

                clientSelect = new TomSelect('#js-client-select', {
                    controlInput: null,
                    allowEmptyOption: true,
                });

                statusSelect = new TomSelect('#js-status-select', {
                    controlInput: null,
                });
            }
        }, 500);
    });

    function onShowForm() {
        project = {...projectTpl};
        validateErrors = {
            clientId: '',
            name: '',
            billableRate: '0.00',
        };
        formClientSelect.setValue('', false);
        projectModal.show();
    }

    function onApplyClientBillableRate() {
        project.billableRate = project.client.billableRate;
    }

    async function onSubmitForm() {
        validateErrors = {
            clientId: '',
            name: '',
            billableRate: '0.00',
        };

        const data = {
            clientId: project.client.id,
            name: project.name,
            billableRate: moneyInteger(project.billableRate),
            comment: project.comment,
        };

        validation.validateProject(
            data,
            async () => {
                if (project.id) {
                    const res = await ProjectAPI.update(project.id, data);
                    if (res !== null) {
                        const idx = projectList.findIndex((e) => e.id === project.id);
                        if (idx !== -1) {
                            projectList[idx] = res;
                            projectModal.hide();
                        }
                    }
                } else {
                    const res = await ProjectAPI.create(data);
                    if (res !== null) {
                        projectList = [res, ...projectList];

                        projectModal.hide();
                    }
                }
            },
            (err) => {
                validateErrors = err;
            }
        );
    }

    async function onEditProject(projectId: string) {
        const res = await ProjectAPI.getById(projectId);
        if (res !== null) {
            let client: ClientOption = {...clientTpl};

            const idx = clientList.findIndex((e) => e.id === res.clientId);
            if (idx !== -1) {
                client = clientList[idx];
            }

            project.id = res.id;
            project.client = client;
            project.name = res.name;
            project.billableRate = moneyFormatted(res.billableRate / 100);
            project.comment = res.comment;

            formClientSelect.setValue(client.id, true);

            projectModal.show();
        }
    }

    async function onArchiveProject(projectId: string) {
        const res = await ProjectAPI.archive(projectId);
        if (res !== null) {
            const idx = projectList.findIndex((e) => e.id === projectId);
            if (idx !== -1) {
                projectList[idx] = res;
            }
        }
    }

    async function onUnarchiveProject(projectId: string) {
        const res = await ProjectAPI.unarchive(projectId);
        if (res !== null) {
            const idx = projectList.findIndex((e) => e.id === projectId);
            if (idx !== -1) {
                projectList[idx] = res;
            }
        }
    }

    async function onTasksByProject(projectId: string) {
        projectSelectedForTasks = null;
        const idx = projectList.findIndex((e) => e.id === projectId);
        if (idx !== -1) {
            projectSelectedForTasks = projectList[idx];
        }

        const res = await TaskAPI.getList(projectId);
        if (res !== null) {
            taskList = res;
            tasksModal.show();
        }
    }

    async function onSubmitTaskForm() {
        if (!editingTask) {
            return;
        }

        if (editingTask.id) {
            const res = await TaskAPI.update(editingTask.id, editingTask);
            if (res !== null) {
                const idx = taskList.findIndex((e) => e.id === editingTask!.id);
                if (idx !== -1) {
                    taskList[idx] = res;
                    editingTask = null;
                }
            }
        } else {
            const res = await TaskAPI.create(editingTask);
            if (res !== null) {
                taskList = [res, ...taskList];
                editingTask = null;
            }
        }
    }

    async function onEditTask(taskId: string) {
        const res = await TaskAPI.getById(taskId);
        if (res !== null) {
            editingTask = res;
        }
    }

    async function onCompleteTask(taskId: string) {
        const res = await TaskAPI.complete(taskId);
        if (res !== null) {
            const idx = taskList.findIndex((e) => e.id === taskId);
            if (idx !== -1) {
                taskList[idx] = res;
            }
        }
    }

    async function onIncompleteTask(taskId: string) {
        const res = await TaskAPI.incomplete(taskId);
        if (res !== null) {
            const idx = taskList.findIndex((e) => e.id === taskId);
            if (idx !== -1) {
                taskList[idx] = res;
            }
        }
    }

    async function onArchiveTask(taskId: string) {
        const res = await TaskAPI.archive(taskId);
        if (res !== null) {
            const idx = taskList.findIndex((e) => e.id === taskId);
            if (idx !== -1) {
                taskList[idx] = res;
                editingTask = null;
            }
        }
    }

    async function onUnarchiveTask(taskId: string) {
        const res = await TaskAPI.unarchive(taskId);
        if (res !== null) {
            const idx = taskList.findIndex((e) => e.id === taskId);
            if (idx !== -1) {
                taskList[idx] = res;
                editingTask = null;
            }
        }
    }

    async function onStartTask(taskId: string) {
        const idxTask = taskList.findIndex((e) => e.id === taskId);
        if (idxTask !== -1) {
            const task = taskList[idxTask];

            const idxProject = projectList.findIndex((e) => e.id === task.projectId);
            if (idxProject !== -1) {
                const res = await TimelogAPI.create({
                    projectId: task.projectId,
                    taskId: task.id,
                    date: dayjs().format('YYYY-MM-DD'),
                    timeStart: nowTime(true),
                    timeEnd: null,
                    billableRate: projectList[idxProject].billableRate,
                    comment: task.comment,
                });
                if (res !== null) {
                    tasksModal.hide();
                    push('/');
                }
            }
        }
    }

    async function onReportTask(taskId: string) {
        const res: TaskReport = await TaskAPI.report(taskId);
        if (res !== null) {
            const idx = taskList.findIndex((e) => e.id === taskId);
            if (idx !== -1) {
                taskList[idx].hasReport = true;
                taskList[idx].durationSeconds = res.durationSeconds;
                taskList[idx].billableAmount = res.billableAmount;
                taskList[idx].uniqueDatesCount = res.uniqueDatesCount;
            }
        }
    }

    function findClient(clientId: string): ClientOption {
        let client: ClientOption = {...clientTpl};

        const idx = clientList.findIndex((e) => e.id === clientId);
        if (idx !== -1) {
            client = clientList[idx];
        }

        return client;
    }
</script>

<!-- svelte-ignore a11y_label_has_associated_control -->
<PageHeader title={$t('projects.title')}>
    {#if projectList.length > 0}
        <div class="col-sm-auto pt2 col-client">
            <select id="js-client-select" bind:value={clientSelectedId} class="form-select form-select-sm">
                <option value="">{$t('clients.all')}</option>
                {#each clientList as client (client.id)}
                    <option value={client.id}>{client.name}</option>
                {/each}
            </select>
        </div>
        <div class="col-sm-auto pt2 col-status">
            <select id="js-status-select" bind:value={statusSelected} class="form-select form-select-sm">
                <option value="all">{$t('statuses.all')}</option>
                <option value="active">{$t('statuses.active')}</option>
                <option value="archived">{$t('statuses.archived')}</option>
            </select>
        </div>
        <div class="col-sm-auto pt2">
            <button type="button" onclick={onShowForm} class="btn btn-sm btn-primary">
                <Fa icon={faPlus} />
                <span class="d-none d-md-inline">{$t('projects.add')}</span>
            </button>
        </div>
    {/if}
</PageHeader>

{#if projectList.length > 0}
    {#if projectFilteredList.length > 0}
        <table class="table align-middle table-hover mt-1">
            <thead>
                <tr class="text-uppercase">
                    <td class="w1">&nbsp;</td>
                    <td>{$t('projects.project')}</td>
                    <td>{$t('clients.client')}</td>
                    <td class="w1 text-end text-nowrap">{$t('rate')}</td>
                    <td class="d-none d-md-table-cell">{$t('comment')}</td>
                    <td class="d-none d-md-table-cell w1 text-center text-nowrap">{$t('inArchive')}</td>
                </tr>
            </thead>
            <tbody class="border-top">
                {#each projectFilteredList as project (project.id)}
                    <tr class:text-muted={project.archivedAt}>
                        <td>
                            <div class="dropdown">
                                <button
                                    class="btn {project.archivedAt
                                        ? 'btn-outline-secondary'
                                        : 'btn-outline-primary'} btn-sm dropdown-toggle"
                                    type="button"
                                    data-bs-toggle="dropdown"
                                >
                                    <Fa icon={faEllipsis} />
                                </button>
                                <ul class="dropdown-menu">
                                    <li>
                                        <button
                                            type="button"
                                            onclick={() => onTasksByProject(project.id)}
                                            class="dropdown-item"
                                        >
                                            <Fa fw icon={faListCheck} />
                                            {$t('projects.tasks.title')}
                                        </button>
                                    </li>
                                    <li>
                                        <button
                                            type="button"
                                            onclick={() => onEditProject(project.id)}
                                            class="dropdown-item"
                                        >
                                            <Fa fw icon={faPencil} />
                                            {$t('edit')}
                                        </button>
                                    </li>
                                    {#if project.archivedAt === null}
                                        <li>
                                            <button
                                                type="button"
                                                onclick={() => onArchiveProject(project.id)}
                                                class="dropdown-item"
                                            >
                                                <Fa fw icon={faBoxArchive} />
                                                {$t('toArchive')}
                                            </button>
                                        </li>
                                    {:else}
                                        <li>
                                            <button
                                                type="button"
                                                onclick={() => onUnarchiveProject(project.id)}
                                                class="dropdown-item text-info"
                                            >
                                                <Fa fw icon={faBoxOpen} />
                                                {$t('return')}
                                            </button>
                                        </li>
                                    {/if}
                                </ul>
                            </div>
                        </td>
                        <td>
                            {project.name}
                        </td>
                        <td>
                            {project.clientName}
                        </td>
                        <td class="text-end text-nowrap">
                            <MoneyValue value={project.billableRate} />
                        </td>
                        <td class="d-none d-md-table-cell">
                            {#if project.comment}
                                {project.comment}
                            {/if}
                        </td>
                        <td class="d-none d-md-table-cell text-center text-nowrap">
                            {#if project.archivedAt}
                                {getDate(project.archivedAt)}
                                <small>{getTime(project.archivedAt)}</small>
                            {:else}
                                {$t('no')}
                            {/if}
                        </td>
                    </tr>
                {/each}
            </tbody>
        </table>
    {:else}
        <NothingFound />
    {/if}
{:else}
    <NoContent
        callback={() => {
            onShowForm();
        }}
    />
{/if}

<!-- svelte-ignore a11y_label_has_associated_control -->
<Modal id={projectModalId} size="lg" title={$t('projects.project')}>
    {#snippet body()}
        <div>
            <form
                onsubmit={(e) => {
                    e.preventDefault();
                    onSubmitForm();
                }}
            >
                <div class="mb-3">
                    <label class="form-label">{$t('clients.client')}</label>
                    <select
                        id="js-form-client-select"
                        class="form-select"
                        class:is-invalid={validateErrors.clientId}
                    ></select>
                    {#if validateErrors.clientId}
                        <div class="invalid-feedback">{validateErrors.clientId}</div>
                    {/if}
                </div>
                <div class="mb-3">
                    <label class="form-label">{$t('projects.name')}</label>
                    <input
                        type="text"
                        bind:value={project.name}
                        class="form-control"
                        class:is-invalid={validateErrors.name}
                    />
                    {#if validateErrors.name}
                        <div class="invalid-feedback">{validateErrors.name}</div>
                    {/if}
                </div>
                <div class="mb-3">
                    <div class="row">
                        <div class="col">
                            <label class="form-label">{$t('projects.rate')}</label>
                            <input
                                type="text"
                                bind:value={project.billableRate}
                                use:moneyFormat
                                class="form-control"
                                class:is-invalid={validateErrors.billableRate}
                            />
                            {#if validateErrors.billableRate}
                                <div class="invalid-feedback">{validateErrors.billableRate}</div>
                            {/if}
                        </div>
                        <div class="col">
                            <label class="form-label">{$t('clients.rate')}</label>
                            <div class="input-group">
                                <button
                                    type="button"
                                    onclick={onApplyClientBillableRate}
                                    title={$t('projects.applyClientBillableRate')}
                                    class="btn btn-secondary"
                                >
                                    <Fa icon={faArrowLeft} />
                                </button>
                                <input type="text" value={project.client.billableRate} class="form-control" disabled />
                            </div>
                        </div>
                    </div>
                </div>
                <div>
                    <label class="form-label">{$t('comment')}</label>
                    <textarea bind:value={project.comment} rows="3" class="form-control"></textarea>
                </div>
            </form>
        </div>
    {/snippet}
    {#snippet footer()}
        <div>
            {#if project.id}
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
<Modal
    id={tasksModalId}
    size="lg"
    title="{$t('projects.tasks.title')}{projectSelectedForTasks !== null ? ' - ' + projectSelectedForTasks.name : ''}"
>
    {#snippet body()}
        <div>
            {#if editingTask && editingTask.id === ''}
                <form
                    onsubmit={(e) => {
                        e.preventDefault();
                        onSubmitTaskForm();
                    }}
                >
                    <div class="mb-2">
                        <label class="form-label"><small>{$t('projects.tasks.name')}</small></label>
                        <input type="text" bind:value={editingTask.name} class="form-control form-control-sm" />
                    </div>
                    <div class="mb-2">
                        <label class="form-label"><small>{$t('projects.tasks.description')}</small></label>
                        <textarea
                            bind:value={editingTask.comment}
                            rows="3"
                            class="form-control form-control-sm"
                        ></textarea>
                    </div>
                    <button type="button" onclick={onSubmitTaskForm} class="btn btn-sm btn-primary">
                        <Fa icon={faPlus} />
                        {$t('add')}
                    </button>
                    &nbsp;
                    <button type="button" onclick={() => (editingTask = null)} class="btn btn-sm btn-secondary">
                        {$t('close')}
                    </button>
                </form>
            {:else}
                <button
                    type="button"
                    onclick={() =>
                        (editingTask = {
                            id: '',
                            projectId: projectSelectedForTasks !== null ? projectSelectedForTasks.id : '',
                            name: '',
                            comment: '',
                        })}
                    class="btn btn-sm btn-primary"
                >
                    <Fa icon={faPlus} />
                    {$t('add')}
                </button>
            {/if}
            <table class="table align-middle mt-3">
                <tbody class="border-top">
                    {#each taskList as task (task.id)}
                        {#if editingTask && editingTask.id === task.id}
                            <!-- Inline Form -->
                            <tr>
                                <td colspan="5">
                                    <form
                                        onsubmit={(e) => {
                                            e.preventDefault();
                                            onSubmitTaskForm();
                                        }}
                                    >
                                        <div class="mb-2">
                                            <label class="form-label"><small>{$t('projects.tasks.name')}</small></label>
                                            <input
                                                type="text"
                                                bind:value={editingTask.name}
                                                class="form-control form-control-sm"
                                            />
                                        </div>
                                        <div class="mb-2">
                                            <label class="form-label">
                                                <small>{$t('projects.tasks.description')}</small>
                                            </label>
                                            <textarea
                                                bind:value={editingTask.comment}
                                                rows="3"
                                                class="form-control form-control-sm"
                                            ></textarea>
                                        </div>
                                        <button type="button" onclick={onSubmitTaskForm} class="btn btn-sm btn-success">
                                            <Fa icon={faSave} />
                                            {$t('save')}
                                        </button>
                                        &nbsp;
                                        <button
                                            type="button"
                                            onclick={() => (editingTask = null)}
                                            class="btn btn-sm btn-secondary"
                                        >
                                            {$t('close')}
                                        </button>
                                        {#if editingTask.archivedAt === null}
                                            <button
                                                type="button"
                                                onclick={() => onArchiveTask(editingTask!.id)}
                                                class="btn btn-sm btn-link text-danger"
                                            >
                                                {$t('toArchive')}
                                            </button>
                                        {:else}
                                            <button
                                                type="button"
                                                onclick={() => onUnarchiveTask(editingTask!.id)}
                                                class="btn btn-sm btn-link text-danger"
                                            >
                                                {$t('return')}
                                            </button>
                                        {/if}
                                    </form>
                                </td>
                            </tr>
                        {:else}
                            <!-- Display Row -->
                            <tr class:text-muted={task.archivedAt}>
                                <td class="w1">
                                    {#if task.archivedAt === null}
                                        {#if task.completedAt === null}
                                            <button
                                                type="button"
                                                onclick={() => onCompleteTask(task.id)}
                                                class="btn btn-sm text-secondary"
                                                title={$t('projects.tasks.toComplete')}
                                            >
                                                <Fa fw icon={faCircle} />
                                            </button>
                                        {:else}
                                            <button
                                                type="button"
                                                onclick={() => onIncompleteTask(task.id)}
                                                class="btn btn-sm text-success"
                                                title={$t('projects.tasks.toIncomplete')}
                                            >
                                                <Fa fw icon={faCircleCheck} />
                                            </button>
                                        {/if}
                                    {/if}
                                </td>
                                <td class="w1">
                                    <button
                                        type="button"
                                        onclick={() => onStartTask(task.id)}
                                        class="btn btn-sm btn-outline-success"
                                        title={$t('track.toRun')}
                                    >
                                        <Fa fw icon={faPlay} />
                                    </button>
                                </td>
                                <td class="w1 text-nowrap text-end">
                                    {#if task.hasReport}
                                        <span title="{$t('hours')}: {hoursFromSeconds(task.durationSeconds!, 1)}">
                                            {formatDurationSmart(task.durationSeconds!)}
                                            <Fa fw icon={faClock} />
                                        </span>
                                        <br />
                                        {moneyFormatted(task.billableAmount! / 100)}
                                        <Fa fw icon={faRubleSign} />
                                        <br />
                                        <span title={$t('projects.tasks.reportDays')}>
                                            {task.uniqueDatesCount!}
                                            <Fa fw icon={faCalendarDays} />
                                        </span>
                                    {:else}
                                        <button
                                            type="button"
                                            onclick={() => onReportTask(task.id)}
                                            class="btn btn-sm btn-outline-info"
                                            title={$t('projects.tasks.report')}
                                        >
                                            <Fa fw icon={faBusinessTime} />
                                        </button>
                                    {/if}
                                </td>
                                <td>
                                    <button
                                        type="button"
                                        onclick={() => onEditTask(task.id)}
                                        class="btn btn-sm btn-link"
                                    >
                                        {task.name}
                                    </button>
                                </td>
                                <td class="d-none d-md-table-cell">
                                    {#if task.comment}
                                        {task.comment}
                                    {/if}
                                </td>
                            </tr>
                        {/if}
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

<style>
    .col-client {
        width: 200px;
    }

    .col-status {
        width: 150px;
    }
</style>