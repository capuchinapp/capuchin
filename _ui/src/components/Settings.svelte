<script>
    import Fa from 'svelte-fa';
    import {faSave} from '@fortawesome/free-solid-svg-icons';
    import TomSelect from 'tom-select';
    import {onMount} from 'svelte';
    import {_} from 'svelte-i18n';
    import SettingsAPI from '../services/Settings';
    import {settings as settingsStore} from '../stores';

    let dateFormat = $settingsStore.dateFormat;

    onMount(async () => {
        new TomSelect('#js-dateformat-select', {
            onChange: (value) => (dateFormat = String(value)),
        });
    });

    async function onSave() {
        SettingsAPI.put({
            dateFormat: dateFormat,
        }).then(() => window.location.reload());
    }
</script>

<!-- svelte-ignore a11y-label-has-associated-control -->
<form on:submit|preventDefault={onSave}>
    <div class="mb-3">
        <label class="form-label">{$_('settings.dateFormat')}</label>
        <select id="js-dateformat-select" bind:value={dateFormat} class="form-select">
            <option value="MM/DD/YYYY">MM/DD/YYYY</option>
            <option value="DD/MM/YYYY">DD/MM/YYYY</option>
            <option value="YYYY-MM-DD">YYYY-MM-DD</option>
            <option value="DD.MM.YYYY">DD.MM.YYYY</option>
            <option value="DD-MM-YYYY">DD-MM-YYYY</option>
            <option value="MM-DD-YYYY">MM-DD-YYYY</option>
        </select>
    </div>
    <button type="button" on:click={onSave} class="btn btn-success">
        <Fa icon={faSave} />
        {$_('save')}
    </button>
</form>
