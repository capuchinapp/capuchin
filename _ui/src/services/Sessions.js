import Api from './Api';

const get = async () => Api.get(`/sessions`);

const deleteRequest = async (/** @type {number} */ sessID) => Api.delete(`/sessions/${sessID}`);

export default {
    get,
    delete: deleteRequest,
};
