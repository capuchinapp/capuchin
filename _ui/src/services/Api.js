import axios from 'axios';
import {_} from 'svelte-i18n';
import {push} from 'svelte-spa-router';
import {get as getStore} from 'svelte/store';
import {preloaderCount} from '../stores';
import {toastDanger, toastWarning} from './Toast';

// Create a instance of axios to use the same base url.
const axiosApi = axios.create({
    baseURL: import.meta.env.DEV ? import.meta.env.VITE_API_BASE_URL : `http://${location.host}/api`,
});

// Add a request interceptor
axiosApi.interceptors.request.use(
    function (config) {
        preloaderCount.update((n) => n + 1);

        return config;
    },
    function (error) {
        preloaderCount.update((n) => n - 1);

        return Promise.reject(error);
    }
);

// Add a response interceptor
axiosApi.interceptors.response.use(
    function (response) {
        // Any status code that lie within the range of 2xx cause this function to trigger
        preloaderCount.update((n) => n - 1);

        return response;
    },
    function (error) {
        // Any status codes that falls outside the range of 2xx cause this function to trigger
        preloaderCount.update((n) => n - 1);

        return Promise.reject(error);
    }
);

// implement a method to execute all the request from here.
const apiRequest = async (
    /** @type {string} */ method,
    /** @type {string} */ url,
    /** @type {Array|Object|undefined} */ data
) => {
    try {
        let config = {
            method: method,
            url: url,
        };

        switch (method) {
            case 'post':
            case 'patch':
            case 'put':
            case 'delete':
                config.data = data;
                break;
        }

        const res = await axiosApi(config);

        return await Promise.resolve(res.data);
    } catch (err) {
        let toast = toastDanger;
        let content = getStore(_)('api.unknownError');

        switch (err.response.status) {
            case 400:
                toast = toastWarning;

                switch (err.response.data) {
                    case 'Wrong password':
                        content = getStore(_)('api.wrongPassword');
                        break;

                    default:
                        content = getStore(_)('api.badRequest');
                        break;
                }
                break;

            case 401:
                push('/login');
                break;

            case 404:
                if (err.response.data === 'Not Found') {
                    toast = toastWarning;
                    content = getStore(_)('api.objectNotFound');
                } else {
                    content = getStore(_)('api.routeNotFound');
                }
                break;

            case 500:
                toast = toastWarning;
                content = getStore(_)('api.internalServerError');
                break;

            default:
                console.error(err);
                break;
        }

        toast(String(content), false);

        return await Promise.reject(err);
    }
};

// function to execute the http get request
const get = (/** @type {string} */ url, /** @type {Array|Object|undefined} */ data) => apiRequest('get', url, data);

// function to execute the http delete request
const deleteRequest = (/** @type {string} */ url, /** @type {Array|Object|undefined} */ data) =>
    apiRequest('delete', url, data);

// function to execute the http post request
const post = (/** @type {string} */ url, /** @type {Array|Object|undefined} */ data) => apiRequest('post', url, data);

// function to execute the http put request
const put = (/** @type {string} */ url, /** @type {Array|Object|undefined} */ data) => apiRequest('put', url, data);

// function to execute the http patch request
const patch = (/** @type {string} */ url, /** @type {Array|Object|undefined} */ data) => apiRequest('patch', url, data);

// expose your method to other services or actions
const Api = {
    get,
    delete: deleteRequest,
    post,
    put,
    patch,
};

export default Api;
