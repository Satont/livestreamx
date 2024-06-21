import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import tailwind from "tailwindcss"
import autoprefixer from "autoprefixer"
import path from "node:path"
import { watch } from 'vite-plugin-watch'
import { webUpdateNotice } from "@plugin-web-update-notification/vite";

// https://vitejs.dev/config/
export default defineConfig({
  css: {
    postcss: {
      plugins: [tailwind(), autoprefixer()],
    },
  },
  plugins: [
    vue(),
    watch({
      onInit: true,
      pattern: 'src/**/*.ts',
      command: 'graphql-codegen',
    }),
    webUpdateNotice({
      notificationProps: {
        title: 'New version',
        description: 'An update available, please refresh the page to get latest features and bug fixes!',
        buttonText: 'refresh',
        dismissButtonText: 'cancel',
      },
      checkInterval: 1 * 60 * 1000,
    }),
  ],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
  server: {
    proxy: {
      '/api': {
        target: 'http://127.0.0.1:1337',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, ''),
        ws: true
      },
    }
  }
})
