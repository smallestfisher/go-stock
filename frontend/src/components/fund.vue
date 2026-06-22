<script setup>
import {ref, onBeforeUnmount} from "vue";
import {EventsOn, EventsOff} from "../api/runtime";
import {useRoute} from 'vue-router'
import FundFollow from "./FundFollow.vue";
import FundRanking from "./FundRanking.vue";

const nowTab = ref("基金自选")
const route = useRoute()

nowTab.value = route.query.name || '基金自选'

EventsOn("changeFundTab", async (msg) => {
  nowTab.value = msg.name
})

onBeforeUnmount(() => {
  EventsOff("changeFundTab")
})
</script>

<template>
  <n-card class="fund-page-shell">
    <n-tabs class="fund-mobile-tabs" type="line" animated v-model:value="nowTab" style="">
      <n-tab-pane name="基金自选" display-directive="show">
        <FundFollow/>
      </n-tab-pane>
      <n-tab-pane name="基金排行" display-directive="show">
        <FundRanking/>
      </n-tab-pane>
    </n-tabs>
  </n-card>
</template>

<style scoped>
@media (max-width: 768px) {
  .fund-page-shell {
    margin: 0 8px calc(var(--mobile-bottom-nav-height) + var(--safe-bottom) + 8px);
  }

  .fund-page-shell :deep(.n-card__content) {
    padding: 10px;
  }

  .fund-mobile-tabs :deep(.n-tabs-nav-scroll-content) {
    min-width: max-content;
  }

  .fund-mobile-tabs :deep(.n-tab-pane) {
    padding-top: 8px;
  }
}
</style>
