import {get} from 'svelte/store';
import {settings} from '../stores';
import Api from './Api';

/**
 * @typedef {Object} Settings
 *
 * @property {string} dateFormat
 */

const init = async () => {
    await Api.get(`/settings`).then(
        (res) => {
            res.map((item) => (get(settings)[item['key']] = item['value']));
        },
        (err) => {
            console.error('An error occurred while getting the settings', err);
        }
    );
};

const put = (/** @type {Settings} */ settings) => Api.put(`/settings`, settings);

export default {
    init,
    put,
};
