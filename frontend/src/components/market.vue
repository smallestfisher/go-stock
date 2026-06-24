<script setup>
import * as echarts from "echarts";
import {computed, h, nextTick, onBeforeMount, onBeforeUnmount, onMounted,onUnmounted, ref, watch} from 'vue'
import {
  GetAIResponseResult,
  GetConfig,
  GetIndustryRank,
  GetPromptTemplates,
  GetTelegraphList,
  GlobalStockIndexes,
  ReFleshTelegraphList,
  SaveAIResponseResult,
  SaveAsMarkdown,
  ShareAnalysis,
  SummaryStockNews,
  GetAiConfigs,
} from "../api/app";
import {EventsOff, EventsOn} from "../api/runtime";
import {registerFeed, stopFeed} from "../api/scheduler";
import {anyOpen} from "../api/marketClock";
import NewsList from "./newsList.vue";
import KLineChart from "./KLineChart.vue";
import StockLightweightKlineChart from "./StockLightweightKlineChart.vue";
import { CaretDown, CaretUp, PulseOutline,} from "@vicons/ionicons5";
import {NAvatar, NButton, NFlex, NText, useMessage, useNotification} from "naive-ui";
import {MdPreview} from "md-editor-v3";
import {useRoute} from 'vue-router'
import RankTable from "./rankTable.vue";
import IndustryMoneyRank from "./industryMoneyRank.vue";
import StockResearchReportList from "./StockResearchReportList.vue";
import StockNoticeList from "./StockNoticeList.vue";
import LongTigerRankList from "./LongTigerRankList.vue";
import IndustryResearchReportList from "./IndustryResearchReportList.vue";
import HotStockList from "./HotStockList.vue";
import HotEvents from "./HotEvents.vue";
import HotTopics from "./HotTopics.vue";
import InvestCalendarTimeLine from "./InvestCalendarTimeLine.vue";
import ClsCalendarTimeLine from "./ClsCalendarTimeLine.vue";
import Stockhotmap from "./stockhotmap.vue";
import BKFundFlowChart from "./bkFundFlowChart.vue";
import ConceptFundFlowChart from "./conceptFundFlowChart.vue";
import AnalyzeMartket from "./AnalyzeMartket.vue";

const route = useRoute()
import {useDevice} from "../composables/useDevice";
const {isMobile} = useDevice()
const icon = ref('https://raw.githubusercontent.com/ArvinLovegood/go-stock/master/build/appicon.png');

const message = useMessage()
const notify = useNotification()

// 指数图表高度：桌面按视口算；移动端固定一个可读高度(避免全屏手机算出超高，
// 或减 130 后负数)。移动端用 visualViewport，并夹在 320~560。
function computePanelHeight() {
  const vh = Math.round(window.visualViewport?.height || window.innerHeight)
  if (isMobile.value) {
    return Math.max(320, Math.min(560, vh - 180))
  }
  return Math.max(360, vh - 240)
}
const panelHeight = ref(computePanelHeight())

const telegraphList = ref([])
const sinaNewsList = ref([])
const foreignNewsList = ref([])
const common = ref([])
const america = ref([])
const europe = ref([])
const asia = ref([])
const other = ref([])
const globalStockIndexes = ref(null)
const summaryModal = ref(false)
const summaryBTN = ref(true)
const darkTheme = ref(false)
const httpProxyEnabled = ref(false)
const theme = computed(() => {
  return darkTheme.value ? 'dark' : 'light'
})
const aiSummary = ref(``)
const aiSummaryTime = ref("")
const modelName = ref("")
const chatId = ref("")
const question = ref(``)
const aiConfigId = ref(null)
const sysPromptId = ref(null)
const loading = ref(true)
const analysisStatus = ref('')
const aiConfigs = ref([])
const sysPromptOptions = ref([])
const userPromptOptions = ref([])
const promptTemplates = ref([])
const industryRanks = ref([])
const sort = ref("0")
const nowTab = ref("市场快讯")
// 移动端：分组导航（两级：分类 → 具体功能），避免一行十几个横向滚动标签
const marketMobileGroups = [
  { category: '行情', icon: '📰', tabs: ['市场快讯', '当前热门'] },
  { category: '指数', icon: '🌐', tabs: ['全球股指', '重大指数'] },
  { category: '资金', icon: '💰', tabs: ['行业排名', '个股资金流向', '板块资金流向', '概念资金流向'] },
  { category: '研报', icon: '📊', tabs: ['龙虎榜', '个股研报', '公司公告', '行业研究', '名站优选'] },
]
const marketMobileActiveGroup = ref('行情')
function marketGroupOf(tabName) {
  const g = marketMobileGroups.find(g => g.tabs.includes(tabName))
  return g ? g.category : marketMobileGroups[0].category
}
const marketMobileCurrentTabs = computed(() => {
  const g = marketMobileGroups.find(g => g.category === marketMobileActiveGroup.value)
  return g ? g.tabs : []
})
function selectMarketMobileGroup(category) {
  marketMobileActiveGroup.value = category
  // 切换分类时跳到该分类的第一个子栏，避免内容停留在旧 tab 且无高亮 chip
  const g = marketMobileGroups.find(g => g.category === category)
  if (g && g.tabs.length && !g.tabs.includes(nowTab.value)) {
    updateTab(g.tabs[0])
  }
}
const mdPreviewRef = ref(null)
const aiResultScrollRef = ref(null)
const stockCode= ref('')
const enableTools= ref(true)
const thinkingMode = ref(true)
const treemapRef = ref(null);
let treemapchart =null;

function getIndex() {
  GlobalStockIndexes().then((res) => {
    globalStockIndexes.value = res
    common.value = res["common"]
    america.value = res["america"]
    europe.value = res["europe"]
    asia.value = res["asia"]
    other.value = res["other"]
  })
}

onBeforeMount(() => {
  nowTab.value = route.query.name || "市场快讯"
  marketMobileActiveGroup.value = marketGroupOf(nowTab.value)
  stockCode.value = route.query.stockCode
  GetConfig().then(result => {
    summaryBTN.value = result.openAiEnable
    darkTheme.value = result.darkTheme
    httpProxyEnabled.value = result.httpProxyEnabled
  })
  GetPromptTemplates("", "").then(res => {
    promptTemplates.value = res
    sysPromptOptions.value = promptTemplates.value.filter(item => item.type === '模型系统Prompt')
    userPromptOptions.value = promptTemplates.value.filter(item => item.type === '模型用户Prompt')
  })

  GetAiConfigs().then(res=>{
    aiConfigs.value = res
    aiConfigId.value = res[0].ID
  })
  GetTelegraphList("财联社电报").then((res) => {
    telegraphList.value = res
  })
  GetTelegraphList("新浪财经").then((res) => {
    sinaNewsList.value = res
  })
  GetTelegraphList("外媒").then((res) => {
    foreignNewsList.value = res
  })
  getIndex();
  industryRank();
  // 行情轮询交给统一调度器：任意市场开市(交易时段)才拉取，页面恢复时自动重连重刷。
  registerFeed("market.index", {
    fetch: getIndex,
    intervalMs: 3000,
    activeWhen: () => anyOpen.value,
  });
  registerFeed("market.industry", {
    fetch: () => {
      industryRank();
      ReFlesh("财联社电报");
      ReFlesh("新浪财经");
      ReFlesh("外媒");
    },
    intervalMs: 10000,
    activeWhen: () => anyOpen.value,
  });
})


