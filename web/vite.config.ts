import { defineConfig } from "vite";
import react from "@vitejs/plugin-react-swc";

import { resolve } from "path";
const __dirname = resolve();

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      "@tabler/icons-react": "@tabler/icons-react/dist/esm/icons/index.mjs",

      "@": resolve(__dirname, "src"),
      "@component": resolve(__dirname, "src/components"),
      "@page": resolve(__dirname, "src/pages"),
      "@hook": resolve(__dirname, "src/hooks"),
    },
  },
});
