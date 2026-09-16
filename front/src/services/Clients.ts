import Api from "./Api";
import type { Client, ClientPayload } from "../types";

const getList = async (): Promise<Client[]> => await Api.get("/clients");

const getById = async (clientId: string): Promise<Client> =>
  await Api.get(`/clients/${clientId}`);

const create = async (client: ClientPayload): Promise<Client> =>
  await Api.post(`/clients`, client);

const update = async (
  clientId: string,
  client: ClientPayload,
): Promise<Client> => await Api.patch(`/clients/${clientId}`, client);

const archive = async (clientId: string): Promise<Client> =>
  await Api.post(`/clients/${clientId}/archive`);

const unarchive = async (clientId: string): Promise<Client> =>
  await Api.post(`/clients/${clientId}/unarchive`);

export default {
  getList,
  getById,
  create,
  update,
  archive,
  unarchive,
};
