<script setup>
import { ref } from 'vue'
import MPullRefresh from '../../components/base/MPullRefresh.vue'
import IndustryCard from '../../components/cards/IndustryCard.vue'

const industries = ref([
  { name: '半导体', changePercent: 5.23, leadingStock: '寒武纪 +10.00%' },
  { name: '新能源汽车', changePercent: 3.87, leadingStock: '比亚迪 +6.54%' },
  { name: '人工智能', changePercent: 2.95, leadingStock: '科大讯飞 +5.32%' },
  { name: '白酒', changePercent: 2.34, leadingStock: '贵州茅台 +2.34%' },
  { name: '医药生物', changePercent: 1.45, leadingStock: '恒瑞医药 +3.21%' },
  { name: '银行', changePercent: 0.67, leadingStock: '工商银行 +1.23%' },
  { name: '房地产', changePercent: -1.23, leadingStock: '万科A -0.89%' },
  { name: '钢铁', changePercent: -2.34, leadingStock: '宝钢股份 -1.56%' },
])

async function handleRefresh() {
  return new Promise(resolve => setTimeout(resolve, 1500))
}
</script>

<template>
  <div class="industry-rank-page">
    <MPullRefresh :on-refresh="handleRefresh">
      <div class="container">
        <IndustryCard :industries="industries" :max-show="20" @view-all="() => {}" />
      </div>
    </MPullRefresh>
  </div>
</template>

<style scoped>
.industry-rank-page { height: 100%; overflow: hidden; }
.container { padding: var(--m-space-md); }
</style>
