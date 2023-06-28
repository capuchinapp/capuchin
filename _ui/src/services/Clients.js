import Api from './Api';

/**
 * @typedef {Object} Client
 *
 * @property {string} name
 * @property {number} billableRate
 * @property {string} comment
 */

const getList = async () => {
    try {
        return await Api.get('/clients');
    } catch (error) {
        return [];
    }
};

const getById = async (/** @type {string} */ clientUUID) => {
    try {
        return await Api.get(`/clients/${clientUUID}`);
    } catch (error) {
        return null;
    }
};

const create = async (/** @type {Client} */ client) => {
    try {
        return await Api.post(`/clients`, client);
    } catch (error) {
        return null;
    }
};

const update = async (/** @type {string} */ clientUUID, /** @type {Client} */ client) => {
    try {
        return await Api.patch(`/clients/${clientUUID}`, client);
    } catch (error) {
        return null;
    }
};

const archive = async (/** @type {string} */ clientUUID) => {
    try {
        return await Api.post(`/clients/${clientUUID}/archive`);
    } catch (error) {
        return null;
    }
};

const unarchive = async (/** @type {string} */ clientUUID) => {
    try {
        return await Api.post(`/clients/${clientUUID}/unarchive`);
    } catch (error) {
        return null;
    }
};

export default {
    getList,
    getById,
    create,
    update,
    archive,
    unarchive,
};
