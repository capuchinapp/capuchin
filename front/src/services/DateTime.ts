import { get } from "svelte/store";
import { settings as settingsStore } from "../stores";
import dayjs from "dayjs";

export const nowDate = (): string => {
  return dayjs().format(get(settingsStore).dateFormat);
};

export const nowTime = (withSeconds: boolean = false): string => {
  const format = withSeconds ? "HH:mm:ss" : "HH:mm";

  return dayjs().format(format);
};

export const getDateTime = (datetime: string): string => {
  let dt = parseDateTime(datetime);
  if (dt === null) {
    return "-";
  }

  return dt.format("L HH:mm");
};

export const getDate = (datetime: string): string => {
  let dt = parseDateTime(datetime);
  if (dt === null) {
    return "-";
  }

  return dt.format(get(settingsStore).dateFormat);
};

export const getTime = (datetime: string): string => {
  let dt = parseDateTime(datetime);
  if (dt === null) {
    return "-";
  }

  return dt.format("HH:mm");
};

export const getDiffSeconds = (timeStart: string, timeEnd: string): number => {
  const d1 = dayjs(`1970-01-01 ${timeEnd}:00`);
  const d2 = dayjs(`1970-01-01 ${timeStart}:00`);

  return d1.diff(d2, "second");
};

export const hoursFromSeconds = (
  seconds: number,
  precision: number,
): number => {
  let val = seconds / 60 / 60;

  if (precision < 0) {
    return val;
  }

  return parseFloat(val.toFixed(precision));
};

export const calculateTime = (
  durationSeconds: number,
  withSeconds: boolean,
): string => {
  const format = withSeconds ? "HH:mm:ss" : "HH:mm";

  return dayjs.duration(durationSeconds, "seconds").format(format);
};

export const formatDurationSmart = (seconds: number): string => {
  if (seconds < 60) {
    return `${seconds}с`;
  }

  if (seconds < 3600) {
    return `${Math.floor(seconds / 60)}м`;
  }

  if (seconds < 86400) {
    const hours = Math.floor(seconds / 3600);
    const minutes = Math.floor((seconds % 3600) / 60);

    return minutes > 0 ? `${hours}ч ${minutes}м` : `${hours}ч`;
  }

  const days = Math.floor(seconds / 86400);
  const hours = Math.floor((seconds % 86400) / 3600);

  return hours > 0 ? `${days}д ${hours}ч` : `${days}д`;
};

const parseDateTime = (datetime: string): dayjs.Dayjs | null => {
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