onBeforeUnmount(() => {
  EventsOff("changeMarketTab")
  EventsOff("newTelegraph")
  EventsOff("newSinaNews")
  EventsOff("summaryStockNews")
  stopFeed("market.index");
  stopFeed("market.industry");
})

onUnmounted(() => {

});
EventsOn("changeMarketTab", async (msg) => {
  //message.info(msg.name)
  console.log(msg.name)
  updateTab(msg.name)
})

EventsOn("newTelegraph", (data) => {
  if (data!=null) {
    for (let i = 0; i < data.length; i++) {
      telegraphList.value.pop()
    }
    telegraphList.value.unshift(...data)
  }
})
EventsOn("newSinaNews", (data) => {
  if (data!=null) {
  for (let i = 0; i < data.length; i++) {
    sinaNewsList.value.pop()
  }
  sinaNewsList.value.unshift(...data)
  }
})
EventsOn("tradingViewNews", (data) => {
  if (data!=null) {
    for (let i = 0; i < data.length; i++) {
      foreignNewsList.value.pop()
    }
    foreignNewsList.value.unshift(...data)
  }
})

//获取页面高度（窗口/方向变化时重算，移动端旋转也覆盖）
window.addEventListener('resize', () => { panelHeight.value = computePanelHeight() })
if (window.visualViewport) {
  window.visualViewport.addEventListener('resize', () => { panelHeight.value = computePanelHeight() })
}
watch(isMobile, () => { panelHeight.value = computePanelHeight() })

// ===================== 移动端"指数"板块：网格选择 + 单图 =====================
// type: 'kline' = KLineChart(echarts), 'lw' = StockLightweightKlineChart
const mobileGlobalIndexList = [
  {name:'上证指数', code:'sh000001', type:'kline'},
  {name:'深证成指', code:'sz399001', type:'kline'},
  {name:'创业板指', code:'sz399006', type:'kline'},
  {name:'恒生指数', code:'hkHSI', type:'kline'},
  {name:'纳斯达克', code:'us.IXIC', type:'kline'},
  {name:'道琼斯', code:'us.DJI', type:'kline'},
  {name:'标普500', code:'us.INX', type:'kline'},
]
const mobileMajorIndexList = [
  {name:'上证指数', code:'000001.SH', type:'lw'},
  {name:'深证指数', code:'399001.SZ', type:'lw'},
  {name:'创业板指', code:'399006.SZ', type:'lw'},
  {name:'恒生指数', code:'100.HSI', type:'lw'},
  {name:'道琼斯', code:'100.DJIA', type:'lw'},
  {name:'标普500', code:'100.SPX', type:'lw'},
  {name:'纳斯达克', code:'100.NDX', type:'lw'},
  {name:'沪深300', code:'000300.SH', type:'lw'},
  {name:'上证50', code:'000016.SH', type:'lw'},
  {name:'中证A500', code:'000510.SH', type:'lw'},
  {name:'中证1000', code:'000852.SH', type:'lw'},
  {name:'科创50', code:'000688.SH', type:'lw'},
  {name:'中证银行', code:'399986.SZ', type:'lw'},
  {name:'中证白酒', code:'399997.SZ', type:'lw'},
]
const mobileGlobalIndex = ref(mobileGlobalIndexList[0])
const mobileMajorIndex = ref(mobileMajorIndexList[0])

function getAreaName(code) {
  switch (code) {
    case "america":
      return "美洲"
    case "europe":
      return "欧洲"
    case "asia":
      return "亚洲"
    case "common":
      return "常用"
    case "other":
      return "其他"
  }
}

// ===================== 移动端"市场快讯"专属数据 =====================
// ① 大盘速览：从全球股指里精选核心指数，按 地区→指数 展平成横滑卡片
const mobileHotIndexCodes = ['sh000001','sz399001','sz399006','hkHSI','us.IXIC','us.DJI','us.INX','us.DX']
const mobileIndexList = computed(() => {
  const g = globalStockIndexes.value
  if (!g) return []
  const all = []
  Object.keys(g).forEach(area => {
    ;(g[area] || []).forEach(item => {
      if (mobileHotIndexCodes.includes(item.qtcode) || mobileHotIndexCodes.includes(item.code)) {
        all.push({...item, area: getAreaName(area)})
      }
    })
  })
  return all
})

// ③ 新闻流来源筛选：全部 / 财联社 / 新浪 / 外媒
const mobileNewsSource = ref('all')
const mobileNewsList = computed(() => {
  const src = mobileNewsSource.value
  const pick = (arr, name) => (arr || []).map(i => ({...i, __source: name}))
  if (src === 'cls') return pick(telegraphList.value, '财联社')
  if (src === 'sina') return pick(sinaNewsList.value, '新浪')
  if (src === 'foreign') return pick(foreignNewsList.value, '外媒')
  // all：混合三源，按时间倒序
  return [
    ...pick(telegraphList.value, '财联社'),
    ...pick(sinaNewsList.value, '新浪'),
    ...pick(foreignNewsList.value, '外媒'),
  ].sort((a, b) => {
    const ta = a.dataTime || a.time || ''
    const tb = b.dataTime || b.time || ''
    return tb.localeCompare(ta)
  })
})
const mobileNewsSourceOptions = computed(() => {
  const opts = [{label:'全部', value:'all'}]
  if (telegraphList.value.length) opts.push({label:`财联社 ${telegraphList.value.length}`, value:'cls'})
  if (sinaNewsList.value.length) opts.push({label:`新浪 ${sinaNewsList.value.length}`, value:'sina'})
  if (foreignNewsList.value.length) opts.push({label:`外媒 ${foreignNewsList.value.length}`, value:'foreign'})
  return opts
})
const sourceTagType = (s) => s === '财联社' ? 'success' : s === '新浪' ? 'info' : 'warning'


function changeIndustryRankSort() {
  if (sort.value === "0") {
    sort.value = "1"
  } else {
    sort.value = "0"
  }
  industryRank()
}

function industryRank() {

  GetIndustryRank(sort.value, 150).then(result => {
    if (result.length > 0) {
      //console.log(result)
      industryRanks.value = result
    } else {
      message.info("暂无数据")
    }
  })
}

function reAiSummary() {
  aiSummary.value = ""
  summaryModal.value = true
  loading.value = true
  analysisStatus.value = "正在连接AI服务..."
  SummaryStockNews(question.value,aiConfigId.value, sysPromptId.value,enableTools.value,thinkingMode.value,"summaryStockNews","")
}

