<script setup>

import {CaretDown, CaretUp, RefreshCircleOutline} from "@vicons/ionicons5";
import {NText,useMessage} from "naive-ui";
import {computed, onBeforeUnmount, onMounted, ref} from "vue";
import {GetMoneyRankSina} from "../api/app";
import KLineChart from "./KLineChart.vue";
import BottomSheet from "./mobile/BottomSheet.vue";
import {useDevice} from "../composables/useDevice";
import {registerFeed, stopFeed} from "../api/scheduler";

const props = defineProps({
  headerTitle: {
    type: String,
    default: '净流入额排名'
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
// 移动端点击查看 K 线的当前标的
const activeKline = ref(null)
const klineVisible = ref(false)

onMounted(()=>{
  sort.value=props.sort
  GetMoneyRankSinaData()
  // 轮询交给统一调度器，页面恢复时自动重刷。
  registerFeed("rankTable." + props.sort, { fetch: GetMoneyRankSinaData, intervalMs: 1000 * 60 })
})
onBeforeUnmount(()=>{
  stopFeed("rankTable." + props.sort)
})
function GetMoneyRankSinaData(){
  message.loading("正在刷新数据...")
  GetMoneyRankSina(sort.value).then(result => {
    if(result.length>0){
      dataList.value = result
    }
  })
}

// 金额统一换算成“万”
function wan(v){
  return (Number(v||0)/10000).toFixed(2)
}
// 比率统一换算成百分号
function pct(v){
  return (Number(v||0)*100).toFixed(2)
}
// 涨跌着色：>0 红(error)，否则绿(success)
function trendColor(v){
  return Number(v||0) > 0 ? 'error' : 'success'
}

// 按排序字段决定明细分组：r0_* 主力、r3_* 散群、其余整体资金
const flowGroup = computed(() => {
  const s = sort.value || ''
  if (s.startsWith('r0')) return 'main'
  if (s.startsWith('r3')) return 'retail'
  return 'all'
})

// 各分组在卡片底部展示的明细指标
const detailCells = computed(() => (item) => {
  const g = flowGroup.value
  if (g === 'main') {
    return [
      {label: '主力流入/万', value: wan(item.r0_in), type: 'error'},
      {label: '主力流出/万', value: wan(item.r0_out), type: 'success'},
      {label: '主力净流入/万', value: wan(item.r0_net), type: trendColor(item.r0_net)},
      {label: '主力净流入率', value: pct(item.r0_ratio) + '%', type: trendColor(item.r0_ratio)},
    ]
  }
  if (g === 'retail') {
    return [
      {label: '散户流入/万', value: wan(item.r3_in), type: 'error'},
      {label: '散户流出/万', value: wan(item.r3_out), type: 'success'},
      {label: '散户净流入/万', value: wan(item.r3_net), type: trendColor(item.r3_net)},
      {label: '散户净流入率', value: pct(item.r3_ratio) + '%', type: trendColor(item.r3_ratio)},
    ]
  }
  return [
    {label: '流入/万', value: wan(item.inamount), type: 'error'},
    {label: '流出/万', value: wan(item.outamount), type: 'success'},
    {label: '成交额/万', value: wan(item.amount), type: 'info'},
    {label: '换手率', value: pct(item.turnover) + '%', type: Number(item.turnover||0) > 500 ? 'error' : 'info'},
  ]
})

function openKline(item){
  activeKline.value = item
  klineVisible.value = true
}
</script>

<template>
  <!-- 桌面端：宽表格 -->
  <n-table v-if="!isMobile" striped size="small">
    <n-thead>
      <n-tr>
        <n-th>代码</n-th>
        <n-th>名称</n-th>
        <n-th>最新价</n-th>
        <n-th>涨跌幅</n-th>
        <n-th>换手率</n-th>
        <n-th>成交额/万</n-th>
        <n-th>流出资金/万</n-th>
        <n-th>流入资金/万</n-th>
        <n-th>净流入/万</n-th>
        <n-th>净流入率</n-th>
        <n-th v-if="sort === 'r0_net'||sort==='r0_out'">主力流出/万</n-th>
        <n-th v-if="sort === 'r0_net'">主力流入/万</n-th>
        <n-th v-if="sort === 'r0_net'">主力净流入/万</n-th>
        <n-th >主力净流入率</n-th>
        <n-th v-if="sort === 'r3_net'||sort==='r3_out'">散户流出/万</n-th>
        <n-th v-if="sort === 'r3_net'">散户流入/万</n-th>
        <n-th v-if="sort === 'r3_net'">散户净流入/万</n-th>
        <n-th >散户净流入率</n-th>
      </n-tr>
    </n-thead>
    <n-tbody>
      <n-tr v-for="item in dataList" :key="item.symbol">
        <n-td><n-tag :bordered=false type="info">{{ item.symbol }}</n-tag></n-td>
        <n-td>
          <n-popover trigger="hover" placement="right">
            <template #trigger>
              <n-button tag="a"  text :type="item.changeratio>0?'error':'success'" :bordered=false >{{ item.name }}</n-button>
            </template>
            <k-line-chart style="width: 800px" :code="item.symbol" :chart-height="500" :stockName="item.name" :k-days="20" :dark-theme="true"></k-line-chart>
          </n-popover>
        </n-td>
        <n-td><n-text :type="item.changeratio>0?'error':'success'">{{item.trade}}</n-text></n-td>
        <n-td><n-text :type="item.changeratio>0?'error':'success'">{{(item.changeratio*100).toFixed(2)}}%</n-text></n-td>
        <n-td><n-text :type="item.turnover>500?'error':'info'">{{(item.turnover/100).toFixed(2)}}%</n-text></n-td>
        <n-td><n-text type="info">{{(item.amount/10000).toFixed(2)}}</n-text></n-td>
        <n-td><n-text type="info"> {{(item.outamount/10000).toFixed(2)}}</n-text></n-td>
        <n-td><n-text type="info"> {{(item.inamount/10000).toFixed(2)}}</n-text></n-td>
        <n-td><n-text type="info"> {{(item.netamount/10000).toFixed(2)}}</n-text></n-td>
        <n-td><n-text :type="item.ratioamount>0?'error':'success'"> {{(item.ratioamount*100).toFixed(2)}}%</n-text></n-td>
        <n-td v-if="sort === 'r0_net'||sort==='r0_out'"><n-text  type="success"> {{(item.r0_out/10000).toFixed(2)}}</n-text></n-td>
        <n-td v-if="sort === 'r0_net'"><n-text  type="error"> {{(item.r0_in/10000).toFixed(2)}}</n-text></n-td>
        <n-td v-if="sort === 'r0_net'"><n-text :type="item.r0_net>0?'error':'success'"> {{(item.r0_net/10000).toFixed(2)}}</n-text></n-td>
        <n-td ><n-text :type="item.r0_ratio>0?'error':'success'"> {{(item.r0_ratio*100).toFixed(2)}}%</n-text></n-td>
        <n-td v-if="sort === 'r3_net'||sort==='r3_out'"><n-text  type="success"> {{(item.r3_out/10000).toFixed(2)}}</n-text></n-td>
        <n-td v-if="sort === 'r3_net'"><n-text  type="error"> {{(item.r3_in/10000).toFixed(2)}}</n-text></n-td>
        <n-td v-if="sort === 'r3_net'"><n-text :type="item.r3_net>0?'error':'success'"> {{(item.r3_net/10000).toFixed(2)}}</n-text></n-td>
        <n-td ><n-text :type="item.r3_ratio>0?'error':'success'"> {{(item.r3_ratio*100).toFixed(2)}}%</n-text></n-td>
      </n-tr>
    </n-tbody>
  </n-table>

  <!-- 移动端：资金流卡片列表，点击查看 K 线 -->
  <div v-else class="rank-mobile">
    <div class="rank-mobile__hint">点击卡片查看 K 线 · {{ headerTitle }}</div>
    <button
        v-for="item in dataList"
        :key="item.symbol"
        type="button"
        class="rank-card"
        :class="'rank-card--' + trendColor(item.netamount)"
        @click="openKline(item)"
    >
      <div class="rank-card__head">
        <span class="rank-card__name" :class="'text-' + trendColor(item.changeratio)">{{ item.name }}</span>
        <n-tag size="tiny" :bordered="false">{{ item.symbol }}</n-tag>
        <span class="rank-card__quote" :class="'text-' + trendColor(item.changeratio)">
          <span class="rank-card__price">{{ item.trade }}</span>
          <span class="rank-card__change"
                :class="trendColor(item.changeratio) === 'error' ? 'is-up' : 'is-down'">{{ pct(item.changeratio) }}%</span>
        </span>
      </div>

      <!-- 核心资金指标 -->
      <div class="rank-card__core">
        <div class="rank-metric rank-metric--lead">
          <span class="rank-metric__label">净流入/万</span>
          <span class="rank-metric__value" :class="'text-' + trendColor(item.netamount)">{{ wan(item.netamount) }}</span>
        </div>
        <div class="rank-metric">
          <span class="rank-metric__label">净流入率</span>
          <span class="rank-metric__value" :class="'text-' + trendColor(item.ratioamount)">{{ pct(item.ratioamount) }}%</span>
        </div>
        <div class="rank-metric">
          <span class="rank-metric__label">换手率</span>
          <span class="rank-metric__value text-info">{{ pct(item.turnover) }}%</span>
        </div>
      </div>

      <!-- 分组明细 -->
      <div class="rank-card__detail">
        <div v-for="(cell, ci) in detailCells(item)" :key="ci" class="rank-sub">
          <span class="rank-sub__label">{{ cell.label }}</span>
          <n-text :type="cell.type" class="rank-sub__value">{{ cell.value }}</n-text>
        </div>
      </div>
    </button>
    <n-empty v-if="!dataList.length" description="暂无资金流向数据" style="padding: 40px 0" />

    <!-- K 线抽屉：复用桌面同款 KLineChart，传入相同的 code 格式 -->
    <BottomSheet :show="klineVisible" :title="activeKline ? activeKline.name : ''" height="78vh" @update:show="(v) => klineVisible = v">
      <div v-if="activeKline" class="rank-kline-wrap">
        <k-line-chart
            :key="activeKline.symbol"
            :code="activeKline.symbol"
            :chart-height="420"
            :stockName="activeKline.name"
            :k-days="20"
            :dark-theme="true"
        />
      </div>
    </BottomSheet>
  </div>
</template>

<style scoped>
/* ============ 移动端资金流卡片 ============ */
.rank-mobile {
  padding: 6px 4px calc(var(--safe-bottom) + 8px);
}

.rank-mobile__hint {
  color: var(--n-text-color-3, #999);
  font-size: 12px;
  padding: 2px 8px 8px;
}

.rank-card {
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

.rank-card:active {
  background: var(--n-color-hover, #f8fafc);
}

.rank-card--error {
  border-left-color: #d03050;
}

.rank-card--success {
  border-left-color: #18a058;
}

.rank-card__head {
  align-items: center;
  display: flex;
  gap: 8px;
}

.rank-card__name {
  flex: 0 1 auto;
  font-size: 16px;
  font-weight: 700;
  max-width: 46%;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.rank-card__quote {
  align-items: baseline;
  display: flex;
  flex: 1 1 auto;
  gap: 8px;
  justify-content: flex-end;
}

.rank-card__price {
  font-size: 17px;
  font-weight: 800;
}

.rank-card__change {
  border-radius: 4px;
  color: #fff;
  font-size: 13px;
  font-weight: 700;
  min-width: 62px;
  padding: 2px 7px;
  text-align: center;
}

.rank-card__change.is-up {
  background: #d03050;
}

.rank-card__change.is-down {
  background: #18a058;
}

.rank-card__core {
  display: grid;
  gap: 8px;
  grid-template-columns: 1fr 1fr 1fr;
}

.rank-metric {
  align-items: flex-start;
  background: var(--n-color-target, #f5f7fa);
  border-radius: 8px;
  display: flex;
  flex-direction: column;
  gap: 2px;
  padding: 7px 9px;
}

.rank-metric--lead {
  background: rgba(208, 48, 80, 0.06);
}

.rank-card--success .rank-metric--lead {
  background: rgba(24, 160, 88, 0.08);
}

.rank-metric__label {
  color: var(--n-text-color-3, #98a2b3);
  font-size: 11px;
}

.rank-metric__value {
  font-size: 15px;
  font-weight: 700;
}

.rank-card__detail {
  border-top: 1px dashed var(--n-divider-color, #eef0f4);
  display: grid;
  gap: 6px 12px;
  grid-template-columns: 1fr 1fr;
  padding-top: 8px;
}

.rank-sub {
  align-items: center;
  display: flex;
  font-size: 12px;
  justify-content: space-between;
}

.rank-sub__label {
  color: var(--n-text-color-3, #98a2b3);
}

.rank-sub__value {
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

.rank-kline-wrap {
  padding: 4px 8px 8px;
}
</style>
