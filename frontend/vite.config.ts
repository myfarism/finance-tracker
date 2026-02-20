import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

export default defineConfig({
  plugins: [react()],
  build: {
    outDir: "dist",
    chunkSizeWarningLimit: 1000,
    rollupOptions: {
      output: {
        manualChunks: {
          // React core
          "vendor-react": ["react", "react-dom", "react-router-dom"],
          // State management + form
          "vendor-store": ["zustand", "react-hook-form", "@hookform/resolvers", "zod"],
          // Charts
          "vendor-charts": ["recharts"],
          // HTTP
          "vendor-axios": ["axios"],
        },
      },
    },
  },
});
