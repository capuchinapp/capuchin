import {writable} from 'svelte/store';
import {getLocaleFromNavigator} from 'svelte-i18n';

const locale = getLocaleFromNavigator();

export let currentLocale = writable(locale);
export let preloaderCount = writable(0);
export let toasts = writable([]);
export let timeTotalSeconds = writable(0);
export let timeTotalClass = writable('');
export let capuchin = writable({
    authMode: false,
    isAuth: false,
});
export let settings = writable({
    dateFormat: locale === 'ru' ? 'DD.MM.YYYY' : 'MM/DD/YYYY',
});
