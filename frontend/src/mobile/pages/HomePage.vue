<script setup>
import { ref } from 'vue'
import MCard from '../components/base/MCard.vue'
import MButton from '../components/base/MButton.vue'
import MSheet from '../components/base/MSheet.vue'
import MTabs from '../components/base/MTabs.vue'
import MLoading from '../components/base/MLoading.vue'
import MEmpty from '../components/base/MEmpty.vue'
import MPullRefresh from '../components/base/MPullRefresh.vue'

const deviceInfo = ref({
  width: window.innerWidth,
  height: window.innerHeight
})

// Sheet 测试
const sheetVisible = ref(false)

// Tabs 测试
const activeTab = ref('1')
const tabs = [
  { label: '全部', value: '1' },
  { label: '沪深', value: '2' },
  { label: '港股', value: '3' },
  { label: '美股', value: '4' },
  { label: '已禁用', value: '5', disabled: true },
]

// Loading 测试
const loading = ref(false)

// 下拉刷新测试
const refreshCount = ref(0)
async function handleRefresh() {
  return new Promise(resolve => {
    setTimeout(() => {
      refreshCount.value++
      resolve()
    }, 1500)
  })
}

function showLoading() {
  loading.value = true
  setTimeout(() => {
    loading.value = false
  }, 2000)
}
</script>

<template>
  <MPullRefresh :on-refresh="handleRefresh">
    <div class="home-page">
      <!-- 欢迎卡片 -->
      <MCard>
        <h1 class="welcome-title">Hello Mobile! 🎉</h1>
        <p class="welcome-subtitle">移动端基础组件已就绪</p>
        <div class="device-info">
          <div class="info-item">
            <span>设备宽度:</span>
            <span>{{ deviceInfo.width }}px</span>
          </div>
          <div class="info-item">
            <span>设备高度:</span>
            <span>{{ deviceInfo.height }}px</span>
          </div>
          <div class="info-item">
            <span>刷新次数:</span>
            <span class="m-rise">{{ refreshCount }} 次</span>
          </div>
        </div>
      </MCard>

      <!-- MButton 演示 -->
      <MCard>
        <h2 class="section-title">MButton - 按钮组件</h2>
        <div class="button-group">
          <MButton>默认按钮</MButton>
          <MButton type="primary">主要按钮</MButton>
          <MButton type="success">成功按钮</MButton>
          <MButton type="danger">危险按钮</MButton>
          <MButton type="text">文本按钮</MButton>
        </div>
        <div class="button-group">
          <MButton size="small">小按钮</MButton>
          <MButton size="medium">中按钮</MButton>
          <MButton size="large">大按钮</MButton>
        </div>
        <div class="button-group">
          <MButton disabled>禁用按钮</MButton>
          <MButton :loading="loading" @click="showLoading">加载按钮</MButton>
        </div>
        <MButton type="primary" block round>块级圆角按钮</MButton>
      </MCard>

      <!-- MTabs 演示 -->
      <MCard padding="none">
        <h2 class="section-title" style="padding: 16px 16px 0;">MTabs - 标签页组件</h2>
        <MTabs v-model="activeTab" :tabs="tabs" />
        <div style="padding: 16px;">
          <p>当前选中: <strong class="m-rise">{{ tabs.find(t => t.value === activeTab)?.label }}</strong></p>
        </div>
      </MCard>

      <!-- MLoading 演示 -->
      <MCard>
        <h2 class="section-title">MLoading - 加载组件</h2>
        <div class="loading-group">
          <MLoading size="small" />
          <MLoading size="medium" />
          <MLoading size="large" />
        </div>
        <div class="loading-group">
          <MLoading text="加载中..." />
          <MLoading text="请稍候..." vertical />
        </div>
      </MCard>

      <!-- MEmpty 演示 -->
      <MCard>
        <h2 class="section-title">MEmpty - 空状态组件</h2>
        <MEmpty description="暂无数据">
          <MButton type="primary" size="small">重新加载</MButton>
        </MEmpty>
      </MCard>

      <!-- MSheet 演示 -->
      <MCard>
        <h2 class="section-title">MSheet - 底部抽屉组件</h2>
        <MButton type="primary" @click="sheetVisible = true">打开底部抽屉</MButton>
      </MCard>

      <!-- 完成清单 -->
      <MCard clickable>
        <h2 class="section-title">✅ 阶段 1 完成清单</h2>
        <div class="checklist">
          <div class="checklist-item">✅ MCard - 卡片容器</div>
          <div class="checklist-item">✅ MButton - 按钮组件</div>
          <div class="checklist-item">✅ MSheet - 底部抽屉</div>
          <div class="checklist-item">✅ MTabs - 横向标签页</div>
          <div class="checklist-item">✅ MLoading - 加载状态</div>
          <div class="checklist-item">✅ MEmpty - 空状态</div>
          <div class="checklist-item">✅ MPullRefresh - 下拉刷新</div>
          <div class="checklist-item">✅ usePullRefresh - 下拉刷新逻辑</div>
        </div>
      </MCard>

      <!-- 下一步 -->
      <MCard>
        <h2 class="section-title">🚀 阶段 2：首页信息流</h2>
        <p class="next-desc">接下来开发首页信息流，实现真实数据展示</p>
        <ul class="next-list">
          <li>自选概览卡片</li>
          <li>市场快讯卡片</li>
          <li>热点话题卡片</li>
          <li>异动预警卡片</li>
        </ul>
      </MCard>
    </div>
  </MPullRefresh>

  <!-- 底部抽屉 -->
  <MSheet
    v-model:show="sheetVisible"
    title="底部抽屉示例"
    height="60vh"
  >
    <div style="padding: 20px;">
      <h3>这是一个底部抽屉</h3>
      <p>可以放置任何内容</p>
      <MButton type="primary" block @click="sheetVisible = false" style="margin-top: 20px;">
        关闭抽屉
      </MButton>
    </div>

    <template #footer>
      <MButton block type="primary">底部操作按钮</MButton>
    </template>
  </MSheet>
