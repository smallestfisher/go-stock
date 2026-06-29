<script setup>
import { computed, onBeforeMount, onBeforeUnmount, ref, watch } from 'vue'
import PageHeader from '../../components/widgets/PageHeader.vue'
import MPullRefresh from '../../components/base/MPullRefresh.vue'
import MEmpty from '../../components/base/MEmpty.vue'
import MLoading from '../../components/base/MLoading.vue'
import MSheet from '../../components/base/MSheet.vue'
import MiniSparkline from '../../components/charts/MiniSparkline.vue'
import { registerFeed, stopFeed } from '../../../api/scheduler'
import { anyOpen } from '../../../api/marketClock'
import { toast } from '../../composables/useToast'
import {
  GetFollowedFundPaged,
  GetFundRanking,
  GetFundTop10Holdings,
  GetFundHistoryNetValue,
  SearchFundCodes,
  FollowFund,
  UnFollowFund,
} from '../../../api/app'

// 两个 Tab：基金自选 / 基金排行（对齐桌面端 fund.vue 的 FundFollow + FundRanking）
const tabs = [
  { label: '自选', value: 'follow' },
  { label: '排行', value: 'ranking' },
]
const activeTab = ref('follow')

// ---------- 通用格式化 ----------
function fmtPct(value) {
  const n = Number(value)
  if (!Number.isFinite(n)) return '--'
  return `${n > 0 ? '+' : ''}${n.toFixed(2)}%`
}
function pctClass(value) {
  const n = Number(value)
  if (!Number.isFinite(n) || n === 0) return 'm-flat'
  return n > 0 ? 'm-rise' : 'm-fall'
}
function fmtNum(value, digits = 4) {
  const n = Number(value)
  if (!Number.isFinite(n)) return '--'
  return n.toFixed(digits)
}

// ================= 基金自选 =================
const followLoading = ref(false)
const followList = ref([])
const followKeyword = ref('')
const followPage = ref(1)
const followPageSize = 20
const followTotal = ref(0)
const followTotalPages = ref(1)
let followSearchTimer = null

// 关注基金的「实际/估算净值」与涨跌幅（对齐桌面 FundFollow 取值优先级）
function followNet(info) {
  // 实际净值优先，否则估算净值，否则单位净值
  if (info.netActualRate != null && info.netUnitValue != null) {
    return { value: info.netUnitValue, rate: info.netActualRate, label: '实际净值', date: info.netUnitValueDate }
  }
  const estUnit = info.netEstimatedUnit ?? info.fundBasic?.netEstimatedUnit
  const estRate = info.netEstimatedRate ?? info.fundBasic?.netEstimatedRate
  if (estUnit != null) {
    return { value: estUnit, rate: estRate, label: '估算净值', date: info.netEstimatedUnitTime || '' }
  }
  const unit = info.netUnitValue ?? info.fundBasic?.netUnitValue
  return { value: unit, rate: null, label: '单位净值', date: info.netUnitValueDate ?? info.fundBasic?.netUnitValueDate ?? '' }
}

// 关注基金的多周期涨幅标签
function growthTags(info) {
  const b = info.fundBasic || {}
  return [
    { label: '近1月', val: b.netGrowth1 },
    { label: '近3月', val: b.netGrowth3 },
    { label: '近6月', val: b.netGrowth6 },
    { label: '近1年', val: b.netGrowth12 },
    { label: '今年', val: b.netGrowthYTD },
  ].filter(t => t.val != null)
}

async function loadFollow(reset = false) {
  if (reset) followPage.value = 1
  followLoading.value = true
  try {
    const res = await GetFollowedFundPaged(followPage.value, followPageSize, followKeyword.value || '')
    const items = Array.isArray(res?.items) ? res.items : []
    followList.value = items.filter(it => it && typeof it === 'object')
    followTotal.value = Number(res?.totalCount) || 0
    followTotalPages.value = Number(res?.totalPages) || 1
  } catch (e) {
    console.error('加载关注基金失败:', e)
    followList.value = []
  } finally {
    followLoading.value = false
  }
}

function onFollowSearchInput() {
  if (followSearchTimer) clearTimeout(followSearchTimer)
  followSearchTimer = setTimeout(() => loadFollow(true), 300)
}

