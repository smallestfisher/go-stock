<script setup lang="ts">
import {onBeforeMount, ref} from 'vue'
import {LongTigerRank} from "../api/app";
import {BrowserOpenURL} from "../api/runtime";
import {ArrowDownOutline} from "@vicons/ionicons5";
import _ from "lodash";
import KLineChart from "./KLineChart.vue";
import MoneyTrend from "./moneyTrend.vue";
import {NButton, NText, useMessage} from "naive-ui";
import {useDevice} from "../composables/useDevice";
import BottomSheet from "./mobile/BottomSheet.vue";
const message = useMessage()
const {isMobile} = useDevice()

// 移动端：点击展开 K 线 / 资金流向的当前标的
const activeKline = ref(null as any)
const klineVisible = ref(false)
const activeMoney = ref(null as any)
const moneyVisible = ref(false)

// 龙虎榜代码统一格式：SECUCODE = "000001.SZ" → "sz000001"
function emCode(secucode: string){
  if(!secucode) return ''
  const parts = String(secucode).split('.')
  if(parts.length < 2) return parts[0]
  return parts[1].toLowerCase() + parts[0]
}

const lhbList=  ref([])
const EXPLANATIONs = ref([])

const today = new Date();
const year = today.getFullYear();
const month = String(today.getMonth() + 1).padStart(2, '0'); // 月份从0开始，需要+1
const day = String(today.getDate()).padStart(2, '0');

// 常见格式：YYYY-MM-DD
const formattedDate = `${year}-${month}-${day}`;

const SearchForm=  ref({
  dateValue:  formattedDate,
  EXPLANATION:null,
})

onBeforeMount(() => {
  longTiger(formattedDate);
})
function longTiger_old(date) {
  if(date) {
    SearchForm.value.dateValue = date
  }
  let loading1=message.loading("正在获取龙虎榜数据...",{
    duration: 0,
  })
  LongTigerRank(date).then(res => {
    lhbList.value = res
    loading1.destroy()
    if (res.length === 0) {
      message.info("暂无数据,请切换日期")
    }
    EXPLANATIONs.value=_.uniqBy(_.map(lhbList.value,function (item){
      return {
        label: item['EXPLANATION'],
        value: item['EXPLANATION'],
      }
    }),'label');
  })
}

function longTiger(date) {
  if (date) {
    SearchForm.value.dateValue = date;
  }

  let loading1 = message.loading("正在获取龙虎榜数据...", {
    duration: 0,
  });

  const fetchDate = (currentDate, retryCount = 0) => {
    if (retryCount > 7) { // 防止无限循环，最多尝试7次
      lhbList.value = [];
      EXPLANATIONs.value = [];
      loading1.destroy();
      message.info("暂无历史数据");
      return;
    }

    LongTigerRank(currentDate).then(res => {
      if (res.length === 0) {
        const previousDate = new Date(currentDate);
        previousDate.setDate(previousDate.getDate() - 1);

        const year = previousDate.getFullYear();
        const month = String(previousDate.getMonth() + 1).padStart(2, '0');
        const day = String(previousDate.getDate()).padStart(2, '0');
        const prevFormattedDate = `${year}-${month}-${day}`;

        message.info(`当前日期 ${currentDate} 暂无数据，尝试查询前一日：${prevFormattedDate}`);

        SearchForm.value.dateValue = prevFormattedDate;
        fetchDate(prevFormattedDate, retryCount + 1); // 递归调用
      } else {
        lhbList.value = res;
        loading1.destroy();
        EXPLANATIONs.value = _.uniqBy(_.map(lhbList.value, function (item) {
          return {
            label: item['EXPLANATION'],
            value: item['EXPLANATION'],
          };
        }), 'label');
      }
    }).catch(err => {
      loading1.destroy();
      message.error("获取数据失败，请重试");
      console.error(err);
    });
  };

  fetchDate(date || formattedDate);
}

function handleEXPLANATION(value, option){
  SearchForm.value.EXPLANATION = value
  if(value){
    LongTigerRank(SearchForm.value.dateValue).then(res => {
      lhbList.value=_.filter(res, function(o) { return o['EXPLANATION']===value; });
      if (res.length === 0) {
        message.info("暂无数据,请切换日期")
      }
    })
  }else{
    longTiger(SearchForm.value.dateValue)
  }
}

