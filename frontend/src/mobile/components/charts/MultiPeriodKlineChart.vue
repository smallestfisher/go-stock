<script setup>
import { ref, computed, onMounted, onBeforeUnmount, watch, nextTick, shallowRef } from 'vue'
import {
  bollingerBands, sarValues, supertrendValues, donchianChannelValues, ichimokuValues,
  macdBundle, kdjBundle, rsiBundle, cciValues, williamsRValues, atrValues,
  obvValues, mfiValues, adxValues,
  emaFinite, kamaValues, hullMaValues, demaValues, temaValues, vwapValues,
  keltnerChannelValues, alligatorValues, vwapBandsValues, pivotPointsValues, zigzagValues,
  stochRsiValues, cmoValues, trixValues, rocValues, coppockValues, smiValues, aoValues,
  cmfValues, adValues, forceIndexValues, chaikinOscValues,
  chopValues, massIndexValues, ulcerIndexValues, elderRayValues, satsValues,
  emaLeadingNull, ttmSqueezeValues, aroonValues, smaValues,
} from '../../../components/kline/calc'

const props = defineProps({
  // K线数据：[{ day, open, close, high, low, volume, amount, changePercent, ... }]
  data: {
    type: Array,
    default: () => []
  },
  height: {
    type: Number,
    default: 360
  },
  period: {
    type: String,
    default: 'day'
  },
  // 已叠加的指标（由指标面板控制）
  indicators: {
    type: Array,
    default: () => ['MA', 'VOL']
  },
  // OHLC 信息条显示的字段（窄屏精简用）。默认全集（兼容个股 K 线分析）。
  // 可选项：open/close/high/low/changePercent/changeValue/volume/amount/amplitude/avgAmp5/avgAmp10/avgAmp20/turnoverRate/volumeRatio
  fields: {
    type: Array,
    default: () => ['open', 'close', 'high', 'low', 'changePercent', 'changeValue', 'volume', 'amount', 'amplitude', 'avgAmp5', 'avgAmp10', 'avgAmp20', 'turnoverRate', 'volumeRatio']
  }
})

// 副图指标清单（顺序即从上到下排列）
const SUB_INDICATORS = ['MACD', 'KDJ', 'RSI', 'CCI', 'WR', 'STOCHRSI', 'CMO', 'TRIX', 'ROC', 'COPPOCK', 'SMI', 'AO', 'OBV', 'MFI', 'CMF', 'AD', 'FI', 'CHAIKINOSC', 'ATR', 'TTM', 'AVGAMP', 'MASSINDEX', 'ULCER', 'SATS', 'ADX', 'AROON', 'CHOP', 'ELDERRAY', 'SIGNALRATIO']
// 当前生效的副图
const activeSubs = computed(() => SUB_INDICATORS.filter(c => props.indicators.includes(c)))

const chartRef = ref(null)
let chart = null
const echartsLib = shallowRef(null)

// 光标悬停的 K 线索引；null 表示显示最新一根
const hoverIndex = ref(null)

function cssColor(name, fallback) {
  if (typeof window === 'undefined') return fallback
  return getComputedStyle(document.documentElement).getPropertyValue(name).trim() || fallback
}

// 成交量/成交额中文化（对齐桌面端 formatVolumeCn / formatAmountCn）
function fmtVolCn(v) {
  const n = Number(v)
  if (!Number.isFinite(n)) return '--'
  if (n >= 1e8) return (n / 1e8).toFixed(2) + '亿'
  if (n >= 1e4) return (n / 1e4).toFixed(2) + '万'
  return String(n)
}
function fmtAmountCn(v) {
  const n = Number(v)
  if (!Number.isFinite(n)) return '--'
  if (n >= 1e8) return (n / 1e8).toFixed(2) + '亿'
  if (n >= 1e4) return (n / 1e4).toFixed(2) + '万'
  return String(n)
}
function fmtPct(v) {
  const n = Number(v)
  return Number.isFinite(n) ? n.toFixed(2) + '%' : '--'
}
function fmtSigned2(v) {
  const n = Number(v)
  if (!Number.isFinite(n)) return '--'
  return (n > 0 ? '+' : '') + n.toFixed(2)
}

// 均幅 N：基于到当前索引(含)的振幅序列
function avgAmp(arr, period) {
  if (!arr || arr.length < period) return NaN
  let s = 0, cnt = 0
  for (let i = arr.length - period; i < arr.length; i++) {
    if (Number.isFinite(arr[i])) { s += arr[i]; cnt++ }
  }
  return cnt === period ? s / cnt : NaN
}

// ---- 当前激活的 K 线（光标悬停 or 最新），含均幅5/10/20 ----
const activeBar = computed(() => {
  const arr = props.data || []
  if (!arr.length) return null
  const idx = (hoverIndex.value != null && hoverIndex.value < arr.length) ? hoverIndex.value : arr.length - 1
  const d = arr[idx]
  const chg = Number(d.changePercent)
  const sign = chg > 0 ? 1 : chg < 0 ? -1 : 0
  const ohlcC = sign > 0 ? 'var(--m-color-rise)' : sign < 0 ? 'var(--m-color-fall)' : 'var(--m-text-secondary)'

  // 均幅序列（到 idx 含）
  const amps = []
  for (let i = 0; i <= idx; i++) {
    const row = arr[i]
    const rawAmp = Number(row.amplitude)
    const o = Number(row.open), h = Number(row.high), l = Number(row.low)
    if (Number.isFinite(rawAmp)) amps.push(rawAmp)
    else if (Number.isFinite(o) && o > 0) amps.push((h - l) / o * 100)
    else amps.push(NaN)
  }
  const a5 = avgAmp(amps, 5), a10 = avgAmp(amps, 10), a20 = avgAmp(amps, 20)

  return {
    day: d.day || '',
    isLatest: hoverIndex.value == null,
    open: Number(d.open).toFixed(2),
    close: Number(d.close).toFixed(2),
    high: Number(d.high).toFixed(2),
    low: Number(d.low).toFixed(2),
    changePercent: fmtPct(d.changePercent),
    changeValue: fmtSigned2(d.changeValue),
    volume: fmtVolCn(d.volume),
    amount: fmtAmountCn(d.amount),
    amplitude: fmtPct(d.amplitude),
    avgAmp5: Number.isFinite(a5) ? a5.toFixed(2) + '%' : '--',
    avgAmp10: Number.isFinite(a10) ? a10.toFixed(2) + '%' : '--',
    avgAmp20: Number.isFinite(a20) ? a20.toFixed(2) + '%' : '--',
    turnoverRate: fmtPct(d.turnoverRate),
    volumeRatio: (() => { const n = Number(d.volumeRatio); return Number.isFinite(n) ? n.toFixed(2) : '--' })(),
    cOpenClose: ohlcC,
    cChg: ohlcC,
  }
})

