<script setup lang="ts">
import {onBeforeMount, onBeforeUnmount, ref, computed} from 'vue'
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

// 移动端：榜单摘要（几涨几跌、最大涨幅/跌幅）
const mobileSummary = computed(() => {
  const arr = list.value || []
  if (!arr.length) return null
  const up = arr.filter(i => Number(i.percent) > 0).length
  const down = arr.filter(i => Number(i.percent) < 0).length
  const flat = arr.length - up - down
  const sorted = [...arr].sort((a, b) => Number(b.percent) - Number(a.percent))
  const top = sorted[0]
  const bottom = sorted[sorted.length - 1]
  return { total: arr.length, up, down, flat, top, bottom }
})

// 移动端：热度归一化(0~100)，用于热度条宽度
function heatRatio(item) {
  const arr = list.value || []
  if (!arr.length) return 0
  const max = Math.max(...arr.map(i => Number(i.value) || 0))
  if (!max) return 0
  return Math.min(100, Math.round((Number(item.value) || 0) / max * 100))
}

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

    <!-- 移动端：紧凑热门榜单（一屏多看，点击看 K线） -->
    <div v-else class="hot-mobile">
      <!-- 榜单摘要 -->
      <div v-if="mobileSummary" class="hot-summary">
        <div class="hot-summary__item hot-summary__item--up">
          <span class="hot-summary__num">{{ mobileSummary.up }}</span>
          <span class="hot-summary__label">上涨</span>
        </div>
        <div class="hot-summary__item hot-summary__item--flat">
          <span class="hot-summary__num">{{ mobileSummary.flat }}</span>
          <span class="hot-summary__label">平</span>
        </div>
        <div class="hot-summary__item hot-summary__item--down">
          <span class="hot-summary__num">{{ mobileSummary.down }}</span>
          <span class="hot-summary__label">下跌</span>
        </div>
        <div class="hot-summary__extremes">
          <div v-if="mobileSummary.top" class="hot-extreme">
            <span class="hot-extreme__tag bg-error">最热</span>
            <span class="hot-extreme__name">{{ mobileSummary.top.name }}</span>
            <span class="text-error">{{ mobileSummary.top.percent }}%</span>
          </div>
          <div v-if="mobileSummary.bottom" class="hot-extreme">
            <span class="hot-extreme__tag bg-success">最弱</span>
            <span class="hot-extreme__name">{{ mobileSummary.bottom.name }}</span>
            <span class="text-success">{{ mobileSummary.bottom.percent }}%</span>
          </div>
        </div>
      </div>

      <!-- 榜单行 -->
      <div class="hot-rank-list">
        <button
            v-for="(item, idx) in list"
            :key="item.code"
            type="button"
            class="hot-row"
            @click="openKline(item)"
        >
          <span class="hot-row__rank" :class="{'hot-row__rank--top': idx < 3}">{{ idx + 1 }}</span>

          <div class="hot-row__main">
            <div class="hot-row__name-line">
              <span class="hot-row__name">{{ item.name }}</span>
              <span class="hot-row__code">{{ item.code }}</span>
            </div>
            <!-- 热度条 + 热度值 -->
            <div class="hot-row__heat">
              <div class="hot-row__heat-bar">
                <div class="hot-row__heat-fill" :style="{width: heatRatio(item) + '%'}"></div>
              </div>
              <span class="hot-row__heat-val">🔥 {{ item.value }}</span>
            </div>
          </div>

          <!-- 右侧价格 + 涨跌幅色块 -->
          <div class="hot-row__right">
            <span class="hot-row__price text-info">{{ item.current }}</span>
            <span class="hot-row__change" :class="'bg-' + trendColor(item.percent)">
              {{ item.percent }}%
            </span>
          </div>
        </button>
      </div>

      <!-- 明细提示 + 变化箭头(榜单底部小注) -->
      <div v-if="list.length" class="hot-mobile__hint">点击行查看 K 线 · 🔥为热度 · 数字旁箭头为热度/排名变化</div>
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
/* ============ 移动端紧凑热门榜单 ============ */
.hot-mobile {
  padding: 6px 0 calc(var(--safe-bottom) + 8px);
}

