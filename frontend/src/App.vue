<script setup>
import {
  EventsEmit,
  EventsOff,
  EventsOn
} from './api/runtime'
import {defineAsyncComponent, h, onBeforeMount, onBeforeUnmount, onMounted, ref, computed, watch} from "vue";
import {RouterLink, useRoute, useRouter} from 'vue-router'
import {createDiscreteApi,darkTheme,lightTheme , NIcon, NText,NButton,dateZhCN,zhCN} from 'naive-ui'
import {
  AlarmOutline,
  AnalyticsOutline,
  BarChartSharp, Bonfire, BonfireOutline, DiamondOutline, EaselSharp,
  ExpandOutline, Flag,
  Flame, FlameSharp, FlaskOutline, GlobeOutline, InformationOutline,
  LogoGithub,
  ChatbubblesOutline,
  NewspaperOutline,
  NewspaperSharp, Notifications,
  PowerOutline, Pulse,
  ReorderTwoOutline,
  SettingsOutline, ServerOutline, Skull, SkullOutline, SkullSharp,
  SparklesOutline, FlashOutline, Star,
  StarOutline,
  StatsChartOutline,
  Wallet, WarningOutline, TimeOutline, SearchOutline,
} from '@vicons/ionicons5'
import {AnalyzeSentiment, GetConfig, GetGroupList, GetVersionInfo} from "./api/app";
import {useDevice} from "./composables/useDevice";
import {cnOpen, hkOpen, usOpen} from "./api/marketClock";
import {registerFeed, stopFeed} from "./api/scheduler";
import {Dragon, Fire, FirefoxBrowser, Gripfire, Robot} from "@vicons/fa";
import {Prompt, ReportAnalytics, ReportMoney, ReportSearch, TrendingUp} from "@vicons/tabler";
import {LocalFireDepartmentRound} from "@vicons/material";
import {AppsList20Regular, BoxSearch20Regular,SlideHide24Filled, CommentNote20Filled} from "@vicons/fluent";
import {FireFilled, MoneyCollectOutlined, NotificationFilled, StockOutlined} from "@vicons/antd";




const route = useRoute()
const router = useRouter()
const FloatingAgentAssistant = defineAsyncComponent(() => import("./components/FloatingAgentAssistant.vue"))
const loading = ref(true)
const loadingMsg = ref("加载数据中...")
const enableNews = ref(false)
const enableFund = ref(false)
const enableAgent = ref(false)
const enableDarkTheme = ref(null)
const content = ref('未经授权,禁止商业目的!\n\n数据来源于网络,仅供参考;投资有风险,入市需谨慎')
const isFullscreen = ref(false)
const activeKey = ref('stock')
const containerRef = ref({})
const realtimeProfit = ref(0)
const telegraph = ref([])
const groupList = ref([])
const officialStatement= ref("")
const marketStatus = ref('')
// 全局单例设备状态，替代组件内 matchMedia 监听
const {isMobile} = useDevice()
const mobileMenuVisible = ref(false)

const mobileBottomNavItems = [
  { key: 'stock', label: '自选', icon: StarOutline, route: { name: 'stock', query: { groupName: '全部', groupId: 0 } } },
  { key: 'market', label: '市场', icon: NewspaperOutline, route: { name: 'market', query: { name: '市场快讯' } } },
  { key: 'klineAnalysis', label: 'K线', icon: AnalyticsOutline, route: { name: 'klineAnalysis' } },
  { key: 'promptPlaza', label: '提示词', icon: GlobeOutline, route: { name: 'research', query: { name: '提示词广场' } } },
  { key: 'more', label: '更多', icon: ReorderTwoOutline },
]

const routeActiveKeyMap = {
  stock: 'stock',
  market: 'market',
  klineAnalysis: 'klineAnalysis',
  fund: 'fund',
  agent: 'agent',
  research: 'research',
  cronTasks: 'research',
  mcpServers: 'research',
  settings: 'settings',
  about: 'about',
}

const investmentMottos = [
  "投资有风险，入市需谨慎",
  "别人贪婪我恐惧，别人恐惧我贪婪",
  "股市有风险，投资需谨慎",
  "不要把所有鸡蛋放在一个篮子里",
  "时间是优秀企业的朋友",
  "买股票就是买公司",
  "市场短期是投票机，长期是称重机",
  "保住本金是投资的第一要务",
  "在别人恐慌时贪婪，在别人贪婪时恐慌",
  "风险来自于你不知道自己在做什么",
  "价格是你付出的，价值是你得到的",
  "投资最重要的品质是耐心",
  "机会总是留给有准备的人",
  "知行合一，方能致远",
  "顺势而为，逆势而思",
  "投资是一场马拉松，不是百米冲刺",
  "独立思考是投资成功的关键",
  "市场永远在波动，但价值终将回归",
  "控制风险比追求收益更重要",
  "学习是最好的投资",
]
const currentMotto = ref(investmentMottos[Math.floor(Math.random() * investmentMottos.length)])

function refreshMotto() {
  currentMotto.value = investmentMottos[Math.floor(Math.random() * investmentMottos.length)]
}

