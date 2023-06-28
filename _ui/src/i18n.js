import {register, init, locale} from 'svelte-i18n';
import {get} from 'svelte/store';
import {currentLocale} from './stores';

register('en', () => import('./locales/en.json'));
register('ru', () => import('./locales/ru.json'));

init({
    fallbackLocale: 'en',
    initialLocale: get(currentLocale),
});

locale.set(get(currentLocale));
