<script>
    import Fa from 'svelte-fa';
    import {faEllipsis, faPencil, faBoxArchive, faBoxOpen, faPlus, faSave} from '@fortawesome/free-solid-svg-icons';
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
    import {moneyFormat, moneyInteger, moneyFormatted} from '../services/MoneyHelper';
    import {getDate, getTime} from '../services/DateTime';
    import Validation from '../services/Validation';

    const clientModalID = 'js-client-form';
    const clientTpl = {
        uuid: '',
        name: '',
        billableRate: '0.00',
        comment: '',
    };

    let /** @type {TomSelect} */ statusSelect;
    let /** @type {BSModal} */ clientModal;
    let /** @type {Validation} */ validation;

    let selectedStatus = 'active';
    let validateErrors = {};
    let clientList = [];
    let clientFilteredList = [];
    let client = clientTpl;

    onMount(async () => {
        validation = new Validation();
        clientModal = new BSModal(document.getElementById(clientModalID));
        clientList = await ClientAPI.getList();

        const _waiter = setInterval(() => {
            if (!statusSelect && document.getElementById('js-status-select')) {
                clearInterval(_waiter);

                statusSelect = new TomSelect('#js-status-select', {
                    controlInput: null,
                });
            }
        }, 500);
    });

    $: {
        clientFilteredList = clientList;

        switch (selectedStatus) {
            case 'active':
                clientFilteredList = clientList.filter((c) => !c.archivedAt);
                break;

            case 'archived':
                clientFilteredList = clientList.filter((c) => c.archivedAt);
                break;
        }
    }

    async function onSubmitForm() {
        validateErrors = {};

        const data = {
            name: client.name,
            billableRate: moneyInteger(client.billableRate),
            comment: client.comment,
        };

        validation.validateClient(
            data,
            async () => {
                if (client.uuid) {
                    const res = await ClientAPI.update(client.uuid, data);
                    if (res !== null) {
                        const idx = clientList.findIndex((e) => e.uuid === client.uuid);
                        if (idx !== -1) {
                            clientList[idx] = res;
                            clientModal.hide();
                        }
                    }
                } else {
                    const res = await ClientAPI.create(data);
                    if (res !== null) {
                        clientList = [res, ...clientList];

                        clientModal.hide();
                    }
                }
            },
            (err) => {
                validateErrors = err;
            }
        );
    }

    function onShowForm() {
        client = clientTpl;
        validateErrors = {};
        clientModal.show();
    }

    async function onEditClient(/** @type {string} */ clientUUID) {
        const res = await ClientAPI.getById(clientUUID);
        if (res !== null) {
            client.uuid = res.uuid;
            client.name = res.name;
            client.billableRate = moneyFormatted(res.billableRate / 100);
            client.comment = res.comment;

            clientModal.show();
        }
    }

    async function onArchiveClient(/** @type {string} */ clientUUID) {
        const res = await ClientAPI.archive(clientUUID);
        if (res !== null) {
            const idx = clientList.findIndex((e) => e.uuid === clientUUID);
            if (idx !== -1) {
                clientList[idx] = res;
            }
        }
    }

    async function onUnarchiveClient(/** @type {string} */ clientUUID) {
        const res = await ClientAPI.unarchive(clientUUID);
        if (res !== null) {
            const idx = clientList.findIndex((e) => e.uuid === clientUUID);
            if (idx !== -1) {
                clientList[idx] = res;
            }
        }
    }
</script>

<!-- svelte-ignore a11y-label-has-associated-control -->
<PageHeader title={$_('clients.title')}>
    {#if clientList.length > 0}
        <div class="col-auto pt2 col-status">
            <select id="js-status-select" bind:value={selectedStatus} class="form-select form-select-sm">
                <option value="all">{$_('statuses.all')}</option>
                <option value="active">{$_('statuses.active')}</option>
                <option value="archived">{$_('statuses.archived')}</option>
            </select>
        </div>
        <div class="col-auto pt2">
            <button type="button" on:click={onShowForm} class="btn btn-sm btn-primary">
                <Fa icon={faPlus} />
                {$_('clients.add')}
            </button>
        </div>
    {/if}
</PageHeader>

{#if clientList.length > 0}
    {#if clientFilteredList.length > 0}
        <table class="table table-hover mt-1">
            <thead>
                <tr class="text-uppercase">
                    <td class="w1">&nbsp;</td>
                    <td>{$_('clients.client')}</td>
                    <td class="w1 text-end text-nowrap">{$_('rate')}</td>
                    <td>{$_('comment')}</td>
                    <td class="w1 text-center text-nowrap">{$_('inArchive')}</td>
                </tr>
            </thead>
            <tbody class="border-top">
                {#each clientFilteredList as client (client.uuid)}
                    <tr class:text-muted={client.archivedAt}>
                        <td>
                            <div class="dropdown">
                                <button
                                    class="btn {client.archivedAt
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
                                            on:click={() => onEditClient(client.uuid)}
                                            class="dropdown-item"
                                        >
                                            <Fa fw icon={faPencil} />
                                            {$_('edit')}
                                        </button>
                                    </li>
                                    {#if client.archivedAt === null}
                                        <li>
                                            <button
                                                type="button"
                                                on:click={() => onArchiveClient(client.uuid)}
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
                                                on:click={() => onUnarchiveClient(client.uuid)}
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
                            {client.name}
                        </td>
                        <td class="text-end text-nowrap">
                            <MoneyValue value={client.billableRate} />
                        </td>
                        <td>
                            {#if client.comment}
                                {client.comment}
                            {/if}
                        </td>
                        <td class="text-center text-nowrap">
                            {#if client.archivedAt}
                                {getDate(client.archivedAt)}
                                <small>{getTime(client.archivedAt)}</small>
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
<Modal id={clientModalID} size="lg" title={$_('clients.client')}>
    <div slot="body">
        <form on:submit|preventDefault={onSubmitForm}>
            <div class="mb-3">
                <label class="form-label">{$_('clients.name')}</label>
                <input
                    type="text"
                    bind:value={client.name}
                    class="form-control"
                    class:is-invalid={validateErrors.name}
                />
                {#if validateErrors.name}
                    <div class="invalid-feedback">{validateErrors.name}</div>
                {/if}
            </div>
            <div class="mb-3">
                <label class="form-label">{$_('rate')}</label>
                <input
                    type="text"
                    bind:value={client.billableRate}
                    use:moneyFormat
                    class="form-control"
                    class:is-invalid={validateErrors.billableRate}
                />
                {#if validateErrors.billableRate}
                    <div class="invalid-feedback">{validateErrors.billableRate}</div>
                {/if}
            </div>
            <div class="mb-3">
                <label class="form-label">{$_('comment')}</label>
                <textarea bind:value={client.comment} rows="3" class="form-control" />
            </div>
            <div class="text-end">
                {#if client.uuid}
                    <button type="submit" class="btn btn-success">
                        <Fa icon={faSave} />
                        {$_('save')}
                    </button>
                {:else}
                    <button type="submit" class="btn btn-primary">
                        <Fa icon={faPlus} />
                        {$_('add')}
                    </button>
                {/if}
                &nbsp;
                <button type="button" data-bs-dismiss="modal" class="btn btn-secondary">{$_('close')}</button>
            </div>
        </form>
    </div>
</Modal>

<style>
    .col-status {
        width: 150px;
    }
</style>
