import Api from "./Api";
import type { Project, ProjectPayload } from "../types";

const getList = async (filterArchivedClients: boolean): Promise<Project[]> => {
  let fac = "";
  if (filterArchivedClients) {
    fac = "?filter_archived_clients=1";
  }

  return await Api.get(`/projects${fac}`);
};

const getById = async (projectId: string): Promise<Project> =>
  await Api.get(`/projects/${projectId}`);

const create = async (project: ProjectPayload): Promise<Project> =>
  await Api.post(`/projects`, project);

const update = async (
  projectId: string,
  project: ProjectPayload,
): Promise<Project> => await Api.patch(`/projects/${projectId}`, project);

const archive = async (projectId: string): Promise<Project> =>
  await Api.post(`/projects/${projectId}/archive`);

const unarchive = async (projectId: string): Promise<Project> =>
  await Api.post(`/projects/${projectId}/unarchive`);

export default {
  getList,
  getById,
  create,
  update,
  archive,
  unarchive,
};
