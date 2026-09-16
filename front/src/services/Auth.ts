import Api from "./Api";
import type {
  ActivatePayload,
  ApplyCodePayload,
  LoginPayload,
  RegisterPayload,
} from "../types";

const register = async (payload: RegisterPayload) =>
  await Api.post(`/auth/register`, payload);

const activate = async (payload: ActivatePayload) =>
  await Api.post(`/auth/activate`, payload);

const login = async (payload: LoginPayload) =>
  await Api.post(`/auth/login`, payload);

const applyCode = async (payload: ApplyCodePayload) =>
  await Api.post(`/auth/apply_code`, payload);

const logout = async () => await Api.post(`/auth/logout`);

export default {
  register,
  activate,
  login,
  applyCode,
  logout,
};
