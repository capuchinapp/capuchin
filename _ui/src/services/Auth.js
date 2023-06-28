import Api from './Api';

/**
 * @typedef {Object} Payload
 *
 * @property {string} password
 */

const login = async (/** @type {Payload} */ payload) => Api.post(`/auth/login`, payload);

const logout = async () => Api.post(`/auth/logout`);

export default {
    login,
    logout,
};
