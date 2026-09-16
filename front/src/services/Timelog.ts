import Api from "./Api";
import type {
  Timelog,
  TimelogCreatePayload,
  TimelogUpdatePayload,
} from "../types";
import { capuchin } from "../stores";

const getList = async (dateFrom: string, dateTo: string): Promise<Timelog[]> =>
  await Api.get(`/timelogs?date_from=${dateFrom}&date_to=${dateTo}`);

const getListByClient = async (
  dateFrom: string,
  dateTo: string,
  clientId: string,
): Promise<Timelog[]> =>
  await Api.get(
    `/timelogs?date_from=${dateFrom}&date_to=${dateTo}&client_id=${clientId}`,
  );

const getListByProject = async (
  dateFrom: string,
  dateTo: string,
  projectId: string,
): Promise<Timelog[]> =>
  await Api.get(
    `/timelogs?date_from=${dateFrom}&date_to=${dateTo}&project_id=${projectId}`,
  );

const getById = async (timelogId: string): Promise<Timelog> =>
  await Api.get(`/timelogs/${timelogId}`);

const create = async (timelog: TimelogCreatePayload): Promise<Timelog> => {
  let res = await Api.post<Timelog>(`/timelogs`, timelog);

  capuchin.update((c) => {
    c.runningTimelogDatetime = null;
    return c;
  });

  return res;
};

const update = async (
  timelogId: string,
  timelog: TimelogUpdatePayload,
): Promise<Timelog> => await Api.patch(`/timelogs/${timelogId}`, timelog);

const deleteRequest = async (timelogId: string): Promise<Timelog> =>
  await Api.delete(`/timelogs/${timelogId}`);

const stop = async (
  timelogId: string,
  date: string,
  timeEnd: string,
): Promise<Timelog> => {
  let res = await Api.post<Timelog>(`/timelogs/${timelogId}/stop`, {
    date: date,
    timeEnd: timeEnd,
  });

  capuchin.update((c) => {
    c.runningTimelogDatetime = null;
    return c;
  });

  return res;
};

const getLastN = async (n: number): Promise<Timelog[]> =>
  await Api.get(`/timelogs/last/${n}`);

export default {
  getList,
  getListByClient,
  getListByProject,
  getById,
  create,
  update,
  delete: deleteRequest,
  stop,
  getLastN,
};
