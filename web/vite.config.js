import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import qiankun from 'vite-plugin-qiankun'

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    qiankun('raven-mail', {
      useDevMode: true
    }),
    vue(),
  ],
  server: {
    port: 5173,
    cors: true,
    origin: 'http://localhost:5173',
  },
})
