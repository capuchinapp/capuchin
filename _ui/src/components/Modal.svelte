<script>
    import {onMount} from 'svelte';
    import {_} from 'svelte-i18n';

    /**
     * ID of the modal window
     *
     * @type {string}
     */
    export let id;

    /**
     * Modal window size (xl, lg, md, sm)
     *
     * @type {string}
     */
    export let size;

    /**
     * Name of the modal window
     *
     * @type {string}
     */
    export let title;

    let sizeClass;

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
                <slot name="header">
                    <h5 class="modal-title">{title}</h5>
                    <button type="button" class="btn-close" data-bs-dismiss="modal" />
                </slot>
            </div>
            <div class="modal-body">
                <slot name="body" />
            </div>
            {#if $$slots.footer}
                <div class="modal-footer">
                    <slot name="footer" />
                </div>
            {/if}
        </div>
    </div>
</div>
