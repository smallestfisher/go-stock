<script setup>
import { computed, onBeforeMount, ref } from 'vue'
import PageHeader from '../../components/widgets/PageHeader.vue'
import MPullRefresh from '../../components/base/MPullRefresh.vue'
import MCard from '../../components/base/MCard.vue'
import MEmpty from '../../components/base/MEmpty.vue'
import MLoading from '../../components/base/MLoading.vue'
import { GetFollowedFund, GetFundRanking } from '../../../api/app'

const loading = ref(false)
const followedFunds = ref([])
const rankingFunds = ref([])

const visibleFollowedFunds = computed(() => followedFunds.value.slice(0, 4))
const visibleRankingFunds = computed(() => rankingFunds.value.slice(0, 8))

function normalizePercent(value) {
  const n = Number(value)
  if (!Number.isFinite(n)) return '--'
  return `${n > 0 ? '+' : ''}${n.toFixed(2)}%`
}

function percentClass(value) {
  const n = Number(value)
  if (!Number.isFinite(n) || n === 0) return 'm-flat'
  return n > 0 ? 'm-rise' : 'm-fall'
}

function getFundCode(fund) {
  return fund.fundCode || fund.code || fund.FundCode || fund.Code || '--'
}

function getFundName(fund) {
  return fund.fundName || fund.name || fund.FundName || fund.Name || '未命名基金'
}

function getFundYield(fund) {
  return fund.jnzf ?? fund.yield ?? fund.rate ?? fund.ChangePercent ?? fund.changePercent
}

async function loadData() {
  loading.value = true
  try {
    const [followed, ranking] = await Promise.all([
      GetFollowedFund(),
      GetFundRanking('all', 'all', 'jnzf', 'desc', 1, 20)
    ])
    followedFunds.value = Array.isArray(followed) ? followed : []
    rankingFunds.value = Array.isArray(ranking?.items) ? ranking.items : []
  } catch (error) {
    console.error('加载基金数据失败:', error)
    followedFunds.value = []
    rankingFunds.value = []
  } finally {
    loading.value = false
  }
}

onBeforeMount(loadData)
</script>

<template>
  <div class="fund-page">
    <PageHeader title="基金中心" />

    <MPullRefresh :on-refresh="loadData">
      <div class="fund-content">
        <MLoading v-if="loading" text="加载基金数据..." />

        <template v-else>
          <MCard>
            <div class="section-header">
              <div>
                <h2 class="section-title">关注基金</h2>
                <p class="section-subtitle">已关注 {{ followedFunds.length }} 只</p>
              </div>
            </div>

            <div v-if="visibleFollowedFunds.length" class="fund-list">
              <div
                v-for="fund in visibleFollowedFunds"
                :key="getFundCode(fund)"
                class="fund-row"
              >
                <div class="fund-main">
                  <div class="fund-name">{{ getFundName(fund) }}</div>
                  <div class="fund-code">{{ getFundCode(fund) }}</div>
                </div>
                <div class="fund-yield" :class="percentClass(getFundYield(fund))">
                  {{ normalizePercent(getFundYield(fund)) }}
                </div>
              </div>
            </div>

            <MEmpty v-else description="还没有关注基金" image-size="88px" />
          </MCard>

          <MCard>
            <div class="section-header">
              <div>
                <h2 class="section-title">基金排行</h2>
                <p class="section-subtitle">按近年收益排序</p>
              </div>
            </div>

            <div v-if="visibleRankingFunds.length" class="fund-list">
              <div
                v-for="(fund, index) in visibleRankingFunds"
                :key="getFundCode(fund)"
                class="fund-row"
              >
                <div class="rank-index">{{ index + 1 }}</div>
                <div class="fund-main">
                  <div class="fund-name">{{ getFundName(fund) }}</div>
                  <div class="fund-code">{{ getFundCode(fund) }}</div>
                </div>
                <div class="fund-yield" :class="percentClass(getFundYield(fund))">
                  {{ normalizePercent(getFundYield(fund)) }}
                </div>
              </div>
            </div>

            <MEmpty v-else description="暂无基金排行数据" image-size="88px" />
          </MCard>
        </template>
      </div>
    </MPullRefresh>
  </div>
</template>

<style scoped>
.fund-page {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: var(--m-bg-primary);
}

.fund-content {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-md);
  padding: var(--m-space-md);
  padding-bottom: calc(var(--m-space-2xl) + var(--m-safe-bottom));
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--m-space-md);
}

.section-title {
  font-size: var(--m-font-lg);
  font-weight: var(--m-font-weight-bold);
  color: var(--m-text-primary);
}

.section-subtitle {
  margin-top: var(--m-space-xs);
  font-size: var(--m-font-xs);
  color: var(--m-text-tertiary);
}

.fund-list {
  display: flex;
  flex-direction: column;
}

.fund-row {
  display: flex;
  align-items: center;
  min-height: 58px;
  border-bottom: 1px solid var(--m-divider-color);
}

.fund-row:last-child {
  border-bottom: none;
}

.rank-index {
  width: 28px;
  font-size: var(--m-font-md);
  font-weight: var(--m-font-weight-bold);
  color: var(--m-text-tertiary);
}

.fund-main {
  flex: 1;
  min-width: 0;
}

.fund-name {
  overflow: hidden;
  color: var(--m-text-primary);
  font-size: var(--m-font-md);
  font-weight: var(--m-font-weight-medium);
  text-overflow: ellipsis;
  white-space: nowrap;
}

.fund-code {
  margin-top: var(--m-space-xs);
  color: var(--m-text-tertiary);
  font-size: var(--m-font-xs);
}

.fund-yield {
  min-width: 72px;
  text-align: right;
  font-size: var(--m-font-md);
  font-weight: var(--m-font-weight-bold);
}
</style>
