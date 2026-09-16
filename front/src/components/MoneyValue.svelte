<script lang="ts">
    import Fa from 'svelte-fa';
    import {faRuble} from '@fortawesome/free-solid-svg-icons';
    import {moneyFormatted} from '../services/MoneyHelper';
    import type {Snippet} from 'svelte';

    type Props = {
        value: number;
        children?: Snippet;
    };

    let {value, children}: Props = $props();
</script>

{#if children}{@render children()}{:else}
    {moneyFormatted(value / 100)}
    <Fa icon={faRuble} />
{/if}