<script setup>
import { ref, computed, watch, onBeforeMount, onBeforeUnmount } from 'vue'
import {
  GetFollowList, GetGroupList, GetStockList, Follow, UnFollow, Greet,
  AddGroup, RemoveGroup, AddStockGroup, RemoveStockGroup
} from '../../../api/app'
import { registerFeed, stopFeed } from '../../../api/scheduler'
import { EventsOn, EventsOff } from '../../../api/runtime'
import PageHeader from '../../components/widgets/PageHeader.vue'
import MPullRefresh from '../../components/base/MPullRefresh.vue'
import MTabs from '../../components/base/MTabs.vue'
import MEmpty from '../../components/base/MEmpty.vue'
import MButton from '../../components/base/MButton.vue'
import StockCard from '../../components/cards/StockCard.vue'
import VirtualList from '../../components/widgets/VirtualList.vue'
import StockDetailSheet from '../../components/sheets/StockDetailSheet.vue'
import MIcon from '../../components/base/MIcon.vue'

// 字段容错映射（对齐 HomePage loadStockData，应对后端字段大小写不一致）
function pick(item, keys, fallback) {
  for (const key of keys) {
    if (item && item[key] !== undefined && item[key] !== null && item[key] !== '') return item[key]
  }
  return fallback
}

// 判断两个股票列表是否内容相同（避免空列表轮询时无谓重渲染导致页面闪烁）
function sameStocks(a, b) {
  if (!Array.isArray(a) || !Array.isArray(b) || a.length !== b.length) return false
  return a.every((s, i) => s.code === b[i].code)
}

// 分组列表（含"全部"，value=0 对齐桌面端 stock.vue 默认全部分组）
const groups = ref([{ label: '全部', value: 0 }])
const activeGroup = ref(0)
// 原始分组（不含"全部"）：管理分组（删除）和"加入分组"菜单需要 ID/name
const rawGroups = ref([])      // [{ id, name }]

// 自选股票列表（已映射成 StockCard 字段）
const stocksData = ref([])
const currentStocks = computed(() => stocksData.value)

const loading = ref(false)
const selectedStock = ref(null)
const detailVisible = ref(false)

// 加载分组列表
async function loadGroups() {
  try {
    const result = await GetGroupList()
    if (Array.isArray(result)) {
      // 原始分组（含 ID/name，用于管理分组的删除、归类操作）
      rawGroups.value = result.map(g => ({
        id: g.ID ?? g.id ?? 0,
        name: pick(g, ['Name', 'name', '名称'], '未命名')
      }))
      groups.value = [
        { label: '全部', value: 0 },
        ...rawGroups.value.map(g => ({ label: g.name, value: g.id }))
      ]
    }
  } catch (e) {
    console.error('加载分组失败:', e)
  }
}

