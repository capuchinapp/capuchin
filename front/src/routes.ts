import { push } from "svelte-spa-router";
import { wrap } from "svelte-spa-router/wrap";
import { get as getStore } from "svelte/store";
import Login from "./routes/Login.svelte";
import Activate from "./routes/Activate.svelte";
import NotFound from "./routes/NotFound.svelte";
import Loading from "./components/Loading.svelte";
import { capuchin } from "./stores";

// eslint-disable-next-line @typescript-eslint/no-explicit-any
const LoadingAny: any = Loading;

const checkAuth = (): boolean => {
  if (getStore(capuchin).isAuth) {
    return true;
  }

  push("/login");

  return false;
};

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export const routes: Record<string, any> = {
  "/": wrap({
    asyncComponent: () => import("./routes/Track.svelte"),
    loadingComponent: LoadingAny,
    conditions: checkAuth,
  }) as any,

  "/login": Login,
  "/activate/:uid/:code": Activate,

  "/clients": wrap({
    asyncComponent: () => import("./routes/Clients.svelte"),
    loadingComponent: LoadingAny,
    conditions: checkAuth,
  }) as any,
  "/projects": wrap({
    asyncComponent: () => import("./routes/Projects.svelte"),
    loadingComponent: LoadingAny,
    conditions: checkAuth,
  }) as any,
  "/reports": wrap({
    asyncComponent: () => import("./routes/Reports.svelte"),
    loadingComponent: LoadingAny,
    conditions: checkAuth,
  }) as any,
  "/profile": wrap({
    asyncComponent: () => import("./routes/Profile.svelte"),
    loadingComponent: LoadingAny,
    conditions: checkAuth,
  }) as any,

  "*": NotFound,
};
