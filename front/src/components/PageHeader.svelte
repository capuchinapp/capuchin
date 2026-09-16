<script lang="ts">
    import type {Snippet} from 'svelte';

    type Props = {
        title: string;
        firstColClass?: string;
        children?: Snippet;
    };

    let {title, firstColClass = 'col', children}: Props = $props();
</script>

<div class="row">
    <div class={firstColClass}>
        <h3 class="m-0">{title}</h3>
    </div>
    {@render children?.()}
</div>