// 加载自选股票列表（数据未变化时不重新赋值，避免空列表轮询导致闪烁）
// silent=true 时为后台轮询，不翻转 loading（否则空态 MEmpty 会每 10s 隐藏/显示 = 闪烁）
async function loadStocks(silent = false) {
  if (!silent) loading.value = true
  try {
    const result = await GetFollowList(activeGroup.value)
    if (Array.isArray(result)) {
      const next = result.map(stock => ({
        code: pick(stock, ['StockCode', 'stockCode', 'code'], ''),
        name: pick(stock, ['Name', 'StockName', 'stockName', 'name'], '未命名股票'),
        price: Number(pick(stock, ['Price', 'price', 'currentPrice', '当前价格'], 0)) || 0,
        changePercent: Number(pick(stock, ['ChangePercent', 'changePercent', 'change_percent'], 0)) || 0,
        changeAmount: Number(pick(stock, ['PriceChange', 'ChangePrice', 'changeAmount', 'priceChange', '涨跌额'], 0)) || 0,
        open: Number(pick(stock, ['Open', '今日开盘价'], 0)) || 0,
        preClose: Number(pick(stock, ['PreClose', '昨日收盘价'], 0)) || 0,
        high: Number(pick(stock, ['High', '今日最高价'], 0)) || 0,
        low: Number(pick(stock, ['Low', '今日最低价'], 0)) || 0,
        volume: Number(pick(stock, ['Volume', 'volume', '成交的股票数'], 0)) || 0,
        turnover: Number(pick(stock, ['Turnover', 'Amount', '成交金额'], 0)) || 0,
        time: pick(stock, ['Time', '时间'], ''),
        profitToday: 0,
        // 设置弹窗预填字段（来自 GetFollowList 的关注记录）
        costPrice: Number(pick(stock, ['CostPrice'], 0)) || 0,
        holdVolume: Number(pick(stock, ['Volume'], 0)) || 0,
        alarmChangePercent: Number(pick(stock, ['AlarmChangePercent'], 0)) || 0,
        alarmPrice: Number(pick(stock, ['AlarmPrice'], 0)) || 0,
        entryPrice: Number(pick(stock, ['EntryPrice'], 0)) || 0,
        takeProfitPrice: Number(pick(stock, ['TakeProfitPrice'], 0)) || 0,
        stopLossPrice: Number(pick(stock, ['StopLossPrice'], 0)) || 0,
        sort: Number(pick(stock, ['Sort'], 0)) || 0
      }))
      // 内容相同（同样的 code 集合）则保留现有对象引用，只更新行情字段；
      // 不同（增删自选）才整体替换。这样既避免空列表闪烁，又保证 fetchRealtime
      // 改的是视图实际绑定的响应式对象。
      if (sameStocks(next, stocksData.value)) {
        // 沿用旧对象，但同步基本字段（名称等可能变）
        stocksData.value.forEach((s, i) => {
          s.name = next[i].name
        })
      } else {
        stocksData.value = next
      }
      // 拉取实时行情填充价格/涨跌（GetFollowList 只存关注时刻快照）
      // 操作 stocksData.value（视图绑定的响应式对象），更新才生效
      fetchRealtime(stocksData.value)
    }
  } catch (e) {
    console.error('加载自选失败:', e)
  } finally {
    if (!silent) loading.value = false
  }
}

// 切换分组时重载列表
watch(activeGroup, () => { loadStocks() })

// 下拉刷新
async function handleRefresh() {
  await loadStocks()
}

// 拉取实时行情并增量更新列表（对齐桌面端 stock.vue 的 Greet + updateData）
// Greet 返回中文 key 的行情对象；拿不到（如非交易时段/抓取受限）则保留快照值不覆盖
async function fetchRealtime(list) {
  if (!Array.isArray(list) || !list.length) return
  await Promise.all(list.map(async (stock) => {
    if (!stock.code) return
    try {
      const rt = await Greet(stock.code)
      if (!rt) return
      const price = Number(rt['当前价格'] || rt['卖一报价']) || 0
      // 只有真正拿到行情才更新（避免空响应把快照价格清零）
      if (price > 0) {
        stock.price = price
        stock.changePercent = Number(rt.changePercent) || 0
        stock.changeAmount = Number(rt.changePrice ?? rt['涨跌额']) || 0
        stock.open = Number(rt['今日开盘价']) || stock.open
        stock.preClose = Number(rt['昨日收盘价']) || stock.preClose
        stock.high = Number(rt['今日最高价']) || stock.high
        stock.low = Number(rt['今日最低价']) || stock.low
        stock.volume = Number(rt['成交的股票数']) || stock.volume
        stock.turnover = Number(rt['成交金额']) || stock.turnover
        stock.time = rt['时间'] || stock.time
        // 持仓盈亏（需持仓成本，对齐桌面端 Greet 返回字段）
        stock.costPrice = Number(rt.costPrice) || stock.costPrice || 0
        stock.costVolume = Number(rt.costVolume) || 0
        stock.profit = Number(rt.profit) || 0           // 总盈亏率 %
        stock.profitAmount = Number(rt.profitAmount) || 0 // 总盈亏额 ¥
        stock.profitToday = Number(rt.profitAmountToday) || 0 // 今日盈亏额 ¥
      }
    } catch (e) {
      // 单只拉取失败不打断其他
    }
  }))
}

