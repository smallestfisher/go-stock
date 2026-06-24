<script setup lang="ts">
import {onBeforeMount, onBeforeUnmount, ref} from 'vue'
import {HotStock} from "../api/app";
import KLineChart from "./KLineChart.vue";
import {ArrowDown, ArrowUp} from "@vicons/ionicons5";
import {registerFeed, stopFeed} from "../api/scheduler";
import {anyOpen} from "../api/marketClock";
import {useDevice} from "../composables/useDevice";
import BottomSheet from "./mobile/BottomSheet.vue";

const {isMobile} = useDevice()
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

// 移动端点击查看 K 线的当前标的
const activeKline = ref(null as any)
const klineVisible = ref(false)

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

function trendColor(v){ return Number(v||0) > 0 ? 'error' : 'success' }

function openKline(item){
  activeKline.value = {code: getMarketCode(item), name: item.name}
  klineVisible.value = true
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
    <n-table v-if="!isMobile" striped size="small">
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

    <!-- 移动端：热门股卡片，点击查看 K 线 -->
    <div v-else class="hot-mobile">
      <div class="hot-mobile__hint">点击卡片查看 K 线</div>
      <button
          v-for="item in list"
          :key="item.code"
          type="button"
          class="hot-card"
          :class="'hot-card--' + trendColor(item.percent)"
          @click="openKline(item)"
      >
        <div class="hot-card__head">
          <span class="hot-card__name">{{ item.name }}</span>
          <n-tag size="tiny" :bordered="false">{{ item.code }}</n-tag>
          <span class="hot-card__change" :class="'hot-card__change--' + trendColor(item.percent)">
            {{ item.percent }}%
          </span>
        </div>

        <div class="hot-card__core">
          <div class="hot-metric">
            <span class="hot-metric__label">当前价格</span>
            <span class="hot-metric__value text-info">{{ item.current }}</span>
          </div>
          <div class="hot-metric hot-metric--lead">
            <span class="hot-metric__label">热度</span>
            <span class="hot-metric__value">{{ item.value }}</span>
          </div>
        </div>

        <div class="hot-card__detail">
          <div class="hot-sub">
            <span class="hot-sub__label">热度变化</span>
            <n-text :type="trendColor(item.increment)" class="hot-sub__value">
              {{ item.increment }}
              <n-icon v-if="item.increment>0" :component="ArrowUp"/>
              <n-icon v-else :component="ArrowDown"/>
            </n-text>
          </div>
          <div class="hot-sub">
            <span class="hot-sub__label">排名变化</span>
            <n-text :type="trendColor(item.rank_change)" class="hot-sub__value">
              {{ item.rank_change }}
              <n-icon v-if="item.rank_change>0" :component="ArrowUp"/>
              <n-icon v-else-if="item.rank_change<0" :component="ArrowDown"/>
            </n-text>
          </div>
        </div>
      </button>
      <n-empty v-if="!loading && list.length === 0" description="暂无热门股票数据" style="padding: 40px 0" />

      <!-- K 线抽屉 -->
      <BottomSheet :show="klineVisible" :title="activeKline ? activeKline.name : ''" height="76vh" @update:show="(v) => klineVisible = v">
        <div v-if="activeKline" class="hot-kline-wrap">
          <k-line-chart
              :key="activeKline.code"
              style="width: 100%"
              :code="activeKline.code"
              :chart-height="420"
              :stockName="activeKline.name"
              :k-days="20"
              :dark-theme="true"
          />
        </div>
      </BottomSheet>
    </div>
  </n-spin>
</template>

<style scoped>
/* ============ 移动端热门股卡片 ============ */
.hot-mobile {
  padding: 6px 4px calc(var(--safe-bottom) + 8px);
}

.hot-mobile__hint {
  color: var(--n-text-color-3, #999);
  font-size: 12px;
  padding: 2px 8px 8px;
}

.hot-card {
  appearance: none;
  background: var(--n-color, #fff);
  border: 1px solid var(--n-border-color, #edf0f5);
  border-left: 4px solid #d0d5dd;
  border-radius: 10px;
  box-shadow: 0 1px 2px rgba(15, 23, 42, 0.04);
  color: inherit;
  display: flex;
  flex-direction: column;
  font: inherit;
  gap: 10px;
  margin-bottom: 8px;
  padding: 11px 12px;
  text-align: left;
  width: 100%;
}

.hot-card:active {
  background: var(--n-color-hover, #f8fafc);
}

.hot-card--error {
  border-left-color: #d03050;
}

.hot-card--success {
  border-left-color: #18a058;
}

.hot-card__head {
  align-items: center;
  display: flex;
  gap: 8px;
}

.hot-card__name {
  flex: 1 1 auto;
  font-size: 16px;
  font-weight: 700;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.hot-card__change {
  border-radius: 4px;
  color: #fff;
  flex: 0 0 auto;
  font-size: 13px;
  font-weight: 700;
  min-width: 62px;
  padding: 2px 8px;
  text-align: center;
}

.hot-card__change--error {
  background: #d03050;
}

.hot-card__change--success {
  background: #18a058;
}

.hot-card__core {
  display: grid;
  gap: 8px;
  grid-template-columns: 1fr 1fr;
}

.hot-metric {
  align-items: flex-start;
  background: var(--n-color-target, #f5f7fa);
  border-radius: 8px;
  display: flex;
  flex-direction: column;
  gap: 2px;
  padding: 7px 9px;
}

.hot-metric--lead {
  background: rgba(32, 128, 240, 0.08);
}

.hot-metric__label {
  color: var(--n-text-color-3, #98a2b3);
  font-size: 11px;
}

.hot-metric__value {
  font-size: 16px;
  font-weight: 700;
}

.hot-card__detail {
  border-top: 1px dashed var(--n-divider-color, #eef0f4);
  display: grid;
  gap: 6px 12px;
  grid-template-columns: 1fr 1fr;
  padding-top: 8px;
}

.hot-sub {
  align-items: center;
  display: flex;
  font-size: 12px;
  gap: 4px;
  justify-content: space-between;
}

.hot-sub__label {
  color: var(--n-text-color-3, #98a2b3);
}

.hot-sub__value {
  font-weight: 600;
}

/* 涨跌色 */
.text-error {
  color: #d03050;
}

.text-success {
  color: #18a058;
}

.text-info {
  color: #2080f0;
}

.hot-kline-wrap {
  padding: 4px 8px 8px;
}
</style>
