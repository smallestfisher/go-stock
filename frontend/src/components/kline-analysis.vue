<script setup>
import { GetStockList, GetConfig } from '../api/app'
import { EventsOn } from '../api/runtime'
import StockLightweightKlineChart from './StockLightweightKlineChart.vue'
import { NAutoComplete, NButton, NFlex, NText, NInputGroup } from 'naive-ui'
import { useMessage } from 'naive-ui'
import { onBeforeMount, onMounted, onBeforeUnmount, ref } from 'vue'
import {useDevice} from '../composables/useDevice'

const message = useMessage()
const {isMobile} = useDevice()
const searchQuery = ref('')
const selectedCode = ref('000001.SH')
const selectedName = ref('上证指数')
const stockList = ref([])
const options = ref([])
const darkTheme = ref(false)
const chartHeight = ref(window.innerHeight - 230)
const recentStocks = ref([])
const unsupportedCode = ref(false)
const mobileSearchVisible = ref(false)
let stockChangeHandler = null

function toEastMoneyCode(code) {
  if (!code) return ''
  const c = String(code).trim()
  if (/\.(SH|SZ|BJ|HK|US|SS)$/i.test(c)) return c.toUpperCase()
  const lower = c.toLowerCase()
  if (lower.startsWith('sh')) return lower.slice(2) + '.SH'
  if (lower.startsWith('sz')) return lower.slice(2) + '.SZ'
  if (lower.startsWith('bj')) return lower.slice(2) + '.BJ'
  if (lower.startsWith('hk')) return lower.slice(2).toUpperCase() + '.HK'
  if (lower.startsWith('us')) return lower.slice(2).toUpperCase() + '.US'
  if (lower.startsWith('gb_')) return lower.slice(3).toUpperCase() + '.US'
  if (/^\d+$/.test(c)) {
    const d = c[0]
    if (d === '6') return c + '.SH'
    if (d === '0' || d === '3') return c + '.SZ'
    if (d === '8' || d === '9') return c + '.BJ'
    return c + '.SZ'
  }
  // 纯字母代码视为美股（如 AAPL → AAPL.US）
  if (/^[a-zA-Z]+$/.test(c)) return c.toUpperCase() + '.US'
  return ''
}

function findStockList(val) {
  if (!val || !val.trim()) {
    options.value = []
    return
  }
  const q = val.trim().toLowerCase()
  const filtered = stockList.value.filter(item =>
    item.name.toLowerCase().includes(q) ||
    item.ts_code.toLowerCase().includes(q) ||
    String(item.symbol || '').toLowerCase().includes(q) ||
    String(item.cnspell || '').toLowerCase().includes(q)
  ).slice(0, 30)
  options.value = filtered.map(item => ({
    label: item.name + ' - ' + item.ts_code,
    value: item.ts_code,
  }))
}

function normalizeSearchValue(value) {
  const raw = String(value || '').trim()
  if (!raw) return ''
  const parts = raw.split(' - ')
  return parts.length > 1 ? parts[parts.length - 1].trim() : raw
}

function findStockByQuery(value) {
  const raw = normalizeSearchValue(value)
  if (!raw) return null
  const q = raw.toLowerCase()
  return stockList.value.find(item =>
    String(item.ts_code || '').toLowerCase() === q ||
    String(item.name || '').toLowerCase() === q ||
    String(item.symbol || '').toLowerCase() === q ||
    String(item.cnspell || '').toLowerCase() === q
  ) || stockList.value.find(item =>
    String(item.ts_code || '').toLowerCase().startsWith(q + '.') ||
    String(item.ts_code || '').toLowerCase().startsWith(q)
  ) || null
}

function applySearch(value) {
  const raw = normalizeSearchValue(value)
  const found = findStockByQuery(raw)
  const code = found ? found.ts_code : raw
  const emCode = toEastMoneyCode(code)
  if (!emCode) {
    unsupportedCode.value = true
    return
  }
  unsupportedCode.value = false
  selectedCode.value = emCode
  selectedName.value = found ? found.name : ''
  addToRecent(code, selectedName.value)
  searchQuery.value = found ? `${found.name} - ${found.ts_code}` : code
  closeMobileSearch()
}