function getAiSummary() {
  summaryModal.value = true
  loading.value = true
  GetAIResponseResult("市场资讯").then(result => {
    loading.value = false
    if (result.content) {
      aiSummary.value = result.content
      question.value = result.question
      loading.value = false

      const date = new Date(result.CreatedAt);
      const year = date.getFullYear();
      const month = String(date.getMonth() + 1).padStart(2, '0');
      const day = String(date.getDate()).padStart(2, '0');
      const hours = String(date.getHours()).padStart(2, '0');
      const minutes = String(date.getMinutes()).padStart(2, '0');
      const seconds = String(date.getSeconds()).padStart(2, '0');
      aiSummaryTime.value = `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`
      modelName.value = result.modelName
    } else {
      aiSummaryTime.value = ""
      aiSummary.value = ""
      modelName.value = ""
      //SummaryStockNews(question.value, sysPromptId.value,enableTools.value)
    }
  })
}

function updateTab(name) {
  if (!name) {
    name = "市场快讯"
  }
  summaryBTN.value = (name === "市场快讯");
  nowTab.value = name
  marketMobileActiveGroup.value = marketGroupOf(name)
}

EventsOn("summaryStockNews", async (msg) => {
  if (msg === "DONE") {
    await SaveAIResponseResult("市场资讯", "市场资讯", aiSummary.value, chatId.value, question.value,aiConfigId.value)
    loading.value = false
    analysisStatus.value = "分析完成"
    message.destroyAll()
    notify.success({
      title: 'AI分析完成',
      content: '市场资讯分析已完成',
      duration: 3000,
    })
    setTimeout(() => {
      analysisStatus.value = ""
    }, 3000)
  } else {
    if (msg.chatId) {
      chatId.value = msg.chatId
    }
    if (msg.question) {
      question.value = msg.question
    }
    if (msg.content || msg.reasoning_content || msg.extraContent) {
      if (!aiSummary.value) {
        analysisStatus.value = "AI正在分析中..."
      }
      loading.value = false
    }
    if (msg.content) {
      aiSummary.value = aiSummary.value + msg.content
    }
    if (msg.reasoning_content) {
      aiSummary.value = aiSummary.value + msg.reasoning_content
    }
    if (msg.extraContent) {
      aiSummary.value = aiSummary.value + msg.extraContent
    }
    if (msg.model) {
      modelName.value = msg.model
    }
    if (msg.time) {
      aiSummaryTime.value = msg.time
    }
    scrollToAiResultBottom()
  }
})

function scrollToAiResultBottom() {
  nextTick(() => {
    requestAnimationFrame(() => {
      const el = aiResultScrollRef.value
      if (el) {
        el.scrollTop = el.scrollHeight
      }
    })
  })
}

async function copyToClipboard() {
  try {
    await navigator.clipboard.writeText(aiSummary.value);
    message.success('分析结果已复制到剪切板');
  } catch (err) {
    message.error('复制失败: ' + err);
  }
}

function saveAsMarkdown() {
  SaveAsMarkdown('市场资讯', '市场资讯').then(result => {
    message.success(result)
  })
}

function share() {
  ShareAnalysis('市场资讯', '市场资讯').then(msg => {
    //message.info(msg)
    notify.info({
      avatar: () =>
          h(NAvatar, {
            size: 'small',
            round: false,
            src: icon.value
          }),
      title: '分享到社区',
      duration: 1000 * 30,
      content: () => {
        return h('div', {
          style: {
            'text-align': 'left',
            'font-size': '14px',
          }
        }, {default: () => msg})
      },
    })
  })
}

function ReFlesh(source) {
  //console.log("ReFlesh:", source)
  ReFleshTelegraphList(source).then(res => {
    if (source === "财联社电报") {
      telegraphList.value = res
    }
    if (source === "新浪财经") {
      sinaNewsList.value = res
    }
    if (source === "外媒") {
      foreignNewsList.value = res
    }
  })
}
</script>

