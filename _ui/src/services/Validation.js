import {_} from 'svelte-i18n';
import {get} from 'svelte/store';
import * as yup from 'yup';

export default class Validation {
    #msgRequiredInput = '';
    #msgRequiredSelect = '';

    /** @type {yup.ObjectSchema} */ #clientSchema;
    /** @type {yup.ObjectSchema} */ #projectSchema;
    /** @type {yup.ObjectSchema} */ #timelogSchema;

    /** @type {yup.StringSchema} */ #requiredPassword;
    /** @type {yup.StringSchema} */ #requiredName;
    /** @type {yup.StringSchema} */ #requiredSelectUUID;
    /** @type {yup.NumberSchema} */ #requiredBillableRate;

    constructor() {
        this.#msgRequiredInput = get(_)('validate.required.input');
        this.#msgRequiredSelect = get(_)('validate.required.select');
    }

    /**
     * @param {Object} data
     * @param {Function} callbackSuccess
     * @param {Function} callbackError
     */
    validateLogin(data, callbackSuccess, callbackError) {
        if (!this.#clientSchema) {
            this.#clientSchema = yup.object().shape({
                password: this.#getRequiredPassword(),
            });
        }

        this.#validateSchema(this.#clientSchema, data, callbackSuccess, callbackError);
    }

    /**
     * @param {Object} data
     * @param {Function} callbackSuccess
     * @param {Function} callbackError
     */
    validateClient(data, callbackSuccess, callbackError) {
        if (!this.#clientSchema) {
            this.#clientSchema = yup.object().shape({
                name: this.#getRequiredName(),
                billableRate: this.#getRequiredBillableRate(),
            });
        }

        this.#validateSchema(this.#clientSchema, data, callbackSuccess, callbackError);
    }

    /**
     * @param {Object} data
     * @param {Function} callbackSuccess
     * @param {Function} callbackError
     */
    validateProject(data, callbackSuccess, callbackError) {
        if (!this.#projectSchema) {
            this.#projectSchema = yup.object().shape({
                clientUUID: this.#getRequiredSelectUUID(),
                name: this.#getRequiredName(),
                billableRate: this.#getRequiredBillableRate(),
            });
        }

        this.#validateSchema(this.#projectSchema, data, callbackSuccess, callbackError);
    }

    /**
     * @param {Object} data
     * @param {Function} callbackSuccess
     * @param {Function} callbackError
     */
    validateTimelog(data, callbackSuccess, callbackError) {
        if (!this.#timelogSchema) {
            const time = yup.string().matches(/\d{0,23}:\d{0,59}/, {message: get(_)('validate.time')});

            this.#timelogSchema = yup.object().shape({
                projectUUID: this.#getRequiredSelectUUID(),
                date: yup.date().required(this.#msgRequiredInput),
                timeStart: time.required(this.#msgRequiredInput),
                timeEnd: time.nullable(),
                billableRate: this.#getRequiredBillableRate(),
            });
        }

        this.#validateSchema(this.#timelogSchema, data, callbackSuccess, callbackError);
    }

    /**
     * @param {yup.ObjectSchema} scheme
     * @param {Object} data
     * @param {Function} callbackSuccess
     * @param {Function} callbackError
     */
    #validateSchema(scheme, data, callbackSuccess, callbackError) {
        scheme
            .validate(data, {abortEarly: false})
            .then(() => callbackSuccess())
            //.then(() => console.log(data)) // debug
            .catch((err) => callbackError(this.#extractErrors(err)));
    }

    /**
     * @returns {yup.StringSchema}
     */
    #getRequiredPassword() {
        if (!this.#requiredPassword) {
            this.#requiredPassword = yup.string().required(this.#msgRequiredInput).min(8, get(_)('validate.password'));
        }

        return this.#requiredPassword;
    }

    /**
     * @returns {yup.StringSchema}
     */
    #getRequiredName() {
        if (!this.#requiredName) {
            this.#requiredName = yup.string().required(this.#msgRequiredInput);
        }

        return this.#requiredName;
    }

    /**
     * @returns {yup.StringSchema}
     */
    #getRequiredSelectUUID() {
        if (!this.#requiredSelectUUID) {
            this.#requiredSelectUUID = yup.string().uuid().required(this.#msgRequiredSelect);
        }

        return this.#requiredSelectUUID;
    }

    /**
     * @returns {yup.NumberSchema}
     */
    #getRequiredBillableRate() {
        if (!this.#requiredBillableRate) {
            this.#requiredBillableRate = yup
                .number()
                .required(this.#msgRequiredInput)
                .min(0, get(_)('validate.number.gtezero'));
        }

        return this.#requiredBillableRate;
    }

    #extractErrors({inner}) {
        return inner.reduce((acc, err) => {
            return {...acc, [err.path]: err.message};
        }, {});
    }
}
