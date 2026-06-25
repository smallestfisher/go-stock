<script setup>
import { ref } from 'vue'
import MPullRefresh from '../../components/base/MPullRefresh.vue'
import IndustryCard from '../../components/cards/IndustryCard.vue'
import MCard from '../../components/base/MCard.vue'

const industries = ref([
  { name: '半导体', changePercent: 5.23, leadingStock: '寒武纪 +10.00%', analysis: '受益于AI需求，行业景气度持续上升' },
  { name: '新能源汽车', changePercent: 3.87, leadingStock: '比亚迪 +6.54%', analysis: '销量持续超预期，产业链持续景气' },
  { name: '人工智能', changePercent: 2.95, leadingStock: '科大讯飞 +5.32%', analysis: 'ChatGPT热度不减，应用端加速落地' },
])

const industryDetail = ref({
  name: '半导体行业',
  summary: 'AI算力需求激增，半导体行业迎来新一轮景气周期。国产替代加速，政策支持力度加大。',
  opportunities: ['AI芯片', '存储芯片', '汽车芯片'],
  risks: ['国际贸易摩擦', '技术壁垒', '产能过剩风险']
})

async function handleRefresh() {
  return new Promise(resolve => setTimeout(resolve, 1500))
}
</script>

<template>
  <div class="page">
    <MPullRefresh :on-refresh="handleRefresh">
      <div class="container">
        <!-- 行业排名 -->
        <IndustryCard :industries="industries" :max-show="20" />

        <!-- 行业详情 -->
        <MCard>
          <h3 class="section-title">{{ industryDetail.name }}深度分析</h3>
          <p class="summary">{{ industryDetail.summary }}</p>
          
          <div class="detail-section">
            <h4 class="detail-title">投资机会</h4>
            <div class="tag-list">
              <span v-for="(item, i) in industryDetail.opportunities" :key="i" class="tag tag--rise">
                {{ item }}
              </span>
            </div>
          </div>

          <div class="detail-section">
            <h4 class="detail-title">风险提示</h4>
            <div class="tag-list">
              <span v-for="(item, i) in industryDetail.risks" :key="i" class="tag tag--fall">
                {{ item }}
              </span>
            </div>
          </div>
        </MCard>
      </div>
    </MPullRefresh>
  </div>
</template>

<style scoped>
.page { height: 100%; overflow: hidden; }
.container { padding: var(--m-space-md); display: flex; flex-direction: column; gap: var(--m-space-md); }
.section-title { font-size: var(--m-font-lg); font-weight: var(--m-font-weight-medium); margin-bottom: var(--m-space-md); }
.summary { color: var(--m-text-secondary); line-height: var(--m-line-height-loose); margin-bottom: var(--m-space-lg); }
.detail-section { margin-bottom: var(--m-space-lg); }
.detail-section:last-child { margin-bottom: 0; }
.detail-title { font-weight: var(--m-font-weight-medium); margin-bottom: var(--m-space-md); }
.tag-list { display: flex; flex-wrap: wrap; gap: var(--m-space-sm); }
.tag { padding: var(--m-space-xs) var(--m-space-md); border-radius: var(--m-radius-sm); font-size: var(--m-font-sm); }
.tag--rise { background: var(--m-color-rise-light); color: var(--m-color-rise); }
.tag--fall { background: var(--m-color-fall-light); color: var(--m-color-fall); }
</style>
