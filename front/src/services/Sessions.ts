import Api from "./Api";
import type { Session } from "../types";

const get = async (): Promise<Session[]> => await Api.get(`/sessions`);

const deleteRequest = async (sessionId: string): Promise<Session> =>
  await Api.delete(`/sessions/${sessionId}`);

export default {
  get,
  delete: deleteRequest,
};