function updateMarketStatus() {
  // 交易时段状态由 marketClock 统一维护（收口原散落各处的 IsTradingTime().catch(()=>false)，
  // 并修复其 bug：网络失败不再被误判为"收盘→停轮询"）。
  const parts = []
  parts.push(cnOpen.value ? 'A股交易中' : 'A股休市')
  parts.push(hkOpen.value ? '港股交易中' : '港股休市')
  parts.push(usOpen.value ? '美股交易中' : '美股休市')
  marketStatus.value = parts.join(' | ')
  document.title = "go-stock " + marketStatus.value
}

// 交易时段状态变化时同步标题（marketClock 自带 60s 轮询 + 页面恢复立即重判）。
watch([cnOpen, hkOpen, usOpen], updateMarketStatus, { immediate: true })

// 内容区高度随 isMobile 自动响应（设备状态由 useDevice 单例维护）
const contentStyle = computed(() => isMobile.value
    ? "height: calc(100dvh - var(--mobile-bottom-nav-height) - env(safe-area-inset-bottom));overflow: auto"
    : "height: calc(100vh - var(--desktop-bottom-menu-height));overflow: auto")

function handleMobileNav(item) {
  if (item.key === 'more') {
    mobileMenuVisible.value = true
    return
  }
  activeKey.value = item.key === 'promptPlaza' ? 'research' : item.key
  if (item.key === 'stock') {
    EventsEmit("changeTab", {ID: 0, name: '全部'})
  }
  if (item.key === 'market') {
    EventsEmit("changeMarketTab", {ID: 0, name: '市场快讯'})
  }
  if (item.key === 'promptPlaza') {
    setTimeout(() => {
      EventsEmit("changeResearchTab", {ID: 10, name: '提示词广场'})
    }, 100)
  }
  router.push(item.route)
  mobileMenuVisible.value = false
}

function syncActiveKeyFromRoute(routeName) {
  const key = routeActiveKeyMap[String(routeName || '')]
  if (key) {
    activeKey.value = key
  }
}

function isMobileBottomNavActive(item) {
  if (item.key === 'promptPlaza') {
    return activeKey.value === 'research'
  }
  if (item.key === 'more') {
    return !['stock', 'market', 'klineAnalysis', 'research'].includes(activeKey.value)
  }
  return activeKey.value === item.key
}

// "更多"抽屉的功能网格：分组 + tile。每个 tile 的 onActivate 原样复制 menuOptions 对应项的
// EventsEmit + setTimeout(100) 协调，route 复制对应的 router.push 目标，保证子 Tab 同步不丢失。
const mobileMoreSections = computed(() => {
  const sections = []

  // 自选管理：全部 + 动态分组
  const stockTiles = [
    {
      label: '股票自选', icon: StarOutline,
      route: {name: 'stock', query: {groupName: '全部', groupId: 0}},
      onActivate: () => { activeKey.value = 'stock'; EventsEmit("changeTab", {ID: 0, name: '全部'}) },
    },
  ]
  groupList.value.forEach(g => {
    stockTiles.push({
      label: g.name, icon: StarOutline,
      route: {name: 'stock', query: {groupName: g.name, groupId: g.ID}},
      onActivate: () => { activeKey.value = 'stock'; setTimeout(() => { EventsEmit("changeTab", g) }, 100) },
    })
  })
  sections.push({category: '自选管理', tiles: stockTiles})

  // 市场行情
  const marketTabs = ['市场快讯', '全球股指', '重大指数', '行业排名', '个股资金流向', '板块资金流向', '概念资金流向', '龙虎榜', '个股研报', '公司公告', '行业研究', '当前热门', '名站优选']
  sections.push({
    category: '市场行情',
    tiles: marketTabs.map(name => ({
      label: name, icon: NewspaperOutline,
      route: {name: 'market', query: {name}},
      onActivate: () => { activeKey.value = 'market'; EventsEmit("changeMarketTab", {ID: 0, name}) },
    })),
  })

  // 研究分析
  const researchTiles = [
    {label: 'AI分析报告', icon: ReportAnalytics, tab: {ID: 0, name: 'AI分析报告'}},
    {label: '股票推荐记录', icon: TrendingUp, tab: {ID: 1, name: '股票推荐记录'}},
    {label: '异动监控', icon: WarningOutline, tab: {ID: 2, name: '异动监控'}},
    {label: '涨停梯队', icon: Flame, tab: {ID: 9, name: '涨停梯队'}},
    {label: '提示词模板', icon: SparklesOutline, tab: {ID: 3, name: '提示词模板'}},
    {label: '提示词广场', icon: GlobeOutline, tab: {ID: 10, name: '提示词广场'}},
    {label: '问答广场', icon: ChatbubblesOutline, tab: {ID: 11, name: '问答广场'}},
    {label: '形态选股', icon: SearchOutline, tab: {ID: 3, name: '形态选股'}},
    {label: '指标选股', icon: BoxSearch20Regular, tab: {ID: 0, name: '指标选股'}},
    {label: '定时任务', icon: TimeOutline, tab: {ID: 5, name: '定时任务'}},
    {label: '交易日志', icon: Wallet, tab: {ID: 6, name: '交易日志'}},
    {label: 'MCP服务', icon: ServerOutline, tab: {ID: 7, name: 'MCP服务'}},
    {label: '技能管理', icon: AppsList20Regular, tab: {ID: 8, name: '技能管理'}},
  ]
  sections.push({
    category: '研究分析',
    tiles: researchTiles.map(t => ({
      label: t.label, icon: t.icon,
      route: {name: 'research', query: {name: t.label}},
      onActivate: () => { activeKey.value = 'research'; setTimeout(() => { EventsEmit("changeResearchTab", t.tab) }, 100) },
    })),
  })

  // 系统设置（基金/AI 智能体按配置显隐）
  const settingTiles = [
    {
      label: 'K线分析', icon: StatsChartOutline,
      route: {name: 'klineAnalysis'},
      onActivate: () => { activeKey.value = 'klineAnalysis' },
    },
  ]
  if (enableFund.value) {
    settingTiles.push({
      label: '基金自选', icon: SparklesOutline,
      route: {name: 'fund', query: {name: '基金自选'}},
      onActivate: () => { activeKey.value = 'fund'; EventsEmit("changeFundTab", {name: '基金自选'}) },
    })
    settingTiles.push({
      label: '基金排行', icon: TrendingUp,
      route: {name: 'fund', query: {name: '基金排行'}},
      onActivate: () => { activeKey.value = 'fund'; EventsEmit("changeFundTab", {name: '基金排行'}) },
    })
  }
  if (enableAgent.value) {
    settingTiles.push({
      label: 'Ai智能体', icon: Robot,
      route: {name: 'agent', query: {name: 'Ai智能体'}},
      onActivate: () => { activeKey.value = 'agent' },
    })
  }
  settingTiles.push({
    label: '设置', icon: SettingsOutline,
    route: {name: 'settings', query: {name: '设置'}},
    onActivate: () => { activeKey.value = 'settings' },
  })
  settingTiles.push({
    label: '关于', icon: InformationOutline,
    route: {name: 'about', query: {name: '关于'}},
    onActivate: () => { activeKey.value = 'about' },
  })
  sections.push({category: '系统设置', tiles: settingTiles})

  return sections
})

