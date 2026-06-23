<script setup>

import {CaretDown, CaretUp, RefreshCircleOutline} from "@vicons/ionicons5";
import {NText,useMessage} from "naive-ui";
import {computed, onBeforeUnmount, onMounted, ref} from "vue";
import {GetIndustryMoneyRankSina} from "../api/app";
import KLineChart from "./KLineChart.vue";
import BottomSheet from "./mobile/BottomSheet.vue";
import {useDevice} from "../composables/useDevice";

const props = defineProps({
  headerTitle: {
    type: String,
    default: '行业资金排名(净流入)'
  },
  fenlei: {
    type: String,
    default: '0'
  },
  sort: {
    type: String,
    default: 'netamount'
  },
})
const {isMobile} = useDevice()
const message = useMessage()
const dataList= ref([])
const sort = ref(props.sort)
const fenlei= ref(props.fenlei)
const activeKline = ref(null)
const klineVisible = ref(false)

const interval = ref(null)
onMounted(()=>{
  sort.value=props.sort
  fenlei.value=props.fenlei
  GetRankData()
  interval.value=setInterval(()=>{
    GetRankData()
  },1000*60)
})
onBeforeUnmount(()=>{
  clearInterval(interval)
})
function GetRankData(){
  message.loading("正在刷新数据...")
  GetIndustryMoneyRankSina(fenlei.value,sort.value).then(result => {
    if(result.length>0){
      dataList.value = result
    }
  })
}

function wan(v){
  return (Number(v||0)/10000).toFixed(2)
}
function pct(v){
  return (Number(v||0)*100).toFixed(2)
}
function trendColor(v){
  return Number(v||0) > 0 ? 'error' : 'success'
}

function openKline(item){
  // 领涨股 K 线
  activeKline.value = {symbol: item.ts_symbol, name: item.ts_name}
  klineVisible.value = true
}
</script>

<template>
  <!-- 桌面端：宽表格 -->
  <n-table v-if="!isMobile" striped size="small">
    <n-thead>
      <n-tr>
        <n-th>板块名称</n-th>
        <n-th>涨跌幅</n-th>
        <n-th>流入资金/万</n-th>
        <n-th>流出资金/万</n-th>
        <n-th>净流入/万<n-icon v-if="sort==='0'" :component="CaretDown"/><n-icon  v-if="sort==='1'" :component="CaretUp"/></n-th>
        <n-th>净流入率</n-th>
        <n-th>领涨股</n-th>
        <n-th>涨跌幅</n-th>
        <n-th>最新价</n-th>
        <n-th>净流入率</n-th>
      </n-tr>
    </n-thead>
    <n-tbody>
      <n-tr v-for="item in dataList" :key="item.category">
        <n-td><n-tag :bordered=false type="info">{{item.name}}</n-tag></n-td>
        <n-td> <n-text :type="item.avg_changeratio>0?'error':'success'">{{(item.avg_changeratio*100).toFixed(2)}}%</n-text></n-td>
        <n-td><n-text type="info">{{(item.inamount/10000).toFixed(2)}}</n-text></n-td>
        <n-td><n-text type="info">{{(item.outamount/10000).toFixed(2)}}</n-text></n-td>
        <n-td><n-text :type="item.netamount>0?'error':'success'">{{(item.netamount/10000).toFixed(2)}}</n-text></n-td>
        <n-td><n-text  :type="item.ratioamount>0?'error':'success'">{{(item.ratioamount*100).toFixed(2)}}%</n-text></n-td>
        <n-td>
          <n-popover trigger="hover" placement="right">
            <template #trigger>
              <n-button tag="a"  text :type="item.ts_changeratio>0?'error':'success'" :bordered=false >{{ item.ts_name }}</n-button>
            </template>
            <k-line-chart style="width: 800px" :code="item.ts_symbol" :chart-height="500" :name="item.ts_name" :k-days="20" :dark-theme="true"></k-line-chart>
          </n-popover>
        </n-td>
        <n-td><n-text :type="item.ts_changeratio>0?'error':'success'">{{(item.ts_changeratio*100).toFixed(2)}}%</n-text></n-td>
        <n-td><n-text type="info">{{item.ts_trade}}</n-text></n-td>
        <n-td><n-text :type="item.ts_ratioamount>0?'error':'success'">{{(item.ts_ratioamount*100).toFixed(2)}}%</n-text></n-td>
      </n-tr>
    </n-tbody>
  </n-table>

  <!-- 移动端：板块资金卡片，点击领涨股查看 K 线 -->
  <div v-else class="imr-mobile">
    <div class="imr-mobile__hint">点击领涨股查看 K 线 · {{ headerTitle }}</div>
    <div
        v-for="item in dataList"
        :key="item.category"
        class="imr-card"
        :class="'imr-card--' + trendColor(item.netamount)"
    >
      <div class="imr-card__head">
        <n-tag size="small" :bordered="false" type="info">{{ item.name }}</n-tag>
        <span class="imr-card__change" :class="'imr-card__change--' + trendColor(item.avg_changeratio)">
          {{ pct(item.avg_changeratio) }}%
        </span>
      </div>

      <div class="imr-card__core">
        <div class="imr-metric imr-metric--lead">
          <span class="imr-metric__label">净流入/万</span>
          <span class="imr-metric__value" :class="'text-' + trendColor(item.netamount)">{{ wan(item.netamount) }}</span>
        </div>
        <div class="imr-metric">
          <span class="imr-metric__label">净流入率</span>
          <span class="imr-metric__value" :class="'text-' + trendColor(item.ratioamount)">{{ pct(item.ratioamount) }}%</span>
        </div>
        <div class="imr-metric">
          <span class="imr-metric__label">流入/万</span>
          <span class="imr-metric__value text-info">{{ wan(item.inamount) }}</span>
        </div>
      </div>

      <button v-if="item.ts_name" type="button" class="imr-leader" @click="openKline(item)">
        <span class="imr-leader__label">领涨</span>
        <span class="imr-leader__name" :class="'text-' + trendColor(item.ts_changeratio)">{{ item.ts_name }}</span>
        <span class="imr-leader__price">{{ item.ts_trade }}</span>
        <span class="imr-leader__ratio" :class="'text-' + trendColor(item.ts_changeratio)">{{ pct(item.ts_changeratio) }}%</span>
        <span class="imr-leader__arrow">›</span>
      </button>
    </div>
    <n-empty v-if="!dataList.length" description="暂无资金流向数据" style="padding: 40px 0" />

    <BottomSheet :show="klineVisible" :title="activeKline ? activeKline.name : ''" height="76vh" @update:show="(v) => klineVisible = v">
      <div v-if="activeKline" class="imr-kline-wrap">
        <k-line-chart
            :key="activeKline.symbol"
            style="width: 100%"
            :code="activeKline.symbol"
            :chart-height="420"
            :name="activeKline.name"
            :k-days="20"
            :dark-theme="true"
        />
      </div>
    </BottomSheet>
  </div>