</template>

<style scoped>
.home-page {
  padding: var(--m-space-md);
  display: flex;
  flex-direction: column;
  gap: var(--m-space-md);
}

.welcome-title {
  font-size: var(--m-font-2xl);
  font-weight: var(--m-font-weight-bold);
  color: var(--m-color-rise);
  margin-bottom: var(--m-space-sm);
  text-align: center;
}

.welcome-subtitle {
  font-size: var(--m-font-lg);
  color: var(--m-text-secondary);
  text-align: center;
  margin-bottom: var(--m-space-lg);
}

.device-info {
  background: var(--m-bg-primary);
  border-radius: var(--m-radius-sm);
  padding: var(--m-space-md);
}

.info-item {
  display: flex;
  justify-content: space-between;
  padding: var(--m-space-sm) 0;
  font-size: var(--m-font-sm);
  color: var(--m-text-secondary);
}

.info-item span:last-child {
  color: var(--m-text-primary);
  font-weight: var(--m-font-weight-medium);
}

.section-title {
  font-size: var(--m-font-lg);
  font-weight: var(--m-font-weight-medium);
  color: var(--m-text-primary);
  margin-bottom: var(--m-space-lg);
}

.button-group {
  display: flex;
  flex-wrap: wrap;
  gap: var(--m-space-md);
  margin-bottom: var(--m-space-md);
}

.loading-group {
  display: flex;
  justify-content: space-around;
  align-items: center;
  padding: var(--m-space-lg) 0;
  margin-bottom: var(--m-space-md);
}

.checklist {
  display: flex;
  flex-direction: column;
  gap: var(--m-space-sm);
}

.checklist-item {
  padding: var(--m-space-md);
  background: var(--m-color-rise-light);
  color: var(--m-color-rise);
  border-radius: var(--m-radius-sm);
  font-size: var(--m-font-md);
  font-weight: var(--m-font-weight-medium);
}

.next-desc {
  color: var(--m-text-secondary);
  margin-bottom: var(--m-space-md);
  line-height: var(--m-line-height-normal);
}

.next-list {
  margin-left: var(--m-space-lg);
  color: var(--m-text-secondary);
  line-height: var(--m-line-height-loose);
}
</style>