function onTileSelect(tile) {
  tile.onActivate?.()
  if (tile.route) {
    router.push(tile.route)
  }
  mobileMenuVisible.value = false
}

watch(
  () => route.name,
  syncActiveKeyFromRoute,
  { immediate: true },
)
const menuOptions = ref([
  {
    label: () =>
        h(
            RouterLink,
            {
              to: {
                name: 'stock',
                query: {
                  groupName: '全部',
                  groupId: 0,
                },
                params: {},
              },
              onClick: () => {
                activeKey.value = 'stock'
              },
            },
            {default: () => '股票自选',}
        ),
    key: 'stock',
    icon: renderIcon(StarOutline),
    children: [
      {
        label: () =>
            h(
                'a',
                {
                  href: '#',
                  type: 'info',
                  onClick: () => {
                    activeKey.value = 'stock'
                    //console.log("push",item)
                    router.push({
                      name: 'stock',
                      query: {
                        groupName: '全部',
                        groupId: 0,
                      },
                    })
                    EventsEmit("changeTab", {ID: 0, name: '全部'})
                  },
                  to: {
                    name: 'stock',
                    query: {
                      groupName: '全部',
                      groupId: 0,
                    },
                  }
                },
                {default: () => '全部',}
            ),
        key: 0,
      }
    ],
  },
  {
    label: () =>
        h(
            RouterLink,
            {
              href: '#',
              to: {
                name: 'market',
                params: {}
              },
              onClick: () => {
                activeKey.value = 'market'
                EventsEmit("changeMarketTab", {ID: 0, name: '市场快讯'})
              },
            },
            {default: () => '市场行情'}
        ),
    key: 'market',
    icon: renderIcon(NewspaperOutline),
    children: [
      {
        label: () =>
            h(
                RouterLink,
                {
                  href: '#',
                  to: {
                    name: 'market',
                    query: {
                      name: "市场快讯",
                    }
                  },
                  onClick: () => {
                    activeKey.value = 'market'
                    EventsEmit("changeMarketTab", {ID: 0, name: '市场快讯'})
                  },
                },
                {default: () => '市场快讯',}
            ),
        key: 'market1',
        icon: renderIcon(NewspaperSharp),
      },
      {
        label: () =>
            h(
                RouterLink,
                {
                  href: '#',
                  to: {
                    name: 'market',
                    query: {
                      name: "全球股指",
                    },
                  },
                  onClick: () => {
                    activeKey.value = 'market'
                    EventsEmit("changeMarketTab", {ID: 0, name: '全球股指'})
                  },
                },
                {default: () => '全球股指',}
            ),
        key: 'market2',
        icon: renderIcon(BarChartSharp),
      },
      {
        label: () =>
            h(
                RouterLink,
                {
                  href: '#',
                  to: {
                    name: 'market',
                    query: {
                      name: "重大指数",
                    }
                  },
                  onClick: () => {
                    activeKey.value = 'market'
                    EventsEmit("changeMarketTab", {ID: 0, name: '重大指数'})
                  },
                },
                {default: () => '重大指数',}
            ),
        key: 'market3',
        icon: renderIcon(AnalyticsOutline),
      },
      {
        label: () =>
            h(
                RouterLink,
                {
                  href: '#',
                  to: {
                    name: 'market',
                    query: {
                      name: "行业排名",
                    }
                  },
                  onClick: () => {
                    activeKey.value = 'market'
                    EventsEmit("changeMarketTab", {ID: 0, name: '行业排名'})
                  },
                },
                {default: () => '行业排名',}
            ),
        key: 'market4',
        icon: renderIcon(Flag),
      },
      {
        label: () =>
            h(
                RouterLink,
                {
                  href: '#',
                  to: {
                    name: 'market',
                    query: {
                      name: "个股资金流向",
                    }
                  },
                  onClick: () => {
                    activeKey.value = 'market'
                    EventsEmit("changeMarketTab", {ID: 0, name: '个股资金流向'})
                  },
                },
                {default: () => '个股资金流向',}
            ),
        key: 'market5',
        icon: renderIcon(Pulse),
      },
      {
        label: () =>
            h(
                RouterLink,
                {
                  href: '#',
                  to: {
                    name: 'market',
                    query: {
                      name: "板块资金流向",
                    }
                  },
                  onClick: () => {
                    activeKey.value = 'market'
                    EventsEmit("changeMarketTab", {ID: 0, name: '板块资金流向'})
                  },
                },
                {default: () => '板块资金流向',}
            ),
        key: 'market5_1',
        icon: renderIcon(ReportMoney),
      },
      {
        label: () =>
            h(
                RouterLink,
                {
                  href: '#',
                  to: {
                    name: 'market',
                    query: {
                      name: "概念资金流向",
                    }
                  },
                  onClick: () => {
                    activeKey.value = 'market'
                    EventsEmit("changeMarketTab", {ID: 0, name: '概念资金流向'})
                  },
                },
                {default: () => '概念资金流向',}
            ),
        key: 'market5_2',
        icon: renderIcon(TrendingUp),
      },
      {
        label: () =>
            h(
                RouterLink,
                {
                  href: '#',
                  to: {
                    name: 'market',
                    query: {
                      name: "龙虎榜",
                    }
                  },
                  onClick: () => {
                    activeKey.value = 'market'
                    EventsEmit("changeMarketTab", {ID: 0, name: '龙虎榜'})
                  },
                },
                {default: () => '龙虎榜',}
            ),
        key: 'market6',
        icon: renderIcon(Dragon),
      },
      {
        label: () =>
            h(
                RouterLink,
                {
                  href: '#',
                  to: {
                    name: 'market',
                    query: {
                      name: "个股研报",
                    }
                  },
                  onClick: () => {
                    activeKey.value = 'market'
                    EventsEmit("changeMarketTab", {ID: 0, name: '个股研报'})
                  },
                },
                {default: () => '个股研报',}
            ),
        key: 'market7',
        icon: renderIcon(StockOutlined),
      },
      {
        label: () =>
            h(
                RouterLink,
                {
                  href: '#',
                  to: {
                    name: 'market',
                    query: {
                      name: "公司公告",
                    }
                  },
                  onClick: () => {
                    activeKey.value = 'market'
                    EventsEmit("changeMarketTab", {ID: 0, name: '公司公告'})
                  },
                },
                {default: () => '公司公告',}
            ),
        key: 'market8',
        icon: renderIcon(NotificationFilled),
      },
      {
        label: () =>
            h(
                RouterLink,
                {
                  href: '#',
                  to: {
                    name: 'market',
                    query: {
                      name: "行业研究",
                    }
                  },
                  onClick: () => {
                    activeKey.value = 'market'
                    EventsEmit("changeMarketTab", {ID: 0, name: '行业研究'})
                  },
                },
                {default: () => '行业研究',}
            ),
        key: 'market9',
        icon: renderIcon(ReportSearch),
      },
      {
        label: () =>
            h(
                RouterLink,
                {
                  href: '#',
                  to: {
                    name: 'market',
                    query: {
                      name: "当前热门",
                    }
                  },
                  onClick: () => {
                    activeKey.value = 'market'
                    EventsEmit("changeMarketTab", {ID: 0, name: '当前热门'})
                  },
                },
                {default: () => '当前热门',}
            ),
        key: 'market10',
        icon: renderIcon(Gripfire),
      },
      {
        label: () =>
            h(
                RouterLink,
                {
                  href: '#',
                  to: {
                    name: 'market',
                    query: {
                      name: "名站优选",
                    }
                  },
                  onClick: () => {
                    activeKey.value = 'market'
                    EventsEmit("changeMarketTab", {ID: 0, name: '名站优选'})
                  },
                },
                {default: () => '名站优选',}
            ),
        key: 'market11',
        icon: renderIcon(FirefoxBrowser),
      },
    ]
  },
  {
    label: () =>
        h(
            RouterLink,
            {
              to: {
                name: 'klineAnalysis',
              },
              onClick: () => {
                activeKey.value = 'klineAnalysis'
              },
            },
            {default: () => 'K线分析'}
        ),
    key: 'klineAnalysis',
    icon: renderIcon(StatsChartOutline),
  },
  {
    label: () =>
        h(
            RouterLink,
            {
              to: {
                name: 'fund',
                query: {
                  name: '基金自选',
                },
              },
              onClick: () => {
                activeKey.value = 'fund'
              },
            },
            {default: () => '基金自选',}
        ),
    show: enableFund.value,
    key: 'fund',
    icon: renderIcon(SparklesOutline),
    children: [
      {
        label: () =>
            h(
                RouterLink,
                {
                  to: {name: 'fund', query: {name: '基金自选'}},
                  onClick: () => {
                    activeKey.value = 'fund'
                    EventsEmit("changeFundTab", {name: '基金自选'})
                  },
                },
                {default: () => '基金自选'}
            ),
        key: 'fundFollow',
        icon: renderIcon(StarOutline),
      },
      {
        label: () =>
            h(
                RouterLink,
                {
                  to: {name: 'fund', query: {name: '基金排行'}},
                  onClick: () => {
                    activeKey.value = 'fund'
                    EventsEmit("changeFundTab", {name: '基金排行'})
                  },
                },
                {default: () => '基金排行'}
            ),
        key: 'fundRanking',
        icon: renderIcon(TrendingUp),
      },
    ]
  },
  {
    label: () =>
        h(
            RouterLink,
            {
              to: {
                name: 'agent',
                query: {
                  name:"Ai智能体",
                },
                onClick: () => {
                  activeKey.value = 'agent'
                },
              }
            },
            {default: () => 'Ai智能体'}
        ),
    key: 'agent',
    show:enableAgent.value,
    icon: renderIcon(Robot),
  },
    {
      label: () =>
          h(
              RouterLink,
              {
                to: {
                  name: 'research',
                  query: {
                    name:"研究中心",
                  },
                },
                onClick: () => {
                  activeKey.value = 'research'
                  setTimeout(() => {
                    EventsEmit("changeResearchTab", {ID: 0, name: 'AI分析报告'})
                  }, 100)
                },
              },
              {default: () => '研究中心'}
          ),
      key: 'research',
      icon: renderIcon(FlaskOutline),
      children:[
          {
            label: () =>
                h(
                    RouterLink,
                    {
                      to: {
                        name: 'research',
                        query: {
                          name:"AI分析报告",
                        },
                      },
                      onClick: () => {
                        activeKey.value = 'research'
                        setTimeout(() => {
                          EventsEmit("changeResearchTab", {ID: 0, name: 'AI分析报告'})
                        }, 100)
                      },
                    },
                    {default: () => 'AI分析报告'}
                ),
            key: 'research1',
            icon: renderIcon(ReportAnalytics),
          },
        {
          label: () =>
              h(
                  RouterLink,
                  {
                    to: {
                      name: 'research',
                      query: {
                        name:"股票推荐记录",
                      },
                    },
                    onClick: () => {
                      activeKey.value = 'research'
                      setTimeout(() => {
                        EventsEmit("changeResearchTab", {ID: 1, name: '股票推荐记录'})
                      }, 100)
                    },
                  },
                  {default: () => '股票推荐记录'}
              ),
          key: 'research2',
          icon: renderIcon(Star),
        },
        {
          label: () =>
              h(
                  RouterLink,
                  {
                    to: {
                      name: 'research',
                      query: {
                        name:"异动监控",
                      },
                    },
                    onClick: () => {
                      activeKey.value = 'research'
                      setTimeout(() => {
                        EventsEmit("changeResearchTab", {ID: 2, name: '异动监控'})
                      }, 100)
                    },
                  },
                  {default: () => '异动监控'}
              ),
          key: 'stockChanges',
          icon: renderIcon(TrendingUp),
        },
        {
          label: () =>
              h(
                  RouterLink,
                  {
                    to: {
                      name: 'research',
                      query: {
                        name:"涨停梯队",
                      },
                    },
                    onClick: () => {
                      activeKey.value = 'research'
                      setTimeout(() => {
                        EventsEmit("changeResearchTab", {ID: 9, name: '涨停梯队'})
                      }, 100)
                    },
                  },
                  {default: () => '涨停梯队'}
              ),
          key: 'uplimitLadder',
          icon: renderIcon(LocalFireDepartmentRound),
        },
        {
          label: () =>
              h(
                  RouterLink,
                  {
                    to: {
                      name: 'research',
                      query: {
                        name:"提示词模板",
                      },
                    },
                    onClick: () => {
                      activeKey.value = 'research'
                      setTimeout(() => {
                        EventsEmit("changeResearchTab", {ID: 3, name: '提示词模板'})
                      }, 100)
                    },
                  },
                  {default: () => '提示词模板'}
              ),
          key: 'research3',
          icon: renderIcon(Prompt),
        },
        {
          label: () =>
              h(
                  RouterLink,
                  {
                    to: {
                      name: 'research',
                      query: {
                        name:"提示词广场",
                      },
                    },
                    onClick: () => {
                      activeKey.value = 'research'
                      setTimeout(() => {
                        EventsEmit("changeResearchTab", {ID: 10, name: '提示词广场'})
                      }, 100)
                    },
                  },
                  {default: () => '提示词广场'}
              ),
          key: 'promptPlaza',
          icon: renderIcon(GlobeOutline),
        },
        {
          label: () =>
              h(
                  RouterLink,
                  {
                    to: {
                      name: 'research',
                      query: {
                        name:"问答广场",
                      },
                    },
                    onClick: () => {
                      activeKey.value = 'research'
                      setTimeout(() => {
                        EventsEmit("changeResearchTab", {ID: 11, name: '问答广场'})
                      }, 100)
                    },
                  },
                  {default: () => '问答广场'}
              ),
          key: 'promptQa',
          icon: renderIcon(ChatbubblesOutline),
        },
        {
          label: () =>
              h(
                  RouterLink,
                  {
                    to: {
                      name: 'research',
                      query: {
                        name:"形态选股",
                      },
                    },
                    onClick: () => {
                      activeKey.value = 'research'
                      setTimeout(() => {
                        EventsEmit("changeResearchTab", {ID: 3, name: '形态选股'})
                      }, 100)
                    },
                  },
                  {default: () => '形态选股'}
              ),
          key: 'research4',
          icon: renderIcon(SearchOutline),
        },
        {
          label: () =>
              h(
                  RouterLink,
                  {
                    to: {
                      name: 'research',
                      query: {
                        name:"指标选股",
                      },
                    },
                    onClick: () => {
                      activeKey.value = 'research'
                      setTimeout(() => {
                        EventsEmit("changeResearchTab", {ID: 0, name: '指标选股'})
                      }, 100)
                    },
                  },
                  {default: () => '指标选股'}
              ),
          key: 'research_select_stock',
          icon: renderIcon(BoxSearch20Regular),
        },
        {
          label: () =>
              h(
                  RouterLink,
                  {
                    to: {
                      name: 'research',
                      query: {
                        name:"定时任务",
                      },
                    },
                    onClick: () => {
                      activeKey.value = 'research'
                      setTimeout(() => {
                        EventsEmit("changeResearchTab", {ID: 5, name: '定时任务'})
                      }, 100)
                    },
                  },
                  {default: () => '定时任务'}
              ),
          key: 'research5',
          icon: renderIcon(TimeOutline),
        },
        {
          label: () =>
              h(
                  RouterLink,
                  {
                    to: {
                      name: 'research',
                      query: {
                        name:"交易日志",
                      },
                    },
                    onClick: () => {
                      activeKey.value = 'research'
                      setTimeout(() => {
                        EventsEmit("changeResearchTab", {ID: 6, name: '交易日志'})
                      }, 100)
                    },
                  },
                  {default: () => '交易日志(beta)'}
              ),
          key: 'research6',
          icon: renderIcon(MoneyCollectOutlined),
        },
        {
          label: () =>
              h(
                  RouterLink,
                  {
                    to: {
                      name: 'research',
                    },
                    onClick: () => {
                      activeKey.value = 'research'
                      setTimeout(() => {
                        EventsEmit("changeResearchTab", {ID: 7, name: 'MCP服务'})
                      }, 100)
                    },
                  },
                  {default: () => 'MCP服务'}
              ),
          key: 'mcpServers',
          icon: renderIcon(ServerOutline),
        },
        {
          label: () =>
              h(
                  RouterLink,
                  {
                    to: {
                      name: 'research',
                    },
                    onClick: () => {
                      activeKey.value = 'research'
                      setTimeout(() => {
                        EventsEmit("changeResearchTab", {ID: 8, name: '技能管理'})
                      }, 100)
                    },
                  },
                  {default: () => '技能管理'}
              ),
          key: 'skills',
          icon: renderIcon(FlashOutline),
          show: false,
        },
      ],
    },
  {
    label: () =>
        h(
            RouterLink,
            {
              to: {
                name: 'settings',
                query: {
                  name:"设置",
                },
                onClick: () => {
                  activeKey.value = 'settings'
                },
              }
            },
            {default: () => '设置'}
        ),
    key: 'settings',
    icon: renderIcon(SettingsOutline),
  },
  {
    label: () =>
        h(
            RouterLink,
            {
              to: {
                name: 'about',
                query: {
                  name:"关于",
                }
              },
              onClick: () => {
                activeKey.value = 'about'
              },
            },
            {default: () => '关于'}
        ),
    key: 'about',
    icon: renderIcon(LogoGithub),
    show: true,
  },
  {
    show:false,
    label: () => h("a", {
      href: '#',
      onClick: toggleFullscreen,
      title: '全屏 Ctrl+F 退出全屏 Esc',
    }, {default: () => isFullscreen.value ? '取消全屏' : '全屏'}),
    key: 'full',
    icon: renderIcon(ExpandOutline),
  },
  // {
  //   label: ()=> h("a", {
  //     href: 'javascript:void(0)',
  //     style: 'cursor: move;',
  //     onClick: toggleStartMoveWindow,
  //   }, { default: () => '移动' }),
  //   key: 'move',
  //   icon: renderIcon(MoveOutline),
  // },
])

