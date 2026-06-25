<script setup>
import { ref } from 'vue'
import MPullRefresh from '../../components/base/MPullRefresh.vue'
import MCard from '../../components/base/MCard.vue'
import MEmpty from '../../components/base/MEmpty.vue'
import MButton from '../../components/base/MButton.vue'

// 定时任务列表（对齐桌面端 cron-task-manager.vue）
// 字段参考：{ name, type, cron, enabled, nextRunTime }
// 数据后续接 GetCronTaskList 填充，操作接 EnableCronTask/ExecuteCronTaskNow/DeleteCronTask
const tasks = ref([])

async function handleRefresh() {
  // 后续接真实 API
}

function toggleEnabled(task) {
  // 后续接 EnableCronTask(id, enabled)
}

function executeNow(task) {
  // 后续接 ExecuteCronTaskNow(id)
}

function handleAdd() {
  // 后续打开新增抽屉，接 CreateCronTask
}
</script>

<template>
  <div class="page">
    <MPullRefresh :on-refresh="handleRefresh">
      <div class="container">
        <MCard v-if="tasks.length" padding="none">
          <div class="task-list">
            <div v-for="(task, i) in tasks" :key="i" class="task-item">
              <div class="task-main">
                <div class="task-name">{{ task.name }}</div>
                <div class="task-meta">
                  <span v-if="task.type" class="task-type">{{ task.type }}</span>
                  <span class="task-cron">⏰ {{ task.cron }}</span>
                </div>
                <div v-if="task.nextRunTime" class="task-next">下次: {{ task.nextRunTime }}</div>
              </div>
              <div class="task-actions">
                <span class="task-status" :class="task.enabled ? 'task-status--on' : 'task-status--off'">
                  {{ task.enabled ? '启用' : '停用' }}
                </span>
              </div>
            </div>
          </div>
        </MCard>
        <MEmpty v-else description="还没有定时任务">
          <MButton type="primary" @click="handleAdd">新建任务</MButton>
        </MEmpty>
      </div>
    </MPullRefresh>
  </div>
</template>

<style scoped>
.page { height: 100%; overflow: hidden; }
.container { padding: var(--m-space-md); }
.task-list { padding: var(--m-space-md); display: flex; flex-direction: column; gap: var(--m-space-md); }
.task-item { display: flex; justify-content: space-between; align-items: flex-start; padding: var(--m-space-md); background: var(--m-bg-primary); border-radius: var(--m-radius-sm); }
.task-main { flex: 1; }
.task-name { font-size: var(--m-font-md); font-weight: var(--m-font-weight-medium); color: var(--m-text-primary); margin-bottom: var(--m-space-xs); }
.task-meta { display: flex; gap: var(--m-space-md); font-size: var(--m-font-xs); color: var(--m-text-secondary); margin-bottom: var(--m-space-xs); }
.task-type { color: var(--m-color-rise); }
.task-next { font-size: var(--m-font-xs); color: var(--m-text-tertiary); }
.task-actions { display: flex; flex-direction: column; align-items: flex-end; gap: var(--m-space-xs); }
.task-status { font-size: var(--m-font-xs); padding: 2px var(--m-space-sm); border-radius: var(--m-radius-sm); }
.task-status--on { background: var(--m-color-rise-light); color: var(--m-color-rise); }
.task-status--off { background: var(--m-divider-color); color: var(--m-text-tertiary); }
</style>
