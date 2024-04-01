import { defineConfig } from "vite";
import million from "million/compiler";
import path from "path";
import react from "@vitejs/plugin-react-swc";
import { TanStackRouterVite } from "@tanstack/router-vite-plugin";
// https://vitejs.dev/config/
export default defineConfig({
  plugins: [million.vite({ auto: true }), react(), TanStackRouterVite()],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
});