function renderIcon(icon) {
  return () => h(NIcon, null, {default: () => h(icon)})
}

function toggleFullscreen(e) {
  activeKey.value = 'full'
  if (document.fullscreenElement) {
    document.exitFullscreen?.()
    isFullscreen.value = false
    return
  }
  document.documentElement.requestFullscreen?.()
  isFullscreen.value = true
}

EventsOn("realtime_profit", (data) => {
  realtimeProfit.value = data
})
EventsOn("telegraph", (data) => {
  telegraph.value = data
})

EventsOn("loadingMsg", (data) => {
  if(data==="done"){
    loadingMsg.value = "加载完成..."
    EventsEmit("loadingDone", "app")
    loading.value  = false
  }else{
    loading.value  = true
    loadingMsg.value = data
  }
})

setTimeout(() => {
  if (loading.value) {
    loading.value = false
    loadingMsg.value = "加载完成..."
    EventsEmit("loadingDone", "app")
  }
}, 8000)

onBeforeUnmount(() => {
  stopFeed("app.motto")
  EventsOff("realtime_profit")
  EventsOff("loadingMsg")
  EventsOff("telegraph")
  EventsOff("newsPush")
})

window.onerror = function (msg, source, lineno, colno, error) {
  // 将错误信息发送给后端
  EventsEmit("frontendError", {
    page: "App.vue",
    message: msg,
    source: source,
    lineno: lineno,
    colno: colno,
    error: error ? error.stack : null,
  });
  return true;
};

