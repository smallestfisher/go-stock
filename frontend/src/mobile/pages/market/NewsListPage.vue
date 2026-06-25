<script setup>
import { ref } from 'vue'
import MPullRefresh from '../../components/base/MPullRefresh.vue'
import MCard from '../../components/base/MCard.vue'
import NewsCard from '../../components/cards/NewsCard.vue'
import MEmpty from '../../components/base/MEmpty.vue'

// 模拟新闻数据
const newsData = ref([
  { title: '央行宣布降准0.5个百分点，释放长期资金约1万亿元', summary: '此次降准将于12月15日实施，旨在保持流动性合理充裕。', time: new Date(Date.now() - 600000), source: '财联社' },
  { title: 'A股三大指数集体低开，半导体板块领跌', summary: '沪指低开0.3%，深成指低开0.5%，创业板指低开0.8%。', time: new Date(Date.now() - 1800000), source: '证券时报' },
  { title: '外资净流入50亿元，连续五日加仓A股', summary: '北向资金持续流入，主要流向白酒、新能源等板块。', time: new Date(Date.now() - 3600000), source: '第一财经' },
  { title: '美联储宣布维持利率不变，市场反应平淡', summary: '市场普遍预期美联储将在明年开始降息。', time: new Date(Date.now() - 7200000), source: '华尔街日报' },
  { title: '新能源汽车销量创新高，产业链持续景气', summary: '11月新能源汽车销量同比增长40%，渗透率超过35%。', time: new Date(Date.now() - 10800000), source: '中国汽车工业协会' },
  { title: 'ChatGPT概念股大涨，AI芯片需求激增', summary: '多家AI芯片公司订单爆满，产能持续紧张。', time: new Date(Date.now() - 14400000), source: '科创板日报' },
  { title: '房地产政策再放松，多地取消限购', summary: '一线城市外围区域陆续取消限购，市场活跃度提升。', time: new Date(Date.now() - 18000000), source: '经济观察报' },
  { title: '半导体行业回暖，多家公司上调业绩预期', summary: '受益于AI需求，半导体行业景气度明显回升。', time: new Date(Date.now() - 21600000), source: '集微网' },
])

const refreshCount = ref(0)
const loadingMore = ref(false)

// 下拉刷新
async function handleRefresh() {
  return new Promise(resolve => {
    setTimeout(() => {
      refreshCount.value++
      // TODO: 调用真实API刷新数据
      resolve()
    }, 1500)
  })
}

// 加载更多
function handleLoadMore() {
  if (loadingMore.value) return
  loadingMore.value = true
  setTimeout(() => {
    // TODO: 加载更多数据
    loadingMore.value = false
  }, 1000)
}

// 点击新闻
function handleNewsClick(news) {
  console.log('点击新闻:', news)
  // TODO: 打开新闻详情
}
</script>

<template>
  <div class="news-list-page">
    <MPullRefresh :on-refresh="handleRefresh">
      <div class="news-list-container">
        <MCard v-if="newsData.length" padding="none">
          <div class="news-list">
            <NewsCard
              v-for="(news, index) in newsData"
              :key="index"
              :news="news"
              show-summary
              @click="handleNewsClick"
            />
          </div>

          <!-- 加载更多 -->
          <div class="load-more" @click="handleLoadMore">
            <span v-if="!loadingMore">加载更多...</span>
            <span v-else>加载中...</span>
          </div>
        </MCard>

        <MEmpty v-else description="暂无快讯" />
      </div>
    </MPullRefresh>
  </div>
</template>

<style scoped>
.news-list-page {
  height: 100%;
  overflow: hidden;
}

.news-list-container {
  padding: var(--m-space-md);
}

.news-list {
  padding: var(--m-space-md);
}

.load-more {
  text-align: center;
  padding: var(--m-space-lg);
  color: var(--m-text-secondary);
  font-size: var(--m-font-sm);
  cursor: pointer;
  border-top: 1px solid var(--m-divider-color);
}

.load-more:active {
  opacity: 0.6;
}
</style>