// ---- OHLC 信息条字段配置（受 fields prop 控制，窄屏可精简） ----
// color: 'oc' = 开收同色(涨跌)，'chg' = 涨跌色，'rise' = 固定红，'fall' = 固定绿
const FIELD_DEFS = [
  { key: 'open', label: '开', color: 'oc' },
  { key: 'close', label: '收', color: 'oc' },
  { key: 'high', label: '高', color: 'rise' },
  { key: 'low', label: '低', color: 'fall' },
  { key: 'changePercent', label: '涨跌幅', color: 'chg' },
  { key: 'changeValue', label: '涨跌额', color: 'chg' },
  { key: 'volume', label: '成交量' },
  { key: 'amount', label: '成交额' },
  { key: 'amplitude', label: '振幅' },
  { key: 'avgAmp5', label: '均幅5' },
  { key: 'avgAmp10', label: '均幅10' },
  { key: 'avgAmp20', label: '均幅20' },
  { key: 'turnoverRate', label: '换手率' },
  { key: 'volumeRatio', label: '量比', color: 'chg' },
]
const visibleFields = computed(() => FIELD_DEFS.filter(f => props.fields.includes(f.key)))
function fieldStyle(f) {
  if (!f.color) return null
  if (f.color === 'oc') return { color: activeBar.value?.cOpenClose }
  if (f.color === 'chg') return { color: activeBar.value?.cChg }
  return null
}

// ---- 数据规范化 ----
const norm = computed(() => {
  const arr = props.data || []
  const category = []
  const values = []
  const volumes = []
  for (let i = 0; i < arr.length; i++) {
    const d = arr[i]
    const open = Number(d.open) || 0
    const close = Number(d.close) || 0
    const low = Number(d.low) || 0
    const high = Number(d.high) || 0
    const volume = Number(d.volume) || 0
    category.push(d.day || '')
    values.push([open, close, low, high])
    volumes.push([i, +(volume / 10000).toFixed(2), close >= open ? 1 : -1])
  }
  return { category, values, volumes }
})

function calcMA(period, values) {
  const result = []
  for (let i = 0; i < values.length; i++) {
    if (i < period - 1) { result.push('-'); continue }
    let sum = 0
    for (let j = 0; j < period; j++) sum += values[i - j][1]
    result.push(+(sum / period).toFixed(2))
  }
  return result
}

// 当前叠加指标的颜色图例（说明每条线代表什么 + 当前值）
const MA_COLORS = { MA5: '#3b82f6', MA10: '#22c55e', MA20: '#f59e0b', MA30: '#ef4444' }
const legendItems = computed(() => {
  const items = []
  if (props.indicators.includes('MA')) {
    const idx = (hoverIndex.value != null && hoverIndex.value < norm.value.values.length)
      ? hoverIndex.value : norm.value.values.length - 1
    for (const p of [5, 10, 20, 30]) {
      const arr = calcMA(p, norm.value.values)
      const v = arr[idx]
      items.push({ label: `MA${p}`, color: MA_COLORS[`MA${p}`], value: (v != null && v !== '-') ? Number(v).toFixed(2) : '--' })
    }
  }
  return items
})

const zoomStart = computed(() => {
  const n = norm.value.category.length
  if (n <= 60) return 0
  return Math.max(0, Math.round((1 - 60 / n) * 100))
})

// 图表垂直布局常量（buildOption 与 chartHeight 共用，保证一致）
// 模型：PAD_TOP → 主图 → [SUB_GAP + 副图]×N → ZOOM_GAP → slider(ZOOM_H) → PAD_BOTTOM
const PAD_TOP = 8
const ZOOM_H = 16             // dataZoom slider 高度
const ZOOM_BOTTOM = 6         // slider 距 canvas 底部
const ZOOM_GAP = 18           // 最后一个图(主图/副图)底边 与 slider 顶部 的间距
const ZOOM_AREA = ZOOM_GAP + ZOOM_H + ZOOM_BOTTOM  // 底部固定区总高(32)
const SUB_GAP = 6             // 副图之间间距
const SUB_H = 58              // 单个副图高度
const MAIN_MIN = 150          // 主图最小高度

// 实际所需 canvas 高度：随副图数增长；父页面 chart-area 可滚动。
const chartHeight = computed(() => {
  const showVol = props.indicators.includes('VOL')
  const subCount = activeSubs.value.length + (showVol ? 1 : 0)
  const subsTotal = subCount * SUB_H + Math.max(0, subCount) * SUB_GAP
  return Math.max(props.height, PAD_TOP + MAIN_MIN + subsTotal + ZOOM_AREA)
})

// OHLC 数值数组（供指标计算）
const closes = computed(() => norm.value.values.map(v => v[1]))
const highs = computed(() => norm.value.values.map(v => v[3]))
const lows = computed(() => norm.value.values.map(v => v[2]))
const vols = computed(() => (props.data || []).map(d => Number(d.volume) || 0))

// 均幅序列：每根 K 的振幅（缺失则用 (高-低)/开*100），供「均幅」副图
function avgAmpSeries(cl, hi, lo) {
  const arr = (props.data || []).map(d => {
    const raw = Number(d.amplitude)
    const o = Number(d.open), h = Number(d.high), l = Number(d.low)
    if (Number.isFinite(raw)) return raw
    if (Number.isFinite(o) && o > 0) return (h - l) / o * 100
    return null
  })
  return arr
}

