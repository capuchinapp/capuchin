import globals from "globals";
import js from "@eslint/js";

export default [
  {
    ignores: ["dist/**"],
  },
  {
    files: ["**/*.js", "**/*.cjs", "**/*.mjs"],
    ...js.configs.recommended,
    languageOptions: {
      globals: globals.browser,
    },
  },
];
