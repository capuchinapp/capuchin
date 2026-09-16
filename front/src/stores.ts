import { writable } from "svelte/store";
import { appVersion } from "./assets/version";

export interface CapuchinState {
  isAuth: boolean;
  appVersionBack: string;
  appVersionFront: string;
  runningTimelogDatetime: string | null;
}

export interface SettingsState {
  dateFormat: string;
  workingDays: string;
  [key: string]: string;
}

export const preloaderCount = writable(0);
export const toasts = writable<unknown[]>([]);
export const timeTotalSeconds = writable(0);
export const timeTotalClass = writable("");
export const capuchin = writable<CapuchinState>({
  isAuth: false,
  appVersionBack: "",
  appVersionFront: appVersion,
  runningTimelogDatetime: null,
});
export const settings = writable<SettingsState>({
  dateFormat: "DD.MM.YYYY",
  workingDays: "",
});