// 信号比例副图（对齐桌面端 signalRatio）：逐根 K 线评估全部指标，
// 输出 看多%/看空%/净% 三条曲线。预计算指标数组一次再逐根判定，O(N)。
function buildSignalRatioSeries(gridIdx, c) {
  const cl = closes.value, hi = highs.value, lo = lows.value, vl = vols.value
  const n = cl.length
  const bullish = new Array(n).fill(null)
  const bearish = new Array(n).fill(null)
  const net = new Array(n).fill(null)
  if (n < 2) return []

  // 预计算所有指标数组（只算一次）
  const ma5 = smaValues(cl, 5), ma10 = smaValues(cl, 10), ma20 = smaValues(cl, 20), ma60 = smaValues(cl, 60)
  const ema12 = emaFinite(cl, 12), ema21 = emaFinite(cl, 21)
  const { upper: bollU, mid: bollM, lower: bollL } = bollingerBands(cl, 20, 2)
  const vw = vwapValues(hi, lo, cl, vl, 20)
  const dema21 = demaValues(cl, 21), tema21 = temaValues(cl, 21)
  const kama10 = kamaValues(cl, 10, 2, 30), hull9 = hullMaValues(cl, 9)
  const { upper: kU, lower: kL } = keltnerChannelValues(hi, lo, cl, 20, 10, 1.5)
  const { direction: stDir } = supertrendValues(hi, lo, cl, 10, 3)
  const { tenkan: ichTen, kijun: ichKij, spanA: ichSA, senkouB: ichSB } = ichimokuValues(hi, lo, cl)
  const { direction: sarDir } = sarValues(hi, lo, cl, 0.02, 0.2)
  const { upper: dcU, lower: dcL } = donchianChannelValues(hi, lo, 20)
  const { jaw: agJ, teeth: agT, lips: agL } = alligatorValues(hi, lo, cl)
  const { directions: zzDir } = zigzagValues(hi, lo, cl, 5)
  const { direction: satsDir } = satsValues(hi, lo, cl, vl)
  const { pp: pivPP, s1: pivS1, r1: pivR1 } = pivotPointsValues(hi, lo, cl)
  const { vwap: vbM, upper: vbU, lower: vbL } = vwapBandsValues(hi, lo, cl, vl)
  const { dif: macdDif, dea: macdDea, hist: macdHist } = macdBundle(cl)
  const rsi14 = rsiBundle(cl, 14)
  const { K: kdjK, D: kdjD, J: kdjJ } = kdjBundle(hi, lo, cl, 9)
  const cci20 = cciValues(hi, lo, cl, 20), wr14 = williamsRValues(hi, lo, cl, 14)
  const { k: stochK, d: stochD } = stochRsiValues(cl, 14, 14, 3, 3)
  const { adx: adxVal, diP: adxP, diM: adxM } = adxValues(hi, lo, cl, 14)
  const { up: arUp, down: arDown } = aroonValues(hi, lo, 25)
  const cmo14 = cmoValues(cl, 14), trix15 = trixValues(cl, 15), trixSig = emaLeadingNull(trix15, 9)
  const roc12 = rocValues(cl, 12), coppock = coppockValues(cl)
  const { smi: smiD, signal: smiS } = smiValues(hi, lo, cl)
  const ao534 = aoValues(hi, lo), obv = obvValues(cl, vl), mfi14 = mfiValues(hi, lo, cl, vl, 14)
  const cmf20 = cmfValues(hi, lo, cl, vl, 20), adLine = adValues(hi, lo, cl, vl)
  const fi13 = forceIndexValues(cl, vl, 13), co310 = chaikinOscValues(hi, lo, cl, vl, 3, 10)
  const atr14 = atrValues(hi, lo, cl, 14), chop14 = chopValues(hi, lo, cl, 14)
  const miVal = massIndexValues(hi, lo), uiVal = ulcerIndexValues(cl)
  const { squeeze: ttmSq, momentum: ttmMo } = ttmSqueezeValues(hi, lo, cl)
  const { bullPower: erBull, bearPower: erBear } = elderRayValues(hi, lo, cl, 13)

  // ZigZag 方向缓存：取最近非零方向
  const zzLastDir = new Array(n).fill(0)
  let lastNonZero = 0
  for (let i = 0; i < zzDir.length; i++) {
    if (zzDir[i] === 1 || zzDir[i] === -1) lastNonZero = zzDir[i]
    zzLastDir[i] = lastNonZero
  }

  const v = (arr, i) => (i >= 0 && i < arr.length && arr[i] != null && Number.isFinite(arr[i])) ? arr[i] : null
  const vPrev = (arr, i) => v(arr, i - 1)

  // 逐根评估（判定逻辑与桌面端严格一致）
  for (let i = 0; i < n; i++) {
    const cc = cl[i]
    if (cc == null || !Number.isFinite(cc)) continue
    let bull = 0, bear = 0, neut = 0, osci = 0, cnt = 0
    // MA
    { const a5=v(ma5,i),a10=v(ma10,i),a20=v(ma20,i),a60=v(ma60,i)
      if(a5!=null&&a10!=null&&a20!=null&&a60!=null){cnt++;if(a5>a10&&a10>a20&&a20>a60)bull++;else if(a5<a10&&a10<a20&&a20<a60)bear++;else if((a5>a20&&a10<a60)||(a5<a20&&a10>a60))osci++;else neut++;} }
    // EMA
    { const e12=v(ema12,i),e21=v(ema21,i)
      if(e12!=null&&e21!=null){cnt++;if(e12>e21)bull++;else if(e12<e21)bear++;else neut++;} }
    // BOLL
    { const bu=v(bollU,i),bm=v(bollM,i),bl=v(bollL,i)
      if(bu!=null&&bm!=null&&bl!=null){cnt++;if(cc>bu)bull++;else if(cc<bl)bear++;else if(cc>bm)osci++;else neut++;} }
    // VWAP
    { const vwv=v(vw,i);if(vwv!=null){cnt++;if(cc>vwv)bull++;else if(cc<vwv)bear++;else neut++;} }
    // DEMA
    { const dv=v(dema21,i);if(dv!=null){cnt++;if(cc>dv)bull++;else if(cc<dv)bear++;else neut++;} }
    // TEMA
    { const tv2=v(tema21,i);if(tv2!=null){cnt++;if(cc>tv2)bull++;else if(cc<tv2)bear++;else neut++;} }
    // KAMA
    { const kv=v(kama10,i),kp=vPrev(kama10,i)
      if(kv!=null&&kp!=null){cnt++;if(cc>kv&&kv>kp)bull++;else if(cc<kv&&kv<kp)bear++;else neut++;} }
    // HullMA
    { const hv=v(hull9,i),hp=vPrev(hull9,i)
      if(hv!=null&&hp!=null){cnt++;if(hv>hp)bull++;else if(hv<hp)bear++;else neut++;} }
    // Keltner
    { const ku=v(kU,i),kl=v(kL,i)
      if(ku!=null&&kl!=null){cnt++;if(cc>ku)bull++;else if(cc<kl)bear++;else osci++;} }
    // SuperTrend
    { const sd=v(stDir,i);if(sd!=null){cnt++;if(sd===1)bull++;else if(sd===-1)bear++;else neut++;} }
    // Ichimoku
    { const it=v(ichTen,i),ik=v(ichKij,i),isa=v(ichSA,i),isb=v(ichSB,i)
      if(it!=null&&ik!=null&&isa!=null&&isb!=null){const ct=Math.max(isa,isb),cb=Math.min(isa,isb);cnt++;if(cc>ct&&it>ik)bull++;else if(cc<cb&&it<ik)bear++;else if(cc>=cb&&cc<=ct)osci++;else neut++;} }
    // SAR
    { const sd2=v(sarDir,i);if(sd2!=null){cnt++;if(sd2===1)bull++;else if(sd2===-1)bear++;else neut++;} }
    // Donchian
    { const du=v(dcU,i),dl=v(dcL,i)
      if(du!=null&&dl!=null){cnt++;if(cc>=du)bull++;else if(cc<=dl)bear++;else osci++;} }
    // Alligator
    { const aj=v(agJ,i),at=v(agT,i),al=v(agL,i)
      if(aj!=null&&at!=null&&al!=null){cnt++;if(al>at&&at>aj)bull++;else if(al<at&&at<aj)bear++;else osci++;} }
    // ZigZag
    { const zd=zzLastDir[i];if(zd!==0){cnt++;if(zd===-1)bull++;else if(zd===1)bear++;else neut++;} }
    // SATS
    { const sd3=v(satsDir,i);if(sd3!=null){cnt++;if(sd3===1)bull++;else if(sd3===-1)bear++;else neut++;} }
    // Pivot
    { const pp=v(pivPP,i),s1=v(pivS1,i),r1=v(pivR1,i)
      if(pp!=null&&r1!=null&&s1!=null){cnt++;if(cc>r1)bull++;else if(cc<s1)bear++;else if(cc>pp)osci++;else neut++;} }
    // VWAPBands
    { const vu=v(vbU,i),vl2=v(vbL,i)
      if(vu!=null&&vl2!=null){cnt++;if(cc>vu)bull++;else if(cc<vl2)bear++;else osci++;} }
    // MACD
    { const md=v(macdDif,i),me=v(macdDea,i),mh=v(macdHist,i)
      if(md!=null&&me!=null&&mh!=null){cnt++;if(md>me&&mh>0)bull++;else if(md<me&&mh<0)bear++;else if((md>0&&mh<0)||(md<0&&mh>0))osci++;else neut++;} }
    // RSI
    { const rv=v(rsi14,i);if(rv!=null){cnt++;if(rv>70)osci++;else if(rv<30)osci++;else if(rv>50)bull++;else bear++;} }
    // KDJ
    { const kk=v(kdjK,i),kd=v(kdjD,i),kj=v(kdjJ,i)
      if(kk!=null&&kd!=null&&kj!=null){cnt++;if(kj>kk&&kk>kd&&kk<80)bull++;else if(kj<kk&&kk<kd&&kk>20)bear++;else if(kk>80)bear++;else if(kk<20)bull++;else osci++;} }
    // CCI
    { const cv=v(cci20,i);if(cv!=null){cnt++;if(cv>100)bull++;else if(cv<-100)bear++;else osci++;} }
    // W%R
    { const wv=v(wr14,i);if(wv!=null){cnt++;if(wv<-80)bull++;else if(wv>-20)bear++;else osci++;} }
    // StochRSI
    { const sk=v(stochK,i),sd4=v(stochD,i)
      if(sk!=null&&sd4!=null){cnt++;if(sk<20&&sd4<20&&sk>sd4)bull++;else if(sk>80&&sd4>80&&sk<sd4)bear++;else osci++;} }
    // ADX
    { const av=v(adxVal,i),ap=v(adxP,i),am=v(adxM,i)
      if(av!=null&&ap!=null&&am!=null){cnt++;if(av>25&&ap>am)bull++;else if(av>25&&ap<am)bear++;else osci++;} }
    // Aroon
    { const au=v(arUp,i),ad2=v(arDown,i)
      if(au!=null&&ad2!=null){cnt++;if(au>70&&ad2<30)bull++;else if(ad2>70&&au<30)bear++;else osci++;} }
    // CMO
    { const cv2=cmo14[i];if(cv2!=null&&Number.isFinite(cv2)){cnt++;if(cv2>50)bull++;else if(cv2<-50)bear++;else osci++;} }
    // TRIX
    { const tv3=v(trix15,i),ts=v(trixSig,i)
      if(tv3!=null&&ts!=null){cnt++;if(tv3>ts)bull++;else if(tv3<ts)bear++;else neut++;} }
    // ROC
    { const rv2=v(roc12,i);if(rv2!=null){cnt++;if(rv2>0)bull++;else if(rv2<0)bear++;else neut++;} }
    // Coppock
    { const cv3=v(coppock,i),cp=vPrev(coppock,i)
      if(cv3!=null&&cp!=null){cnt++;if(cv3>0&&cp<=0)bull++;else if(cv3<0)bear++;else neut++;} }
    // SMI
    { const sv=v(smiD,i),ss=v(smiS,i)
      if(sv!=null&&ss!=null){cnt++;if(sv>ss&&sv>0)bull++;else if(sv<ss&&sv<0)bear++;else osci++;} }
    // AO
    { const av2=v(ao534,i),ap2=vPrev(ao534,i)
      if(av2!=null&&ap2!=null){cnt++;if(av2>0&&av2>ap2)bull++;else if(av2<0&&av2<ap2)bear++;else if(av2>0&&av2<ap2)osci++;else neut++;} }
    // OBV
    { const ov=v(obv,i),op=vPrev(obv,i)
      if(ov!=null&&op!=null){cnt++;if(ov>op)bull++;else if(ov<op)bear++;else neut++;} }
    // MFI
    { const mv=v(mfi14,i);if(mv!=null){cnt++;if(mv>80)osci++;else if(mv<20)osci++;else if(mv>50)bull++;else bear++;} }
    // CMF
    { const cv4=v(cmf20,i);if(cv4!=null){cnt++;if(cv4>0.05)bull++;else if(cv4<-0.05)bear++;else osci++;} }
    // A/D
    { const av3=v(adLine,i),ap3=vPrev(adLine,i)
      if(av3!=null&&ap3!=null){cnt++;if(av3>ap3)bull++;else if(av3<ap3)bear++;else neut++;} }
    // FI
    { const fv=v(fi13,i);if(fv!=null){cnt++;if(fv>0)bull++;else if(fv<0)bear++;else neut++;} }
    // ChaikinOsc
    { const cv5=v(co310,i),cp2=vPrev(co310,i)
      if(cv5!=null&&cp2!=null){cnt++;if(cv5>0&&cv5>cp2)bull++;else if(cv5<0&&cv5<cp2)bear++;else osci++;} }
    // ATR
    { const av4=v(atr14,i),ap4=vPrev(atr14,i)
      if(av4!=null&&ap4!=null){cnt++;if(av4>ap4)osci++;else neut++;} }
    // CHOP
    { const cv6=v(chop14,i);if(cv6!=null){cnt++;if(cv6>61.8)osci++;else neut++;} }
    // MassIndex
    { const mv2=v(miVal,i),mp=vPrev(miVal,i)
      if(mv2!=null&&mp!=null){cnt++;if(mp>27&&mv2<27)bull++;else neut++;} }
    // UlcerIndex
    { const uv=v(uiVal,i);if(uv!=null){cnt++;if(uv<5)bull++;else if(uv>15)bear++;else neut++;} }
    // TTM
    { const sq=v(ttmSq,i),mo=v(ttmMo,i)
      if(sq!=null&&mo!=null){cnt++;if(!sq&&mo>0)bull++;else if(!sq&&mo<0)bear++;else if(sq)osci++;else neut++;} }
    // ElderRay
    { const bp=v(erBull,i),brp=v(erBear,i)
      if(bp!=null&&brp!=null){cnt++;if(bp>0&&bp>brp)bull++;else if(brp<0&&brp<bp)bear++;else osci++;} }
    // 汇总
    if (cnt > 0) {
      bullish[i] = +(bull / cnt * 100).toFixed(2)
      bearish[i] = +(bear / cnt * 100).toFixed(2)
      net[i] = +((bull - bear) / cnt * 100).toFixed(2)
    }
  }

  const mkSubLine = (name, data, color) => ({
    name, type: 'line', data, xAxisIndex: gridIdx, yAxisIndex: gridIdx,
    showSymbol: false, lineStyle: { width: 1, color },
  })
  // 看多=涨红、看空=跌绿、净=蓝（对齐移动端涨红跌绿习惯）
  return [
    mkSubLine('看多%', bullish, c.rise),
    mkSubLine('看空%', bearish, c.fall),
    mkSubLine('净%', net, '#3b82f6'),
  ]
}