// SSE 实时价格增量更新（对齐 HomePage 的 stock_price 处理）
function onStockPrice(data) {
  if (!data) return
  const code = data['股票代码']
  const price = data['当前价格'] || data['卖一报价']
  if (!code) return
  const target = stocksData.value.find(s => s.code === code)
  if (target && price > 0) {
    target.price = Number(price) || target.price
    target.changePercent = data.changePercent || 0
    target.changeAmount = Number(data.changePrice ?? data['涨跌额']) || 0
    target.high = Number(data['今日最高价']) || target.high
    target.low = Number(data['今日最低价']) || target.low
    target.time = data['时间'] || target.time
  }
}

// 点击股票
function handleStockClick(stock) {
  selectedStock.value = stock
  detailVisible.value = true
}

// 设置保存后：重新拉取列表（成本/持仓等已变）
function handleSaved() {
  loadStocks()
}

// ========== 长按操作菜单：取消关注 ==========
const actionVisible = ref(false)      // 操作菜单弹层
const actionStock = ref(null)         // 当前长按选中的股票
const removingCode = ref('')          // 正在取消关注的 code（防重复点击）
let pressTimer = null                 // 长按计时器
let pressTriggered = false            // 长按是否已触发（用于抑制随后的 click）

// 长按开始：500ms 后弹出操作菜单
function onPressStart(stock) {
  pressTriggered = false
  clearTimeout(pressTimer)
  pressTimer = setTimeout(() => {
    pressTriggered = true
    actionStock.value = stock
    actionVisible.value = true
  }, 500)
}

// 长按取消（手指移开/抬起，未到时长）
function onPressEnd() {
  clearTimeout(pressTimer)
}

// 卡片点击：若刚触发过长按则吞掉这次 click（避免长按后又打开详情）
function onCardClick(stock) {
  if (pressTriggered) {
    pressTriggered = false
    return
  }
  handleStockClick(stock)
}

function closeActionMenu() {
  actionVisible.value = false
  actionStock.value = null
}

// 取消关注：调后端 UnFollow，成功后刷新列表
async function confirmUnfollow() {
  const stock = actionStock.value
  if (!stock || !stock.code || removingCode.value) return
  removingCode.value = stock.code
  try {
    const res = await UnFollow(stock.code)
    // 后端返回提示串（如"取消关注成功"）；无论文案都以重新拉取列表为准
    await loadStocks()
    closeActionMenu()
  } catch (e) {
    console.error('取消关注失败:', e)
    toast.error('取消关注失败')
  } finally {
    removingCode.value = ''
  }
}

// ========== 分组归类：把长按选中的股票加入/移出分组 ==========
const groupPickerVisible = ref(false)   // 「设置分组」选择弹层
const groupBusy = ref(false)            // 归类操作进行中

// 打开分组选择（从长按菜单进入）
function openGroupPicker() {
  actionVisible.value = false
  groupPickerVisible.value = true
}

function closeGroupPicker() {
  groupPickerVisible.value = false
}

// 把股票加入指定分组（对齐桌面端 AddStockGroupInfo）
async function addToGroup(groupId) {
  const stock = actionStock.value
  if (!stock || !stock.code || groupBusy.value) return
  groupBusy.value = true
  try {
    let code = stock.code
    // 对齐桌面端：美股代码前缀转换
    if (code.startsWith('gb_')) code = 'us' + code.replace('gb_', '').toLowerCase()
    await AddStockGroup(groupId, code)
    closeGroupPicker()
    closeActionMenu()
    await loadGroups()
  } catch (e) {
    console.error('加入分组失败:', e)
    toast.error('加入分组失败')
  } finally {
    groupBusy.value = false
  }
}

// 把股票移出当前分组（对齐桌面端 delStockGroup，仅在非「全部」分组下可用）
async function removeFromGroup() {
  const stock = actionStock.value
  if (!stock || !stock.code || groupBusy.value || !activeGroup.value) return
  groupBusy.value = true
  try {
    await RemoveStockGroup(stock.code, stock.name || '', activeGroup.value)
    closeActionMenu()
    await loadStocks()
  } catch (e) {
    console.error('移出分组失败:', e)
    toast.error('移出分组失败')
  } finally {
    groupBusy.value = false
  }
}

