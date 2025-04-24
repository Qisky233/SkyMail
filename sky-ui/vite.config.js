import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { fileURLToPath, URL } from 'node:url'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  },
  // VitePluginLess
  server: {
    proxy: {
      // 匹配所有以 /api 开头的路径
      '/api': {
        target: 'http://localhost:7010',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, ''),
        secure: false, // 如果是HTTPS目标，则设置为true
      }
    }
  }
})
