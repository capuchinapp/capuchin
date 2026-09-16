import { get } from "svelte/store";
import { settings } from "../stores";
import Api from "./Api";
import type { Setting, SettingsPayload } from "../types";

const init = async () => {
  await Api.get<Setting[]>(`/settings`).then(
    (res) => {
      res.map((item) => (get(settings)[item.key] = item.value));
    },
    (err) => {
      console.error("An error occurred while getting the settings", err);
    },
  );
};

const put = async (payload: SettingsPayload) =>
  await Api.put(`/settings`, payload);

function workingDaysFromString(workingDaysStr: string): string[] {
  if (workingDaysStr === "") {
    return ["1", "2", "3", "4", "5"];
  }

  return workingDaysStr.split(",");
}

function workingDaysToString(workingDaysArr: string[]): string {
  return workingDaysArr.join(",");
}

export default {
  init,
  put,
  workingDaysFromString,
  workingDaysToString,
};