function changeFollowPage(delta) {
  const next = followPage.value + delta
  if (next < 1 || next > followTotalPages.value) return
  followPage.value = next
  loadFollow()
}

// ---- 添加关注（搜索基金） ----
const addSheetVisible = ref(false)
const addKeyword = ref('')
const addResults = ref([])
const addSearching = ref(false)
let addSearchTimer = null

function openAddSheet() {
  addKeyword.value = ''
  addResults.value = []
  addSheetVisible.value = true
}

function onAddSearchInput() {
  if (addSearchTimer) clearTimeout(addSearchTimer)
  const kw = addKeyword.value.trim()
  if (!kw) { addResults.value = []; return }
  addSearchTimer = setTimeout(async () => {
    addSearching.value = true
    try {
      const res = await SearchFundCodes(kw)
      addResults.value = Array.isArray(res) ? res : []
    } catch (e) {
      console.error('搜索基金失败:', e)
      addResults.value = []
    } finally {
      addSearching.value = false
    }
  }, 300)
}

async function doFollow(code) {
  try {
    const res = await FollowFund(code)
    toast(res || '已关注')
    addSheetVisible.value = false
    loadFollow(true)
  } catch (e) {
    toast('关注失败')
  }
}

async function doUnFollow(code) {
  try {
    const res = await UnFollowFund(code)
    toast(res || '已取消关注')
    loadFollow()
  } catch (e) {
    toast('取消关注失败')
  }
}

// ================= 基金排行 =================
const rankLoading = ref(false)
const rankList = ref([])
const rankPage = ref(1)
const rankPageSize = 30
const rankTotalPages = ref(1)

// 基金类型筛选（对齐桌面 FundRanking fundType，移动端精简为常用类型；marketType 固定场外 kf）
const rankTypeTabs = [
  { label: '全部', value: 'all' },
  { label: '股票型', value: 'gp' },
  { label: '混合型', value: 'hh' },
  { label: '债券型', value: 'zq' },
  { label: '指数型', value: 'zs' },
  { label: 'QDII', value: 'qdii' },
]
const rankType = ref('all')

// 排序字段（value 严格对齐桌面 FundRanking sortFieldOptions）
const rankSortTabs = [
  { label: '今年来', value: 'jnzf' },
  { label: '日涨幅', value: 'rzdf' },
  { label: '近1周', value: '1yzf' },
  { label: '近1月', value: '1mzf' },
  { label: '近3月', value: '3mzf' },
  { label: '近1年', value: '1nzf' },
]
const rankSort = ref('jnzf')

// 排序字段 → FundRankingItem 取值（用于展示当前排序维度的涨幅）
const sortFieldGetter = {
  jnzf: it => it.ytdGrowth,
  rzdf: it => it.dailyGrowth,
  '1yzf': it => it.weekGrowth,
  '1mzf': it => it.monthGrowth,
  '3mzf': it => it.threeMonthGrowth,
  '1nzf': it => it.yearGrowth,
}

function rankGrowth(item) {
  const getter = sortFieldGetter[rankSort.value] || sortFieldGetter.jnzf
  return getter(item)
}

async function loadRanking(reset = false) {
  if (reset) rankPage.value = 1
  rankLoading.value = true
  try {
    // 第1参 marketType(场外kf)，第2参 fundType(rankType)，对齐桌面 GetFundRanking 签名
    const res = await GetFundRanking('kf', rankType.value, rankSort.value, 'desc', rankPage.value, rankPageSize)
    const items = Array.isArray(res?.items) ? res.items : []
    rankList.value = items.filter(it => it && typeof it === 'object')
    rankTotalPages.value = Number(res?.totalPages) || 1
  } catch (e) {
    console.error('加载基金排行失败:', e)
    rankList.value = []
  } finally {
    rankLoading.value = false
  }
}

function changeRankPage(delta) {
  const next = rankPage.value + delta
  if (next < 1 || next > rankTotalPages.value) return
  rankPage.value = next
  loadRanking()
}

watch(rankType, () => loadRanking(true))
watch(rankSort, () => loadRanking(true))

// ================= 基金详情抽屉 =================
const detailVisible = ref(false)
const detailLoading = ref(false)
const detailFund = ref(null)
const detailHoldings = ref([])
const detailNetValues = ref([])

