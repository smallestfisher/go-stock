import {defineConfig} from 'vite'
import vue from '@vitejs/plugin-vue'
import AutoImport from 'unplugin-auto-import/vite';
import Components from 'unplugin-vue-components/vite';
import { TDesignResolver } from '@tdesign-vue-next/auto-import-resolver';
import { NaiveUiResolver } from 'unplugin-vue-components/resolvers';
import { VitePWA } from 'vite-plugin-pwa';

const tDesignChatResolver = TDesignResolver({
    library: 'chat'
})

function manualChunks(id) {
    if (!id.includes('node_modules')) {
        return
    }
    if (
        id.includes('/vue/') ||
        id.includes('/vue-router/') ||
        id.includes('/naive-ui/') ||
        id.includes('/vdirs/') ||
        id.includes('/vooks/') ||
        id.includes('/vueuc/')
    ) {
        return 'ui-vendor'
    }
    if (id.includes('/echarts/') || id.includes('/zrender/') || id.includes('/lightweight-charts/')) {
        return 'chart-vendor'
    }
    if (id.includes('/md-editor-v3/') || id.includes('/@lezer/') || id.includes('/@vavt/')) {
        return 'markdown-vendor'
    }
    if (id.includes('/html2canvas/') || id.includes('/html-docx-js-typescript/')) {
        return 'export-vendor'
    }
    if (id.includes('/@tdesign-vue-next/') || id.includes('/tdesign-icons-vue-next/')) {
        return 'tdesign-vendor'
    }
    if (id.includes('/@vicons/')) {
        return 'icon-vendor'
    }
}

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
      vue(),
      AutoImport({
          resolvers: [tDesignChatResolver],
      }),
      Components({
          resolvers: [NaiveUiResolver(), tDesignChatResolver],
      }),
      VitePWA({
          registerType: 'autoUpdate',
          includeAssets: ['favicon.ico', 'appicon.png'],
          manifest: {
              name: 'go-stock: AI赋能股票分析',
              short_name: 'go-stock',
              description: '基于大语言模型的AI赋能股票分析工具，支持A股、港股、美股',
              theme_color: '#18a058',
              background_color: '#ffffff',
              display: 'standalone',
              orientation: 'portrait-primary',
              scope: '/',
              start_url: '/',
              icons: [
                  { src: 'pwa-192x192.png', sizes: '192x192', type: 'image/png' },
                  { src: 'pwa-512x512.png', sizes: '512x512', type: 'image/png' },
                  { src: 'pwa-512x512.png', sizes: '512x512', type: 'image/png', purpose: 'any maskable' }
              ]
          },
          workbox: {
              globPatterns: ['**/*.{js,css,html,ico,png,svg,woff,woff2}'],
              navigateFallback: null,
              runtimeCaching: [
                  {
                      urlPattern: /^https:\/\/.*\.(?:png|jpg|jpeg|svg|gif)$/,
                      handler: 'CacheFirst',
                      options: {
                          cacheName: 'images-cache',
                          expiration: { maxEntries: 60, maxAgeSeconds: 30 * 24 * 60 * 60 }
                      }
                  }
              ]
          },
          devOptions: { enabled: false }
      }),
  ],
  server: {
      port: 5173,
      proxy: {
          '/api': {
              target: 'http://localhost:18888',
              changeOrigin: true
          }
      }
  },
  build: {
      chunkSizeWarningLimit: 2000,
      rollupOptions: {
          output: {
              manualChunks,
          },
      },
  },
})
