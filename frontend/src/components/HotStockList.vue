<script setup lang="ts">
import {onBeforeMount, onBeforeUnmount, ref} from 'vue'
import {HotStock} from "../api/app";
import KLineChart from "./KLineChart.vue";
import {ArrowDown, ArrowUp} from "@vicons/ionicons5";
import {registerFeed, stopFeed} from "../api/scheduler";
import {anyOpen} from "../api/marketClock";
const {marketType}=defineProps(
    {
      marketType: {
        type: String,
        default: '10'
      }
    }
)
const list  = ref([])
const loading = ref(false)
const errorMsg = ref('')


async function fetchHotStock() {
  try {
    loading.value = true
    errorMsg.value = ''
    const res = await HotStock(marketType)
    list.value = res || []
    if (list.value.length === 0) {
      errorMsg.value = '暂无数据，可能为API请求失败或休市'
    }
  } catch (e) {
    console.error('fetchHotStock error:', e)
    errorMsg.value = '请求失败: ' + (e?.message || e || '未知错误')
  } finally {
    loading.value = false
  }
}

onBeforeMount(async () => {
  await fetchHotStock()
  // 任意市场开市(交易时段)才轮询；页面恢复时自动重刷。
  registerFeed("hotStock." + marketType, {
    fetch: fetchHotStock,
    intervalMs: 5000,
    activeWhen: () => anyOpen.value,
  })
})

onBeforeUnmount(()=>{
  stopFeed("hotStock." + marketType)
})

function getMarketCode(item) {
  if (item.exchange	 === 'SZ') {
    return item.code.toLowerCase()
  }
  if (item.exchange	 === 'SH') {
    return item.code.toLowerCase()
  }
  if (item.exchange	 === 'HK') {
    return (item.exchange + item.code).toLowerCase()
  }
  return ("gb_"+item.code).toLowerCase()
}

</script>

<template>
  <n-spin :show="loading" size="small">
    <n-alert v-if="errorMsg && list.length === 0" type="warning" style="margin-bottom: 10px" :bordered="false">
      <template #header>
        {{ errorMsg }}
        <n-button size="tiny" type="primary" @click="fetchHotStock" style="margin-left: 10px">重试</n-button>
      </template>
    </n-alert>

    <!-- 桌面端：宽表格 -->
    <n-table striped size="small">
      <n-thead>
        <n-tr>
          <n-th>股票名称</n-th>
          <n-th>涨跌幅</n-th>
          <n-th>当前价格</n-th>
          <n-th>热度</n-th>
          <n-th>热度变化</n-th>
          <n-th>排名变化</n-th>
        </n-tr>
      </n-thead>
      <n-tbody>
        <n-tr v-for="item in list" :key="item.code">
          <n-td><n-text type="info">
            <n-popover trigger="hover" placement="right">
              <template #trigger>
                <n-tag type="info"  :bordered="false">  {{item.name}} {{item.code}}</n-tag>
              </template>
              <k-line-chart style="width: 800px" :code="getMarketCode(item)" :chart-height="500" :stockName="item.name" :k-days="20" :dark-theme="true"></k-line-chart>
            </n-popover>
          </n-text></n-td>
          <n-td><n-text :type="item.percent>0?'error':'success'">{{item.percent}}%</n-text></n-td>
          <n-td><n-text type="info">{{item.current}}</n-text></n-td>
          <n-td><n-text type="info">{{item.value}}</n-text></n-td>
          <n-td><n-text  :type="item.increment>0?'error':'success'">
            {{item.increment}}
            <n-icon v-if="item.increment>0" :component="ArrowUp"/>
            <n-icon v-else :component="ArrowDown"/>
          </n-text></n-td>
          <n-td>
            <n-text  :type="item.rank_change>0?'error':'success'">
              {{item.rank_change}}
              <n-icon v-if="item.rank_change>0" :component="ArrowUp"/>
              <n-text v-else-if="item.rank_change==0" ></n-text>
              <n-icon v-else :component="ArrowDown"/>
            </n-text>
          </n-td>
        </n-tr>
      </n-tbody>
    </n-table>

  </n-spin>
</template>

<style scoped>
</style>
