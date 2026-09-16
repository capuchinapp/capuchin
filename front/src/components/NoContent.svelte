<script lang="ts">
    import Fa from 'svelte-fa';
    import {faPlus} from '@fortawesome/free-solid-svg-icons';
    import {t} from '../i18n';
    import EmptyPage from './EmptyPage.svelte';
    import imgInsertRecord from '../assets/images/insert_record.png';

    type Props = {
        callback: () => void;
    };

    let {callback}: Props = $props();
</script>

<EmptyPage>
    {#snippet image()}
        <img src={imgInsertRecord} alt="No content" />
    {/snippet}
    {#snippet title()}
        <span>{$t('noContent')}</span>
    {/snippet}
    {#snippet footer()}
        <button onclick={callback} type="button" class="btn btn-primary">
            <Fa icon={faPlus} />
            {$t('addFirstEntry')}
        </button>
    {/snippet}
</EmptyPage>