<template>
  <n-card class="market-page-shell">
    <!-- 移动端：分类 + 功能两级菜单（替换原生横向滚动标签栏） -->
    <div v-if="isMobile" class="market-mobile-nav">
      <div class="market-mobile-nav__groups">
        <button
            v-for="g in marketMobileGroups"
            :key="g.category"
            type="button"
            class="market-group-tile"
            :class="{ 'market-group-tile--active': marketMobileActiveGroup === g.category }"
            @click="selectMarketMobileGroup(g.category)"
        >
          <span class="market-group-tile__icon">{{ g.icon }}</span>
          <span class="market-group-tile__name">{{ g.category }}</span>
        </button>
      </div>
      <div class="market-mobile-nav__tabs">
        <button
            v-for="tab in marketMobileCurrentTabs"
            :key="tab"
            type="button"
            class="market-tab-chip"
            :class="{ 'market-tab-chip--active': nowTab === tab }"
            @click="updateTab(tab)"
        >
          {{ tab }}
        </button>
      </div>
    </div>
    <n-tabs :class="{ 'market-mobile-tabs--native-hidden': isMobile }" type="line" animated @update-value="updateTab" :value="nowTab" style="">
      <n-tab-pane name="市场快讯" tab="市场快讯">

        <!-- ============ 桌面端：热词 + 三栏新闻 ============ -->
        <n-grid v-if="!isMobile" :cols="1" :y-gap="0">
          <n-gi>
            <div class="market-desktop-heat-panel desktop-only">
              <AnalyzeMartket :dark-theme="darkTheme" :chart-height="300" :kDays="1" :name="'最近24小时热词'" />
            </div>
          </n-gi>
          <n-gi>
            <n-grid class="market-news-grid market-desktop-news-grid" :cols="foreignNewsList.length?3:2" :y-gap="0">
              <n-gi>
                <news-list :newsList="telegraphList" :header-title="'财联社电报'" @update:message="ReFlesh"></news-list>
              </n-gi>
              <n-gi>
                <news-list :newsList="sinaNewsList" :header-title="'新浪财经'" @update:message="ReFlesh"></news-list>
              </n-gi>
              <n-gi v-if="foreignNewsList.length>0">
                <news-list :newsList="foreignNewsList" :header-title="'外媒'" @update:message="ReFlesh"></news-list>
              </n-gi>
            </n-grid>
          </n-gi>
        </n-grid>

        <!-- ============ 移动端：重新设计的快讯布局 ============ -->
        <div v-else class="mkt-brief">

          <!-- ① 大盘速览：核心指数横滑卡片 -->
          <section class="mkt-section">
            <div class="mkt-section__title">📊 大盘速览</div>
            <div v-if="mobileIndexList.length" class="mkt-index-rail">
              <div
                  v-for="item in mobileIndexList"
                  :key="item.code"
                  class="mkt-index-card"
                  :class="'mkt-index-card--' + (item.zdf>0?'up':'down')"
              >
                <div class="mkt-index-card__head">
                  <n-image :src="item.img" :width="16" preview-disabled />
                  <span class="mkt-index-card__name">{{ item.name }}</span>
                </div>
                <div class="mkt-index-card__price" :class="'text-' + (item.zdf>0?'error':'success')">{{ item.zxj }}</div>
                <div class="mkt-index-card__zdf" :class="'bg-' + (item.zdf>0?'error':'success')">
                  <n-number-animation :precision="2" :from="0" :to="item.zdf"/>%
                </div>
                <div class="mkt-index-card__state">{{ item.state === 'open' ? '开市' : '休市' }}</div>
              </div>
            </div>
            <n-empty v-else description="指数加载中" size="small" style="padding: 16px 0" />
          </section>

          <!-- ② 24h 热词：平铺，不再折叠 -->
          <section class="mkt-section">
            <div class="mkt-section__title">🔥 最近24小时热词</div>
            <AnalyzeMartket :dark-theme="darkTheme" :chart-height="260" :kDays="1" :name="'最近24小时热词'" />
          </section>

          <!-- ③ 资讯流：来源吸顶筛选 + 混合时间流 -->
          <section class="mkt-section">
            <div class="mkt-news-sticky">
              <div class="mkt-section__title">📰 资讯</div>
              <div class="mkt-news-source-bar">
                <button
                    v-for="opt in mobileNewsSourceOptions"
                    :key="opt.value"
                    type="button"
                    class="mkt-source-chip"
                    :class="{ 'mkt-source-chip--active': mobileNewsSource === opt.value }"
                    @click="mobileNewsSource = opt.value"
                >{{ opt.label }}</button>
              </div>
            </div>

            <div class="mkt-news-list">
              <article
                  v-for="(item, idx) in mobileNewsList"
                  :key="(item.ID||'') + '-' + idx"
                  class="mkt-news-item"
                  :class="{'mkt-news-item--red': item.isRed}"
              >
                <div class="mkt-news-item__top">
                  <n-tag size="tiny" :bordered="false" :type="sourceTagType(item.__source)">{{ item.__source }}</n-tag>
                  <span v-if="item.time" class="mkt-news-item__time">{{ item.time }}</span>
                  <n-tag v-if="item.sentimentResult" size="tiny" :bordered="false"
                         :type="item.sentimentResult==='看涨'?'error':item.sentimentResult==='看跌'?'success':'info'">
                    {{ item.sentimentResult }}
                  </n-tag>
                </div>
                <div v-if="item.title" class="mkt-news-item__title" :class="{'text-error': item.isRed}">{{ item.title }}</div>
                <div v-if="item.content" class="mkt-news-item__content">{{ item.content }}</div>
                <div v-if="item.subjects || item.stocks || item.url" class="mkt-news-item__tags">
                  <n-tag v-for="sub in (item.subjects||[])" :key="'s'+sub" :bordered="false" type="success" size="tiny">{{ sub }}</n-tag>
                  <n-tag v-for="sub in (item.stocks||[])" :key="'k'+sub" :bordered="false" type="warning" size="tiny">{{ sub }}</n-tag>
                  <a v-if="item.url" :href="item.url" target="_blank" class="mkt-news-item__link">原文 ›</a>
                </div>
              </article>
              <n-empty v-if="!mobileNewsList.length" description="暂无资讯" size="small" style="padding: 24px 0" />
            </div>
          </section>

        </div>

      </n-tab-pane>
      <n-tab-pane name="全球股指" tab="全球股指">
        <!-- 移动端：指数网格选择 + 选中指数单图全宽 -->
        <div v-if="isMobile" class="mkt-index-mobile">
          <div class="mkt-index-mobile__title">选择指数</div>
          <div class="mkt-index-grid">
            <button
                v-for="item in mobileGlobalIndexList"
                :key="item.code"
                type="button"
                class="mkt-index-grid__tile"
                :class="{'mkt-index-grid__tile--active': mobileGlobalIndex.code === item.code}"
                @click="mobileGlobalIndex = item"
            >{{ item.name }}</button>
          </div>
          <div class="mkt-index-chart-wrap">
            <k-line-chart
                :key="mobileGlobalIndex.code"
                :code="mobileGlobalIndex.code"
                :chart-height="panelHeight"
                :stockName="mobileGlobalIndex.name"
                :k-days="20"
                :dark-theme="true"
            />
          </div>
        </div>
        <n-tabs v-else type="segment" animated>
          <n-tab-pane name="全球指数" tab="全球指数">
            <n-grid class="market-mobile-scroll" :cols="5" :y-gap="0">
              <n-gi v-for="(val, key) in globalStockIndexes" :key="key">
                <n-list class="market-mobile-index-card" bordered>
                  <template #header>
                    {{ getAreaName(key) }}
                  </template>
                  <n-list-item v-for="item in val" :key="item.code">
                    <n-grid :cols="3" :y-gap="0">
                      <n-gi>

                        <n-text :type="item.zdf>0?'error':'success'">
                          <n-image :src="item.img" width="20"/> &nbsp;{{ item.name }}
                        </n-text>
                      </n-gi>
                      <n-gi>
                        <n-text :type="item.zdf>0?'error':'success'">{{ item.zxj }}</n-text>&nbsp;
                        <n-text :type="item.zdf>0?'error':'success'">
                          <n-number-animation :precision="2" :from="0" :to="item.zdf"/>
                          %
                        </n-text>

                      </n-gi>
                      <n-gi>
                        <n-text :type="item.state === 'open' ? 'success' : 'warning'">{{
                            item.state === 'open' ? '开市' : '休市'
                          }}
                        </n-text>
                      </n-gi>
                    </n-grid>
                  </n-list-item>
                </n-list>
              </n-gi>
            </n-grid>
          </n-tab-pane>
          <n-tab-pane name="上证指数" tab="上证指数">
            <k-line-chart code="sh000001" :chart-height="panelHeight" stockName="上证指数" :k-days="20"
                          :dark-theme="true"></k-line-chart>
          </n-tab-pane>
          <n-tab-pane name="深证成指" tab="深证成指">
            <k-line-chart code="sz399001" :chart-height="panelHeight" stockName="深证成指" :k-days="20"
                          :dark-theme="true"></k-line-chart>
          </n-tab-pane>
          <n-tab-pane name="创业板指" tab="创业板指">
            <k-line-chart code="sz399006" :chart-height="panelHeight" stockName="创业板指" :k-days="20"
                          :dark-theme="true"></k-line-chart>
          </n-tab-pane>
          <n-tab-pane name="恒生指数" tab="恒生指数">
            <k-line-chart code="hkHSI" :chart-height="panelHeight" stockName="恒生指数" :k-days="20"
                          :dark-theme="true"></k-line-chart>
          </n-tab-pane>
          <n-tab-pane name="纳斯达克" tab="纳斯达克">
            <k-line-chart code="us.IXIC" :chart-height="panelHeight" stockName="纳斯达克" :k-days="20"
                          :dark-theme="true"></k-line-chart>
          </n-tab-pane>
          <n-tab-pane name="道琼斯" tab="道琼斯">
            <k-line-chart code="us.DJI" :chart-height="panelHeight" stockName="道琼斯" :k-days="20"
                          :dark-theme="true"></k-line-chart>
          </n-tab-pane>
          <n-tab-pane name="标普500" tab="标普500">
            <k-line-chart code="us.INX" :chart-height="panelHeight" stockName="标普500" :k-days="20"
                          :dark-theme="true"></k-line-chart>
          </n-tab-pane>
        </n-tabs>
      </n-tab-pane>
      <n-tab-pane name="重大指数" tab="重大指数">
        <!-- 移动端：指数网格选择 + 选中指数单图全宽 -->
        <div v-if="isMobile" class="mkt-index-mobile">
          <div class="mkt-index-mobile__title">选择指数</div>
          <div class="mkt-index-grid">
            <button
                v-for="item in mobileMajorIndexList"
                :key="item.code"
                type="button"
                class="mkt-index-grid__tile"
                :class="{'mkt-index-grid__tile--active': mobileMajorIndex.code === item.code}"
                @click="mobileMajorIndex = item"
            >{{ item.name }}</button>
          </div>
          <div class="mkt-index-chart-wrap">
            <StockLightweightKlineChart
                :key="mobileMajorIndex.code"
                :code="mobileMajorIndex.code"
                :chart-height="panelHeight"
                :stock-name="mobileMajorIndex.name"
                :dark-theme="true"
            />
          </div>
        </div>
        <n-tabs v-else type="segment" animated>

