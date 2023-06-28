import {toast} from '@zerodevx/svelte-toast';

export const toastSuccess = (/** @type {string} */ content, /** @type {boolean} */ autohide) =>
    show('toast-success', content, autohide);

export const toastDanger = (/** @type {string} */ content, /** @type {boolean} */ autohide) =>
    show('toast-danger', content, autohide);

export const toastWarning = (/** @type {string} */ content, /** @type {boolean} */ autohide) =>
    show('toast-warning', content, autohide);

export const toastInfo = (/** @type {string} */ content, /** @type {boolean} */ autohide) =>
    show('toast-info', content, autohide);

const show = (/** @type {string} */ className, /** @type {string} */ content, /** @type {boolean} */ autohide) => {
    toast.push(content, {
        initial: autohide ? 1 : 0,
        classes: ['toast-item', className],
    });
};
