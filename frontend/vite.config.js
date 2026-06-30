import {defineConfig} from 'vite'
import vue from '@vitejs/plugin-vue'
import { VitePWA } from 'vite-plugin-pwa';

function manualChunks(id) {
    if (!id.includes('node_modules')) {
        return
    }
    if (id.includes('/echarts/') || id.includes('/zrender/')) {
        return 'chart-vendor'
    }
    if (id.includes('/md-editor-v3/') || id.includes('/@lezer/') || id.includes('/@vavt/')) {
        return 'markdown-vendor'
    }
}

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
      vue(),
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
              // markdown-vendor(md-editor-v3) 约 2.3MB，超过默认 2MiB 预缓存上限会导致 build 失败，放宽到 4MiB。
              maximumFileSizeToCacheInBytes: 4 * 1024 * 1024,
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
      host: '0.0.0.0',
      port: 5173,
      allowedHosts: [
          '192.168.5.20',
          'lanbeilvdou.dpdns.org',
      ],
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