// 构建某个副图指标的 series（返回数组）
function buildSubSeries(code, gridIdx, c) {
  const cl = closes.value, hi = highs.value, lo = lows.value, vl = vols.value
  const mkLine = (name, data, color) => ({
    name, type: 'line', data, xAxisIndex: gridIdx, yAxisIndex: gridIdx,
    showSymbol: false, lineStyle: { width: 1, color },
  })
  switch (code) {
    case 'MACD': {
      const { dif, dea, hist } = macdBundle(cl)
      return [
        { name: 'MACD柱', type: 'bar', data: hist, xAxisIndex: gridIdx, yAxisIndex: gridIdx,
          itemStyle: { color: (v) => v >= 0 ? c.rise : c.fall } },
        mkLine('DIF', dif, '#3b82f6'),
        mkLine('DEA', dea, '#f59e0b'),
      ]
    }
    case 'KDJ': {
      const { K, D, J } = kdjBundle(hi, lo, cl, 9)
      return [mkLine('K', K, '#3b82f6'), mkLine('D', D, '#f59e0b'), mkLine('J', J, '#ef4444')]
    }
    case 'RSI':
      return [mkLine('RSI', rsiBundle(cl, 14), '#8b5cf6')]
    case 'CCI':
      return [mkLine('CCI', cciValues(hi, lo, cl, 20), '#10b981')]
    case 'WR':
      return [mkLine('WR', williamsRValues(hi, lo, cl, 14), '#f59e0b')]
    case 'STOCHRSI': {
      const { k, d } = stochRsiValues(cl, 14, 14, 3, 3)
      return [mkLine('K', k, '#3b82f6'), mkLine('D', d, '#f59e0b')]
    }
    case 'CMO':
      return [mkLine('CMO', cmoValues(cl, 14), '#10b981')]
    case 'TRIX': {
      const trix = trixValues(cl, 15)
      return [mkLine('TRIX', trix, '#8b5cf6'), mkLine('信号', emaLeadingNull(trix, 9), '#f59e0b')]
    }
    case 'ROC':
      return [mkLine('ROC', rocValues(cl, 12), '#3b82f6')]
    case 'COPPOCK':
      return [mkLine('Coppock', coppockValues(cl), '#ef4444')]
    case 'SMI': {
      const { smi: s, signal } = smiValues(hi, lo, cl)
      return [mkLine('SMI', s, '#8b5cf6'), mkLine('信号', signal, '#f59e0b')]
    }
    case 'AO':
      return [mkLine('AO', aoValues(hi, lo), '#3b82f6')]
    case 'OBV':
      return [mkLine('OBV', obvValues(cl, vl), '#3b82f6')]
    case 'MFI':
      return [mkLine('MFI', mfiValues(hi, lo, cl, vl, 14), '#8b5cf6')]
    case 'CMF':
      return [mkLine('CMF', cmfValues(hi, lo, cl, vl, 20), '#10b981')]
    case 'AD':
      return [mkLine('A/D', adValues(hi, lo, cl, vl), '#f59e0b')]
    case 'FI':
      return [mkLine('FI', forceIndexValues(cl, vl, 13), '#ef4444')]
    case 'CHAIKINOSC':
      return [mkLine('Chaikin', chaikinOscValues(hi, lo, cl, vl, 3, 10), '#8b5cf6')]
    case 'ATR':
      return [mkLine('ATR', atrValues(hi, lo, cl, 14), '#ef4444')]
    case 'TTM': {
      const { squeeze, momentum } = ttmSqueezeValues(hi, lo, cl)
      return [
        { name: 'TTM动量', type: 'bar', data: momentum, xAxisIndex: gridIdx, yAxisIndex: gridIdx,
          itemStyle: { color: (v) => v >= 0 ? c.rise : c.fall } },
      ]
    }
    case 'AVGAMP':
      return [mkLine('均幅', avgAmpSeries(cl, hi, lo), '#3b82f6')]
    case 'MASSINDEX':
      return [mkLine('Mass', massIndexValues(hi, lo), '#f59e0b')]
    case 'ULCER':
      return [mkLine('Ulcer', ulcerIndexValues(cl), '#ef4444')]
    case 'SATS': {
      const { direction } = satsValues(hi, lo, cl, vl)
      return [mkLine('SATS', direction, '#8b5cf6')]
    }
    case 'ADX': {
      const { adx, diP, diM } = adxValues(hi, lo, cl, 14)
      return [mkLine('ADX', adx, '#ef4444'), mkLine('+DI', diP, '#10b981'), mkLine('-DI', diM, '#3b82f6')]
    }
    case 'AROON': {
      const { up, down } = aroonValues(hi, lo, 25)
      return [mkLine('Aroon上', up, '#10b981'), mkLine('Aroon下', down, '#ef4444')]
    }
    case 'CHOP':
      return [mkLine('CHOP', chopValues(hi, lo, cl, 14), '#10b981')]
    case 'ELDERRAY': {
      const { bullPower, bearPower } = elderRayValues(hi, lo, cl, 13)
      return [mkLine('多头力', bullPower, '#ef4444'), mkLine('空头力', bearPower, '#10b981')]
    }
    case 'SIGNALRATIO':
      return buildSignalRatioSeries(gridIdx, c)
  }
  return []
}

