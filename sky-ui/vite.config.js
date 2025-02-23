import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';
import path from 'path';

export default defineConfig({
  plugins: [vue()],
  base: './',
  publicDir: 'public',
  assetsInclude: ['**/*.svg', '**/*.png'],
  build: {
    assetsInlineLimit: 4096,
    cssCodeSplit: false  // 移动端建议关闭 CSS 代码分割
  },
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src') // 设置 @ 别名指向 src 目录
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
});