<script>
    import Fa from 'svelte-fa';
    import {
        faArrowLeft,
        faBoxArchive,
        faBoxOpen,
        faEllipsis,
        faPencil,
        faPlus,
        faSave,
    } from '@fortawesome/free-solid-svg-icons';
    import TomSelect from 'tom-select';
    import {Modal as BSModal} from 'bootstrap/dist/js/bootstrap.esm';
    import {onMount} from 'svelte';
    import {_} from 'svelte-i18n';
    import PageHeader from '../components/PageHeader.svelte';
    import MoneyValue from '../components/MoneyValue.svelte';
    import NoContent from '../components/NoContent.svelte';
    import NothingFound from '../components/NothingFound.svelte';
    import Modal from '../components/Modal.svelte';
    import ClientAPI from '../services/Clients';
    import ProjectAPI from '../services/Projects';
    import {moneyFormat, moneyInteger, moneyFormatted} from '../services/MoneyHelper';
    import {getDate, getTime} from '../services/DateTime';
    import Validation from '../services/Validation';

    const projectModalID = 'js-project-form';
    const clientTpl = {
        uuid: '',
        name: '',
        billableRate: '0.00',
    };
    const projectTpl = {
        uuid: '',
        client: clientTpl,
        name: '',
        billableRate: '0.00',
        comment: '',
    };

    let /** @type {TomSelect} */ clientSelect;
    let /** @type {TomSelect} */ statusSelect;
    let /** @type {TomSelect} */ formClientSelect;

    let /** @type {BSModal} */ projectModal;
    let /** @type {Validation} */ validation;

    let statusSelected = 'active';
    let validateErrors = {};
    let clientSelectedUUID = '';
    let projectList = [];
    let projectFilteredList = [];
    let clientList = [];
    let project = projectTpl;

    onMount(async () => {
        validation = new Validation();
        projectModal = new BSModal(document.getElementById(projectModalID));
        projectList = await ProjectAPI.getList(true);

        await ClientAPI.getList().then(
            (res) =>
                (clientList = res
                    .filter((e) => !e.archivedAt)
                    .map((e) => ({
                        uuid: e.uuid,
                        name: e.name,
                        billableRate: moneyFormatted(e.billableRate / 100),
                    })))
        );

        formClientSelect = new TomSelect('#js-form-client-select', {
            options: clientList,
            items: [''],
            valueField: 'uuid',
            labelField: 'name',
            sortField: 'name',
            searchField: ['name'],
            onChange: (value) => (project.client = findClient(String(value))),
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

    $: {
        projectFilteredList = projectList;

        switch (statusSelected) {
            case 'active':
                projectFilteredList = projectList.filter((c) => !c.archivedAt);
                break;

            case 'archived':
                projectFilteredList = projectList.filter((c) => c.archivedAt);
                break;
        }

        if (clientSelectedUUID !== '') {
            projectFilteredList = projectFilteredList.filter((c) => c.clientUUID === clientSelectedUUID);
        }
    }

    function onShowForm() {
        project = projectTpl;
        validateErrors = {};
        formClientSelect.setValue('', false);
        projectModal.show();
    }

    function onApplyClientBillableRate() {
        project.billableRate = project.client.billableRate;
    }

    async function onSubmitForm() {
        validateErrors = {};

        const data = {
            clientUUID: project.client.uuid,
            name: project.name,
            billableRate: moneyInteger(project.billableRate),
            comment: project.comment,
        };

        validation.validateProject(
            data,
            async () => {
                if (project.uuid) {
                    const res = await ProjectAPI.update(project.uuid, data);
                    if (res !== null) {
                        const idx = projectList.findIndex((e) => e.uuid === project.uuid);
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

    async function onEditProject(/** @type {string} */ projectUUID) {
        const res = await ProjectAPI.getById(projectUUID);
        if (res !== null) {
            let client = clientTpl;

            const idx = clientList.findIndex((e) => e.uuid === res.clientUUID);
            if (idx !== -1) {
                client = clientList[idx];
            }

            project.uuid = res.uuid;
            project.client = client;
            project.name = res.name;
            project.billableRate = moneyFormatted(res.billableRate / 100);
            project.comment = res.comment;

            formClientSelect.setValue(client.uuid, true);

            projectModal.show();
        }
    }

    async function onArchiveProject(/** @type {string} */ projectUUID) {
        const res = await ProjectAPI.archive(projectUUID);
        if (res !== null) {
            const idx = projectList.findIndex((e) => e.uuid === projectUUID);
            if (idx !== -1) {
                projectList[idx] = res;
            }
        }
    }

    async function onUnarchiveProject(/** @type {string} */ projectUUID) {
        const res = await ProjectAPI.unarchive(projectUUID);
        if (res !== null) {
            const idx = projectList.findIndex((e) => e.uuid === projectUUID);
            if (idx !== -1) {
                projectList[idx] = res;
            }
        }
    }

    /**
     * Поиск клиента в списке
     *
     * @param {string} clientUUID
     */
    function findClient(clientUUID) {
        let client = clientTpl;

        const idx = clientList.findIndex((e) => e.uuid === clientUUID);
        if (idx !== -1) {
            client = clientList[idx];
        }

        return client;
    }
</script>

<!-- svelte-ignore a11y-label-has-associated-control -->
<PageHeader title={$_('projects.title')}>
    {#if projectList.length > 0}
        <div class="col-auto pt2 col-client">
            <select id="js-client-select" bind:value={clientSelectedUUID} class="form-select form-select-sm">
                <option value="">{$_('clients.all')}</option>
                {#each clientList as client (client.uuid)}
                    <option value={client.uuid}>{client.name}</option>
                {/each}
            </select>
        </div>
        <div class="col-auto pt2 col-status">
            <select id="js-status-select" bind:value={statusSelected} class="form-select form-select-sm">
                <option value="all">{$_('statuses.all')}</option>
                <option value="active">{$_('statuses.active')}</option>
                <option value="archived">{$_('statuses.archived')}</option>
            </select>
        </div>
        <div class="col-auto pt2">
            <button type="button" on:click={onShowForm} class="btn btn-sm btn-primary">
                <Fa icon={faPlus} />
                {$_('projects.add')}
            </button>
        </div>
    {/if}
</PageHeader>

{#if projectList.length > 0}
    {#if projectFilteredList.length > 0}
        <table class="table table-hover mt-1">
            <thead>
                <tr class="text-uppercase">
                    <td class="w1">&nbsp;</td>
                    <td>{$_('clients.client')}</td>
                    <td>{$_('projects.project')}</td>
                    <td class="w1 text-end text-nowrap">{$_('rate')}</td>
                    <td>{$_('comment')}</td>
                    <td class="w1 text-center text-nowrap">{$_('inArchive')}</td>
                </tr>
            </thead>
            <tbody class="border-top">
                {#each projectFilteredList as project (project.uuid)}
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
                                            on:click={() => onEditProject(project.uuid)}
                                            class="dropdown-item"
                                        >
                                            <Fa fw icon={faPencil} />
                                            {$_('edit')}
                                        </button>
                                    </li>
                                    {#if project.archivedAt === null}
                                        <li>
                                            <button
                                                type="button"
                                                on:click={() => onArchiveProject(project.uuid)}
                                                class="dropdown-item"
                                            >
                                                <Fa fw icon={faBoxArchive} />
                                                {$_('toArchive')}
                                            </button>
                                        </li>
                                    {:else}
                                        <li>
                                            <button
                                                type="button"
                                                on:click={() => onUnarchiveProject(project.uuid)}
                                                class="dropdown-item text-info"
                                            >
                                                <Fa fw icon={faBoxOpen} />
                                                {$_('return')}
                                            </button>
                                        </li>
                                    {/if}
                                </ul>
                            </div>
                        </td>
                        <td>
                            {project.clientName}
                        </td>
                        <td>
                            {project.name}
                        </td>
                        <td class="text-end text-nowrap">
                            <MoneyValue value={project.billableRate} />
                        </td>
                        <td>
                            {#if project.comment}
                                {project.comment}
                            {/if}
                        </td>
                        <td class="text-center text-nowrap">
                            {#if project.archivedAt}
                                {getDate(project.archivedAt)}
                                <small>{getTime(project.archivedAt)}</small>
                            {:else}
                                {$_('no')}
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

<!-- svelte-ignore a11y-label-has-associated-control -->
<Modal id={projectModalID} size="lg" title={$_('projects.project')}>
    <div slot="body">
        <form on:submit|preventDefault={onSubmitForm}>
            <div class="mb-3">
                <label class="form-label">{$_('clients.client')}</label>
                <select id="js-form-client-select" class="form-select" class:is-invalid={validateErrors.clientUUID} />
                {#if validateErrors.clientUUID}
                    <div class="invalid-feedback">{validateErrors.clientUUID}</div>
                {/if}
            </div>
            <div class="mb-3">
                <label class="form-label">{$_('projects.name')}</label>
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
                        <label class="form-label">{$_('projects.rate')}</label>
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
                        <label class="form-label">{$_('clients.rate')}</label>
                        <div class="input-group">
                            <button
                                type="button"
                                on:click={onApplyClientBillableRate}
                                title={$_('projects.applyClientBillableRate')}
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
                <label class="form-label">{$_('comment')}</label>
                <textarea bind:value={project.comment} rows="3" class="form-control" />
            </div>
        </form>
    </div>
    <div slot="footer">
        {#if project.uuid}
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
    .col-client {
        width: 200px;
    }

    .col-status {
        width: 150px;
    }
</style>