</template>

<style scoped>
/* ============ 移动端板块资金卡片 ============ */
.imr-mobile {
  padding: 6px 4px calc(var(--safe-bottom) + 8px);
}

.imr-mobile__hint {
  color: var(--n-text-color-3, #999);
  font-size: 12px;
  padding: 2px 8px 8px;
}

.imr-card {
  background: var(--n-color, #fff);
  border: 1px solid var(--n-border-color, #edf0f5);
  border-left: 4px solid #d0d5dd;
  border-radius: 10px;
  box-shadow: 0 1px 2px rgba(15, 23, 42, 0.04);
  display: flex;
  flex-direction: column;
  gap: 10px;
  margin-bottom: 8px;
  padding: 11px 12px;
}

.imr-card--error {
  border-left-color: #d03050;
}

.imr-card--success {
  border-left-color: #18a058;
}

.imr-card__head {
  align-items: center;
  display: flex;
  justify-content: space-between;
}

.imr-card__change {
  border-radius: 4px;
  color: #fff;
  font-size: 13px;
  font-weight: 700;
  min-width: 70px;
  padding: 2px 8px;
  text-align: center;
}

.imr-card__change--error {
  background: #d03050;
}

.imr-card__change--success {
  background: #18a058;
}

.imr-card__core {
  display: grid;
  gap: 8px;
  grid-template-columns: 1fr 1fr 1fr;
}

.imr-metric {
  align-items: flex-start;
  background: var(--n-color-target, #f5f7fa);
  border-radius: 8px;
  display: flex;
  flex-direction: column;
  gap: 2px;
  padding: 7px 9px;
}

.imr-metric--lead {
  background: rgba(208, 48, 80, 0.06);
}

.imr-card--success .imr-metric--lead {
  background: rgba(24, 160, 88, 0.08);
}

.imr-metric__label {
  color: var(--n-text-color-3, #98a2b3);
  font-size: 11px;
}

.imr-metric__value {
  font-size: 15px;
  font-weight: 700;
}

/* 领涨股可点击行 */
.imr-leader {
  align-items: center;
  appearance: none;
  background: var(--n-color-target, #f5f7fa);
  border: 0;
  border-radius: 8px;
  color: inherit;
  display: flex;
  font: inherit;
  gap: 8px;
  padding: 8px 10px;
  text-align: left;
  width: 100%;
}

.imr-leader:active {
  background: var(--n-color-hover, #eef2f7);
}

.imr-leader__label {
  background: var(--n-color, #fff);
  border-radius: 4px;
  color: var(--n-text-color-3, #98a2b3);
  flex: 0 0 auto;
  font-size: 11px;
  padding: 1px 6px;
}

.imr-leader__name {
  flex: 0 1 auto;
  font-size: 14px;
  font-weight: 700;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.imr-leader__price {
  color: var(--n-text-color-2, #555);
  flex: 1 1 auto;
  font-size: 13px;
  text-align: right;
}

.imr-leader__ratio {
  font-size: 13px;
  font-weight: 700;
}

.imr-leader__arrow {
  color: var(--n-text-color-3, #c0c4cc);
  font-size: 18px;
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

.imr-kline-wrap {
  padding: 4px 8px 8px;
}
</style>
