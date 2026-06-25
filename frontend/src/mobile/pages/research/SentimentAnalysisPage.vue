<script setup>
import { ref } from 'vue'
import MPullRefresh from '../../components/base/MPullRefresh.vue'
import MCard from '../../components/base/MCard.vue'

const marketSentiment = ref({
  overall: '谨慎乐观',
  score: 65,
  analysis: '市场整体情绪偏乐观，但短期波动较大，建议控制仓位。'
})

const sentimentIndicators = ref([
  { name: '恐慌贪婪指数', value: 62, status: '贪婪' },
  { name: '涨跌家数比', value: 1.8, status: '强势' },
  { name: '成交量比', value: 1.2, status: '活跃' },
  { name: '融资融券余额', value: 18500, status: '增长' },
])

const hotWords = ref([
  { word: 'AI芯片', heat: 95 },
  { word: '新能源', heat: 88 },
  { word: 'ChatGPT', heat: 82 },
  { word: '半导体', heat: 76 },
  { word: '锂电池', heat: 70 },
])

async function handleRefresh() {
  return new Promise(resolve => setTimeout(resolve, 1500))
}
</script>

<template>
  <div class="page">
    <MPullRefresh :on-refresh="handleRefresh">
      <div class="container">
        <!-- 市场情绪 -->
        <MCard>
          <h3 class="section-title">市场情绪</h3>
          <div class="sentiment-overall">
            <div class="sentiment-score">{{ marketSentiment.score }}</div>
            <div class="sentiment-info">
              <div class="sentiment-status">{{ marketSentiment.overall }}</div>
              <div class="sentiment-analysis">{{ marketSentiment.analysis }}</div>
            </div>
          </div>

          <!-- 情绪指标 -->
          <div class="indicator-list">
            <div v-for="(item, i) in sentimentIndicators" :key="i" class="indicator-item">
              <div class="indicator-name">{{ item.name }}</div>
              <div class="indicator-value">{{ item.value }}
                <span class="indicator-status">{{ item.status }}</span>
              </div>
            </div>
          </div>
        </MCard>

        <!-- 热词排行 -->
        <MCard>
          <h3 class="section-title">热词排行</h3>
          <div class="hotword-list">
            <div v-for="(item, i) in hotWords" :key="i" class="hotword-item">
              <div class="hotword-rank">{{ i + 1 }}</div>
              <div class="hotword-name">{{ item.word }}</div>
              <div class="hotword-bar">
                <div class="hotword-bar-fill" :style="{ width: `${item.heat}%` }" />
              </div>
              <div class="hotword-heat">{{ item.heat }}</div>
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
.section-title { font-size: var(--m-font-lg); font-weight: var(--m-font-weight-medium); margin-bottom: var(--m-space-lg); }
.sentiment-overall { display: flex; gap: var(--m-space-lg); align-items: center; margin-bottom: var(--m-space-xl); }
.sentiment-score { font-size: 48px; font-weight: var(--m-font-weight-bold); color: var(--m-color-rise); }
.sentiment-info { flex: 1; }
.sentiment-status { font-size: var(--m-font-lg); font-weight: var(--m-font-weight-medium); margin-bottom: var(--m-space-xs); }
.sentiment-analysis { font-size: var(--m-font-sm); color: var(--m-text-secondary); line-height: var(--m-line-height-normal); }
.indicator-list { display: grid; grid-template-columns: 1fr 1fr; gap: var(--m-space-md); }
.indicator-item { padding: var(--m-space-md); background: var(--m-bg-primary); border-radius: var(--m-radius-sm); }
.indicator-name { font-size: var(--m-font-sm); color: var(--m-text-secondary); margin-bottom: var(--m-space-xs); }
.indicator-value { font-size: var(--m-font-lg); font-weight: var(--m-font-weight-bold); }
.indicator-status { font-size: var(--m-font-xs); color: var(--m-color-rise); margin-left: var(--m-space-xs); }
.hotword-list { display: flex; flex-direction: column; gap: var(--m-space-md); }
.hotword-item { display: grid; grid-template-columns: 24px 80px 1fr 40px; gap: var(--m-space-sm); align-items: center; }
.hotword-rank { font-weight: var(--m-font-weight-bold); color: var(--m-text-tertiary); }
.hotword-name { font-weight: var(--m-font-weight-medium); }
.hotword-bar { height: 8px; background: var(--m-bg-primary); border-radius: var(--m-radius-full); overflow: hidden; }
.hotword-bar-fill { height: 100%; background: var(--m-color-rise); transition: width var(--m-duration-normal); }
.hotword-heat { font-size: var(--m-font-sm); color: var(--m-text-secondary); text-align: right; }
</style>
