<script setup>
import { ref, computed, onBeforeMount, onBeforeUnmount } from 'vue'
import { GlobalStockIndexes } from '../../../api/app'
import { registerFeed, stopFeed } from '../../../api/scheduler'
import { anyOpen } from '../../../api/marketClock'
import MPullRefresh from '../../components/base/MPullRefresh.vue'
import MTabs from '../../components/base/MTabs.vue'
import MEmpty from '../../components/base/MEmpty.vue'
import MLoading from '../../components/base/MLoading.vue'
import PercentTag from '../../components/widgets/PercentTag.vue'
import MIcon from '../../components/base/MIcon.vue'

// 对齐桌面端 market.vue「全球股指 → 全球指数」子页：
// GlobalStockIndexes() 返回 { common, asia, america, europe, other } 五组，
// 每项字段：name/location/zxj(最新点位)/zdf(涨跌幅)/state(open|close|break)/img(国旗)/code/qtcode。
// 桌面端按区域分栏；移动端窄屏改为「区域分段 Tab + 列表」，区域顺序与桌面 getAreaName 一致。
const REGIONS = [
  { label: '常用', value: 'common' },
  { label: '亚洲', value: 'asia' },
  { label: '美洲', value: 'america' },
  { label: '欧洲', value: 'europe' },
  { label: '其他', value: 'other' },
]

const activeRegion = ref('common')
const globalData = ref({}) // { regionKey: item[] }
const loading = ref(false)
// 国旗图加载失败的回退标记（按 code/qtcode/name 去重）
const imgFailed = ref({})

const currentList = computed(() => {
  const list = globalData.value[activeRegion.value]
  if (!Array.isArray(list)) return []
  return list.filter(it => it && typeof it === 'object')
})

async function loadData(silent = false) {
  try {
    const res = await GlobalStockIndexes()
    globalData.value = res && typeof res === 'object' ? res : {}
  } catch (e) {
    console.error('加载全球指数失败:', e)
  } finally {
    if (!silent) loading.value = false
  }
}

async function handleRefresh() {
  loading.value = true
  await loadData(true)
}

// 数字解析：zdf/zxj 原始可能是字符串或数字，统一转 number
function toNum(v) {
  const n = Number(v)
  return Number.isFinite(n) ? n : NaN
}

// 涨跌方向：red 涨 / green 跌 / gray 平
function dirClass(zdf) {
  const n = toNum(zdf)
  if (isNaN(n) || n === 0) return 'm-flat'
  return n > 0 ? 'm-rise' : 'm-fall'
}

// 开休市：open→开市(绿)，close→收盘(灰)，break→休市(灰)；对齐桌面端开市/休市语义并细化
function stateInfo(s) {
  switch (String(s || '').toLowerCase()) {
    case 'open': return { text: '开市', open: true }
    case 'close': return { text: '收盘', open: false }
    case 'break': return { text: '休市', open: false }
    default: return { text: '--', open: false }
  }
}

// 最新点位格式化：千分位 + 2 位小数
function fmtPrice(v) {
  const n = toNum(v)
  if (isNaN(n)) return v ?? '--'
  return n.toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
}

function flagKey(item) {
  return item.code || item.qtcode || item.name || Math.random()
}

function onImgError(item) {
  imgFailed.value[flagKey(item)] = true
}

onBeforeMount(() => {
  loading.value = true
  loadData()
  // 对齐桌面端 market.index feed：开市时段 3s 轮询（静默，不触发 loading 闪烁）
  registerFeed('mobile-global-index', {
    fetch: () => loadData(true),
    intervalMs: 3000,
    activeWhen: () => anyOpen.value,
  })
})

onBeforeUnmount(() => {
  stopFeed('mobile-global-index')
})
</script>