function buildOption() {
  const n = norm.value
  const rise = cssColor('--m-color-rise', '#d03050')
  const fall = cssColor('--m-color-fall', '#18a058')
  const textColor = cssColor('--m-text-secondary', '#666')
  const gridColor = cssColor('--m-divider-color', '#eee')
  const c = { rise, fall }

  const subs = activeSubs.value
  const showVol = props.indicators.includes('VOL')
  // 副图总数（含成交量）决定垂直布局
  const subCount = subs.length + (showVol ? 1 : 0)

  // 自上而下累加布局：每个 grid 用绝对 { top, height }，
  // slider 用 top 锚定在所有图下方，物理上不会重叠。
  const totalH = chartHeight.value
  const sliderTop = totalH - ZOOM_BOTTOM - ZOOM_H   // slider 顶边位置
  let cursor = PAD_TOP                               // 当前 grid 顶边游标

  const grids = []
  const xAxes = []
  const yAxes = []
  const series = []
  const allXIdx = [0]

  // 主图 grid[0]：填满副图区之上的剩余空间，底边停在 sliderTop - ZOOM_GAP
  const subBlock = subCount * (SUB_H + SUB_GAP)
  const mainHeight = Math.max(MAIN_MIN, sliderTop - ZOOM_GAP - cursor - subBlock)
  grids.push({ left: 8, right: 52, top: cursor, height: mainHeight })
  cursor += mainHeight
  xAxes.push({
    type: 'category', data: n.category, boundaryGap: false,
    axisLine: { onZero: false, lineStyle: { color: gridColor } },
    axisLabel: { fontSize: 9, color: textColor },
    splitLine: { show: false }, min: 'dataMin', max: 'dataMax',
    axisPointer: { z: 100 },
  })
  yAxes.push({
    scale: true, position: 'right',
    axisLabel: { fontSize: 9, color: textColor, margin: 4 },
    axisLine: { show: false }, axisTick: { show: false },
    splitLine: { lineStyle: { color: gridColor, type: 'dashed' } },
  })

  // 日K + 主图叠加指标
  series.push({
    name: '日K', type: 'candlestick', data: n.values, xAxisIndex: 0, yAxisIndex: 0,
    itemStyle: { color: rise, color0: fall, borderColor: rise, borderColor0: fall },
    markPoint: {
      symbol: 'circle', symbolSize: 8,
      data: [
        {
          name: '最高', type: 'max', valueDim: 'highest', itemStyle: { color: rise },
          label: { formatter: '高 {c}', position: 'top', fontSize: 9, color: rise, fontWeight: 'bold',
            backgroundColor: 'rgba(255,255,255,0.85)', padding: [1, 3], borderRadius: 2 },
        },
        {
          name: '最低', type: 'min', valueDim: 'lowest', itemStyle: { color: fall },
          label: { formatter: '低 {c}', position: 'bottom', fontSize: 9, color: fall, fontWeight: 'bold',
            backgroundColor: 'rgba(255,255,255,0.85)', padding: [1, 3], borderRadius: 2 },
        },
      ],
    },
  })

  // MA 均线（主图叠加）
  if (props.indicators.includes('MA')) {
    series.push({ name: 'MA5', type: 'line', data: calcMA(5, n.values), xAxisIndex: 0, yAxisIndex: 0, smooth: true, showSymbol: false, lineStyle: { width: 1, opacity: 0.8, color: '#3b82f6' } })
    series.push({ name: 'MA10', type: 'line', data: calcMA(10, n.values), xAxisIndex: 0, yAxisIndex: 0, smooth: true, showSymbol: false, lineStyle: { width: 1, opacity: 0.8, color: '#22c55e' } })
    series.push({ name: 'MA20', type: 'line', data: calcMA(20, n.values), xAxisIndex: 0, yAxisIndex: 0, smooth: true, showSymbol: false, lineStyle: { width: 1, opacity: 0.8, color: '#f59e0b' } })
    series.push({ name: 'MA30', type: 'line', data: calcMA(30, n.values), xAxisIndex: 0, yAxisIndex: 0, smooth: true, showSymbol: false, lineStyle: { width: 1, opacity: 0.8, color: '#ef4444' } })
  }
  // BOLL（主图叠加）
  if (props.indicators.includes('BOLL')) {
    const { upper, mid, lower } = bollingerBands(closes.value, 20, 2)
    series.push({ name: 'BOLL上', type: 'line', data: upper, xAxisIndex: 0, yAxisIndex: 0, showSymbol: false, lineStyle: { width: 1, opacity: 0.6, color: '#8b5cf6' } })
    series.push({ name: 'BOLL中', type: 'line', data: mid, xAxisIndex: 0, yAxisIndex: 0, showSymbol: false, lineStyle: { width: 1, opacity: 0.6, color: '#8b5cf6', type: 'dashed' } })
    series.push({ name: 'BOLL下', type: 'line', data: lower, xAxisIndex: 0, yAxisIndex: 0, showSymbol: false, lineStyle: { width: 1, opacity: 0.6, color: '#8b5cf6' } })
  }
  // SAR（主图叠加，点状）
  if (props.indicators.includes('SAR')) {
    const { sar } = sarValues(highs.value, lows.value, closes.value, 0.02, 0.2)
    series.push({ name: 'SAR', type: 'scatter', data: sar, xAxisIndex: 0, yAxisIndex: 0, symbolSize: 3, itemStyle: { color: fall } })
  }
  // EMA（主图叠加）
  if (props.indicators.includes('EMA')) {
    series.push({ name: 'EMA12', type: 'line', data: emaFinite(closes.value, 12), xAxisIndex: 0, yAxisIndex: 0, showSymbol: false, lineStyle: { width: 1, opacity: 0.8, color: '#a855f7' } })
    series.push({ name: 'EMA21', type: 'line', data: emaFinite(closes.value, 21), xAxisIndex: 0, yAxisIndex: 0, showSymbol: false, lineStyle: { width: 1, opacity: 0.8, color: '#14b8a6' } })
  }
  // KAMA / Hull / DEMA / TEMA（主图叠加，单线）
  if (props.indicators.includes('KAMA')) {
    series.push({ name: 'KAMA', type: 'line', data: kamaValues(closes.value, 10, 2, 30), xAxisIndex: 0, yAxisIndex: 0, showSymbol: false, lineStyle: { width: 1, opacity: 0.8, color: '#0ea5e9' } })
  }
  if (props.indicators.includes('HULL')) {
    series.push({ name: 'HullMA', type: 'line', data: hullMaValues(closes.value, 9), xAxisIndex: 0, yAxisIndex: 0, showSymbol: false, lineStyle: { width: 1, opacity: 0.8, color: '#f97316' } })
  }
  if (props.indicators.includes('DEMA')) {
    series.push({ name: 'DEMA', type: 'line', data: demaValues(closes.value, 21), xAxisIndex: 0, yAxisIndex: 0, showSymbol: false, lineStyle: { width: 1, opacity: 0.8, color: '#84cc16' } })
  }
  if (props.indicators.includes('TEMA')) {
    series.push({ name: 'TEMA', type: 'line', data: temaValues(closes.value, 21), xAxisIndex: 0, yAxisIndex: 0, showSymbol: false, lineStyle: { width: 1, opacity: 0.8, color: '#ec4899' } })
  }
  // SuperTrend（主图叠加）
  if (props.indicators.includes('SUPERTREND')) {
    const { supertrend } = supertrendValues(highs.value, lows.value, closes.value, 10, 3)
    series.push({ name: 'SuperTrend', type: 'line', data: supertrend, xAxisIndex: 0, yAxisIndex: 0, showSymbol: false, lineStyle: { width: 1, opacity: 0.8, color: '#6366f1' } })
  }
  // Donchian / Keltner / VWAPBands（主图通道）
  if (props.indicators.includes('DONCHIAN')) {
    const { upper, lower } = donchianChannelValues(highs.value, lows.value, 20)
    series.push({ name: 'Don上', type: 'line', data: upper, xAxisIndex: 0, yAxisIndex: 0, showSymbol: false, lineStyle: { width: 1, opacity: 0.5, color: '#64748b' } })
    series.push({ name: 'Don下', type: 'line', data: lower, xAxisIndex: 0, yAxisIndex: 0, showSymbol: false, lineStyle: { width: 1, opacity: 0.5, color: '#64748b' } })
  }
  if (props.indicators.includes('KELTNER')) {
    const { upper, mid, lower } = keltnerChannelValues(highs.value, lows.value, closes.value, 20, 10, 1.5)
    series.push({ name: 'Kelt上', type: 'line', data: upper, xAxisIndex: 0, yAxisIndex: 0, showSymbol: false, lineStyle: { width: 1, opacity: 0.5, color: '#0891b2' } })
    series.push({ name: 'Kelt中', type: 'line', data: mid, xAxisIndex: 0, yAxisIndex: 0, showSymbol: false, lineStyle: { width: 1, opacity: 0.5, color: '#0891b2', type: 'dashed' } })
    series.push({ name: 'Kelt下', type: 'line', data: lower, xAxisIndex: 0, yAxisIndex: 0, showSymbol: false, lineStyle: { width: 1, opacity: 0.5, color: '#0891b2' } })
  }
  // Ichimoku / Alligator（主图多线）
  if (props.indicators.includes('ICHIMOKU')) {
    const { tenkan, kijun, spanA, senkouB } = ichimokuValues(highs.value, lows.value, closes.value)
    series.push({ name: '转换', type: 'line', data: tenkan, xAxisIndex: 0, yAxisIndex: 0, showSymbol: false, lineStyle: { width: 1, opacity: 0.7, color: '#3b82f6' } })
    series.push({ name: '基准', type: 'line', data: kijun, xAxisIndex: 0, yAxisIndex: 0, showSymbol: false, lineStyle: { width: 1, opacity: 0.7, color: '#ef4444' } })
    series.push({ name: '先行A', type: 'line', data: spanA, xAxisIndex: 0, yAxisIndex: 0, showSymbol: false, lineStyle: { width: 1, opacity: 0.5, color: '#10b981' } })
    series.push({ name: '先行B', type: 'line', data: senkouB, xAxisIndex: 0, yAxisIndex: 0, showSymbol: false, lineStyle: { width: 1, opacity: 0.5, color: '#f59e0b' } })
  }
  if (props.indicators.includes('ALLIGATOR')) {
    const { jaw, teeth, lips } = alligatorValues(highs.value, lows.value, closes.value)
    series.push({ name: '颚', type: 'line', data: jaw, xAxisIndex: 0, yAxisIndex: 0, showSymbol: false, lineStyle: { width: 1, opacity: 0.7, color: '#3b82f6' } })
    series.push({ name: '齿', type: 'line', data: teeth, xAxisIndex: 0, yAxisIndex: 0, showSymbol: false, lineStyle: { width: 1, opacity: 0.7, color: '#ef4444' } })
    series.push({ name: '唇', type: 'line', data: lips, xAxisIndex: 0, yAxisIndex: 0, showSymbol: false, lineStyle: { width: 1, opacity: 0.7, color: '#10b981' } })
  }
  // VWAP / VWAPBands（主图）
  if (props.indicators.includes('VWAP')) {
    series.push({ name: 'VWAP', type: 'line', data: vwapValues(highs.value, lows.value, closes.value, vols.value, 20), xAxisIndex: 0, yAxisIndex: 0, showSymbol: false, lineStyle: { width: 1, opacity: 0.8, color: '#8b5cf6' } })
  }
  if (props.indicators.includes('VWAPBANDS')) {
    const { vwap, upper, lower } = vwapBandsValues(highs.value, lows.value, closes.value, vols.value)
    series.push({ name: 'VWAP', type: 'line', data: vwap, xAxisIndex: 0, yAxisIndex: 0, showSymbol: false, lineStyle: { width: 1, opacity: 0.7, color: '#8b5cf6' } })
    series.push({ name: 'VWAP上', type: 'line', data: upper, xAxisIndex: 0, yAxisIndex: 0, showSymbol: false, lineStyle: { width: 1, opacity: 0.4, color: '#8b5cf6' } })
    series.push({ name: 'VWAP下', type: 'line', data: lower, xAxisIndex: 0, yAxisIndex: 0, showSymbol: false, lineStyle: { width: 1, opacity: 0.4, color: '#8b5cf6' } })
  }
  // Pivot（主图支撑压力线）
  if (props.indicators.includes('PIVOT')) {
    const { pp, s1, r1 } = pivotPointsValues(highs.value, lows.value, closes.value)
    series.push({ name: '枢轴', type: 'line', data: pp, xAxisIndex: 0, yAxisIndex: 0, showSymbol: false, lineStyle: { width: 1, opacity: 0.5, color: '#64748b', type: 'dashed' } })
    series.push({ name: '支撑', type: 'line', data: s1, xAxisIndex: 0, yAxisIndex: 0, showSymbol: false, lineStyle: { width: 1, opacity: 0.5, color: '#10b981', type: 'dashed' } })
    series.push({ name: '阻力', type: 'line', data: r1, xAxisIndex: 0, yAxisIndex: 0, showSymbol: false, lineStyle: { width: 1, opacity: 0.5, color: '#ef4444', type: 'dashed' } })
  }
  // ZigZag（主图折线）
  if (props.indicators.includes('ZIGZAG')) {
    const { zigzag } = zigzagValues(highs.value, lows.value, closes.value, 5)
    series.push({ name: 'ZigZag', type: 'line', data: zigzag, xAxisIndex: 0, yAxisIndex: 0, showSymbol: false, lineStyle: { width: 1, opacity: 0.7, color: '#f59e0b' } })
  }

  // 成交量副图
  let gridIdx = 1
  if (showVol) {
    cursor += SUB_GAP
    grids.push({ left: 8, right: 52, top: cursor, height: SUB_H })
    cursor += SUB_H
    xAxes.push({ type: 'category', gridIndex: gridIdx, data: n.category, boundaryGap: false, axisLine: { lineStyle: { color: gridColor } }, axisTick: { show: false }, axisLabel: { show: false }, splitLine: { show: false }, min: 'dataMin', max: 'dataMax' })
    yAxes.push({ scale: true, gridIndex: gridIdx, splitNumber: 2, position: 'right', axisLabel: { show: false }, axisLine: { show: false }, axisTick: { show: false }, splitLine: { show: false } })
    series.push({ name: '成交量', type: 'bar', xAxisIndex: gridIdx, yAxisIndex: gridIdx, data: n.volumes })
    allXIdx.push(gridIdx)
    gridIdx++
  }

  // 指标副图
  for (const code of subs) {
    cursor += SUB_GAP
    grids.push({ left: 8, right: 52, top: cursor, height: SUB_H })
    cursor += SUB_H
    xAxes.push({ type: 'category', gridIndex: gridIdx, data: n.category, boundaryGap: false, axisLine: { lineStyle: { color: gridColor } }, axisTick: { show: false }, axisLabel: { fontSize: 9, color: textColor }, splitLine: { show: false }, min: 'dataMin', max: 'dataMax' })
    yAxes.push({ scale: true, gridIndex: gridIdx, splitNumber: 2, position: 'right', axisLabel: { fontSize: 9, color: textColor }, axisLine: { show: false }, axisTick: { show: false }, splitLine: { show: false } })
    const subSeries = buildSubSeries(code, gridIdx, c)
    series.push(...subSeries)
    allXIdx.push(gridIdx)
    gridIdx++
  }

  return {
    animation: false,
    legend: { show: false },
    tooltip: { showContent: false, trigger: 'axis',
      axisPointer: { type: 'cross', lineStyle: { color: '#888', width: 1, opacity: 0.8 } } },
    axisPointer: { link: [{ xAxisIndex: 'all' }], label: { backgroundColor: '#6b7077' } },
    grid: grids,
    xAxis: xAxes,
    yAxis: yAxes,
    dataZoom: [
      { type: 'inside', xAxisIndex: allXIdx, start: zoomStart.value, end: 100 },
      { show: true, type: 'slider', xAxisIndex: allXIdx, top: sliderTop, height: ZOOM_H,
        start: zoomStart.value, end: 100, handleSize: '120%',
        fillerColor: 'rgba(135,158,222,0.18)', borderColor: gridColor,
        textStyle: { fontSize: 9, color: textColor } },
    ],
    series,
  }
}

