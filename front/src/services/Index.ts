import Api from "./Api";
import type { IndexResponse } from "../types";

const get = async (): Promise<IndexResponse> => await Api.get(`/`);

export default {
  get,
};