<!--          <n-tab-pane name="西部数据" tab="西部数据">-->
<!--            <StockLightweightKlineChart code="105.WDC" :chart-height="panelHeight" stock-name="西部数据"-->
<!--                                        :dark-theme="true"></StockLightweightKlineChart>-->
<!--          </n-tab-pane>-->

          <n-tab-pane name="上证指数" tab="上证指数"  >
            <StockLightweightKlineChart code="000001.SH" :chart-height="panelHeight-130" stock-name="上证指数" :dark-theme="true"></StockLightweightKlineChart>
          </n-tab-pane>
          <n-tab-pane name="深证指数" tab="深证指数"  >
            <StockLightweightKlineChart code="399001.SZ" :chart-height="panelHeight-130" stock-name="深证指数" :dark-theme="true"></StockLightweightKlineChart>
          </n-tab-pane>
          <n-tab-pane name="创业板指" tab="创业板指"  >
            <StockLightweightKlineChart code="399006.SZ" :chart-height="panelHeight-130" stock-name="创业板指" :dark-theme="true"></StockLightweightKlineChart>
          </n-tab-pane>

          <n-tab-pane name="恒生指数" tab="恒生指数">
            <StockLightweightKlineChart code="100.HSI" :chart-height="panelHeight" stock-name="恒生指数"
                                        :dark-theme="true"></StockLightweightKlineChart>
          </n-tab-pane>
          <n-tab-pane name="道琼斯" tab="道琼斯">
            <StockLightweightKlineChart code="100.DJIA" :chart-height="panelHeight" stock-name="道琼斯"
                                        :dark-theme="true"></StockLightweightKlineChart>
          </n-tab-pane>
          <n-tab-pane name="标普500" tab="标普500">
            <StockLightweightKlineChart code="100.SPX" :chart-height="panelHeight" stock-name="标普500"
                                        :dark-theme="true"></StockLightweightKlineChart>
          </n-tab-pane>
          <n-tab-pane name="纳斯达克" tab="纳斯达克">
            <StockLightweightKlineChart code="100.NDX" :chart-height="panelHeight" stock-name="纳斯达克"
                                        :dark-theme="true"></StockLightweightKlineChart>
          </n-tab-pane>

          <n-tab-pane name="沪深300" tab="沪深300">
            <StockLightweightKlineChart code="000300.SH" :chart-height="panelHeight-130" stock-name="沪深 300" :dark-theme="true"></StockLightweightKlineChart>
          </n-tab-pane>
          <n-tab-pane name="上证50" tab="上证50">
            <StockLightweightKlineChart code="000016.SH" :chart-height="panelHeight-130" stock-name="上证 50" :dark-theme="true"></StockLightweightKlineChart>
          </n-tab-pane>
          <n-tab-pane name="中证A500" tab="中证A500">
            <StockLightweightKlineChart code="000510.SH" :chart-height="panelHeight-130" stock-name="中证 A500" :dark-theme="true"></StockLightweightKlineChart>
          </n-tab-pane>
          <n-tab-pane name="中证1000" tab="中证1000">
            <StockLightweightKlineChart code="000852.SH" :chart-height="panelHeight-130" stock-name="中证 1000" :dark-theme="true"></StockLightweightKlineChart>
          </n-tab-pane>

          <n-tab-pane name="科创50" tab="科创50"  >
            <StockLightweightKlineChart code="000688.SH" :chart-height="panelHeight-130" stock-name="科创 50" :dark-theme="true"></StockLightweightKlineChart>
          </n-tab-pane>
          <n-tab-pane name="科创芯片" tab="科创芯片"  >
            <StockLightweightKlineChart code="000685.SH" :chart-height="panelHeight-130" stock-name="科创芯片" :dark-theme="true"></StockLightweightKlineChart>
          </n-tab-pane>
          <n-tab-pane name="证券龙头" tab="证券龙头"  >
            <StockLightweightKlineChart code="399437.SZ" :chart-height="panelHeight-130" stock-name="证券龙头" :dark-theme="true"></StockLightweightKlineChart>
          </n-tab-pane>
          <n-tab-pane name="高端装备" tab="高端装备"  >
            <StockLightweightKlineChart code="399437.SZ" :chart-height="panelHeight-130" stock-name="高端装备" :dark-theme="true"></StockLightweightKlineChart>
          </n-tab-pane>
          <n-tab-pane name="中证银行" tab="中证银行">
            <StockLightweightKlineChart code="399986.SZ" :chart-height="panelHeight-130" stock-name="中证银行" :dark-theme="true"></StockLightweightKlineChart>
          </n-tab-pane>
          <n-tab-pane name="上证医药" tab="上证医药">
            <StockLightweightKlineChart code="000037.SH" :chart-height="panelHeight-130" stock-name="上证医药" :dark-theme="true"></StockLightweightKlineChart>
          </n-tab-pane>
          <n-tab-pane name="中证白酒" tab="中证白酒">
            <StockLightweightKlineChart code="399997.SZ" :chart-height="panelHeight-130" stock-name="中证白酒" :dark-theme="true"></StockLightweightKlineChart>
          </n-tab-pane>
          <n-tab-pane name="富时中国三倍做多" tab="富时中国三倍做多">
            <k-line-chart code="usYINN.AM" :chart-height="panelHeight" stockName="富时中国三倍做多" :k-days="20"
                          :dark-theme="true"></k-line-chart>
          </n-tab-pane>
          <n-tab-pane name="VIX恐慌指数" tab="VIX恐慌指数">
            <k-line-chart code="usUVXY.AM" :chart-height="panelHeight" stockName="VIX恐慌指数" :k-days="20"
                          :dark-theme="true"></k-line-chart>
          </n-tab-pane>
        </n-tabs>
      </n-tab-pane>
      <n-tab-pane name="行业排名" tab="行业排名">
        <n-tabs type="card" animated>
          <n-tab-pane name="行业涨幅排名" tab="行业涨幅排名">
            <n-table striped>
              <n-thead>
                <n-tr>
                  <n-th>行业名称</n-th>
                  <n-th @click="changeIndustryRankSort">行业涨幅
                    <n-icon v-if="sort==='0'" :component="CaretDown"/>
                    <n-icon v-if="sort==='1'" :component="CaretUp"/>
                  </n-th>
                  <n-th>行业5日涨幅</n-th>
                  <n-th>行业20日涨幅</n-th>
                  <n-th>领涨股</n-th>
                  <n-th>涨幅</n-th>
                  <n-th>最新价</n-th>
                </n-tr>
              </n-thead>
              <n-tbody>
                <n-tr v-for="item in industryRanks" :key="item.bd_code">
                  <n-td>
                    <n-tag :bordered=false type="info">{{ item.bd_name }}</n-tag>
                  </n-td>
                  <n-td>
                    <n-text :type="item.bd_zdf>0?'error':'success'">{{ item.bd_zdf }}%</n-text>
                  </n-td>
                  <n-td>
                    <n-text :type="item.bd_zdf5>0?'error':'success'">{{ item.bd_zdf5 }}%</n-text>
                  </n-td>
                  <n-td>
                    <n-text :type="item.bd_zdf20>0?'error':'success'">{{ item.bd_zdf20 }}%</n-text>
                  </n-td>
                  <n-td>
                    <n-text :type="item.nzg_zdf>0?'error':'success'"> {{ item.nzg_name }}
                      <n-text type="info">{{ item.nzg_code }}</n-text>
                    </n-text>
                  </n-td>
                  <n-td>
                    <n-text :type="item.nzg_zdf>0?'error':'success'"> {{ item.nzg_zdf }}%</n-text>
                  </n-td>
                  <n-td>
                    <n-text :type="item.nzg_zdf>0?'error':'success'">{{ item.nzg_zxj }}</n-text>
                  </n-td>
                </n-tr>
              </n-tbody>
            </n-table>
            <n-table striped>
              <n-thead>
                <n-tr>
                  <n-th>行业名称</n-th>
                  <n-th @click="changeIndustryRankSort">行业涨幅
                    <n-icon v-if="sort==='0'" :component="CaretDown"/>
                    <n-icon v-if="sort==='1'" :component="CaretUp"/>
                  </n-th>
                  <n-th>行业5日涨幅</n-th>
                  <n-th>行业20日涨幅</n-th>
                  <n-th>领涨股</n-th>
                  <n-th>涨幅</n-th>
                  <n-th>最新价</n-th>
                </n-tr>
              </n-thead>
              <n-tbody>
                <n-tr v-for="item in industryRanks" :key="item.bd_code">
                  <n-td>
                    <n-tag :bordered=false type="info">{{ item.bd_name }}</n-tag>
                  </n-td>
                  <n-td>
                    <n-text :type="item.bd_zdf>0?'error':'success'">{{ item.bd_zdf }}%</n-text>
                  </n-td>
                  <n-td>
                    <n-text :type="item.bd_zdf5>0?'error':'success'">{{ item.bd_zdf5 }}%</n-text>
                  </n-td>
                  <n-td>
                    <n-text :type="item.bd_zdf20>0?'error':'success'">{{ item.bd_zdf20 }}%</n-text>
                  </n-td>
                  <n-td>
                    <n-text :type="item.nzg_zdf>0?'error':'success'"> {{ item.nzg_name }}
                      <n-text type="info">{{ item.nzg_code }}</n-text>
                    </n-text>
                  </n-td>
                  <n-td>
                    <n-text :type="item.nzg_zdf>0?'error':'success'"> {{ item.nzg_zdf }}%</n-text>
                  </n-td>
                  <n-td>
                    <n-text :type="item.nzg_zdf>0?'error':'success'">{{ item.nzg_zxj }}</n-text>
                  </n-td>
                </n-tr>
              </n-tbody>
            </n-table>
          </n-tab-pane>
          <n-tab-pane name="行业资金排名(净流入)" tab="行业资金排名">
            <industryMoneyRank :fenlei="'0'" :header-title="'行业资金排名(净流入)'" :sort="'netamount'"/>
          </n-tab-pane>
          <n-tab-pane name="证监会行业资金排名(净流入)" tab="证监会行业资金排名">
            <industryMoneyRank :fenlei="'2'" :header-title="'证监会行业资金排名(净流入)'" :sort="'netamount'"/>
          </n-tab-pane>
          <n-tab-pane name="概念板块资金排名(净流入)" tab="概念板块资金排名">
            <industryMoneyRank :fenlei="'1'" :header-title="'概念板块资金排名(净流入)'" :sort="'netamount'"/>
          </n-tab-pane>
        </n-tabs>
      </n-tab-pane>
      <n-tab-pane name="个股资金流向" tab="个股资金流向">
        <n-tabs type="card" animated>
          <n-tab-pane name="netamount" tab="净流入额排名">
            <RankTable :header-title="'净流入额排名'" :sort="'netamount'"/>
          </n-tab-pane>
          <n-tab-pane name="outamount" tab="流出资金排名">
            <RankTable :header-title="'流出资金排名'" :sort="'outamount'"/>
          </n-tab-pane>
          <n-tab-pane name="ratioamount" tab="净流入率排名">
            <RankTable :header-title="'净流入率排名'" :sort="'ratioamount'"/>
          </n-tab-pane>
          <n-tab-pane name="r0_net" tab="主力净流入额排名">
            <RankTable :header-title="'主力净流入额排名'" :sort="'r0_net'"/>
          </n-tab-pane>
          <n-tab-pane name="r0_out" tab="主力流出排名">
            <RankTable :header-title="'主力流出排名'" :sort="'r0_out'"/>
          </n-tab-pane>
          <n-tab-pane name="r0_ratio" tab="主力净流入率排名">
            <RankTable :header-title="'主力净流入率排名'" :sort="'r0_ratio'"/>
          </n-tab-pane>
          <n-tab-pane name="r3_net" tab="散户净流入额排名">
            <RankTable :header-title="'散户净流入额排名'" :sort="'r3_net'"/>
          </n-tab-pane>
          <n-tab-pane name="r3_out" tab="散户流出排名">
            <RankTable :header-title="'散户流出排名'" :sort="'r3_out'"/>
          </n-tab-pane>
          <n-tab-pane name="r3_ratio" tab="散户净流入率排名">
            <RankTable :header-title="'散户净流入率排名'" :sort="'r3_ratio'"/>
          </n-tab-pane>
        </n-tabs>
      </n-tab-pane>
      <n-tab-pane name="板块资金流向" tab="板块资金流向">
        <BKFundFlowChart :dark-theme="darkTheme" :chart-height="600"/>
      </n-tab-pane>
      <n-tab-pane name="概念资金流向" tab="概念资金流向">
        <ConceptFundFlowChart :dark-theme="darkTheme" :chart-height="600"/>
      </n-tab-pane>
      <n-tab-pane name="龙虎榜" tab="龙虎榜">
        <LongTigerRankList />
      </n-tab-pane>
      <n-tab-pane name="个股研报" tab="个股研报">
        <StockResearchReportList :stock-code="stockCode"/>
      </n-tab-pane>
      <n-tab-pane name="公司公告" tab="公司公告 ">
        <StockNoticeList :stock-code="stockCode" />
      </n-tab-pane>
      <n-tab-pane name="行业研究" tab="行业研究 ">
        <IndustryResearchReportList/>
      </n-tab-pane>
      <n-tab-pane name="当前热门" tab="当前热门">
        <n-tabs type="card" animated>
          <n-tab-pane name="全球" tab="全球">
            <HotStockList :market-type="'10'"/>
          </n-tab-pane>
          <n-tab-pane name="沪深" tab="沪深">
            <HotStockList :market-type="'12'"/>
          </n-tab-pane>
          <n-tab-pane name="港股" tab="港股">
            <HotStockList :market-type="'13'"/>
          </n-tab-pane>
          <n-tab-pane name="美股" tab="美股">
            <HotStockList :market-type="'11'"/>
          </n-tab-pane>
          <n-tab-pane name="热门话题" tab="热门话题">
            <n-grid :cols="1" :y-gap="10">
              <n-grid-item>
                <HotTopics/>
              </n-grid-item>
