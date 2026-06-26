<script setup>

import {CaretDown, CaretUp, RefreshCircleOutline} from "@vicons/ionicons5";
import {NText,useMessage} from "naive-ui";
import {onBeforeUnmount, onMounted, ref} from "vue";
import {GetIndustryMoneyRankSina} from "../api/app";
import KLineChart from "./KLineChart.vue";
import {registerFeed, stopFeed} from "../api/scheduler";

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
const message = useMessage()
const dataList= ref([])
const sort = ref(props.sort)
const fenlei= ref(props.fenlei)

onMounted(()=>{
  sort.value=props.sort
  fenlei.value=props.fenlei
  GetRankData()
  // 轮询交给统一调度器，页面恢复时自动重刷。
  registerFeed("industryMoneyRank." + props.fenlei + "." + props.sort, { fetch: GetRankData, intervalMs: 1000 * 60 })
})
onBeforeUnmount(()=>{
  stopFeed("industryMoneyRank." + props.fenlei + "." + props.sort)
})
function GetRankData(){
  message.loading("正在刷新数据...")
  GetIndustryMoneyRankSina(fenlei.value,sort.value).then(result => {
    if(result.length>0){
      dataList.value = result
    }
  })
}

</script>

<template>
  <!-- 桌面端：宽表格 -->
  <n-table striped size="small">
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

</template>

<style scoped>
</style>
