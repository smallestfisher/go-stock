<script setup>
import {defineAsyncComponent} from 'vue'

const StockSparkLine = defineAsyncComponent(() => import('../stockSparkLine.vue'))

const props = defineProps({
    result: {type: Object, required: true},
    // 0 = 全部（不显示"移出分组"），>0 = 具体分组
    groupId: {type: Number, default: 0},
    // 全部页不显示迷你走势，分组页显示
    showSparkline: {type: Boolean, default: false},
    openAiEnable: {type: Boolean, default: false},
    groupList: {type: Array, default: () => []},
    // 父组件操作函数集合，避免十几个 emit
    actions: {type: Object, default: () => ({})},
})

function a() {
    return props.actions || {}
}
</script>

<template>
    <n-gi :id="result['股票代码']+'_gi'" style="margin-left: 2px;">
        <n-card class="stock-mobile-card" :data-sort="result.sort" :id="result['股票代码']" :data-code="result['股票代码']" :bordered="true"
                :title="result['股票名称']" :closable="false"
                @close="a().removeMonitor && a().removeMonitor(result['股票代码'],result['股票名称'],result.key)">
            <n-grid class="stock-mobile-card-body" :cols="showSparkline ? 12 : 1" :y-gap="6"
                    @click="a().showLightweightKline && a().showLightweightKline(result['股票代码'],result['股票名称'])">
                <n-gi :span="showSparkline ? 6 : undefined">
                    <n-text :type="result.type">
                        <n-number-animation :duration="1000" :precision="2" :from="result['上次当前价格']"
                                            :to="Number(result['当前价格'])"/>
                        <n-tag size="small" :type="result.type" :bordered="false" v-if="result['盘前盘后']>0">
                            ({{ result['盘前盘后'] }} {{ result['盘前盘后涨跌幅'] }}%)
                        </n-tag>
                    </n-text>
                    <n-text style="padding-left: 10px;" :type="result.type">
                        <n-number-animation :duration="1000" :precision="3" :from="0" :to="result.changePercent"/>
                        %
                    </n-text>&nbsp;
                    <n-text size="small" v-if="result.costVolume>0" :type="result.type">
                        <n-number-animation :duration="1000" :precision="2" :from="0" :to="result.profitAmountToday"/>
                    </n-text>
                </n-gi>
                <n-gi v-if="showSparkline" :span="6">
                    <stock-spark-line :last-price="Number(result['当前价格'])" :open-price="Number(result['昨日收盘价'])"
                                      :stock-code="result['股票代码']" :stock-name="result['股票名称']"></stock-spark-line>
                </n-gi>
            </n-grid>
            <n-grid :cols="2" :y-gap="4" :x-gap="4">
                <n-gi>
                    <n-text :type="'info'">{{ "最高 " + result["今日最高价"] + " " + result.highRate }}%</n-text>
                </n-gi>
                <n-gi>
                    <n-text :type="'info'">{{ "最低 " + result["今日最低价"] + " " + result.lowRate }}%</n-text>
                </n-gi>
                <n-gi>
                    <n-text :type="'info'">{{ "昨收 " + result["昨日收盘价"] }}</n-text>
                </n-gi>
                <n-gi>
                    <n-text :type="'info'">{{ "今开 " + result["今日开盘价"] }}</n-text>
                </n-gi>
            </n-grid>
            <n-collapse accordion v-if="result['买一报价']>0">
                <n-collapse-item title="盘口" name="1" v-if="result['买一报价']>0">
                    <template #header-extra>
                        <n-flex justify="space-between">
                            <n-text :type="'info'">{{ "买一 " + result["买一报价"] + '(' + result["买一申报"] + ")" }}</n-text>
                            <n-text :type="'info'">{{ "卖一 " + result["卖一报价"] + '(' + result["卖一申报"] + ")" }}</n-text>
                        </n-flex>
                    </template>
                    <n-grid :cols="2" :y-gap="4" :x-gap="4">
                        <n-gi v-if="result['买一报价']>0">
                            <n-text :type="'info'">{{ "买一 " + result["买一报价"] + '(' + result["买一申报"] + ")" }}</n-text>
                        </n-gi>
                        <n-gi v-if="result['卖一报价']>0">
                            <n-text :type="'info'">{{ "卖一 " + result["卖一报价"] + '(' + result["卖一申报"] + ")" }}</n-text>
                        </n-gi>

                        <n-gi v-if="result['买二报价']>0">
                            <n-text :type="'info'">{{ "买二 " + result["买二报价"] + '(' + result["买二申报"] + ")" }}</n-text>
                        </n-gi>
                        <n-gi v-if="result['卖二报价']>0">
                            <n-text :type="'info'">{{ "卖二 " + result["卖二报价"] + '(' + result["卖二申报"] + ")" }}</n-text>
                        </n-gi>

                        <n-gi v-if="result['买三报价']>0">
                            <n-text :type="'info'">{{ "买三 " + result["买三报价"] + '(' + result["买三申报"] + ")" }}</n-text>
                        </n-gi>
                        <n-gi v-if="result['卖三报价']>0">
                            <n-text :type="'info'">{{ "卖三 " + result["卖三报价"] + '(' + result["卖三申报"] + ")" }}</n-text>
                        </n-gi>

                        <n-gi v-if="result['买四报价']>0">
                            <n-text :type="'info'">{{ "买四 " + result["买四报价"] + '(' + result["买四申报"] + ")" }}</n-text>
                        </n-gi>
                        <n-gi v-if="result['卖四报价']>0">
                            <n-text :type="'info'">{{ "卖四 " + result["卖四报价"] + '(' + result["卖四申报"] + ")" }}</n-text>
                        </n-gi>

                        <n-gi v-if="result['买五报价']>0">
                            <n-text :type="'info'">{{ "买五 " + result["买五报价"] + '(' + result["买五申报"] + ")" }}</n-text>
                        </n-gi>
                        <n-gi v-if="result['卖五报价']>0">
                            <n-text :type="'info'">{{ "卖五 " + result["卖五报价"] + '(' + result["卖五申报"] + ")" }}</n-text>
                        </n-gi>
                    </n-grid>
                </n-collapse-item>
            </n-collapse>
            <template #header-extra>
                <n-flex class="stock-desktop-card-extra" align="center" :size="4">
                    <n-tag size="small" :bordered="false">{{ result['股票代码'] }}</n-tag>
                    <n-button size="tiny" secondary type="primary"
                              @click="a().removeMonitor && a().removeMonitor(result['股票代码'],result['股票名称'],result.key)">
                        取消关注
                    </n-button>
                    <n-button size="tiny" v-if="openAiEnable" secondary type="warning"
                              @click="a().aiCheckStock && a().aiCheckStock(result['股票名称'],result['股票代码'])">
                        AI分析
                    </n-button>
                    <n-button v-if="groupId > 0" secondary type="error" size="tiny"
                              @click="a().delStockGroup && a().delStockGroup(result['股票代码'],result['股票名称'],groupId)">移出分组
                    </n-button>
                </n-flex>
            </template>
            <template #footer>
                <n-flex vertical :size="8">
                    <n-flex justify="center">
                        <n-text :type="'info'">{{ result["日期"] + " " + result["时间"] }}</n-text>
                        <n-tag size="small" v-if="result.volume>0" :type="result.profitType">{{ result.volume + "股" }}</n-tag>
                        <n-tag size="small" v-if="result.costPrice>0" :type="result.profitType">
                            {{
                                "成本:" + result.costPrice + "*" + result.costVolume + " " + result.profit + "%" + " ( " + result.profitAmount + " ¥ )"
                            }}
                        </n-tag>
                    </n-flex>
                    <n-flex justify="center">
                        <n-button size="tiny" type="primary" secondary
                                  @click="a().showLightweightKline && a().showLightweightKline(result['股票代码'],result['股票名称'])">
                            多周期K线
                        </n-button>
                    </n-flex>
                </n-flex>
            </template>
            <template #action>
                <n-flex class="stock-mobile-primary-actions mobile-only" justify="space-between">
                    <n-button size="small" type="warning" secondary v-if="openAiEnable"
                              @click.stop="a().aiCheckStock && a().aiCheckStock(result['股票名称'],result['股票代码'])">AI</n-button>
                    <n-button size="small" type="error" secondary
                              @click.stop="a().showFenshi && a().showFenshi(result['股票代码'],result['股票名称'],result.changePercent)">分时</n-button>
                    <n-button size="small" type="primary" secondary
                              @click.stop="a().showLightweightKline && a().showLightweightKline(result['股票代码'],result['股票名称'])">K线</n-button>
                    <n-button size="small" type="warning" secondary
                              @click.stop="a().setStock && a().setStock(result['股票代码'],result['股票名称'])">成本</n-button>
                    <n-dropdown class="stock-mobile-more-menu" trigger="click"
                                :options="a().getStockMobileMoreOptions ? a().getStockMobileMoreOptions(result, groupId) : []"
                                @select="(key) => a().handleStockMobileMoreSelect && a().handleStockMobileMoreSelect(key, result, groupId)">
                        <n-button size="small" secondary>更多</n-button>
                    </n-dropdown>
                </n-flex>
                <n-flex class="stock-mobile-actions stock-desktop-actions" justify="left">
                    <n-button size="tiny" type="warning" @click="a().setStock && a().setStock(result['股票代码'],result['股票名称'])"> 成本
                    </n-button>
                    <n-button size="tiny" type="error"
                              @click="a().showFenshi && a().showFenshi(result['股票代码'],result['股票名称'],result.changePercent)"> 分时
                    </n-button>
                    <n-button size="tiny" type="error" @click="a().showK && a().showK(result['股票代码'],result['股票名称'])"> 日K</n-button>
                    <n-button size="tiny" type="error" v-if="result['买一报价']>0"
                              @click="a().showMoney && a().showMoney(result['股票代码'],result['股票名称'])"> 资金
                    </n-button>
                    <n-button size="tiny" type="success" @click="a().search && a().search(result['股票代码'],result['股票名称'])"> 详情
                    </n-button>
                    <n-button v-if="result['买一报价']>0" size="tiny" type="success"
                              @click="a().searchNotice && a().searchNotice(result['股票代码'])"> 公告
                    </n-button>
                    <n-button v-if="result['买一报价']>0" size="tiny" type="success"
                              @click="a().searchStockReport && a().searchStockReport(result['股票代码'])"> 研报
                    </n-button>
                    <n-flex justify="right">
                        <n-dropdown trigger="click" :options="groupList" key-field="ID" label-field="name"
                                    @select="(gId) => a().AddStockGroupInfo && a().AddStockGroupInfo(gId,result['股票代码'],result['股票名称'])">
                            <n-button type="warning" size="tiny">设置分组</n-button>
                        </n-dropdown>
                    </n-flex>
                </n-flex>
            </template>
        </n-card>
    </n-gi>
</template>