<!--              <n-grid-item>-->
<!--                <HotEvents/>-->
<!--              </n-grid-item>-->
            </n-grid>
          </n-tab-pane>
          <n-tab-pane name="重大事件时间轴" tab="重大事件时间轴">
            <InvestCalendarTimeLine />
          </n-tab-pane>
          <n-tab-pane name="财经日历" tab="财经日历">
            <ClsCalendarTimeLine />
          </n-tab-pane>
        </n-tabs>
      </n-tab-pane>
      <n-tab-pane name="名站优选" tab="名站优选">
        <Stockhotmap />
      </n-tab-pane>
    </n-tabs>
  </n-card>
  <n-modal class="market-summary-modal" transform-origin="center" v-model:show="summaryModal" preset="card" style="width: 800px;max-width: calc(100vw - 32px);"
           :title="'AI市场资讯总结'">
    <n-spin size="small" :show="loading && !aiSummary">
      <div ref="aiResultScrollRef" style="height: 440px;max-height: 60vh;text-align: left;overflow-y: auto;">
        <MdPreview ref="mdPreviewRef" :modelValue="aiSummary" :theme="theme"/>
      </div>
    </n-spin>
    <template #footer>
      <n-flex justify="space-between" ref="tipsRef">
        <n-text type="info" v-if="aiSummaryTime">
          <n-tag v-if="modelName" type="warning" round :title="chatId" :bordered="false">{{ modelName }}</n-tag>
          {{ aiSummaryTime }}
        </n-text>
        <n-text type="success" v-if="analysisStatus">{{ analysisStatus }}</n-text>
        <n-text type="error">*AI分析结果仅供参考，请以实际行情为准。投资需谨慎，风险自担。</n-text>
      </n-flex>
    </template>
    <template #action>
      <n-flex class="market-summary-modal__switches" justify="left" style="margin-bottom: 10px">
        <n-switch v-model:value="enableTools" :round="false">
          <template #checked>
            工具调用
          </template>
          <template #unchecked>
            非工具调用
          </template>
        </n-switch>
        <n-switch v-model:value="thinkingMode" :round="false">
          <template #checked>
            思考模式
          </template>
          <template #unchecked>
            非思考模式
          </template>
        </n-switch>


        <n-gradient-text type="error" style="margin-left: 10px">*AI函数工具调用可以增强AI获取数据的能力,但会消耗更多tokens。</n-gradient-text>
      </n-flex>
      <n-flex class="market-summary-modal__selectors" justify="space-between" style="margin-bottom: 10px">
        <n-select style="width: 32%" v-model:value="aiConfigId" label-field="name" value-field="ID"
                  :options="aiConfigs" placeholder="请选择AI模型服务配置"/>
        <n-select style="width: 32%" v-model:value="sysPromptId" label-field="name" value-field="ID"
                  :options="sysPromptOptions" placeholder="请选择系统提示词"/>
        <n-select style="width: 32%" v-model:value="question" label-field="name" value-field="content"
                  :options="userPromptOptions" placeholder="请选择用户提示词"/>
      </n-flex>
      <n-flex class="market-summary-modal__actions" justify="right">
        <n-input v-model:value="question" style="text-align: left" clearable
                 type="textarea"
                 :show-count="true"
                 placeholder="请输入您的问题:例如 总结和分析股票市场新闻中的投资机会"
                 :autosize="{
              minRows: 2,
              maxRows: 5
            }"
        />
        <n-button size="tiny" type="warning" @click="reAiSummary">再次总结</n-button>
        <n-button size="tiny" type="success" @click="copyToClipboard">复制到剪切板</n-button>
        <n-button size="tiny" type="primary" @click="saveAsMarkdown">保存为Markdown文件</n-button>
        <n-button size="tiny" type="error" @click="share">分享到项目社区</n-button>
      </n-flex>
    </template>
  </n-modal>

  <div class="market-summary-fab market-mobile-summary-action" style="position: fixed;bottom: 18px;right:25px;z-index: 10;" v-if="summaryBTN">
    <n-input-group>
      <n-button type="primary" @click="getAiSummary">
        <n-icon :component="PulseOutline"/> &nbsp;AI总结
      </n-button>
    </n-input-group>
  </div>



