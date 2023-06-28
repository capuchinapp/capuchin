import {get} from 'svelte/store';
import {settings as settingsStore} from '../stores';
import dayjs from 'dayjs';

export const nowDate = () => {
    return dayjs().format(get(settingsStore).dateFormat);
};

export const nowTime = (/** @type {boolean} */ withSeconds) => {
    const format = withSeconds ? 'HH:mm:ss' : 'HH:mm';

    return dayjs().format(format);
};

export const getDateTime = (/** @type {string} */ datetime) => {
    let dt = parseDateTime(datetime);
    if (dt === null) {
        return '-';
    }

    return dt.format('L HH:mm');
};

export const getDate = (/** @type {string} */ datetime) => {
    let dt = parseDateTime(datetime);
    if (dt === null) {
        return '-';
    }

    return dt.format(get(settingsStore).dateFormat);
};

export const getTime = (/** @type {string} */ datetime) => {
    let dt = parseDateTime(datetime);
    if (dt === null) {
        return '-';
    }

    return dt.format('HH:mm');
};

export const getDiffSeconds = (/** @type {string} */ timeStart, /** @type {string} */ timeEnd) => {
    const d1 = dayjs(`1970-01-01 ${timeEnd}:00`);
    const d2 = dayjs(`1970-01-01 ${timeStart}:00`);

    return d1.diff(d2, 'second');
};

export const hoursFromSeconds = (/** @type {number} */ seconds, /** @type {number} */ precision) => {
    let val = seconds / 60 / 60;

    if (precision < 0) {
        return val;
    }

    return parseFloat(val.toFixed(precision));
};

export const calculateTime = (/** @type {number} */ durationSeconds, /** @type {boolean} */ withSeconds) => {
    const format = withSeconds ? 'HH:mm:ss' : 'HH:mm';

    return dayjs.duration(durationSeconds, 'seconds').format(format);
};

const parseDateTime = (/** @type {string} */ datetime) => {
    try {
        let dt = dayjs(datetime);

        if (dt.isValid()) {
            return dt;
        }

        return null;
    } catch (error) {
        console.error(error);

        return null;
    }
};