async function ensureEcharts() {
  if (echartsLib.value) return echartsLib.value
  const mod = await import('echarts')
  echartsLib.value = mod
  return mod
}

function bindCrosshair() {
  if (!chart) return
  // 光标/手指移动时联动顶部数据条
  chart.on('updateAxisPointer', (params) => {
    const ax = params.axesInfo && params.axesInfo.find(a => a.axisDim === 'x')
    if (ax && typeof ax.value === 'number') {
      hoverIndex.value = ax.value
    }
  })
  chart.on('globalout', () => { hoverIndex.value = null })
}

async function render() {
  const el = chartRef.value
  if (!el || !norm.value.category.length) return
  const echarts = await ensureEcharts()
  if (!chart) {
    chart = echarts.init(el)
    bindCrosshair()
  }
  // canvas 高度可能随副图数变化，先 resize 再 setOption，避免按旧尺寸布局
  chart.resize()
  chart.setOption(buildOption(), true)
}

function resize() { chart && chart.resize() }

watch(() => props.data, () => { nextTick(render) }, { deep: true })
watch(() => props.indicators, () => { nextTick(render) }, { deep: true })

onMounted(() => {
  nextTick(render)
  window.addEventListener('resize', resize)
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', resize)
  if (chart) { chart.dispose(); chart = null }
})
</script>

