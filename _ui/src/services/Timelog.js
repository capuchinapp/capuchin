import Api from './Api';

/**
 * @typedef {Object} Timelog
 *
 * @property {string} projectUUID
 * @property {string} date
 * @property {string} timeStart
 * @property {string|null} timeEnd
 * @property {number} billableRate
 * @property {string|null} comment
 */

const getList = async (/** @type {string} */ dateFrom, /** @type {string} */ dateTo) => {
    try {
        return await Api.get(`/timelogs?date_from=${dateFrom}&date_to=${dateTo}`);
    } catch (error) {
        return [];
    }
};

const getListByClient = async (
    /** @type {string} */ dateFrom,
    /** @type {string} */ dateTo,
    /** @type {string} */ clientUUID
) => {
    try {
        return await Api.get(`/timelogs?date_from=${dateFrom}&date_to=${dateTo}&client_uuid=${clientUUID}`);
    } catch (error) {
        return [];
    }
};

const getListByProject = async (
    /** @type {string} */ dateFrom,
    /** @type {string} */ dateTo,
    /** @type {string} */ projectUUID
) => {
    try {
        return await Api.get(`/timelogs?date_from=${dateFrom}&date_to=${dateTo}&project_uuid=${projectUUID}`);
    } catch (error) {
        return [];
    }
};

const getById = async (/** @type {string} */ timelogUUID) => {
    try {
        return await Api.get(`/timelogs/${timelogUUID}`);
    } catch (error) {
        return null;
    }
};

const create = async (/** @type {Timelog} */ timelog) => {
    try {
        return await Api.post(`/timelogs`, timelog);
    } catch (error) {
        return null;
    }
};

const update = async (/** @type {string} */ timelogUUID, /** @type {Timelog} */ timelog) => {
    try {
        return await Api.patch(`/timelogs/${timelogUUID}`, timelog);
    } catch (error) {
        return null;
    }
};

const deleteRequest = async (/** @type {string} */ timelogUUID) => {
    try {
        return await Api.delete(`/timelogs/${timelogUUID}`);
    } catch (error) {
        return null;
    }
};

const stop = async (/** @type {string} */ timelogUUID, /** @type {string} */ date, /** @type {string} */ timeEnd) => {
    try {
        return await Api.patch(`/timelogs/${timelogUUID}/stop`, {date: date, timeEnd: timeEnd});
    } catch (error) {
        return null;
    }
};

export default {
    getList,
    getListByClient,
    getListByProject,
    getById,
    create,
    update,
    delete: deleteRequest,
    stop,
};
