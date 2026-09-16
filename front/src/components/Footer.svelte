<script lang="ts">
    import {Modal as BSModal} from 'bootstrap/dist/js/bootstrap.esm';
    import {marked} from 'marked';
    import {onMount} from 'svelte';
    import {t} from '../i18n';
    import {capuchin} from './../stores';
    import Modal from './Modal.svelte';
    import changelogRuUrl from './../assets/changelog.ru.md?url';

    const changelogModalId = 'js-changelog-form';

    let changelogModal!: BSModal;

    let changelog = $state('');

    onMount(async () => {
        changelogModal = new BSModal(document.getElementById(changelogModalId)!);
    });

    async function onChangelogShow() {
        try {
            const response = await fetch(changelogRuUrl);
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }

            changelog = await marked.parse(await response.text());
            changelogModal.show();
        } catch (error) {
            console.error('There was a problem with the fetch operation:', error);
        }
    }
</script>

<div class="row">
    <div class="col-md-4 d-flex align-items-center">
        <!-- svelte-ignore a11y_invalid_attribute -->
        &copy; DimNS&nbsp;
        <a
            href="javascript:;"
            onclick={(e) => {
                e.preventDefault();
                onChangelogShow();
            }}
        >
            {$capuchin.appVersionFront}
        </a>
    </div>
    <div class="col-md-4 d-md-flex align-items-center justify-content-md-center">
        &nbsp;
    </div>
    <div class="col-md-4 text-md-end">
        <a href="https://www.flaticon.com/free-icons/amazon" title="amazon icons" class="attribution">
            Amazon icons created by Freepik - Flaticon
        </a>
    </div>
</div>

<Modal id={changelogModalId} size="lg" title={$t('changelogTitle')}>
    {#snippet body()}
        <div>
            {@html changelog}
        </div>
    {/snippet}
</Modal>

<style>
    .attribution {
        font-size: 0.8em;
    }
</style>