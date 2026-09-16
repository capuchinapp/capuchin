<script lang="ts">
    import Fa from 'svelte-fa';
    import {faTriangleExclamation} from '@fortawesome/free-solid-svg-icons';
    import {t} from '../i18n';

    type Props = {
        date: string;
    };

    let {date}: Props = $props();
</script>

<div class="alert alert-warning" role="alert">
    <Fa icon={faTriangleExclamation} />
    {date}
    {$t('existRunningTimelog')}
</div>