import {push} from 'svelte-spa-router';
import {wrap} from 'svelte-spa-router/wrap';
import {get as getStore} from 'svelte/store';
import Login from './routes/Login.svelte';
import NotFound from './routes/NotFound.svelte';
import Loading from './components/Loading.svelte';
import {capuchin} from './stores';

const checkAuth = () => {
    if (!getStore(capuchin).authMode || getStore(capuchin).isAuth) {
        return true;
    }

    push('/login');

    return false;
};

export const routes = {
    // First route
    '/': wrap({
        asyncComponent: () => import('./routes/Track.svelte'),
        loadingComponent: Loading,
        conditions: checkAuth,
    }),

    '/login': Login,

    '/clients': wrap({
        asyncComponent: () => import('./routes/Clients.svelte'),
        loadingComponent: Loading,
        conditions: checkAuth,
    }),
    '/projects': wrap({
        asyncComponent: () => import('./routes/Projects.svelte'),
        loadingComponent: Loading,
        conditions: checkAuth,
    }),
    '/reports': wrap({
        asyncComponent: () => import('./routes/Reports.svelte'),
        loadingComponent: Loading,
        conditions: checkAuth,
    }),
    '/profile': wrap({
        asyncComponent: () => import('./routes/Profile.svelte'),
        loadingComponent: Loading,
        conditions: checkAuth,
    }),

    // Last route
    '*': NotFound,
};
