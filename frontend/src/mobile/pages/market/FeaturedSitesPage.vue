<script setup>
import { ref } from 'vue'
import MPullRefresh from '../../components/base/MPullRefresh.vue'
import MCard from '../../components/base/MCard.vue'
import MEmpty from '../../components/base/MEmpty.vue'

const sites = ref([
  { name: '东方财富', url: 'https://www.eastmoney.com', desc: '中国财经门户网站' },
  { name: '雪球', url: 'https://xueqiu.com', desc: '投资者社交平台' },
  { name: '同花顺', url: 'https://www.10jqka.com.cn', desc: '金融数据服务商' },
])

async function handleRefresh() { return new Promise(resolve => setTimeout(resolve, 1500)) }
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
