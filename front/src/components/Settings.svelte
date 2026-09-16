<script lang="ts">
    import Fa from 'svelte-fa';
    import {faSave} from '@fortawesome/free-solid-svg-icons';
    import TomSelect from 'tom-select';
    import {onMount} from 'svelte';
    import {t} from '../i18n';
    import SettingsAPI from '../services/Settings';
    import {settings as settingsStore} from '../stores';

    let dateFormat = $state($settingsStore.dateFormat);
    let workingDays = $state(SettingsAPI.workingDaysFromString($settingsStore.workingDays));
    let daysOfWeek = ['0', '1', '2', '3', '4', '5', '6'];

    onMount(async () => {
        new TomSelect('#js-dateformat-select', {
            onChange: (value: string) => (dateFormat = String(value)),
        });
    });

    async function onSave() {
        SettingsAPI.put({
            dateFormat: dateFormat,
            workingDays: SettingsAPI.workingDaysToString(workingDays),
        }).then(() => window.location.reload());
    }
</script>

<!-- svelte-ignore a11y_label_has_associated_control -->
<form
    onsubmit={(e) => {
        e.preventDefault();
        onSave();
    }}
>
    <div class="mb-3">
        <label class="form-label">{$t('settings.dateFormat')}</label>
        <select id="js-dateformat-select" bind:value={dateFormat} class="form-select">
            <option value="MM/DD/YYYY">MM/DD/YYYY</option>
            <option value="DD/MM/YYYY">DD/MM/YYYY</option>
            <option value="YYYY-MM-DD">YYYY-MM-DD</option>
            <option value="DD.MM.YYYY">DD.MM.YYYY</option>
            <option value="DD-MM-YYYY">DD-MM-YYYY</option>
            <option value="MM-DD-YYYY">MM-DD-YYYY</option>
        </select>
    </div>
    <div class="mb-3">
        <label class="form-label">{$t('settings.workingDays.title')}</label>
        {#each daysOfWeek as day}
            <div class="form-check">
                <input
                    class="form-check-input"
                    type="checkbox"
                    name="workingDays"
                    id="d{day}"
                    value={day}
                    bind:group={workingDays}
                    checked={workingDays.includes(day)}
                />
                <label class="form-check-label" for="d{day}">{$t('settings.workingDays.names.d' + day)}</label>
            </div>
        {/each}
    </div>
    <button type="button" onclick={onSave} class="btn btn-success">
        <Fa icon={faSave} />
        {$t('save')}
    </button>
</form>