function handleSearch(value) {
  applySearch(value)
}

function submitSearch() {
  applySearch(searchQuery.value)
}

function addToRecent(code, name) {
  const list = recentStocks.value.filter(s => s.code !== code)
  list.unshift({ code, name })
  if (list.length > 10) list.length = 10
  recentStocks.value = list
  try {
    localStorage.setItem('kline-recent-stocks', JSON.stringify(list))
  } catch {}
}

function loadRecentStocks() {
  try {
    const raw = localStorage.getItem('kline-recent-stocks')
    if (raw) recentStocks.value = JSON.parse(raw)
  } catch {}
}

function selectRecent(code, name) {
  const emCode = toEastMoneyCode(code)
  if (!emCode) {
    unsupportedCode.value = true
    return
  }
  unsupportedCode.value = false
  selectedCode.value = emCode
  selectedName.value = name
  addToRecent(code, name)
  closeMobileSearch()
}

function updateChartHeight() {
  const viewportHeight = Math.round(window.visualViewport?.height || window.innerHeight)
  chartHeight.value = isMobile.value
    ? Math.max(390, Math.min(560, viewportHeight - 220))
    : Math.max(400, window.innerHeight - 230)
}

function openMobileSearch() {
  mobileSearchVisible.value = true
}

function closeMobileSearch() {
  if (isMobile.value) {
    mobileSearchVisible.value = false
  }
}

onBeforeMount(() => {
  GetStockList('').then(result => {
    stockList.value = result || []
  }).catch(err => { console.error('GetStockList error:', err) })
  GetConfig().then(result => {
    darkTheme.value = !!result.darkTheme
  }).catch(err => { console.error('GetConfig error:', err) })
})

onMounted(() => {
  loadRecentStocks()
  updateChartHeight()
  window.addEventListener('resize', updateChartHeight)
  window.visualViewport?.addEventListener('resize', updateChartHeight)

  stockChangeHandler = (data) => {
    if (data && data.ts_code) {
      const emCode = toEastMoneyCode(data.ts_code)
      if (!emCode) {
        unsupportedCode.value = true
        return
      }
      unsupportedCode.value = false
      selectedCode.value = emCode
      selectedName.value = data.name || ''
      addToRecent(data.ts_code, data.name || '')
    }
  }
  EventsOn('klineSelectStock', stockChangeHandler)
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', updateChartHeight)
  window.visualViewport?.removeEventListener('resize', updateChartHeight)
})
</script>

