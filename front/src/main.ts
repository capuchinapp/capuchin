import { mount } from "svelte";
import { get } from "svelte/store";
import { capuchin } from "./stores";
import App from "./App.svelte";
import IndexAPI from "./services/Index";
import { loadTranslations } from "./i18n";

void (async () => {
  await loadTranslations("ru", "/");

  try {
    const res = await IndexAPI.get();
    const c = get(capuchin);
    c.isAuth = res.isAuth;
    c.appVersionBack = res.appVersion;
    c.runningTimelogDatetime = res.runningTimelogDatetime;

    mount(App, {
      target: document.getElementById("app")!,
    });
  } catch (err) {
    console.error("An error occurred while getting the index", err);
  }
})();