function wan(v){ return (Number(v||0)/10000).toFixed(2) }
function yi(v){ return (Number(v||0)/100000000).toFixed(2) }
function trendColor(v){ return Number(v||0) > 0 ? 'error' : 'success' }

function openKline(item){
  activeKline.value = {code: emCode(item.SECUCODE), name: item.SECURITY_NAME_ABBR}
  klineVisible.value = true
}
function openMoney(item){
  activeMoney.value = {code: emCode(item.SECUCODE), name: item.SECURITY_NAME_ABBR}
  moneyVisible.value = true
}
</script>

<template>
  <!-- ===================== 桌面端 ===================== -->
  <template v-if="!isMobile">
    <n-form :model="SearchForm" >
      <n-grid :cols="24" :x-gap="24">
        <n-form-item-gi  :span="4" label="日期" path="dateValue" label-placement="left">
          <n-date-picker   v-model:formatted-value="SearchForm.dateValue"
                           value-format="yyyy-MM-dd"  type="date"  :on-update:value="(v,v2)=>longTiger(v2)"/>

        </n-form-item-gi>
        <n-form-item-gi :span="8" label="上榜原因" path="EXPLANATION" label-placement="left">
          <n-select  clearable placeholder="上榜原因过滤" v-model:value="SearchForm.EXPLANATION" :options="EXPLANATIONs" :on-update:value="handleEXPLANATION"/>
        </n-form-item-gi>
        <n-form-item-gi :span="10" label=""  label-placement="left">
          <n-text type="error">*当天的龙虎榜数据通常在收盘结束后一小时左右更新</n-text>
        </n-form-item-gi>
      </n-grid>
    </n-form>
    <n-table :single-line="false" striped>
      <n-thead>
        <n-tr>
          <n-th>代码</n-th>
          <n-th width="60px">名称</n-th>
          <n-th>收盘价</n-th>
          <n-th width="60px">涨跌幅</n-th>
          <n-th>龙虎榜净买额(万)</n-th>
          <n-th>龙虎榜买入额(万)</n-th>
          <n-th>龙虎榜卖出额(万)</n-th>
          <n-th>龙虎榜成交额(万)</n-th>
          <n-th width="60px"  data-field="TURNOVERRATE">换手率<n-icon :component="ArrowDownOutline" /></n-th>
          <n-th>流通市值(亿)</n-th>
          <n-th>上榜原因</n-th>
        </n-tr>
      </n-thead>
      <n-tbody>
        <n-tr v-for="(item, index) in lhbList" :key="index">
          <n-td>
            <n-tag :bordered=false type="info">{{ item.SECUCODE.split('.')[1].toLowerCase()+item.SECUCODE.split('.')[0] }}</n-tag>
          </n-td>
          <n-td>
            <n-popover trigger="hover" placement="right">
              <template #trigger>
                <n-button tag="a"  text :type="item.CHANGE_RATE>0?'error':'success'" :bordered=false >{{ item.SECURITY_NAME_ABBR }}</n-button>
              </template>
              <k-line-chart style="width: 800px" :code="item.SECUCODE.split('.')[1].toLowerCase()+item.SECUCODE.split('.')[0]" :chart-height="500" :stockName="item.SECURITY_NAME_ABBR" :k-days="20" :dark-theme="true"></k-line-chart>
            </n-popover>
          </n-td>
          <n-td>
            <n-text :type="item.CHANGE_RATE>0?'error':'success'">{{ item.CLOSE_PRICE }}</n-text>
          </n-td>
          <n-td>
            <n-text :type="item.CHANGE_RATE>0?'error':'success'">{{ (item.CHANGE_RATE).toFixed(2) }}%</n-text>
          </n-td>
          <n-td>
            <n-popover trigger="hover" placement="right">
              <template #trigger>
                <n-button tag="a"  text :type="item.BILLBOARD_NET_AMT>0?'error':'success'" :bordered="false">{{ (item.BILLBOARD_NET_AMT/10000).toFixed(2) }}</n-button>
              </template>
              <money-trend :code="item.SECUCODE.split('.')[1].toLowerCase()+item.SECUCODE.split('.')[0]" :name="item.SECURITY_NAME_ABBR" :days="360" :dark-theme="true" :chart-height="500" style="width: 800px"></money-trend>
            </n-popover>
          </n-td>
          <n-td>
            <n-text :type="'error'">{{ (item.BILLBOARD_BUY_AMT/10000).toFixed(2) }}</n-text>
          </n-td>
          <n-td>
            <n-text :type="'success'">{{ (item.BILLBOARD_SELL_AMT/10000).toFixed(2) }}</n-text>
          </n-td>
          <n-td>
            <n-text :type="'info'">{{ (item.BILLBOARD_DEAL_AMT/10000).toFixed(2) }}</n-text>
          </n-td>
          <n-td>
            <n-text :type="'info'">{{ (item.TURNOVERRATE).toFixed(2) }}%</n-text>
          </n-td>
          <n-td>
            <n-text :type="'info'">{{ (item.FREE_MARKET_CAP/100000000).toFixed(2) }}</n-text>
          </n-td>
          <n-td>
            <n-text :type="'info'">{{ item.EXPLANATION }}</n-text>
          </n-td>
        </n-tr>
      </n-tbody>
    </n-table>
  </template>

  <!-- ===================== 移动端 ===================== -->
  <div v-else class="lhb-mobile">
    <!-- 顶部筛选：日期 + 上榜原因，竖排 -->
    <div class="lhb-mobile__filter">
      <n-date-picker
          v-model:formatted-value="SearchForm.dateValue"
          value-format="yyyy-MM-dd"
          type="date"
          size="small"
          :on-update:value="(v,v2)=>longTiger(v2)"
      />
      <n-select
          clearable
          size="small"
          placeholder="上榜原因过滤"
          v-model:value="SearchForm.EXPLANATION"
          :options="EXPLANATIONs"
          :on-update:value="handleEXPLANATION"
      />
    </div>
    <div class="lhb-mobile__hint">点击名称看 K 线 · 点击净买额看资金</div>

    <article
        v-for="(item, index) in lhbList"
        :key="index"
        class="lhb-card"
        :class="'lhb-card--' + trendColor(item.CHANGE_RATE)"
    >
      <div class="lhb-card__head">
        <div class="lhb-card__title">
          <button type="button" class="lhb-card__name" :class="'text-' + trendColor(item.CHANGE_RATE)" @click="openKline(item)">
            {{ item.SECURITY_NAME_ABBR }}
          </button>
          <n-tag size="tiny" :bordered="false">{{ emCode(item.SECUCODE) }}</n-tag>
        </div>
        <span class="lhb-card__change" :class="'lhb-card__change--' + trendColor(item.CHANGE_RATE)">
          {{ item.CHANGE_RATE.toFixed(2) }}%
        </span>
      </div>

      <!-- 收盘价 + 净买额（可点） -->
      <div class="lhb-card__prices">
        <div class="lhb-price">
          <span class="lhb-price__label">收盘价</span>
          <span class="lhb-price__value" :class="'text-' + trendColor(item.CHANGE_RATE)">{{ item.CLOSE_PRICE }}</span>
        </div>
        <button type="button" class="lhb-price lhb-price--click" @click="openMoney(item)">
          <span class="lhb-price__label">净买额/万</span>
          <span class="lhb-price__value" :class="'text-' + trendColor(item.BILLBOARD_NET_AMT)">{{ wan(item.BILLBOARD_NET_AMT) }} ›</span>
        </button>
      </div>

      <!-- 龙虎榜买卖成交 -->
      <div class="lhb-card__grid">
        <div class="lhb-cell">
          <span class="lhb-cell__label">买入/万</span>
          <n-text type="error" class="lhb-cell__value">{{ wan(item.BILLBOARD_BUY_AMT) }}</n-text>
        </div>
        <div class="lhb-cell">
          <span class="lhb-cell__label">卖出/万</span>
          <n-text type="success" class="lhb-cell__value">{{ wan(item.BILLBOARD_SELL_AMT) }}</n-text>
        </div>
        <div class="lhb-cell">
          <span class="lhb-cell__label">成交/万</span>
          <n-text type="info" class="lhb-cell__value">{{ wan(item.BILLBOARD_DEAL_AMT) }}</n-text>
        </div>
        <div class="lhb-cell">
          <span class="lhb-cell__label">换手率</span>
          <n-text type="info" class="lhb-cell__value">{{ item.TURNOVERRATE.toFixed(2) }}%</n-text>
        </div>
        <div class="lhb-cell">
          <span class="lhb-cell__label">流通市值/亿</span>
          <n-text type="info" class="lhb-cell__value">{{ yi(item.FREE_MARKET_CAP) }}</n-text>
        </div>
      </div>

      <div class="lhb-card__reason">{{ item.EXPLANATION }}</div>
    </article>
    <n-empty v-if="!lhbList.length" description="暂无龙虎榜数据，请切换日期" style="padding: 40px 0" />

    <!-- K 线抽屉 -->
    <BottomSheet :show="klineVisible" :title="activeKline ? activeKline.name : ''" height="76vh" @update:show="(v) => klineVisible = v">
      <div v-if="activeKline" class="lhb-kline-wrap">
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

    <!-- 资金流向抽屉 -->
    <BottomSheet :show="moneyVisible" :title="activeMoney ? (activeMoney.name + ' · 资金流向') : ''" height="76vh" @update:show="(v) => moneyVisible = v">
      <div v-if="activeMoney" class="lhb-kline-wrap">
        <money-trend
            :key="activeMoney.code"
            :code="activeMoney.code"
            :name="activeMoney.name"
            :days="360"
            :dark-theme="true"
            :chart-height="420"
            style="width: 100%"
        />
      </div>
    </BottomSheet>
  </div>
