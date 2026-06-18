import {defineConfig} from 'vite'
import vue from '@vitejs/plugin-vue'
import AutoImport from 'unplugin-auto-import/vite';
import Components from 'unplugin-vue-components/vite';
import { TDesignResolver } from '@tdesign-vue-next/auto-import-resolver';
import { NaiveUiResolver } from 'unplugin-vue-components/resolvers';

const tDesignChatResolver = TDesignResolver({
    library: 'chat'
})

function manualChunks(id) {
    if (!id.includes('node_modules')) {
        return
    }
    if (id.includes('/vue/') || id.includes('/vue-router/')) {
        return 'vue-vendor'
    }
    if (id.includes('/naive-ui/') || id.includes('/vdirs/') || id.includes('/vooks/') || id.includes('/vueuc/')) {
        return 'naive-vendor'
    }
    if (id.includes('/echarts/') || id.includes('/zrender/') || id.includes('/lightweight-charts/')) {
        return 'chart-vendor'
    }
    if (id.includes('/md-editor-v3/') || id.includes('/@lezer/')) {
        return 'markdown-vendor'
    }
    if (id.includes('/html2canvas/') || id.includes('/html-docx-js-typescript/') || id.includes('/@vavt/')) {
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
  ],
  build: {
      rollupOptions: {
          output: {
              manualChunks,
          },
      },
  },
})
