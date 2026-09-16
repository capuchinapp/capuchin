// API domain types (corresponds to backend/internal/domain)

export interface Client {
  id: string;
  name: string;
  billableRate: number;
  comment: string;
  createdAt: string;
  updatedAt: string | null;
  archivedAt: string | null;
}

export interface Project {
  id: string;
  clientId: string;
  clientName: string;
  name: string;
  billableRate: number;
  comment: string;
  createdAt: string;
  updatedAt: string | null;
  archivedAt: string | null;
}

export interface Task {
  id: string;
  projectId: string;
  projectName: string;
  clientId: string;
  clientName: string;
  name: string;
  comment: string;
  createdAt: string;
  updatedAt: string | null;
  archivedAt: string | null;
  completedAt: string | null;
}

export interface TaskReport {
  durationSeconds: number;
  billableAmount: number;
  uniqueDatesCount: number;
}

export interface Timelog {
  id: string;
  projectId: string;
  projectName: string;
  clientId: string;
  clientName: string;
  date: string;
  timeStart: string;
  timeEnd: string | null;
  durationSeconds: number;
  billableRate: number;
  billableAmount: number;
  comment: string;
  createdAt: string;
  updatedAt: string | null;
  taskId: string | null;
  taskName: string | null;
  taskCompletedAt: string | null;
}

export interface Favorite {
  id: string;
  name: string;
  projectId: string;
  projectName: string;
  clientId: string;
  clientName: string;
  taskId: string | null;
  taskName: string | null;
  billableRate: number;
  comment: string;
}

export interface Session {
  id: string;
  checkedAt: string;
  isCurrent: boolean;
}

export interface Setting {
  key: string;
  value: string;
}

export interface IndexResponse {
  name: string;
  isAuth: boolean;
  appVersion: string;
  runningTimelogDatetime: string | null;
}

// UI option types

export interface ClientOption {
  id: string;
  name: string;
  billableRate: string;
}

export interface ProjectOption {
  id: string;
  clientId: string;
  name: string;
  billableRate: string;
}

export interface TaskOption {
  id: string;
  name: string;
}

// Form payload types (corresponds to backend restapi input structs)

export interface RegisterPayload {
  email: string;
}

export interface ActivatePayload {
  userId: string;
  code: string;
}

export interface LoginPayload {
  email: string;
}

export interface ApplyCodePayload {
  email: string;
  code: string;
}

export interface ClientPayload {
  name: string;
  billableRate: number;
  comment: string | null;
}

export interface ProjectPayload {
  clientId: string;
  name: string;
  billableRate: number;
  comment: string | null;
}

export interface TaskPayload {
  projectId: string;
  name: string;
  comment: string | null;
}

export interface TimelogCreatePayload {
  projectId: string;
  taskId: string | null;
  date: string;
  timeStart: string;
  timeEnd: string | null;
  billableRate: number;
  comment: string | null;
}

export interface TimelogUpdatePayload {
  projectId: string;
  taskId: string | null;
  date: string;
  timeStart: string;
  timeEnd: string | null;
  billableRate: number;
  comment: string | null;
}

export interface TimelogStopPayload {
  date: string;
  timeEnd: string;
}

export interface FavoritePayload {
  name: string;
  projectId: string;
  taskId: string | null;
  billableRate: number;
  comment: string | null;
}

export interface SettingsPayload {
  dateFormat: string;
  workingDays: string;
}