// ========== 管理分组：新建 / 删除 ==========
const groupMgrVisible = ref(false)      // 管理分组弹层
const newGroupName = ref('')            // 新建分组名输入
const groupMgrBusy = ref(false)
const pendingDeleteGroupId = ref(null)  // 待二次确认删除的分组 id（再点一次才真正删除）
let pendingDeleteTimer = null

function clearPendingDelete() {
  pendingDeleteGroupId.value = null
  if (pendingDeleteTimer) { clearTimeout(pendingDeleteTimer); pendingDeleteTimer = null }
}

function openGroupManager() {
  groupMgrVisible.value = true
  newGroupName.value = ''
  clearPendingDelete()
}

function closeGroupManager() {
  groupMgrVisible.value = false
  clearPendingDelete()
}

// 新建分组（对齐桌面端 saveTabPane → AddGroup）
async function createGroup() {
  const name = newGroupName.value.trim()
  if (!name || groupMgrBusy.value) return
  groupMgrBusy.value = true
  try {
    // AddGroup 收的是对象 { name, sort }（对齐桌面端 addTabModel），传字符串会建成「未命名」
    await AddGroup({ name, sort: 1 })
    newGroupName.value = ''
    await loadGroups()
  } catch (e) {
    console.error('新建分组失败:', e)
    toast.error('新建分组失败')
  } finally {
    groupMgrBusy.value = false
  }
}

// 删除分组（对齐桌面端 delTab → RemoveGroup）
async function deleteGroup(group) {
  if (groupMgrBusy.value) return
  // 内联二次确认：首次点击进入待确认态，3 秒内再次点击才真正删除（替代原生 confirm）
  if (pendingDeleteGroupId.value !== group.id) {
    pendingDeleteGroupId.value = group.id
    if (pendingDeleteTimer) clearTimeout(pendingDeleteTimer)
    pendingDeleteTimer = setTimeout(() => { pendingDeleteGroupId.value = null }, 3000)
    return
  }
  if (pendingDeleteTimer) { clearTimeout(pendingDeleteTimer); pendingDeleteTimer = null }
  pendingDeleteGroupId.value = null
  groupMgrBusy.value = true
  try {
    await RemoveGroup(Number(group.id))
    // 若删的是当前选中分组，切回「全部」
    if (activeGroup.value === group.id) activeGroup.value = 0
    await loadGroups()
  } catch (e) {
    console.error('删除分组失败:', e)
    toast.error('删除分组失败')
  } finally {
    groupMgrBusy.value = false
  }
}

// ========== 添加自选：底部搜索面板 ==========
// 与桌面端 stock.vue 一致：用本地全量股票库 GetStockList + 前端过滤，
// 不走 SearchStock（那是东财选股语义搜索，需外网+qgqp_b_id 且返回复杂）。
const allStocks = ref([])       // 本地股票库缓存（一次加载，反复过滤）
let allStocksLoaded = false
const searchVisible = ref(false)
const searchKeyword = ref('')
const searchResults = ref([])
const searching = ref(false)
const addingCode = ref('')

// 把后端股票对象规范化为 { code, name }（code 用 ts_code，关注时 Follow 能识别）
function normStock(item) {
  return {
    code: pick(item, ['ts_code', 'code', 'StockCode', 'stockCode', 'symbol'], ''),
    name: pick(item, ['name', 'Name', 'stockName', 'StockName'], '')
  }
}

async function ensureAllStocks() {
  if (allStocksLoaded) return
  searching.value = true
  try {
    const result = await GetStockList('')
    if (Array.isArray(result)) {
      allStocks.value = result.map(normStock).filter(s => s.code && s.name)
      allStocksLoaded = true
    }
  } catch (e) {
    console.error('加载股票库失败:', e)
  } finally {
    searching.value = false
  }
}

function openSearch() {
  searchVisible.value = true
  searchKeyword.value = ''
  searchResults.value = []
  ensureAllStocks()
}

function closeSearch() {
  searchVisible.value = false
}

// 搜索（防抖：输入停止 300ms 后本地过滤）
let searchTimer = null
function onSearchInput() {
  clearTimeout(searchTimer)
  const keyword = searchKeyword.value.trim().toLowerCase()
  if (!keyword) {
    searchResults.value = []
    return
  }
  searchTimer = setTimeout(() => doSearch(keyword), 300)
}