<template>
  <div class="global-index-page">
    <div class="region-tabs">
      <MTabs v-model="activeRegion" :tabs="REGIONS" />
    </div>

    <MPullRefresh class="index-refresh" :on-refresh="handleRefresh">
      <div class="container">
        <div v-if="loading && !currentList.length" class="loading">
          <MLoading text="加载中..." vertical />
        </div>

        <template v-else>
          <div v-if="currentList.length" class="index-list">
            <div v-for="(item, i) in currentList" :key="item.code || item.qtcode || i" class="index-item">
              <!-- 国旗图，加载失败回退 🌐 -->
              <div class="index-flag">
                <img
                  v-if="item.img && !imgFailed[flagKey(item)]"
                  :src="item.img"
                  alt=""
                  loading="lazy"
                  @error="onImgError(item)"
                />
                <span v-else class="flag-fallback"><MIcon name="globe" :size="14" /></span>
              </div>

              <div class="index-main">
                <div class="index-name">{{ item.name || item.code }}</div>
                <div class="index-location">{{ item.location || item.code }}</div>
              </div>

              <div class="index-right">
                <div class="index-price" :class="dirClass(item.zdf)">{{ fmtPrice(item.zxj) }}</div>
                <div class="index-bottom">
                  <PercentTag :value="item.zdf" size="small" />
                  <span class="state-badge" :class="{ 'state-badge--open': stateInfo(item.state).open }">
                    <span class="state-dot" />
                    {{ stateInfo(item.state).text }}
                  </span>
                </div>
              </div>
            </div>
          </div>
          <MEmpty v-else description="该区域暂无数据" />
        </template>
      </div>
    </MPullRefresh>
  </div>
</template>

<style scoped>
.global-index-page {
  height: 100%;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.region-tabs {
  background: var(--m-bg-card);
  border-bottom: 1px solid var(--m-divider-color);
}

/* MPullRefresh 占满 tabs 之外的剩余高度，否则其 height:100% 与顶部 Tab 叠加后
   超出父容器、底部被 overflow:hidden 裁切，列表滚不到底。 */
.index-refresh {
  flex: 1;
  min-height: 0;
}

.container {
  padding: var(--m-space-md);
}

.loading {
  display: flex;
  justify-content: center;
  padding: var(--m-space-2xl) 0;
}

/* 指数列表（卡片样式，对齐 HotStockPage 列表） */
.index-list {
  background: var(--m-bg-card);
  border-radius: var(--m-radius-md);
  overflow: hidden;
}

.index-list > :not(:last-child) {
  border-bottom: 1px solid var(--m-divider-color);
}

.index-item {
  display: flex;
  align-items: center;
  gap: var(--m-space-md);
  padding: var(--m-space-md);
}

.index-flag {
  width: 24px;
  height: 24px;
  border-radius: var(--m-radius-sm);
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--m-bg-primary);
  overflow: hidden;
}

.index-flag img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.flag-fallback {
  font-size: 14px;
  line-height: 1;
}

.index-main {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.index-name {
  font-size: var(--m-font-md);
  font-weight: var(--m-font-weight-medium);
  color: var(--m-text-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.index-location {
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.index-right {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: var(--m-space-xs);
  flex-shrink: 0;
}

.index-price {
  font-size: var(--m-font-lg);
  font-weight: var(--m-font-weight-bold);
  font-variant-numeric: tabular-nums;
  color: var(--m-text-primary);
}

.index-price.m-rise {
  color: var(--m-color-rise);
}

.index-price.m-fall {
  color: var(--m-color-fall);
}

.index-price.m-flat {
  color: var(--m-color-gray);
}

.index-bottom {
  display: flex;
  align-items: center;
  gap: var(--m-space-sm);
}

.state-badge {
  display: inline-flex;
  align-items: center;
  gap: 3px;
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
}

.state-dot {
  width: 6px;
  height: 6px;
  border-radius: var(--m-radius-full);
  background: var(--m-text-tertiary);
}

.state-badge--open {
  color: var(--m-color-fall);
}

.state-badge--open .state-dot {
  background: var(--m-color-fall);
}
</style>
