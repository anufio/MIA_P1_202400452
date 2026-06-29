import { defineConfig, loadEnv } from "vite";
import react from "@vitejs/plugin-react";

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), "");

  return {
    plugins: [react()],
    server: {
      host: "0.0.0.0",
      port: 5173,
      strictPort: true
    },
    preview: {
      host: "0.0.0.0",
      port: 4173,
      strictPort: true
    },
    define: {
      __API_URL__: JSON.stringify(env.VITE_API_URL || "http://localhost:8080/api")
    },
    build: {
      outDir: "dist",
      emptyOutDir: true
    }
  };
});
