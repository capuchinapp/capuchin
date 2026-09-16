<script lang="ts">
    import Fa from 'svelte-fa';
    import {faCogs, faCookie} from '@fortawesome/free-solid-svg-icons';
    import {t} from '../i18n';
    import PageHeader from '../components/PageHeader.svelte';
    import Settings from '../components/Settings.svelte';
    import Sessions from '../components/Sessions.svelte';
</script>

<PageHeader title={$t('profile.title')} />

<div class="card mt-3 mb-3">
    <div class="card-header fw-bold">
        <Fa fw icon={faCogs} />
        {$t('settings.title')}
    </div>
    <div class="card-body">
        <Settings />
    </div>
</div>

<div class="card">
    <div class="card-header fw-bold">
        <Fa fw icon={faCookie} />
        {$t('sessions.title')}
    </div>
    <div class="card-body">
        <Sessions />
    </div>
</div>