</template>
<style scoped>
@media (max-width: 768px) {
  .market-page-shell {
    margin: 0 6px 86px;
    text-align: left;
  }

  :deep(.market-page-shell > .n-card__content) {
    padding: 8px;
  }

  :deep(.market-page-shell .n-tabs-nav-scroll-content) {
    min-width: max-content;
  }

  /* 移动端用自定义两级菜单，隐藏顶层原生标签栏。
     用直接子选择器(>)，避免穿透到嵌套的子 tabs（财联社/新浪/外媒、
     全球指数、行业排名等 segment/card 子标签）把它们也隐藏掉。 */
  .market-mobile-tabs--native-hidden > :deep(.n-tabs-nav) {
    display: none !important;
  }

  .market-news-grid {
    display: block !important;
  }

  .market-desktop-news-grid {
    display: none !important;
  }

  .market-news-grid :deep(.n-grid-item) {
    margin-bottom: 8px;
  }

  .market-mobile-scroll {
    display: grid !important;
    grid-auto-columns: minmax(260px, 82vw);
    grid-auto-flow: column;
    grid-template-columns: none !important;
    overflow-x: auto;
    padding-bottom: 8px;
    scroll-snap-type: x proximity;
  }

  .market-mobile-scroll :deep(.n-grid-item) {
    scroll-snap-align: start;
  }

  .market-mobile-index-card {
    border-radius: 6px;
    min-height: 100%;
    overflow: hidden;
  }

  .market-mobile-index-card :deep(.n-list-item) {
    padding: 8px 10px;
  }

  .market-mobile-index-card :deep(.n-grid) {
    align-items: center;
    grid-template-columns: minmax(108px, 1fr) minmax(82px, auto) 44px !important;
  }

  .market-summary-fab {
    bottom: calc(var(--mobile-bottom-nav-height) + 10px + env(safe-area-inset-bottom)) !important;
    right: 10px !important;
  }

  .market-mobile-summary-action :deep(.n-button) {
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.14);
    height: 40px;
  }

  :deep(.market-summary-modal.n-modal) {
    margin: 0 !important;
    max-width: 100vw !important;
    width: calc(100vw - 12px) !important;
  }

  :deep(.market-summary-modal .n-card) {
    max-height: calc(100dvh - var(--mobile-bottom-nav-height) - var(--safe-bottom) - 12px);
    overflow: auto;
  }

  .market-summary-modal__switches,
  .market-summary-modal__selectors,
  .market-summary-modal__actions {
    align-items: stretch !important;
    flex-wrap: wrap;
    gap: 8px !important;
  }

  .market-summary-modal__selectors :deep(.n-select),
  .market-summary-modal__actions :deep(.n-input) {
    width: 100% !important;
  }

  .market-summary-modal__actions :deep(.n-button) {
    flex: 1 1 calc(50% - 8px);
    min-width: 120px;
  }
}

