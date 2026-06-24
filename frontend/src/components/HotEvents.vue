<script setup lang="ts">
import {onBeforeMount, onUnmounted, ref} from 'vue'
import {HotEvent} from "../api/app";
import {registerFeed, stopFeed} from "../api/scheduler";
const list  = ref([])

async function loadHotEvents() { list.value = await HotEvent(50) }

onBeforeMount(async () => {
  await loadHotEvents()
  // 轮询交给统一调度器，页面恢复时自动重刷。
  registerFeed("hotEvents", { fetch: loadHotEvents, intervalMs: 1000 * 10 })
})

onUnmounted(() => {
  stopFeed("hotEvents")
})
</script>

<template>
  <n-list bordered>
    <template #header>
      雪球热门
    </template>
    <n-list-item v-for="(item, index) in list" :key="index">
        <n-thing :title="item.tag" :description="item.content"  >
          <template v-if="item.pic" #avatar>
            <n-avatar :src="item.pic" :size="60">
            </n-avatar>
          </template>
        </n-thing>
    </n-list-item>
  </n-list>
</template>

<style scoped>

</style>