// 净值走势（取历史净值的 netValue 序列喂 MiniSparkline，时间正序）
const netValueSeries = computed(() => {
  return [...detailNetValues.value].reverse().map(d => Number(d.netValue)).filter(Number.isFinite)
})

async function openDetail(code, name) {
  detailFund.value = { code, name }
  detailHoldings.value = []
  detailNetValues.value = []
  detailVisible.value = true
  detailLoading.value = true
  try {
    const [holdings, netValues] = await Promise.all([
      GetFundTop10Holdings(code),
      GetFundHistoryNetValue(code, 1, 30, '', ''),
    ])
    detailHoldings.value = Array.isArray(holdings) ? holdings : []
    detailNetValues.value = Array.isArray(netValues) ? netValues : []
  } catch (e) {
    console.error('加载基金详情失败:', e)
  } finally {
    detailLoading.value = false
  }
}

// ================= 生命周期 =================
async function handleRefresh() {
  if (activeTab.value === 'follow') await loadFollow(true)
  else await loadRanking(true)
}

watch(activeTab, (tab) => {
  if (tab === 'ranking' && !rankList.value.length) loadRanking()
})

onBeforeMount(() => {
  loadFollow()
  // 关注基金净值/估值随行情变动，开市时段 60s 静默刷新当前 Tab
  registerFeed('mobile-fund', {
    fetch: () => { if (activeTab.value === 'follow') loadFollow() },
    intervalMs: 60_000,
    activeWhen: () => anyOpen.value,
  })
})

onBeforeUnmount(() => {
  stopFeed('mobile-fund')
  if (followSearchTimer) clearTimeout(followSearchTimer)
  if (addSearchTimer) clearTimeout(addSearchTimer)
})
</script>

