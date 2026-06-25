<script setup>
import { ref } from 'vue'
import MPullRefresh from '../../components/base/MPullRefresh.vue'
import MCard from '../../components/base/MCard.vue'
import MEmpty from '../../components/base/MEmpty.vue'
import MButton from '../../components/base/MButton.vue'

// MCP服务列表（对齐桌面端 mcp-server-manager.vue）
// 字段参考：{ name, baseUrl, enabled, toolCount, description }
// 数据后续接 GetMCPServerList 填充，操作接 EnableMCPServer/DeleteMCPServer/CreateMCPServer
const servers = ref([])

async function handleRefresh() {
  // 后续接真实 API
}

function handleAdd() {
  // 后续打开新增抽屉，接 CreateMCPServer
}
</script>

<template>
  <div class="page">
    <MPullRefresh :on-refresh="handleRefresh">
      <div class="container">
        <MCard v-if="servers.length" padding="none">
          <div class="server-list">
            <div v-for="(server, i) in servers" :key="i" class="server-item">
              <div class="server-main">
                <div class="server-header">
                  <div class="server-name">{{ server.name }}</div>
                  <span class="server-status" :class="server.enabled ? 'server-status--on' : 'server-status--off'">
                    {{ server.enabled ? '已启用' : '已停用' }}
                  </span>
                </div>
                <div v-if="server.baseUrl" class="server-url">{{ server.baseUrl }}</div>
                <p v-if="server.description" class="server-desc">{{ server.description }}</p>
                <div v-if="server.toolCount != null" class="server-tools">
                  🛠 {{ server.toolCount }} 个工具
                </div>
              </div>
            </div>
          </div>
        </MCard>
        <MEmpty v-else description="还没有配置 MCP 服务">
          <MButton type="primary" @click="handleAdd">添加服务</MButton>
        </MEmpty>
      </div>
    </MPullRefresh>
  </div>
</template>

<style scoped>
.page { height: 100%; overflow: hidden; }
.container { padding: var(--m-space-md); }
.server-list { padding: var(--m-space-md); display: flex; flex-direction: column; gap: var(--m-space-md); }
.server-item { padding: var(--m-space-md); background: var(--m-bg-primary); border-radius: var(--m-radius-sm); }
.server-main { display: flex; flex-direction: column; gap: var(--m-space-xs); }
.server-header { display: flex; justify-content: space-between; align-items: center; }
.server-name { font-size: var(--m-font-md); font-weight: var(--m-font-weight-medium); color: var(--m-text-primary); }
.server-status { font-size: var(--m-font-xs); padding: 2px var(--m-space-sm); border-radius: var(--m-radius-sm); }
.server-status--on { background: var(--m-color-rise-light); color: var(--m-color-rise); }
.server-status--off { background: var(--m-divider-color); color: var(--m-text-tertiary); }
.server-url { font-size: var(--m-font-xs); color: var(--m-text-secondary); word-break: break-all; }
.server-desc { font-size: var(--m-font-sm); color: var(--m-text-secondary); }
.server-tools { font-size: var(--m-font-xs); color: var(--m-text-tertiary); }
</style>