/* 摘要条 */
.hot-summary {
  align-items: center;
  background: var(--n-color, #fff);
  border: 1px solid var(--n-border-color, #edf0f5);
  border-radius: 10px;
  display: flex;
  flex-wrap: wrap;
  gap: 4px 10px;
  margin: 0 8px 8px;
  padding: 8px 12px;
}

.hot-summary__item {
  align-items: center;
  display: flex;
  gap: 4px;
}

.hot-summary__num {
  font-size: 18px;
  font-weight: 800;
}

.hot-summary__label {
  color: var(--n-text-color-3, #98a2b3);
  font-size: 11px;
}

.hot-summary__item--up .hot-summary__num { color: #d03050; }
.hot-summary__item--down .hot-summary__num { color: #18a058; }
.hot-summary__item--flat .hot-summary__num { color: var(--n-text-color-2, #666); }

.hot-summary__extremes {
  border-left: 1px solid var(--n-border-color, #eef0f4);
  display: flex;
  flex: 1 1 100%;
  flex-direction: column;
  gap: 2px;
  margin-top: 4px;
  padding-left: 0;
}

.hot-extreme {
  align-items: center;
  display: flex;
  font-size: 12px;
  gap: 6px;
}

.hot-extreme__tag {
  border-radius: 3px;
  color: #fff;
  font-size: 10px;
  padding: 0 5px;
}

.hot-extreme__name {
  color: var(--n-text-color-2, #555);
  flex: 1 1 auto;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* 榜单行 */
.hot-rank-list {
  background: var(--n-color, #fff);
  border: 1px solid var(--n-border-color, #edf0f5);
  border-radius: 10px;
  margin: 0 8px;
  overflow: hidden;
}

.hot-row {
  align-items: center;
  appearance: none;
  background: transparent;
  border: 0;
  border-bottom: 1px solid var(--n-border-color, #f0f1f5);
  color: inherit;
  display: flex;
  font: inherit;
  gap: 10px;
  padding: 9px 10px;
  text-align: left;
  width: 100%;
}

.hot-row:last-child {
  border-bottom: 0;
}

.hot-row:active {
  background: var(--n-color-target, #f5f7fa);
}

.hot-row__rank {
  color: var(--n-text-color-3, #98a2b3);
  flex: 0 0 20px;
  font-size: 13px;
  font-weight: 700;
  text-align: center;
}

.hot-row__rank--top {
  color: #fff;
}

.hot-row__rank--top {
  background: #f0a020;
  border-radius: 4px;
  line-height: 20px;
}

.hot-row__main {
  flex: 1 1 auto;
  min-width: 0;
}

.hot-row__name-line {
  align-items: baseline;
  display: flex;
  gap: 6px;
  margin-bottom: 4px;
}

.hot-row__name {
  font-size: 15px;
  font-weight: 700;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.hot-row__code {
  color: var(--n-text-color-3, #98a2b3);
  flex: 0 0 auto;
  font-size: 11px;
}

.hot-row__heat {
  align-items: center;
  display: flex;
  gap: 6px;
}

.hot-row__heat-bar {
  background: var(--n-color-target, #f0f1f5);
  border-radius: 3px;
  flex: 1 1 auto;
  height: 4px;
  overflow: hidden;
}

.hot-row__heat-fill {
  background: linear-gradient(90deg, #ff8a3d, #f0a020);
  border-radius: 3px;
  height: 100%;
  transition: width .3s;
}

.hot-row__heat-val {
  color: var(--n-text-color-3, #98a2b3);
  flex: 0 0 auto;
  font-size: 11px;
}

.hot-row__right {
  align-items: flex-end;
  display: flex;
  flex-direction: column;
  flex: 0 0 auto;
  gap: 3px;
}

.hot-row__price {
  font-size: 13px;
  font-weight: 600;
}

.hot-row__change {
  border-radius: 4px;
  color: #fff;
  font-size: 13px;
  font-weight: 700;
  min-width: 64px;
  padding: 2px 7px;
  text-align: center;
}

.hot-mobile__hint {
  color: var(--n-text-color-3, #98a2b3);
  font-size: 11px;
  padding: 8px 12px 0;
  text-align: center;
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

.bg-error {
  background: #d03050;
}

.bg-success {
  background: #0f7a43;
}

.hot-kline-wrap {
  padding: 4px 8px 8px;
}
</style>
