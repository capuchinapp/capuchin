<script lang="ts">
    import {onMount} from 'svelte';
    import type {Snippet} from 'svelte';

    type Props = {
        id: string;
        size: string;
        title: string;
        header?: Snippet;
        body?: Snippet;
        footer?: Snippet;
    };

    let {id, size, title, header, body, footer}: Props = $props();

    let sizeClass = $state('');

    onMount(() => {
        if (size === 'md') {
            sizeClass = '';
        } else {
            sizeClass = 'modal-' + size;
        }
    });
</script>

<div {id} class="modal fade" data-bs-backdrop="static" data-bs-keyboard="false" tabindex="-1">
    <div class="modal-dialog {sizeClass} modal-dialog-centered modal-dialog-scrollable">
        <div class="modal-content">
            <div class="modal-header">
                {#if header}{@render header()}{:else}
                    <h5 class="modal-title">{title}</h5>
                    <!-- svelte-ignore a11y_consider_explicit_label -->
                    <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
                {/if}
            </div>
            <div class="modal-body">
                {@render body?.()}
            </div>
            {#if footer}
                <div class="modal-footer">
                    {@render footer?.()}
                </div>
            {/if}
        </div>
    </div>
</div>