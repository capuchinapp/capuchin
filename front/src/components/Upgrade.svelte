<script lang="ts">
    import Fa from 'svelte-fa';
    import {faLightbulb} from '@fortawesome/free-solid-svg-icons';
    import {t} from '../i18n';
</script>

<div class="alert alert-info" role="alert">
    <Fa icon={faLightbulb} />
    {$t('appUpdated')}
</div>