function doSearch(keyword) {
  // 本地过滤：名称或代码包含关键词（对齐桌面端 stock.vue 行 802 的过滤逻辑）
  searchResults.value = allStocks.value
    .filter(s => s.name.toLowerCase().includes(keyword) || s.code.toLowerCase().includes(keyword))
    .slice(0, 50)  // 最多展示 50 条，避免渲染过多
}

// 选中并关注
async function selectAndFollow(stock) {
  if (addingCode.value) return
  addingCode.value = stock.code
  try {
    const res = await Follow(stock.code)
    if (res === '关注成功' || res === '已经关注了') {
      await loadStocks()
      closeSearch()
    } else {
      toast.error(res || '关注失败')
    }
  } catch (e) {
    console.error('关注失败:', e)
    toast.error('关注失败')
  } finally {
    addingCode.value = ''
  }
}

// + 按钮与"添加自选"都打开搜索面板
function handleAddStock() {
  openSearch()
}

onBeforeMount(async () => {
  await loadGroups()
  await loadStocks()
  // 10s 轮询刷新（与 HomePage 一致）
  registerFeed('mobile-stock-list', { fetch: loadStocks, intervalMs: 10000 })
  // 订阅实时价格推送（对齐桌面端 stock.vue 的 stock_price 事件）
  EventsOn('stock_price', onStockPrice)
})

onBeforeUnmount(() => {
  stopFeed('mobile-stock-list')
  EventsOff('stock_price')
})
</script>