<template>
  <div class="fund-page">
    <PageHeader title="基金中心" />

    <!-- Tab 切换 -->
    <div class="fund-tabs">
      <button
        v-for="t in tabs"
        :key="t.value"
        type="button"
        class="fund-tab"
        :class="{ 'fund-tab--active': activeTab === t.value }"
        @click="activeTab = t.value"
      >
        {{ t.label }}
      </button>
    </div>

    <!-- ============ 自选 ============ -->
    <template v-if="activeTab === 'follow'">
      <div class="toolbar">
        <input
          v-model="followKeyword"
          class="search-input"
          type="search"
          placeholder="搜索关注基金 名称/代码"
          @input="onFollowSearchInput"
        >
        <button class="add-btn" type="button" @click="openAddSheet">+ 关注</button>
      </div>

      <MPullRefresh class="fund-refresh" :on-refresh="handleRefresh">
        <div class="fund-content">
          <MLoading v-if="followLoading && !followList.length" text="加载关注基金..." vertical />

          <template v-else-if="followList.length">
            <div class="follow-list">
              <div
                v-for="info in followList"
                :key="info.code"
                class="follow-card"
                @click="openDetail(info.code, info.fundBasic?.fullName || info.name)"
              >
                <div class="follow-head">
                  <div class="follow-title">
                    <span class="follow-name">{{ info.fundBasic?.fullName || info.name }}</span>
                    <span class="follow-meta">
                      <span class="tag tag--code">{{ info.code }}</span>
                      <span v-if="info.fundBasic?.type" class="tag tag--type">{{ info.fundBasic.type }}</span>
                    </span>
                  </div>
                  <button class="unfollow-btn" type="button" @click.stop="doUnFollow(info.code)">取消</button>
                </div>

                <div class="follow-net">
                  <div class="net-main">
                    <span class="net-label">{{ followNet(info).label }}</span>
                    <span class="net-value" :class="pctClass(followNet(info).rate)">
                      {{ fmtNum(followNet(info).value) }}
                    </span>
                    <span v-if="followNet(info).rate != null" class="net-rate" :class="pctClass(followNet(info).rate)">
                      {{ fmtPct(followNet(info).rate) }}
                    </span>
                  </div>
                  <span v-if="followNet(info).date" class="net-date">{{ followNet(info).date }}</span>
                </div>

                <div v-if="growthTags(info).length" class="growth-tags">
                  <span
                    v-for="g in growthTags(info)"
                    :key="g.label"
                    class="growth-tag"
                    :class="pctClass(g.val)"
                  >{{ g.label }} {{ fmtPct(g.val) }}</span>
                </div>
              </div>
            </div>

            <!-- 分页 -->
            <div v-if="followTotalPages > 1" class="pager">
              <button :disabled="followPage <= 1" @click="changeFollowPage(-1)">上一页</button>
              <span>{{ followPage }} / {{ followTotalPages }}</span>
              <button :disabled="followPage >= followTotalPages" @click="changeFollowPage(1)">下一页</button>
            </div>
          </template>

          <MEmpty v-else :description="followKeyword ? '没有匹配的关注基金' : '还没有关注基金，点击右上「+ 关注」添加'" />
        </div>
      </MPullRefresh>
    </template>

    <!-- ============ 排行 ============ -->
    <template v-else>
      <div class="rank-filters">
        <div class="chip-scroll">
          <button
            v-for="t in rankTypeTabs"
            :key="t.value"
            type="button"
            class="chip"
            :class="{ 'chip--active': rankType === t.value }"
            @click="rankType = t.value"
          >{{ t.label }}</button>
        </div>
        <div class="chip-scroll">
          <button
            v-for="s in rankSortTabs"
            :key="s.value"
            type="button"
            class="chip chip--sort"
            :class="{ 'chip--active': rankSort === s.value }"
            @click="rankSort = s.value"
          >{{ s.label }}</button>
        </div>
      </div>

      <MPullRefresh class="fund-refresh" :on-refresh="handleRefresh">
        <div class="fund-content">
          <MLoading v-if="rankLoading && !rankList.length" text="加载基金排行..." vertical />

          <template v-else-if="rankList.length">
            <div class="rank-list">
              <div
                v-for="(item, index) in rankList"
                :key="item.code"
                class="rank-row"
                @click="openDetail(item.code, item.name)"
              >
                <div class="rank-idx">{{ (rankPage - 1) * rankPageSize + index + 1 }}</div>
                <div class="rank-main">
                  <div class="rank-name">{{ item.name }}</div>
                  <div class="rank-sub">
                    <span class="tag tag--code">{{ item.code }}</span>
                    <span v-if="item.fundTypeDetail" class="tag tag--type">{{ item.fundTypeDetail }}</span>
                  </div>
                </div>
                <div class="rank-right">
                  <div class="rank-growth" :class="pctClass(rankGrowth(item))">{{ fmtPct(rankGrowth(item)) }}</div>
                  <div v-if="item.netUnitValue != null" class="rank-net">净值 {{ fmtNum(item.netUnitValue) }}</div>
                </div>
              </div>
            </div>

            <div v-if="rankTotalPages > 1" class="pager">
              <button :disabled="rankPage <= 1" @click="changeRankPage(-1)">上一页</button>
              <span>{{ rankPage }} / {{ rankTotalPages }}</span>
              <button :disabled="rankPage >= rankTotalPages" @click="changeRankPage(1)">下一页</button>
            </div>
          </template>

          <MEmpty v-else description="暂无基金排行数据" />
        </div>
      </MPullRefresh>
    </template>

    <!-- ============ 添加关注抽屉 ============ -->
    <MSheet v-model:show="addSheetVisible" title="搜索基金" height="80vh">
      <div class="add-body">
        <input
          v-model="addKeyword"
          class="search-input"
          type="search"
          placeholder="输入基金名称 / 代码"
          @input="onAddSearchInput"
        >
        <MLoading v-if="addSearching" text="搜索中..." vertical />
        <div v-else-if="addResults.length" class="add-results">
          <div
            v-for="r in addResults"
            :key="r.code"
            class="add-item"
          >
            <div class="add-info">
              <div class="add-name">{{ r.name }}</div>
              <div class="add-meta">
                <span class="tag tag--code">{{ r.code }}</span>
                <span v-if="r.type" class="tag tag--type">{{ r.type }}</span>
              </div>
            </div>
            <button class="follow-action" type="button" @click="doFollow(r.code)">关注</button>
          </div>
        </div>
        <MEmpty v-else-if="addKeyword.trim()" description="没有找到匹配的基金" />
        <div v-else class="add-tip">输入基金名称或代码开始搜索</div>
      </div>
    </MSheet>

    <!-- ============ 基金详情抽屉 ============ -->
    <MSheet v-model:show="detailVisible" :title="detailFund?.name || '基金详情'" height="85vh">
      <div class="detail-body">
        <MLoading v-if="detailLoading" text="加载详情..." vertical />
        <template v-else>
          <div class="detail-code">{{ detailFund?.code }}</div>

          <!-- 净值走势 -->
          <div class="detail-section">
            <h4 class="detail-title">近 30 期净值走势</h4>
            <div v-if="netValueSeries.length > 1" class="netvalue-chart">
              <MiniSparkline :data="netValueSeries" :width="320" :height="80" />
            </div>
            <MEmpty v-else description="暂无净值数据" />
          </div>

          <!-- 历史净值列表 -->
          <div v-if="detailNetValues.length" class="detail-section">
            <h4 class="detail-title">历史净值</h4>
            <div class="nv-list">
              <div class="nv-row nv-row--head">
                <span>日期</span><span>单位净值</span><span>日增长</span>
              </div>
              <div v-for="nv in detailNetValues.slice(0, 15)" :key="nv.date" class="nv-row">
                <span>{{ nv.date }}</span>
                <span>{{ fmtNum(nv.netValue) }}</span>
                <span :class="pctClass(nv.dailyGrowth)">{{ fmtPct(nv.dailyGrowth) }}</span>
              </div>
            </div>
          </div>

          <!-- 十大持仓 -->
          <div class="detail-section">
            <h4 class="detail-title">
              十大持仓
              <span v-if="detailHoldings[0]?.quarter" class="quarter">{{ detailHoldings[0].quarter }}</span>
            </h4>
            <div v-if="detailHoldings.length" class="holding-list">
              <div v-for="h in detailHoldings" :key="h.stockCode" class="holding-row">
                <span class="holding-rank">{{ h.rank }}</span>
                <span class="holding-name">{{ h.stockName }}</span>
                <span class="holding-code">{{ h.stockCode }}</span>
                <span v-if="h.changeRate != null" class="holding-rate" :class="pctClass(h.changeRate)">{{ fmtPct(h.changeRate) }}</span>
                <span class="holding-ratio">{{ fmtNum(h.ratio, 2) }}%</span>
              </div>
            </div>
            <MEmpty v-else description="暂无持仓数据" />
          </div>
        </template>
      </div>
    </MSheet>
  </div>
