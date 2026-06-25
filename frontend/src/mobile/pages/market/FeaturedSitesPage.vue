<script setup>
import { ref } from 'vue'
import MPullRefresh from '../../components/base/MPullRefresh.vue'
import MCard from '../../components/base/MCard.vue'
import MEmpty from '../../components/base/MEmpty.vue'

// 名站列表（后续接 GetConfig 里的 embeddedUrls 或后端站点配置填充）
const sites = ref([])

async function handleRefresh() {
  // 后续接真实 API
}
</script>

<template>
  <div class="page">
    <MPullRefresh :on-refresh="handleRefresh">
      <div class="container">
        <MCard v-if="sites.length">
          <div class="list">
            <div v-for="(site, i) in sites" :key="i" class="item">
              <div class="name">{{ site.name }}</div>
              <div class="desc">{{ site.desc }}</div>
            </div>
          </div>
        </MCard>
        <MEmpty v-else description="暂无数据" />
      </div>
    </MPullRefresh>
  </div>
</template>

<style scoped>
.page { height: 100%; overflow: hidden; }
.container { padding: var(--m-space-md); }
.list { display: flex; flex-direction: column; gap: var(--m-space-md); }
.item { padding: var(--m-space-md); background: var(--m-bg-primary); border-radius: var(--m-radius-sm); }
.name { font-weight: var(--m-font-weight-medium); margin-bottom: var(--m-space-xs); }
.desc { font-size: var(--m-font-sm); color: var(--m-text-secondary); }
</style>
