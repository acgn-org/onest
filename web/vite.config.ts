import { defineConfig } from "vite";
import react from "@vitejs/plugin-react-swc";

import { resolve } from "path";
const __dirname = resolve();

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      "@": resolve(__dirname, "src"),
      "@component": resolve(__dirname, "src/components"),
    },
  },
});