</template>

<style scoped>
/* ============ 移动端龙虎榜卡片 ============ */
.lhb-mobile {
  padding: 6px 4px calc(var(--safe-bottom) + 8px);
}

.lhb-mobile__filter {
  display: grid;
  gap: 8px;
  padding: 0 8px 10px;
}

.lhb-mobile__hint {
  color: var(--n-text-color-3, #999);
  font-size: 12px;
  padding: 0 8px 10px;
}

.lhb-card {
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

.lhb-card--error {
  border-left-color: #d03050;
}

.lhb-card--success {
  border-left-color: #18a058;
}

.lhb-card__head {
  align-items: center;
  display: flex;
  justify-content: space-between;
}

.lhb-card__title {
  align-items: center;
  display: flex;
  gap: 8px;
  min-width: 0;
}

.lhb-card__name {
  appearance: none;
  background: transparent;
  border: 0;
  color: inherit;
  font: inherit;
  font-size: 16px;
  font-weight: 700;
  padding: 0;
}

.lhb-card__change {
  border-radius: 4px;
  color: #fff;
  font-size: 13px;
  font-weight: 700;
  min-width: 64px;
  padding: 2px 8px;
  text-align: center;
}

.lhb-card__change--error {
  background: #d03050;
}

.lhb-card__change--success {
  background: #18a058;
}

.lhb-card__prices {
  display: grid;
  gap: 8px;
  grid-template-columns: 1fr 1fr;
}

.lhb-price {
  align-items: flex-start;
  background: var(--n-color-target, #f5f7fa);
  border-radius: 8px;
  display: flex;
  flex-direction: column;
  gap: 2px;
  padding: 7px 9px;
}

.lhb-price--click {
  appearance: none;
  border: 0;
  color: inherit;
  font: inherit;
  text-align: left;
  width: 100%;
}

.lhb-price--click:active {
  background: var(--n-color-hover, #eef2f7);
}

.lhb-price__label {
  color: var(--n-text-color-3, #98a2b3);
  font-size: 11px;
}

.lhb-price__value {
  font-size: 16px;
  font-weight: 700;
}

.lhb-card__grid {
  border-top: 1px dashed var(--n-divider-color, #eef0f4);
  display: grid;
  gap: 6px 12px;
  grid-template-columns: 1fr 1fr 1fr;
  padding-top: 8px;
}

.lhb-cell {
  align-items: center;
  display: flex;
  flex-direction: column;
  gap: 1px;
}

.lhb-cell__label {
  color: var(--n-text-color-3, #98a2b3);
  font-size: 11px;
}

.lhb-cell__value {
  font-size: 13px;
  font-weight: 600;
}

.lhb-card__reason {
  background: var(--n-color-target, #f5f7fa);
  border-radius: 6px;
  color: var(--n-text-color-2, #555);
  font-size: 12px;
  line-height: 1.5;
  padding: 6px 8px;
}

/* 涨跌色 */
.text-error {
  color: #d03050;
}

.text-success {
  color: #18a058;
}

.lhb-kline-wrap {
  padding: 4px 8px 8px;
}
</style>