</template>

<style scoped>
.fund-page {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: var(--m-bg-primary);
  overflow: hidden;
}

/* Tab */
.fund-tabs {
  display: flex;
  gap: var(--m-space-sm);
  padding: var(--m-space-sm) var(--m-space-md);
  background: var(--m-bg-card);
  border-bottom: 1px solid var(--m-divider-color);
}
.fund-tab {
  flex: 1;
  padding: var(--m-space-sm) 0;
  border: 1px solid var(--m-border-color);
  border-radius: var(--m-radius-full);
  background: var(--m-bg-card);
  color: var(--m-text-secondary);
  font: inherit;
  font-size: var(--m-font-md);
  font-weight: var(--m-font-weight-medium);
}
.fund-tab--active {
  background: var(--m-color-rise);
  border-color: var(--m-color-rise);
  color: #fff;
}

/* 工具条 */
.toolbar {
  display: flex;
  gap: var(--m-space-sm);
  padding: var(--m-space-sm) var(--m-space-md);
  background: var(--m-bg-card);
  border-bottom: 1px solid var(--m-divider-color);
}
.search-input {
  flex: 1;
  height: 36px;
  padding: 0 var(--m-space-md);
  border: 1px solid var(--m-border-color);
  border-radius: var(--m-radius-md);
  background: var(--m-bg-primary);
  color: var(--m-text-primary);
  font-size: var(--m-font-sm);
  outline: none;
}
.add-btn {
  flex-shrink: 0;
  padding: 0 var(--m-space-md);
  border: 1px solid var(--m-color-rise);
  border-radius: var(--m-radius-md);
  background: var(--m-color-rise-light);
  color: var(--m-color-rise);
  font: inherit;
  font-size: var(--m-font-sm);
  font-weight: var(--m-font-weight-medium);
}