/* ============ 移动端两级导航（仅在 isMobile 渲染） ============ */
.market-mobile-nav {
  display: flex;
  flex-direction: column;
  gap: 10px;
  margin-bottom: 6px;
}

.market-mobile-nav__groups {
  display: grid;
  gap: 8px;
  grid-template-columns: repeat(4, minmax(0, 1fr));
}

.market-group-tile {
  align-items: center;
  appearance: none;
  background: var(--n-color, #fff);
  border: 1px solid var(--n-border-color, #edf0f5);
  border-radius: 10px;
  color: inherit;
  display: flex;
  flex-direction: column;
  font: inherit;
  gap: 2px;
  padding: 8px 4px;
}

.market-group-tile--active {
  background: var(--n-color-target, rgba(32, 128, 240, 0.1));
  border-color: #2080f0;
  color: #2080f0;
}

.market-group-tile__icon {
  font-size: 18px;
}

.market-group-tile__name {
  font-size: 12px;
  font-weight: 600;
}

.market-mobile-nav__tabs {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.market-tab-chip {
  appearance: none;
  background: var(--n-color, #fff);
  border: 1px solid var(--n-border-color, #edf0f5);
  border-radius: 16px;
  color: var(--n-text-color, #333);
  font: inherit;
  font-size: 13px;
  padding: 6px 14px;
}

.market-tab-chip--active {
  background: #2080f0;
  border-color: #2080f0;
  color: #fff;
  font-weight: 700;
}

/* ============ 移动端"市场快讯"重新设计（仅在 isMobile 渲染） ============ */
.mkt-brief {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 4px 0 calc(var(--safe-bottom) + 8px);
  text-align: left;
}

.mkt-section {
  background: var(--n-color, #fff);
  border: 1px solid var(--n-border-color, #edf0f5);
  border-radius: 12px;
  overflow: hidden;
}

.mkt-section__title {
  font-size: 14px;
  font-weight: 700;
  padding: 10px 12px 8px;
}

/* ① 大盘速览横滑 */
.mkt-index-rail {
  display: flex;
  gap: 8px;
  overflow-x: auto;
  padding: 0 12px 12px;
  scroll-snap-type: x proximity;
  -webkit-overflow-scrolling: touch;
}

.mkt-index-card {
  background: var(--n-color-target, #f5f7fa);
  border-radius: 10px;
  flex: 0 0 108px;
  padding: 9px 10px;
  scroll-snap-align: start;
}

.mkt-index-card--up {
  border-left: 3px solid #d03050;
}

.mkt-index-card--down {
  border-left: 3px solid #18a058;
}

.mkt-index-card__head {
  align-items: center;
  display: flex;
  gap: 4px;
  margin-bottom: 4px;
}

.mkt-index-card__name {
  color: var(--n-text-color-2, #555);
  font-size: 12px;
  font-weight: 600;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.mkt-index-card__price {
  font-size: 15px;
  font-weight: 800;
}

.mkt-index-card__zdf {
  border-radius: 4px;
  color: #fff;
  display: inline-block;
  font-size: 12px;
  font-weight: 700;
  margin-top: 3px;
  padding: 1px 6px;
}

.mkt-index-card__state {
  color: var(--n-text-color-3, #98a2b3);
  font-size: 10px;
  margin-top: 4px;
}

/* ③ 资讯流 */
.mkt-news-sticky {
  background: var(--n-color, #fff);
  position: sticky;
  top: 0;
  z-index: 2;
}

.mkt-news-source-bar {
  display: flex;
  gap: 8px;
  overflow-x: auto;
  padding: 0 12px 10px;
}

.mkt-source-chip {
  appearance: none;
  background: var(--n-color-target, #f5f7fa);
  border: 1px solid transparent;
  border-radius: 16px;
  color: var(--n-text-color-2, #555);
  flex: 0 0 auto;
  font: inherit;
  font-size: 13px;
  padding: 5px 14px;
  white-space: nowrap;
}

.mkt-source-chip--active {
  background: #2080f0;
  color: #fff;
  font-weight: 700;
}

.mkt-news-list {
  display: flex;
  flex-direction: column;
}

.mkt-news-item {
  border-top: 1px solid var(--n-border-color, #f0f1f5);
  padding: 10px 12px;
}

.mkt-news-item:first-child {
  border-top: 0;
}

.mkt-news-item--red {
  background: rgba(208, 48, 80, 0.04);
}

.mkt-news-item__top {
  align-items: center;
  display: flex;
  gap: 8px;
  margin-bottom: 4px;
}

.mkt-news-item__time {
  color: var(--n-text-color-3, #98a2b3);
  font-size: 11px;
}

.mkt-news-item__title {
  font-size: 14px;
  font-weight: 700;
  line-height: 1.4;
  margin-bottom: 2px;
}

.mkt-news-item__content {
  color: var(--n-text-color-2, #555);
  font-size: 13px;
  line-height: 1.5;
  overflow-wrap: anywhere;
}

.mkt-news-item__tags {
  align-items: center;
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
  margin-top: 6px;
}

.mkt-news-item__link {
  color: #f0a020;
  font-size: 12px;
}

/* 涨跌色工具类 */
.text-error {
  color: #d03050;
}

.text-success {
  color: #18a058;
}

.bg-error {
  background: #d03050;
}

.bg-success {
  background: #18a058;
}

/* ============ 移动端"指数"板块：网格选择 + 单图 ============ */
.mkt-index-mobile {
  padding: 4px 0 calc(var(--safe-bottom) + 8px);
}

.mkt-index-mobile__title {
  color: var(--n-text-color-3, #98a2b3);
  font-size: 12px;
  font-weight: 600;
  padding: 6px 12px;
}

.mkt-index-grid {
  display: grid;
  gap: 8px;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  padding: 0 12px 12px;
}

.mkt-index-grid__tile {
  appearance: none;
  background: var(--n-color, #fff);
  border: 1px solid var(--n-border-color, #edf0f5);
  border-radius: 8px;
  color: var(--n-text-color, #333);
  font: inherit;
  font-size: 12px;
  font-weight: 600;
  padding: 8px 4px;
  text-align: center;
}

.mkt-index-grid__tile--active {
  background: #2080f0;
  border-color: #2080f0;
  color: #fff;
}

.mkt-index-chart-wrap {
  background: var(--n-color, #fff);
  border: 1px solid var(--n-border-color, #edf0f5);
  border-radius: 10px;
  margin: 0 8px;
  overflow: hidden;
  padding: 6px;
}
</style>