<template>
  <div class="stock-list-page">
    <!-- 顶部导航（只保留标题） -->
    <PageHeader title="自选股票" />

    <!-- 分组标签 + 管理分组 + 添加 -->
    <div class="stock-list-tabs">
      <MTabs v-model="activeGroup" :tabs="groups" />
      <button class="tab-action-btn" title="管理分组" @click="openGroupManager"><MIcon name="folder" :size="20" /></button>
      <button class="tab-action-btn" title="添加自选" @click="handleAddStock"><MIcon name="plus" :size="20" /></button>
    </div>

    <!-- 股票列表 -->
    <MPullRefresh :on-refresh="handleRefresh">
      <div class="stock-list-container">
        <!-- 有数据：虚拟滚动列表 -->
        <VirtualList
          v-if="currentStocks.length"
          :items="currentStocks"
          :item-height="128"
          class="stock-list"
        >
          <template #default="{ item }">
            <div
              class="stock-item-wrapper"
              @touchstart.passive="onPressStart(item)"
              @touchend="onPressEnd"
              @touchmove="onPressEnd"
              @touchcancel="onPressEnd"
              @mousedown="onPressStart(item)"
              @mouseup="onPressEnd"
              @mouseleave="onPressEnd"
              @contextmenu.prevent
            >
              <StockCard :stock="item" @click="onCardClick(item)" />
            </div>
          </template>
        </VirtualList>

        <!-- 空状态 -->
        <MEmpty v-else-if="!loading" description="还没有添加自选股票">
          <MButton type="primary" @click="handleAddStock">
            添加自选
          </MButton>
        </MEmpty>
      </div>
    </MPullRefresh>

    <!-- 股票详情抽屉 -->
    <StockDetailSheet
      v-model:show="detailVisible"
      :stock="selectedStock"
      @saved="handleSaved"
    />

    <!-- 添加自选：底部搜索面板 -->
    <transition name="search-sheet">
      <div v-if="searchVisible" class="search-mask" @click.self="closeSearch">
        <div class="search-sheet">
          <div class="search-sheet__header">
            <span class="search-sheet__title">添加自选</span>
            <button class="search-sheet__close" @click="closeSearch"><MIcon name="close" :size="18" /></button>
          </div>
          <div class="search-sheet__input-wrap">
            <input
              v-model="searchKeyword"
              class="search-sheet__input"
              type="search"
              placeholder="输入股票名称/代码/拼音"
              @input="onSearchInput"
            >
            <span v-if="searching" class="search-sheet__hint">搜索中...</span>
          </div>
          <div class="search-sheet__results">
            <div
              v-for="item in searchResults"
              :key="item.code"
              class="search-result-item"
              @click="selectAndFollow(item)"
            >
              <div class="search-result-item__main">
                <span class="search-result-item__name">{{ item.name }}</span>
                <span class="search-result-item__code">{{ item.code }}</span>
              </div>
              <button
                class="search-result-item__add"
                :disabled="!!addingCode"
                @click.stop="selectAndFollow(item)"
              >
                {{ addingCode === item.code ? '...' : '+' }}
              </button>
            </div>
            <MEmpty
              v-if="searchKeyword && !searching && !searchResults.length"
              description="没有找到相关股票"
            />
          </div>
        </div>
      </div>
    </transition>

    <!-- 长按操作菜单：取消关注 -->
    <transition name="action-sheet">
      <div v-if="actionVisible" class="action-mask" @click.self="closeActionMenu">
        <div class="action-sheet">
          <div class="action-sheet__title">
            {{ actionStock?.name || '' }}
            <span class="action-sheet__code">{{ actionStock?.code || '' }}</span>
          </div>
          <button
            v-if="rawGroups.length"
            class="action-sheet__item"
            @click="openGroupPicker"
          >
            设置分组
          </button>
          <button
            v-if="activeGroup"
            class="action-sheet__item"
            :disabled="groupBusy"
            @click="removeFromGroup"
          >
            {{ groupBusy ? '处理中...' : '移出当前分组' }}
          </button>
          <button
            class="action-sheet__item action-sheet__item--danger"
            :disabled="!!removingCode"
            @click="confirmUnfollow"
          >
            {{ removingCode ? '取消中...' : '取消关注' }}
          </button>
          <button class="action-sheet__item action-sheet__cancel" @click="closeActionMenu">
            取消
          </button>
        </div>
      </div>
    </transition>

    <!-- 设置分组：选择把当前股票加入哪个分组 -->
    <transition name="action-sheet">
      <div v-if="groupPickerVisible" class="action-mask" @click.self="closeGroupPicker">
        <div class="action-sheet">
          <div class="action-sheet__title">
            加入分组
            <span class="action-sheet__code">{{ actionStock?.name || '' }}</span>
          </div>
          <div v-if="!rawGroups.length" class="action-sheet__empty">
            还没有分组，先在「管理分组」里新建
          </div>
          <button
            v-for="g in rawGroups"
            :key="g.id"
            class="action-sheet__item"
            :disabled="groupBusy"
            @click="addToGroup(g.id)"
          >
            {{ g.name }}
          </button>
          <button class="action-sheet__item action-sheet__cancel" @click="closeGroupPicker">
            取消
          </button>
        </div>
      </div>
    </transition>

    <!-- 管理分组：新建 / 删除 -->
    <transition name="action-sheet">
      <div v-if="groupMgrVisible" class="action-mask" @click.self="closeGroupManager">
        <div class="action-sheet">
          <div class="action-sheet__title">管理分组</div>
          <div class="group-mgr__create">
            <input
              v-model="newGroupName"
              class="group-mgr__input"
              type="text"
              placeholder="输入新分组名称"
              maxlength="20"
              @keyup.enter="createGroup"
            >
            <button
              class="group-mgr__add"
              :disabled="!newGroupName.trim() || groupMgrBusy"
              @click="createGroup"
            >
              新建
            </button>
          </div>
          <div class="group-mgr__list">
            <div v-if="!rawGroups.length" class="action-sheet__empty">
              暂无分组
            </div>
            <div
              v-for="g in rawGroups"
              :key="g.id"
              class="group-mgr__item"
            >
              <span class="group-mgr__name">{{ g.name }}</span>
              <button
                class="group-mgr__del"
                :class="{ 'group-mgr__del--confirm': pendingDeleteGroupId === g.id }"
                :disabled="groupMgrBusy"
                @click="deleteGroup(g)"
              >
                {{ pendingDeleteGroupId === g.id ? '确认删除?' : '删除' }}
              </button>
            </div>
          </div>
          <button class="action-sheet__item action-sheet__cancel" @click="closeGroupManager">
            关闭
          </button>
        </div>
      </div>
    </transition>
  </div>
