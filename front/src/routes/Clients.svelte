<script lang="ts">
    import Fa from 'svelte-fa';
    import {faEllipsis, faPencil, faBoxArchive, faBoxOpen, faPlus, faSave} from '@fortawesome/free-solid-svg-icons';
    import TomSelect from 'tom-select';
    import {Modal as BSModal} from 'bootstrap/dist/js/bootstrap.esm';
    import {onMount} from 'svelte';
    import {t} from '../i18n';
    import PageHeader from '../components/PageHeader.svelte';
    import MoneyValue from '../components/MoneyValue.svelte';
    import NoContent from '../components/NoContent.svelte';
    import NothingFound from '../components/NothingFound.svelte';
    import Modal from '../components/Modal.svelte';
    import ClientAPI from '../services/Clients';
    import {moneyFormat, moneyInteger, moneyFormatted} from '../services/MoneyHelper';
    import {getDate, getTime} from '../services/DateTime';
    import Validation from '../services/Validation';
    import type {Client} from '../types';

    type ValidateErrors = {
        name: string;
        billableRate: string;
    };

    const clientModalId = 'js-client-form';
    const clientTpl = {
        id: '',
        name: '',
        billableRate: '0.00',
        comment: '',
    };

    let statusSelect!: TomSelect;
    let clientModal!: BSModal;
    let validation!: Validation;

    let validateErrors = $state<ValidateErrors>({
        name: '',
        billableRate: '0.00',
    });

    let selectedStatus = $state('active');
    let clientList = $state<Client[]>([]);
    let clientFilteredList = $state<Client[]>([]);
    let client = $state({...clientTpl});

    $effect(() => {
        clientFilteredList = clientList;

        switch (selectedStatus) {
            case 'active':
                clientFilteredList = clientList.filter((c) => !c.archivedAt);
                break;

            case 'archived':
                clientFilteredList = clientList.filter((c) => c.archivedAt);
                break;
        }
    });

    onMount(async () => {
        validation = new Validation();
        clientModal = new BSModal(document.getElementById(clientModalId)!);
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

    async function onSubmitForm() {
        validateErrors = {
            name: '',
            billableRate: '0.00',
        };

        const data = {
            name: client.name,
            billableRate: moneyInteger(client.billableRate),
            comment: client.comment,
        };

        validation.validateClient(
            data,
            async () => {
                if (client.id) {
                    const res = await ClientAPI.update(client.id, data);
                    if (res !== null) {
                        const idx = clientList.findIndex((e) => e.id === client.id);
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
        client = {...clientTpl};
        validateErrors = {
            name: '',
            billableRate: '0.00',
        };
        clientModal.show();
    }

    async function onEditClient(clientId: string) {
        const res = await ClientAPI.getById(clientId);
        if (res !== null) {
            client.id = res.id;
            client.name = res.name;
            client.billableRate = moneyFormatted(res.billableRate / 100);
            client.comment = res.comment;

            clientModal.show();
        }
    }

    async function onArchiveClient(clientId: string) {
        const res = await ClientAPI.archive(clientId);
        if (res !== null) {
            const idx = clientList.findIndex((e) => e.id === clientId);
            if (idx !== -1) {
                clientList[idx] = res;
            }
        }
    }

    async function onUnarchiveClient(clientId: string) {
        const res = await ClientAPI.unarchive(clientId);
        if (res !== null) {
            const idx = clientList.findIndex((e) => e.id === clientId);
            if (idx !== -1) {
                clientList[idx] = res;
            }
        }
    }
</script>

<!-- svelte-ignore a11y_label_has_associated_control -->
<PageHeader title={$t('clients.title')}>
    {#if clientList.length > 0}
        <div class="col-sm-auto pt2 col-status">
            <select id="js-status-select" bind:value={selectedStatus} class="form-select form-select-sm">
                <option value="all">{$t('statuses.all')}</option>
                <option value="active">{$t('statuses.active')}</option>
                <option value="archived">{$t('statuses.archived')}</option>
            </select>
        </div>
        <div class="col-sm-auto pt2">
            <button type="button" onclick={onShowForm} class="btn btn-sm btn-primary">
                <Fa icon={faPlus} />
                <span class="d-none d-sm-inline">{$t('clients.add')}</span>
            </button>
        </div>
    {/if}
</PageHeader>

{#if clientList.length > 0}
    {#if clientFilteredList.length > 0}
        <table class="table align-middle table-hover mt-1">
            <thead>
                <tr class="text-uppercase">
                    <td class="w1">&nbsp;</td>
                    <td>{$t('clients.client')}</td>
                    <td class="w1 text-end text-nowrap">{$t('rate')}</td>
                    <td class="d-none d-md-table-cell">{$t('comment')}</td>
                    <td class="d-none d-md-table-cell w1 text-center text-nowrap">{$t('inArchive')}</td>
                </tr>
            </thead>
            <tbody class="border-top">
                {#each clientFilteredList as client (client.id)}
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
                                            onclick={() => onEditClient(client.id)}
                                            class="dropdown-item"
                                        >
                                            <Fa fw icon={faPencil} />
                                            {$t('edit')}
                                        </button>
                                    </li>
                                    {#if client.archivedAt === null}
                                        <li>
                                            <button
                                                type="button"
                                                onclick={() => onArchiveClient(client.id)}
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
                                                onclick={() => onUnarchiveClient(client.id)}
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
                            {client.name}
                        </td>
                        <td class="text-end text-nowrap">
                            <MoneyValue value={client.billableRate} />
                        </td>
                        <td class="d-none d-md-table-cell">
                            {#if client.comment}
                                {client.comment}
                            {/if}
                        </td>
                        <td class="d-none d-md-table-cell text-center text-nowrap">
                            {#if client.archivedAt}
                                {getDate(client.archivedAt)}
                                <small>{getTime(client.archivedAt)}</small>
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
<Modal id={clientModalId} size="lg" title={$t('clients.client')}>
    {#snippet body()}
        <div>
            <form
                onsubmit={(e) => {
                    e.preventDefault();
                    onSubmitForm();
                }}
            >
                <div class="mb-3">
                    <label class="form-label">{$t('clients.name')}</label>
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
                    <label class="form-label">{$t('rate')}</label>
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
                    <label class="form-label">{$t('comment')}</label>
                    <textarea bind:value={client.comment} rows="3" class="form-control"></textarea>
                </div>
                <div class="text-end">
                    {#if client.id}
                        <button type="submit" class="btn btn-success">
                            <Fa icon={faSave} />
                            {$t('save')}
                        </button>
                    {:else}
                        <button type="submit" class="btn btn-primary">
                            <Fa icon={faPlus} />
                            {$t('add')}
                        </button>
                    {/if}
                    &nbsp;
                    <button type="button" data-bs-dismiss="modal" class="btn btn-secondary">{$t('close')}</button>
                </div>
            </form>
        </div>
    {/snippet}
</Modal>

<style>
    .col-status {
        width: 150px;
    }
</style>