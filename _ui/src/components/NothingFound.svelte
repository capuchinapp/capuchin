<script>
    import {_} from 'svelte-i18n';
    import EmptyPage from './EmptyPage.svelte';
    import image from '../assets/images/nothing_found.png';
</script>

<EmptyPage>
    <img slot="image" src={image} alt="Nothing found" />
    <span slot="title">{$_('nothingFound')}</span>
</EmptyPage>