/* 排行筛选 */
.rank-filters {
  background: var(--m-bg-card);
  border-bottom: 1px solid var(--m-divider-color);
  padding: var(--m-space-sm) 0;
}
.chip-scroll {
  display: flex;
  gap: var(--m-space-sm);
  padding: var(--m-space-xs) var(--m-space-md);
  overflow-x: auto;
  scrollbar-width: none;
}
.chip-scroll::-webkit-scrollbar { display: none; }
.chip {
  flex-shrink: 0;
  padding: var(--m-space-xs) var(--m-space-md);
  border: 1px solid var(--m-border-color);
  border-radius: var(--m-radius-full);
  background: var(--m-bg-card);
  color: var(--m-text-secondary);
  font: inherit;
  font-size: var(--m-font-sm);
  white-space: nowrap;
}
.chip--sort { font-size: var(--m-font-xs); }
.chip--active {
  background: var(--m-color-rise-light);
  border-color: var(--m-color-rise);
  color: var(--m-color-rise);
  font-weight: var(--m-font-weight-medium);
}

.fund-refresh { flex: 1; min-height: 0; }
.fund-content {
  padding: var(--m-space-md);
  padding-bottom: calc(var(--m-space-2xl) + var(--m-safe-bottom));
}

/* 关注卡片 */
.follow-list {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-md);
}
.follow-card {
  background: var(--m-bg-card);
  border-radius: var(--m-radius-md);
  padding: var(--m-space-md);
}
.follow-card:active { background: var(--m-bg-primary); }
.follow-head {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: var(--m-space-sm);
}
.follow-title { min-width: 0; flex: 1; }
.follow-name {
  display: block;
  font-size: var(--m-font-md);
  font-weight: var(--m-font-weight-bold);
  color: var(--m-text-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.follow-meta {
  display: flex;
  gap: var(--m-space-xs);
  margin-top: var(--m-space-xs);
}
.tag {
  padding: 1px var(--m-space-sm);
  border-radius: var(--m-radius-sm);
  font-size: var(--m-font-xs);
}
.tag--code { background: var(--m-bg-primary); color: var(--m-text-secondary); border: 1px solid var(--m-divider-color); }
.tag--type { background: var(--m-color-rise-light); color: var(--m-color-rise); }
.unfollow-btn {
  flex-shrink: 0;
  padding: 2px var(--m-space-sm);
  border: 1px solid var(--m-border-color);
  border-radius: var(--m-radius-sm);
  background: transparent;
  color: var(--m-text-tertiary);
  font: inherit;
  font-size: var(--m-font-xs);
}
.follow-net {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  margin-top: var(--m-space-md);
}
.net-main { display: flex; align-items: baseline; gap: var(--m-space-sm); }
.net-label { font-size: var(--m-font-xs); color: var(--m-text-tertiary); }
.net-value {
  font-size: var(--m-font-xl);
  font-weight: var(--m-font-weight-bold);
  font-variant-numeric: tabular-nums;
}
.net-rate { font-size: var(--m-font-sm); font-weight: var(--m-font-weight-medium); }
.net-date { font-size: var(--m-font-xs); color: var(--m-text-tertiary); }
.growth-tags {
  display: flex;
  flex-wrap: wrap;
  gap: var(--m-space-xs);
  margin-top: var(--m-space-sm);
}
.growth-tag {
  padding: 1px var(--m-space-sm);
  border-radius: var(--m-radius-sm);
  font-size: var(--m-font-xs);
  background: var(--m-bg-primary);
  font-variant-numeric: tabular-nums;
}

/* 排行列表 */
.rank-list {
  background: var(--m-bg-card);
  border-radius: var(--m-radius-md);
  overflow: hidden;
}
.rank-row {
  display: flex;
  align-items: center;
  gap: var(--m-space-md);
  padding: var(--m-space-md);
}
.rank-row:not(:last-child) { border-bottom: 1px solid var(--m-divider-color); }
.rank-row:active { background: var(--m-bg-primary); }
.rank-idx {
  width: 24px;
  text-align: center;
  font-size: var(--m-font-md);
  font-weight: var(--m-font-weight-bold);
  color: var(--m-text-tertiary);
  flex-shrink: 0;
}
.rank-main { flex: 1; min-width: 0; }
.rank-name {
  font-size: var(--m-font-md);
  font-weight: var(--m-font-weight-medium);
  color: var(--m-text-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.rank-sub { display: flex; gap: var(--m-space-xs); margin-top: var(--m-space-xs); }
.rank-right { text-align: right; flex-shrink: 0; }
.rank-growth {
  font-size: var(--m-font-md);
  font-weight: var(--m-font-weight-bold);
  font-variant-numeric: tabular-nums;
}
.rank-net { font-size: var(--m-font-xs); color: var(--m-text-tertiary); margin-top: 2px; }

/* 分页 */
.pager {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--m-space-lg);
  padding: var(--m-space-lg) 0;
  font-size: var(--m-font-sm);
  color: var(--m-text-secondary);
}
.pager button {
  padding: var(--m-space-xs) var(--m-space-md);
  border: 1px solid var(--m-border-color);
  border-radius: var(--m-radius-sm);
  background: var(--m-bg-card);
  color: var(--m-text-primary);
  font: inherit;
  font-size: var(--m-font-sm);
}
.pager button:disabled { opacity: 0.4; }

/* 添加抽屉 */
.add-body { padding: var(--m-space-md); }
.add-results { margin-top: var(--m-space-md); display: flex; flex-direction: column; }
.add-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--m-space-md);
  padding: var(--m-space-md) 0;
  border-bottom: 1px solid var(--m-divider-color);
}
.add-info { min-width: 0; flex: 1; }
.add-name {
  font-size: var(--m-font-md);
  color: var(--m-text-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.add-meta { display: flex; gap: var(--m-space-xs); margin-top: var(--m-space-xs); }
.follow-action {
  flex-shrink: 0;
  padding: var(--m-space-xs) var(--m-space-md);
  border: 1px solid var(--m-color-rise);
  border-radius: var(--m-radius-sm);
  background: var(--m-color-rise);
  color: #fff;
  font: inherit;
  font-size: var(--m-font-sm);
}
.add-tip {
  text-align: center;
  padding: var(--m-space-2xl) 0;
  color: var(--m-text-tertiary);
  font-size: var(--m-font-sm);
}

/* 详情抽屉 */
.detail-body { padding: var(--m-space-md); }
.detail-code { font-size: var(--m-font-sm); color: var(--m-text-tertiary); margin-bottom: var(--m-space-md); }
.detail-section { margin-bottom: var(--m-space-xl); }
.detail-title {
  font-size: var(--m-font-md);
  font-weight: var(--m-font-weight-bold);
  color: var(--m-text-primary);
  margin-bottom: var(--m-space-md);
}
.quarter { font-size: var(--m-font-xs); font-weight: var(--m-font-weight-normal); color: var(--m-text-tertiary); margin-left: var(--m-space-sm); }
.netvalue-chart {
  display: flex;
  justify-content: center;
  padding: var(--m-space-md);
  background: var(--m-bg-primary);
  border-radius: var(--m-radius-sm);
}
.nv-list { font-size: var(--m-font-sm); }
.nv-row {
  display: grid;
  grid-template-columns: 1fr 1fr 1fr;
  gap: var(--m-space-sm);
  padding: var(--m-space-sm) 0;
  border-bottom: 1px solid var(--m-divider-color);
  font-variant-numeric: tabular-nums;
}
.nv-row span:nth-child(2), .nv-row span:nth-child(3) { text-align: right; }
.nv-row--head { color: var(--m-text-tertiary); font-size: var(--m-font-xs); }
.holding-list { display: flex; flex-direction: column; }
.holding-row {
  display: flex;
  align-items: center;
  gap: var(--m-space-sm);
  padding: var(--m-space-sm) 0;
  border-bottom: 1px solid var(--m-divider-color);
  font-size: var(--m-font-sm);
}
.holding-rank { width: 20px; color: var(--m-text-tertiary); flex-shrink: 0; }
.holding-name {
  flex: 1;
  min-width: 0;
  color: var(--m-text-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.holding-code { font-size: var(--m-font-xs); color: var(--m-text-tertiary); flex-shrink: 0; }
.holding-rate { font-size: var(--m-font-xs); flex-shrink: 0; font-variant-numeric: tabular-nums; }
.holding-ratio {
  flex-shrink: 0;
  min-width: 48px;
  text-align: right;
  color: var(--m-text-secondary);
  font-variant-numeric: tabular-nums;
}
</style>
