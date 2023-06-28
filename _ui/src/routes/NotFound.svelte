<script>
    import {link} from 'svelte-spa-router';
    import {_} from 'svelte-i18n';
    import EmptyPage from '../components/EmptyPage.svelte';
    import image from '../assets/images/not_found.png';
</script>

<EmptyPage>
    <img slot="image" src={image} alt="Page not found" />
    <span slot="title">404: {$_('pageNotFound')}</span>
    <a slot="footer" class="btn btn-primary" href="/" use:link>{$_('backToHome')} →</a>
</EmptyPage>
