import Api from "./Api";
import type { Task, TaskPayload, TaskReport } from "../types";

const getList = async (projectId: string): Promise<Task[]> =>
  await Api.get(`/tasks?project_id=${projectId}`);

const getById = async (taskId: string): Promise<Task> =>
  await Api.get(`/tasks/${taskId}`);

const create = async (task: TaskPayload): Promise<Task> =>
  await Api.post(`/tasks`, task);

const update = async (taskId: string, task: TaskPayload): Promise<Task> =>
  await Api.patch(`/tasks/${taskId}`, task);

const archive = async (taskId: string): Promise<Task> =>
  await Api.post(`/tasks/${taskId}/archive`);

const unarchive = async (taskId: string): Promise<Task> =>
  await Api.post(`/tasks/${taskId}/unarchive`);

const complete = async (taskId: string): Promise<Task> =>
  await Api.post(`/tasks/${taskId}/complete`);

const incomplete = async (taskId: string): Promise<Task> =>
  await Api.post(`/tasks/${taskId}/incomplete`);

const report = async (taskId: string): Promise<TaskReport> =>
  await Api.get(`/tasks/${taskId}/report`);

export default {
  getList,
  getById,
  create,
  update,
  archive,
  unarchive,
  complete,
  incomplete,
  report,
};
