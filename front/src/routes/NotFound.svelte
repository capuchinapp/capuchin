<script lang="ts">
    import {link} from 'svelte-spa-router';
    import {t} from '../i18n';
    import EmptyPage from '../components/EmptyPage.svelte';
    import imgNotFound from '../assets/images/not_found.png';
</script>

<EmptyPage>
    {#snippet image()}
        <img src={imgNotFound} alt="Page not found" />
    {/snippet}
    {#snippet title()}
        <span>404: {$t('pageNotFound')}</span>
    {/snippet}
    {#snippet footer()}
        <a class="btn btn-primary" href="/" use:link>{$t('backToHome')} →</a>
    {/snippet}
</EmptyPage>