</template>

<style scoped>
.stock-list-page {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: var(--m-bg-primary);
}

.stock-list-tabs {
  display: flex;
  align-items: center;
  background: var(--m-bg-card);
  border-bottom: 1px solid var(--m-divider-color);
  position: sticky;
  top: 56px;
  z-index: var(--m-z-sticky);
}

/* 让 MTabs 占满宽度 */
.stock-list-tabs :deep(.m-tabs) {
  flex: 1;
  min-width: 0;
}

.tab-action-btn {
  flex-shrink: 0;
  width: var(--m-touch-min);
  height: var(--m-touch-min);
  display: flex;
  align-items: center;
  justify-content: center;
  background: transparent;
  border: none;
  font-size: 20px;
  color: var(--m-text-primary);
  cursor: pointer;
}

.tab-action-btn:active {
  opacity: 0.6;
}

.stock-list-container {
  flex: 1;
  overflow: hidden;
  position: relative;
}

.stock-list {
  height: 100%;
}

.stock-item-wrapper {
  height: 100%;
  padding: var(--m-space-xs) var(--m-space-md);
  box-sizing: border-box;
}

/* ===== 添加自选搜索面板 ===== */
.search-mask {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.4);
  z-index: var(--m-z-modal, 1000);
  display: flex;
  align-items: flex-end;
}

.search-sheet {
  width: 100%;
  max-height: 75%;
  display: flex;
  flex-direction: column;
  background: var(--m-bg-card);
  border-radius: var(--m-radius-lg, 16px) var(--m-radius-lg, 16px) 0 0;
  overflow: hidden;
}

.search-sheet__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--m-space-md);
  border-bottom: 1px solid var(--m-divider-color);
}

.search-sheet__title {
  font-size: var(--m-font-md);
  font-weight: var(--m-font-weight-bold);
  color: var(--m-text-primary);
}

.search-sheet__close {
  width: var(--m-touch-min);
  height: var(--m-touch-min);
  background: transparent;
  border: none;
  font-size: 18px;
  color: var(--m-text-secondary);
}

.search-sheet__input-wrap {
  position: relative;
  padding: var(--m-space-md);
  border-bottom: 1px solid var(--m-divider-color);
}

.search-sheet__input {
  width: 100%;
  height: 40px;
  padding: 0 var(--m-space-md);
  background: var(--m-bg-primary);
  border: 1px solid var(--m-divider-color);
  border-radius: var(--m-radius-full);
  font-size: var(--m-font-md);
  color: var(--m-text-primary);
  outline: none;
}

.search-sheet__hint {
  position: absolute;
  right: var(--m-space-lg);
  top: 50%;
  transform: translateY(-50%);
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
}

.search-sheet__results {
  flex: 1;
  overflow-y: auto;
  -webkit-overflow-scrolling: touch;
}

.search-result-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--m-space-md);
  border-bottom: 1px solid var(--m-divider-color);
}

.search-result-item:active {
  background: var(--m-bg-primary);
}

.search-result-item__main {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}

.search-result-item__name {
  font-size: var(--m-font-md);
  color: var(--m-text-primary);
  font-weight: var(--m-font-weight-medium);
}

.search-result-item__code {
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
}

.search-result-item__add {
  width: 32px;
  height: 32px;
  flex-shrink: 0;
  border-radius: 50%;
  border: none;
  background: var(--m-color-rise);
  color: #fff;
  font-size: 20px;
  line-height: 1;
}

.search-result-item__add:disabled {
  opacity: 0.5;
}

/* 面板滑入动画 */
.search-sheet-enter-active,
.search-sheet-leave-active {
  transition: opacity 0.25s var(--m-ease-out);
}

.search-sheet-enter-active .search-sheet,
.search-sheet-leave-active .search-sheet {
  transition: transform 0.25s var(--m-ease-out);
}

.search-sheet-enter-from,
.search-sheet-leave-to {
  opacity: 0;
}

.search-sheet-enter-from .search-sheet,
.search-sheet-leave-to .search-sheet {
  transform: translateY(100%);
}

