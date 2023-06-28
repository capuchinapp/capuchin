<script>
    import {_} from 'svelte-i18n';
    import Fa from 'svelte-fa';
    import {faPlus} from '@fortawesome/free-solid-svg-icons';
    import EmptyPage from './EmptyPage.svelte';
    import image from '../assets/images/insert_record.png';

    /**
     * Функция на кнопку в подвале
     *
     * @type {any}
     */
    export let callback;
</script>

<EmptyPage>
    <img slot="image" src={image} alt="No content" />
    <span slot="title">{$_('noContent')}</span>
    <button slot="footer" on:click={callback} type="button" class="btn btn-primary">
        <Fa icon={faPlus} /> {$_('addFirstEntry')}
    </button>
</EmptyPage>
