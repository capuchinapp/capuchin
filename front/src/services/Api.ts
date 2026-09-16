import axios from "axios";
import type { AxiosError } from "axios";
import { push } from "svelte-spa-router";
import { get as getStore } from "svelte/store";
import { t } from "../i18n";
import { preloaderCount } from "../stores";
import { toastDanger, toastWarning } from "./Toast";

const axiosApi = axios.create({
  baseURL: getBaseUrl(),
  withCredentials: true,
});

axiosApi.interceptors.request.use(
  function (config) {
    preloaderCount.update((n) => n + 1);

    return config;
  },
  function (error) {
    preloaderCount.update((n) => n - 1);

    return Promise.reject(error);
  },
);

axiosApi.interceptors.response.use(
  function (response) {
    preloaderCount.update((n) => n - 1);

    return response;
  },
  function (error) {
    preloaderCount.update((n) => n - 1);

    return Promise.reject(error);
  },
);

type HttpMethod = "get" | "post" | "patch" | "put" | "delete";

const apiRequest = async <T>(
  method: HttpMethod,
  url: string,
  data?: unknown,
): Promise<T> => {
  try {
    const config: Record<string, unknown> = {
      method: method,
      url: url,
      withCredentials: true,
    };

    switch (method) {
      case "post":
      case "patch":
      case "put":
      case "delete":
        config.data = data;
        break;
    }

    const res = await axiosApi(config);

    return await Promise.resolve(res.data);
  } catch (err) {
    let toast = toastDanger;
    let content: string = getStore(t)("api.unknownError");

    const e = err as AxiosError;

    switch (e.response?.status) {
      case 400:
        toast = toastWarning;
        content = getStore(t)("api.badRequest");
        break;

      case 401:
        push("/login");
        break;

      case 403:
        toast = toastWarning;
        content = getStore(t)("api.forbidden");
        break;

      case 404:
        if (e.response?.headers["x-routenotfound"] === "1") {
          content = getStore(t)("api.routeNotFound");
        } else {
          toast = toastWarning;
          content = getStore(t)("api.objectNotFound");
        }
        break;

      case 409:
        toast = toastWarning;
        switch (e.response?.data) {
          case "user with this email already exists":
            content = getStore(t)("api.userAlreadyExists");
            break;

          case "user not activated":
            content = getStore(t)("api.userNotActivated");
            break;

          case "impossible action":
            content = getStore(t)("api.impossibleAction");
            break;

          default:
            content = getStore(t)("api.objectAlreadyExists");
            break;
        }
        break;

      case 429:
        toast = toastWarning;
        content = getStore(t)("api.tooManyRequests");
        break;

      case 500:
        toast = toastWarning;
        content = getStore(t)("api.internalServerError");
        break;

      default:
        console.error(err);
        break;
    }

    toast(content, false);

    return await Promise.reject(e);
  }
};

function getBaseUrl() {
  if (import.meta.env.DEV) {
    return import.meta.env.VITE_API_BASE_URL;
  }

  return `${location.protocol}//${location.host}/api`;
}

const get = async <T>(url: string, data?: unknown): Promise<T> =>
  await apiRequest<T>("get", url, data);

const deleteRequest = async <T>(url: string, data?: unknown): Promise<T> =>
  await apiRequest<T>("delete", url, data);

const post = async <T>(url: string, data?: unknown): Promise<T> =>
  await apiRequest<T>("post", url, data);

const put = async <T>(url: string, data?: unknown): Promise<T> =>
  await apiRequest<T>("put", url, data);

const patch = async <T>(url: string, data?: unknown): Promise<T> =>
  await apiRequest<T>("patch", url, data);

const Api = {
  get,
  delete: deleteRequest,
  post,
  put,
  patch,
};

export default Api;
