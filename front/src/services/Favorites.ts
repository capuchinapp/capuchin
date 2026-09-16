import Api from "./Api";
import type { Favorite, FavoritePayload } from "../types";

const getList = async (): Promise<Favorite[]> => await Api.get(`/favorites`);

const getById = async (favoriteId: string): Promise<Favorite> =>
  await Api.get(`/favorites/${favoriteId}`);

const create = async (favorite: FavoritePayload): Promise<Favorite> =>
  await Api.post(`/favorites`, favorite);

const update = async (
  favoriteId: string,
  favorite: FavoritePayload,
): Promise<Favorite> => await Api.patch(`/favorites/${favoriteId}`, favorite);

const deleteRequest = async (favoriteId: string): Promise<Favorite> =>
  await Api.delete(`/favorites/${favoriteId}`);

export default {
  getList,
  getById,
  create,
  update,
  delete: deleteRequest,
};
