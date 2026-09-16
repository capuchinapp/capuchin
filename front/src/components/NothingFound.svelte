<script lang="ts">
    import {t} from '../i18n';
    import EmptyPage from './EmptyPage.svelte';
    import imgNothingFound from '../assets/images/nothing_found.png';
</script>

<EmptyPage>
    {#snippet image()}
        <img src={imgNothingFound} alt="Nothing found" />
    {/snippet}
    {#snippet title()}
        <span>{$t('nothingFound')}</span>
    {/snippet}
</EmptyPage>