onBeforeMount(() => {
  GetVersionInfo().then(result => {
    if(result.officialStatement){
      content.value = result.officialStatement+"\n\n"+content.value
    }
    officialStatement.value = result.officialStatement || ""
    updateMarketStatus()
  }).catch(err => {
    console.error("GetVersionInfo error:", err)
  })

  GetGroupList().then(result => {
    groupList.value = result
    menuOptions.value.map((item) => {
      if (item.key === 'stock') {
        item.children.push(...groupList.value.map(item => {
          return {
            label: () =>
                h(
                    'a',
                    {
                      href: '#',
                      type: 'info',
                      onClick: () => {
                        router.push({
                          name: 'stock',
                          query: {
                            groupName: item.name,
                            groupId: item.ID,
                          },
                        })
                        setTimeout(() => {
                          EventsEmit("changeTab", item)
                        }, 100)
                      },
                      to: {
                        name: 'stock',
                        query: {
                          groupName: item.name,
                          groupId: item.ID,
                        },
                      }
                    },
                    {default: () => item.name,}
                ),
            key: item.ID,
          }
        }))
      }
    })
  }).catch(err => {
    console.error("GetGroupList error:", err)
  })


  GetConfig().then((res) => {
    enableFund.value = res.enableFund
    enableAgent.value = res.enableAgent

    menuOptions.value.filter((item) => {
      if (item.key === 'fund') {
        item.show = res.enableFund
      }
      if (item.key === 'agent') {
        item.show = res.enableAgent
      }
    })

    if (res.darkTheme) {
      enableDarkTheme.value = darkTheme
    } else {
      enableDarkTheme.value = null
    }
  }).catch(err => {
    console.error("GetConfig error:", err)
  })
})

