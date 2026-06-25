<script setup>
import { ref } from 'vue'
import MPullRefresh from '../../components/base/MPullRefresh.vue'
import HotTopicCard from '../../components/cards/HotTopicCard.vue'

const topics = ref([
  { title: 'AI芯片', heat: '1.2M', changePercent: 5.67, stocks: 23 },
  { title: '新能源汽车', heat: '980K', changePercent: 3.45, stocks: 45 },
  { title: 'ChatGPT概念', heat: '850K', changePercent: -2.12, stocks: 18 },
  { title: '储能概念', heat: '720K', changePercent: 4.23, stocks: 31 },
  { title: '氢能源', heat: '650K', changePercent: 2.89, stocks: 19 },
])

async function handleRefresh() { return new Promise(resolve => setTimeout(resolve, 1500)) }
</script>

<template>
  <div class="page">
    <MPullRefresh :on-refresh="handleRefresh">
      <div class="container">
        <HotTopicCard :topics="topics" @view-more="() => {}" />
      </div>
    </MPullRefresh>
  </div>
</template>

<style scoped>
.page { height: 100%; overflow: hidden; }
.container { padding: var(--m-space-md); }
</style>
