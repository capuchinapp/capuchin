<script lang="ts">
    import type {Snippet} from 'svelte';

    type Props = {
        image?: Snippet;
        title?: Snippet;
        footer?: Snippet;
    };

    let {image, title, footer}: Props = $props();
</script>

<div class="text-center pb-5">
    <div class="image">
        {@render image?.()}
    </div>
    <h3>
        {@render title?.()}
    </h3>
    <div class="mt-3">
        {@render footer?.()}
    </div>
</div>

<style>
    .image :global(img) {
        height: 30vh;
    }
</style>