/* ===== 长按操作菜单 ===== */
.action-mask {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.4);
  z-index: var(--m-z-modal, 1000);
  display: flex;
  align-items: flex-end;
}

.action-sheet {
  width: 100%;
  display: flex;
  flex-direction: column;
  gap: var(--m-space-sm);
  padding: var(--m-space-md);
  background: var(--m-bg-primary);
}

.action-sheet__title {
  display: flex;
  align-items: baseline;
  gap: var(--m-space-sm);
  padding: var(--m-space-sm) var(--m-space-md);
  font-size: var(--m-font-md);
  font-weight: var(--m-font-weight-bold);
  color: var(--m-text-primary);
}

.action-sheet__code {
  font-size: var(--m-font-xs);
  font-weight: var(--m-font-weight-normal);
  color: var(--m-text-tertiary);
}

.action-sheet__item {
  width: 100%;
  height: 52px;
  border: none;
  border-radius: var(--m-radius-md, 12px);
  background: var(--m-bg-card);
  font-size: var(--m-font-md);
  color: var(--m-text-primary);
  cursor: pointer;
}

.action-sheet__item:active {
  opacity: 0.7;
}

.action-sheet__item--danger {
  color: #e54d42;
  font-weight: var(--m-font-weight-medium);
}

.action-sheet__item--danger:disabled {
  opacity: 0.5;
}

.action-sheet__cancel {
  margin-top: var(--m-space-xs);
  color: var(--m-text-secondary);
}

/* 操作菜单滑入动画 */
.action-sheet-enter-active,
.action-sheet-leave-active {
  transition: opacity 0.2s var(--m-ease-out);
}

.action-sheet-enter-active .action-sheet,
.action-sheet-leave-active .action-sheet {
  transition: transform 0.2s var(--m-ease-out);
}

.action-sheet-enter-from,
.action-sheet-leave-to {
  opacity: 0;
}

.action-sheet-enter-from .action-sheet,
.action-sheet-leave-to .action-sheet {
  transform: translateY(100%);
}

/* ===== 分组选择 / 管理分组 ===== */
.action-sheet__empty {
  padding: var(--m-space-lg) var(--m-space-md);
  text-align: center;
  font-size: var(--m-font-sm);
  color: var(--m-text-tertiary);
}

/* 新建分组输入行 */
.group-mgr__create {
  display: flex;
  gap: var(--m-space-sm);
  padding: var(--m-space-sm) 0 var(--m-space-md);
}

.group-mgr__input {
  flex: 1;
  min-width: 0;
  height: 40px;
  padding: 0 var(--m-space-md);
  background: var(--m-bg-primary);
  border: 1px solid var(--m-divider-color);
  border-radius: var(--m-radius-sm);
  font-size: var(--m-font-md);
  color: var(--m-text-primary);
  outline: none;
}

.group-mgr__add {
  flex-shrink: 0;
  padding: 0 var(--m-space-lg);
  background: var(--m-color-rise-light, rgba(208, 48, 80, 0.12));
  color: var(--m-color-rise);
  border: none;
  border-radius: var(--m-radius-sm);
  font-size: var(--m-font-sm);
  font-weight: var(--m-font-weight-medium);
  cursor: pointer;
}

.group-mgr__add:disabled {
  opacity: 0.4;
}

/* 已有分组列表 */
.group-mgr__list {
  display: flex;
  flex-direction: column;
  max-height: 40vh;
  overflow-y: auto;
}

.group-mgr__item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--m-space-md) var(--m-space-sm);
  border-bottom: 1px solid var(--m-divider-color);
}

.group-mgr__name {
  font-size: var(--m-font-md);
  color: var(--m-text-primary);
}

.group-mgr__del {
  flex-shrink: 0;
  padding: var(--m-space-xs) var(--m-space-md);
  background: transparent;
  border: 1px solid var(--m-color-fall, #18a058);
  color: var(--m-color-fall, #18a058);
  border-radius: var(--m-radius-sm);
  font-size: var(--m-font-sm);
  cursor: pointer;
}

.group-mgr__del:disabled {
  opacity: 0.4;
}
</style>
