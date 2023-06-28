<script>
    import {_} from 'svelte-i18n';
    import Fa from 'svelte-fa';
    import {faCogs, faCookie} from '@fortawesome/free-solid-svg-icons';
    import PageHeader from '../components/PageHeader.svelte';
    import Settings from '../components/Settings.svelte';
    import Sessions from '../components/Sessions.svelte';
    import {capuchin} from '../stores';
</script>

<PageHeader title={$_('profile.title')} />

<div class="card mt-3 mb-3">
    <div class="card-header fw-bold">
        <Fa fw icon={faCogs} />
        {$_('settings.title')}
    </div>
    <div class="card-body">
        <Settings />
    </div>
</div>

{#if $capuchin.authMode}
    <div class="card mb-3">
        <div class="card-header fw-bold">
            <Fa fw icon={faCookie} />
            {$_('sessions.title')}
        </div>
        <div class="card-body">
            <Sessions />
        </div>
    </div>
{/if}