<template>
  <div class="multi-kline">
    <!-- 十字光标联动数据条：17 字段，松手回到最新 -->
    <div v-if="activeBar" class="ohlc-bar">
      <div class="ohlc-row">
        <span class="ohlc-day">{{ activeBar.day }}<span v-if="!activeBar.isLatest" class="ohlc-tag">悬停</span><span v-else class="ohlc-tag">最新</span></span>
      </div>
      <div class="ohlc-grid">
        <span v-for="f in visibleFields" :key="f.key" class="kv">
          <i>{{ f.label }}</i>
          <b :style="fieldStyle(f)" :class="{ 'c-rise': f.color === 'rise', 'c-fall': f.color === 'fall' }">
            {{ activeBar[f.key] }}
          </b>
        </span>
      </div>
    </div>

    <!-- MA 均线图例（说明四条线 + 当前值，跟随光标联动） -->
    <div v-if="legendItems.length" class="legend-bar">
      <span v-for="it in legendItems" :key="it.label" class="legend-item">
        <span class="legend-dot" :style="{ background: it.color }" />
        <span class="legend-label">{{ it.label }}</span>
        <span class="legend-val" :style="{ color: it.color }">{{ it.value }}</span>
      </span>
    </div>

    <div ref="chartRef" class="kline-canvas" :style="{ height: `${chartHeight}px` }" />
    <div v-if="!data.length" class="chart-empty"><p>📊 暂无K线数据</p></div>
  </div>
