import { defineConfig } from "vite";
import { svelte } from "@sveltejs/vite-plugin-svelte";

// https://vitejs.dev/config/
export default defineConfig(({ mode }) => {
  if (mode === "development") {
    const csp =
      "default-src 'self'; script-src 'self'; script-src-elem 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; connect-src 'self' http://localhost:* https://app.capuchin.ru data:; font-src 'self'; object-src 'none'; media-src 'none'; frame-src 'self'; child-src 'none'; form-action 'none'; worker-src 'none'; manifest-src 'self'; frame-ancestors 'none'";

    return {
      plugins: [svelte()],
      server: {
        headers: {
          "Content-Security-Policy": csp,
        },
      },
    };
  }

  return {
    plugins: [svelte()],
  };
});
