import Api from './Api';

/**
 * @typedef {Object} Project
 *
 * @property {string} clientUUID
 * @property {string} name
 * @property {number} billableRate
 * @property {string} comment
 */

const getList = async (/** @type {boolean} */ filterArchivedClients) => {
    try {
        let fac = '';
        if (filterArchivedClients) {
            fac = '?filter_archived_clients=1';
        }

        return await Api.get(`/projects${fac}`);
    } catch (error) {
        return [];
    }
};

const getById = async (/** @type {string} */ projectUUID) => {
    try {
        return await Api.get(`/projects/${projectUUID}`);
    } catch (error) {
        return null;
    }
};

const create = async (/** @type {Project} */ project) => {
    try {
        return await Api.post(`/projects`, project);
    } catch (error) {
        return null;
    }
};

const update = async (/** @type {string} */ projectUUID, /** @type {Project} */ project) => {
    try {
        return await Api.patch(`/projects/${projectUUID}`, project);
    } catch (error) {
        return null;
    }
};

const archive = async (/** @type {string} */ projectUUID) => {
    try {
        return await Api.post(`/projects/${projectUUID}/archive`);
    } catch (error) {
        return null;
    }
};

const unarchive = async (/** @type {string} */ projectUUID) => {
    try {
        return await Api.post(`/projects/${projectUUID}/unarchive`);
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
