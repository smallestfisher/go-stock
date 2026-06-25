<script setup>
import { ref } from 'vue'
import MPullRefresh from '../../components/base/MPullRefresh.vue'
import MCard from '../../components/base/MCard.vue'
import MEmpty from '../../components/base/MEmpty.vue'
import MButton from '../../components/base/MButton.vue'

// 提示词模板列表（后续接 GetPromptTemplateList 填充）
// 字段参考桌面端 promptTemplateList.vue：{ name, content, type, CreatedAt }
const templates = ref([])

async function handleRefresh() {
  // 后续接真实 API
}

function handleAdd() {
  // 后续打开新增抽屉，接 AddPromptTemplate
}
</script>

<template>
  <div class="page">
    <MPullRefresh :on-refresh="handleRefresh">
      <div class="container">
        <MCard v-if="templates.length" padding="none">
          <div class="template-list">
            <div v-for="(item, i) in templates" :key="i" class="template-item">
              <div class="item-header">
                <div class="item-name">{{ item.name }}</div>
                <span v-if="item.type" class="item-type">{{ item.type }}</span>
              </div>
              <p class="item-content">{{ item.content }}</p>
            </div>
          </div>
        </MCard>
        <MEmpty v-else description="还没有提示词模板">
          <MButton type="primary" @click="handleAdd">新建模板</MButton>
        </MEmpty>
      </div>
    </MPullRefresh>
  </div>
</template>

<style scoped>
.page { height: 100%; overflow: hidden; }
.container { padding: var(--m-space-md); }
.template-list { padding: var(--m-space-md); display: flex; flex-direction: column; gap: var(--m-space-lg); }
.template-item { padding-bottom: var(--m-space-lg); border-bottom: 1px solid var(--m-divider-color); }
.template-item:last-child { border-bottom: none; }
.item-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: var(--m-space-sm); }
.item-name { font-size: var(--m-font-md); font-weight: var(--m-font-weight-medium); color: var(--m-text-primary); }
.item-type { font-size: var(--m-font-xs); color: var(--m-color-rise); background: var(--m-color-rise-light); padding: 2px var(--m-space-sm); border-radius: var(--m-radius-sm); }
.item-content { font-size: var(--m-font-sm); color: var(--m-text-secondary); line-height: var(--m-line-height-normal); }
</style>