onMounted(() => {
  // 标题/交易状态由 marketClock + watch 驱动；这里只保留格言的定时刷新。
  registerFeed("app.motto", { fetch: refreshMotto, intervalMs: 60000 })
  GetConfig().then((res) => {
    if (res.enableNews) {
      enableNews.value = true
    }
    enableFund.value = res.enableFund
    enableAgent.value = res.enableAgent
    const {notification } =createDiscreteApi(["notification"], {
      configProviderProps: {
        theme: enableDarkTheme.value ? darkTheme : lightTheme ,
        max: 3,
      },
    })
    EventsOn("newsPush", (data) => {
      //console.log(data)
      if(data.isRed){
        notification.create({
          //type:"error",
         // avatar: () => h(NIcon,{component:Notifications,color:"red"}),
          title: data.time,
          content: () => h('div',{type:"error",style:{
              "text-align":"left",
              "font-size":"14px",
              "color":"#f67979"
            }}, { default: () => data.content }),
          meta: () => h(NText,{type:"warning"}, { default: () => data.source}),
          duration:1000*40,
        })
      }else{
         notification.create({
          //type:"info",
          //avatar: () => h(NIcon,{component:Notifications}),
          title: data.time,
          content: () => h('div',{type:"info",style:{
            "text-align":"left",
              "font-size":"14px",
              "color": data.source==="go-stock"?"#F98C24":"#549EC8"
            }}, { default: () => data.content }),
          meta: () => h(NText,{type:"warning"}, { default: () => data.source}),
          duration:1000*30 ,
        })
      }
    })
  }).catch(err => {
    console.error("GetConfig(onMounted) error:", err)
  })
})
</script>
<template>
  <n-config-provider ref="containerRef" :theme="enableDarkTheme" :locale="zhCN" :date-locale="dateZhCN">
    <n-message-provider>
      <n-notification-provider>
        <n-modal-provider>
          <n-dialog-provider>
            <n-watermark
                :content="''"
                cross
                selectable
                :font-size="16"
                :line-height="16"
                :width="500"
                :height="400"
                :x-offset="50"
                :y-offset="150"
                :rotate="-15"
            >
