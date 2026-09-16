<script lang="ts">
    import Fa from 'svelte-fa';
    import {faCheckCircle, faEllipsis, faTrash} from '@fortawesome/free-solid-svg-icons';
    import {faCircle} from '@fortawesome/free-regular-svg-icons';
    import {onMount} from 'svelte';
    import {t} from '../i18n';
    import SessionsAPI from '../services/Sessions';
    import {getDate, getTime} from '../services/DateTime';
    import type {Session} from '../types';

    let sessionList = $state<Session[]>([]);

    onMount(() => {
        SessionsAPI.get().then((res) => {
            sessionList = res;
        });
    });

    function onDelete(sessionId: string) {
        SessionsAPI.delete(sessionId).then(() => {
            sessionList = sessionList.filter((e) => e.id !== sessionId);
        });
    }
</script>

<table class="table align-middle table-hover mb-0">
    <tbody>
        {#each sessionList as session (session.id)}
            <tr>
                <td class="align-middle w1">
                    {#if !session.isCurrent}
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
                                        onclick={() => onDelete(session.id)}
                                        class="dropdown-item text-danger"
                                    >
                                        <Fa fw icon={faTrash} />
                                        {$t('delete')}
                                    </button>
                                </li>
                            </ul>
                        </div>
                    {/if}
                </td>
                <td class="align-middle text-start">
                    {#if session.isCurrent}
                        <span class="text-success" title="Активная сессия"><Fa fw icon={faCheckCircle} /></span>
                    {:else}
                        <Fa fw icon={faCircle} />
                    {/if}

                    {getDate(session.checkedAt)}
                    <small>{getTime(session.checkedAt)}</small>
                </td>
            </tr>
        {/each}
    </tbody>
</table>