import {get} from 'svelte/store';
import {capuchin} from './stores';
import App from './App.svelte';
import IndexAPI from './services/Index';

let /** @type {App} */ app;

IndexAPI.get().then(
    (res) => {
        get(capuchin)['authMode'] = res.authMode;
        get(capuchin)['isAuth'] = res.isAuth;

        app = new App({
            target: document.getElementById('app'),
        });
    },
    (err) => {
        console.error('An error occurred while getting the index', err);
    }
);

export default app;