<!--              <FloatingAiAssistant />-->
              <FloatingAgentAssistant />
              <n-flex class="app-shell" :class="{ 'app-shell--mobile': isMobile }">
                <n-grid x-gap="12" :cols="1" class="app-shell__grid">
                  <n-gi>
                    <n-spin :show="loading">
                      <template #description>
                        {{ loadingMsg }}
                      </template>
                      <n-marquee :speed="100" style="position: relative;top:0;z-index: 19;width: 100%"
                                 v-if="(telegraph.length>0)&&(enableNews)">
                        <n-tag type="warning" v-for="item in telegraph" style="margin-right: 10px">
                          {{ item }}
                        </n-tag>
                      </n-marquee>
                      <n-scrollbar :style="contentStyle" class="app-content-scroll">
                        <n-skeleton v-if="loading" height="calc(100vh)" />
                        <RouterView/>
                      </n-scrollbar>
                    </n-spin>
                  </n-gi>
                  <n-gi class="desktop-bottom-menu desktop-only" style="position: fixed;bottom:0;z-index: 9;width: 100%;">
                    <n-card size="small" style="">
                      <n-menu style="font-size: 18px;"
                              v-model:value="activeKey"
                              mode="horizontal"
                              :options="menuOptions"
                              responsive
                      />
                    </n-card>
                  </n-gi>
                </n-grid>
              </n-flex>
              <nav class="mobile-bottom-nav mobile-only" aria-label="移动端主导航">
                <button
                    v-for="item in mobileBottomNavItems"
                    :key="item.key"
                    class="mobile-bottom-nav__item"
                    :class="{ 'mobile-bottom-nav__item--active': isMobileBottomNavActive(item) }"
                    type="button"
                    @click="handleMobileNav(item)"
                >
                  <n-icon size="20">
                    <component :is="item.icon" />
                  </n-icon>
                  <span>{{ item.label }}</span>
                </button>
              </nav>
              <n-drawer
                  v-model:show="mobileMenuVisible"
                  class="mobile-menu-drawer"
                  placement="bottom"
                  height="82vh"
              >
                <n-drawer-content title="全部功能" closable>
                  <div class="mobile-more-grid">
                    <div v-for="section in mobileMoreSections" :key="section.category" class="mobile-more-section">
                      <div class="mobile-more-section__title">{{ section.category }}</div>
                      <div class="mobile-more-section__tiles">
                        <button
                            v-for="tile in section.tiles"
                            :key="tile.label"
                            type="button"
                            class="mobile-more-tile"
                            @click="onTileSelect(tile)"
                        >
                          <n-icon size="22">
                            <component :is="tile.icon" />
                          </n-icon>
                          <span>{{ tile.label }}</span>
                        </button>
                      </div>
                    </div>
                  </div>
                </n-drawer-content>
              </n-drawer>
            </n-watermark>
          </n-dialog-provider>
        </n-modal-provider>
      </n-notification-provider>
    </n-message-provider>
  </n-config-provider>
</template>
<style>
.mobile-more-grid {
    padding: 4px 0 calc(var(--safe-bottom) + 8px);
}

.mobile-more-section {
    margin-bottom: 18px;
}

.mobile-more-section__title {
    color: var(--n-text-color-3, #999);
    font-size: 13px;
    font-weight: 600;
    margin: 0 4px 10px;
}

.mobile-more-section__tiles {
    display: grid;
    gap: 10px;
    grid-template-columns: repeat(4, 1fr);
}

.mobile-more-tile {
    align-items: center;
    appearance: none;
    background: var(--n-color-target, rgba(0, 0, 0, 0.03));
    border: 1px solid var(--n-border-color, #efeff5);
    border-radius: 12px;
    color: var(--n-text-color, #333);
    display: flex;
    flex-direction: column;
    font: inherit;
    gap: 6px;
    justify-content: center;
    min-width: 0;
    padding: 12px 4px;
}

.mobile-more-tile:active {
    background: rgba(24, 160, 88, 0.12);
    border-color: var(--n-primary-color, #18a058);
}

.mobile-more-tile span {
    font-size: 12px;
    line-height: 1.2;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    width: 100%;
}
</style>