<template>
  <div class="kline-analysis-page" :class="{ 'kline-analysis-page--dark': darkTheme }">
    <div class="kline-title-bar">
      <NText :depth="darkTheme ? 1 : 3" style="font-size: 15px; font-weight: 700">{{ selectedName }}&nbsp;</NText>
      <NText depth="3" style="font-size: 13px">{{ selectedCode }}</NText>
    </div>
    <div class="kline-mobile-chart-stage">
      <StockLightweightKlineChart
        :key="selectedCode"
        :code="selectedCode"
        :stockName="selectedName"
        :darkTheme="darkTheme"
        :chartHeight="chartHeight"
        :realtimeIntervalMs="60000"
      />
    </div>

    <div class="kline-search-bar mobile-kline-search">
      <n-button class="kline-mobile-search-trigger mobile-only" type="primary" block @click="openMobileSearch">
        换股 / 搜索
      </n-button>
      <n-input-group class="kline-desktop-search-control">
        <n-auto-complete
          v-model:value="searchQuery"
          :options="options"
          placeholder="股票名称/代码搜索..."
          clearable
          :on-select="handleSearch"
          @update:value="findStockList"
          @keydown.enter.prevent="submitSearch"
        />
        <n-button type="primary" @click="submitSearch">
          🔍
        </n-button>
      </n-input-group>
      <NFlex v-if="unsupportedCode" align="center" :size="6" style="margin-top: 4px">
        <NText type="warning" style="font-size: 12px">该股票暂不支持K线图</NText>
      </NFlex>
      <div v-if="recentStocks.length && !selectedCode" class="recent-stocks">
        <NText depth="3" style="font-size: 11px; white-space: nowrap">最近:</NText>
        <n-button
          v-for="s in recentStocks.slice(0, 6)"
          :key="s.code"
          size="tiny"
          secondary
          @click="selectRecent(s.code, s.name)"
        >
          {{ s.name }}
        </n-button>
      </div>
    </div>
    <n-drawer
      v-model:show="mobileSearchVisible"
      class="kline-mobile-search-drawer"
      placement="bottom"
      height="48vh"
    >
      <n-drawer-content title="换股 / 搜索" closable>
        <NFlex vertical :size="12">
          <n-input-group>
            <n-auto-complete
              v-model:value="searchQuery"
              :options="options"
              placeholder="输入股票名称或代码"
              clearable
              :on-select="handleSearch"
              @update:value="findStockList"
              @keydown.enter.prevent="submitSearch"
            />
            <n-button type="primary" @click="submitSearch">
              搜索
            </n-button>
          </n-input-group>
          <NText v-if="unsupportedCode" type="warning" style="font-size: 12px">该股票暂不支持K线图</NText>
          <div v-if="recentStocks.length" class="kline-mobile-recent-list">
            <NText depth="3" style="font-size: 12px">最近查看</NText>
            <div class="kline-mobile-recent-list__items">
              <n-button
                v-for="s in recentStocks.slice(0, 8)"
                :key="s.code"
                size="small"
                secondary
                @click="selectRecent(s.code, s.name)"
              >
                {{ s.name || s.code }}
              </n-button>
            </div>
          </div>
        </NFlex>
      </n-drawer-content>
    </n-drawer>
  </div>
</template>

<style scoped>
.kline-analysis-page {
  width: 100%;
  padding: 4px 8px;
  box-sizing: border-box;

  position: relative;
}
.kline-analysis-page--dark {
  background: #0a0a0a;
  color: #e2e8f0;
}
.kline-title-bar {
  padding: 2px 0 4px 0;
}
.kline-search-bar {
  position: fixed;
  bottom: 18px;
  right: 12px;
  z-index: 10;
  width: 320px;
}
.kline-mobile-chart-stage {
  width: 100%;
}
.recent-stocks {
  display: flex;
  align-items: center;
  gap: 4px;
  margin-top: 4px;
  flex-wrap: wrap;
}

@media (max-width: 768px) {
  .kline-analysis-page {
    min-height: calc(100dvh - var(--mobile-bottom-nav-height));
    padding: 4px 4px calc(var(--mobile-bottom-nav-height) + 76px + env(safe-area-inset-bottom));
  }

  .kline-title-bar {
    align-items: center;
    display: flex;
    gap: 4px;
    justify-content: center;
    min-height: 28px;
  }

  .kline-mobile-chart-stage {
    border-radius: 6px;
    overflow: visible;
  }

  .mobile-kline-search {
    bottom: calc(var(--mobile-bottom-nav-height) + 10px + env(safe-area-inset-bottom));
    left: 8px;
    right: 8px;
    width: auto;
  }

  .kline-desktop-search-control,
  .mobile-kline-search > .n-flex,
  .mobile-kline-search > .recent-stocks {
    display: none !important;
  }

  .kline-mobile-search-trigger {
    height: 42px;
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.14);
  }

  .kline-mobile-search-drawer :deep(.n-drawer-content) {
    border-radius: 12px 12px 0 0;
  }

  .kline-mobile-recent-list {
    display: grid;
    gap: 8px;
  }

  .kline-mobile-recent-list__items {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
  }

  .kline-mobile-recent-list__items :deep(.n-button) {
    flex: 1 1 calc(33.333% - 8px);
    min-width: 86px;
  }

  .recent-stocks {
    max-height: 64px;
    overflow: auto;
  }
}
</style>
