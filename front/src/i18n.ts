import i18n from "sveltekit-i18n";

export const { t, locale, locales, loading, loadTranslations } = new i18n({
  loaders: [
    {
      locale: "ru",
      key: "",
      loader: async () => (await import("./locales/ru.json")).default,
    },
  ],
  fallbackLocale: "ru",
});