</template>

<style scoped>
.multi-kline {
  position: relative;
  width: 100%;
}

.ohlc-bar {
  padding: var(--m-space-xs) var(--m-space-md) var(--m-space-sm);
  background: var(--m-bg-card);
  border-bottom: 1px solid var(--m-divider-color);
}

.ohlc-row {
  margin-bottom: var(--m-space-xs);
}

.ohlc-day {
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
  font-variant-numeric: tabular-nums;
}

.ohlc-tag {
  display: inline-block;
  margin-left: var(--m-space-xs);
  padding: 0 var(--m-space-xs);
  font-size: 10px;
  line-height: 1.5;
  border-radius: var(--m-radius-sm);
  background: var(--m-color-rise-light);
  color: var(--m-color-rise);
}

.ohlc-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: var(--m-space-xs) var(--m-space-sm);
}

.kv {
  display: flex;
  flex-direction: column;
  gap: 1px;
  min-width: 0;
}

.kv i {
  font-style: normal;
  font-size: 10px;
  color: var(--m-text-tertiary);
}

.kv b {
  font-weight: var(--m-font-weight-medium);
  font-size: var(--m-font-xs);
  color: var(--m-text-primary);
  font-variant-numeric: tabular-nums;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.c-rise { color: var(--m-color-rise); }
.c-fall { color: var(--m-color-fall); }

.kline-canvas { width: 100%; }

/* MA 均线图例条 */
.legend-bar {
  display: flex;
  flex-wrap: wrap;
  gap: var(--m-space-sm);
  padding: var(--m-space-xs) var(--m-space-md);
  background: var(--m-bg-card);
  border-bottom: 1px solid var(--m-divider-color);
}

.legend-item {
  display: inline-flex;
  align-items: center;
  gap: 3px;
  font-size: var(--m-font-xs);
}

.legend-dot {
  width: 10px;
  height: 2px;
  border-radius: 1px;
}

.legend-label {
  color: var(--m-text-tertiary);
}

.legend-val {
  font-weight: var(--m-font-weight-medium);
  font-variant-numeric: tabular-nums;
}

.chart-empty {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--m-text-tertiary);
  font-size: var(--m-font-sm);
}
</style>
