<script lang="ts">
    import Fa from 'svelte-fa';
    import {faSpinner} from '@fortawesome/free-solid-svg-icons';
    import image from '../assets/images/loading.png';
    import {t} from '../i18n';
</script>

<div class="text-center pb-5">
    <div class="image">
        <img src={image} alt="Loading" />
    </div>
    <h3>
        <span><Fa pulse icon={faSpinner} /> {$t('loading')}</span>
    </h3>
</div>

<style>
    .image :global(img) {
        height: 30vh;
    }
</style>