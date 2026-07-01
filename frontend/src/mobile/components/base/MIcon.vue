<script setup>
// 轻量本地矢量图标组件：不引第三方图标库（本机内存不足、无法 build 验证依赖），
// 用内联 SVG + currentColor 实现。相比 emoji 的优势：
//   1. 吃 CSS color —— 选中/禁用/涨跌语义色可直接用 color 控制，无需去色 hack；
//   2. 跨平台渲染一致 —— 不受 iOS/Android/系统 emoji 字体差异影响；
//   3. 线条粗细/尺寸可控，视觉统一。
// 图标为 24x24 viewBox 的描边风格（stroke=currentColor）。
// filled=true 时同时用 currentColor 填充，用于「实心/空心」语义（收藏★/☆、点赞♥/♡、色点●）。
// 新增图标只需往 PATHS 里加。
import { computed } from 'vue'

const props = defineProps({
  name: { type: String, required: true },
  size: { type: [Number, String], default: 24 },
  // 描边粗细，可按场景微调（dock 用略粗更醒目）
  strokeWidth: { type: [Number, String], default: 2 },
  // 是否填充：收藏/点赞/色点等「实心态」用 true
  filled: { type: Boolean, default: false },
})

// 每个图标是一段/多段 SVG 子元素。统一 24x24、round 线帽，视觉一致。
// 命名对齐业务语义，方便调用处按含义取用。
const PATHS = {
  // —— 主导航 / 抽屉 ——
  home: '<path d="M3 10.5 12 3l9 7.5"/><path d="M5 9.5V21h14V9.5"/>',
  star: '<path d="M12 3.5l2.6 5.3 5.9.9-4.3 4.2 1 5.9L12 17l-5.2 2.7 1-5.9L3.5 9.7l5.9-.9z"/>',
  chart: '<path d="M4 20V4"/><path d="M4 20h16"/><rect x="7" y="11" width="3" height="6"/><rect x="13" y="7" width="3" height="10"/>',
  research: '<path d="M9 3h6"/><path d="M10 3v6l-4.5 8a2 2 0 0 0 1.8 3h9.4a2 2 0 0 0 1.8-3L14 9V3"/><path d="M8 15h8"/>',
  menu: '<path d="M4 7h16"/><path d="M4 12h16"/><path d="M4 17h16"/>',
  kline: '<path d="M3 3v18h18"/><path d="M7 14l3-4 3 3 4-6"/>',
  fund: '<rect x="3" y="6" width="18" height="13" rx="2"/><path d="M3 10h18"/><path d="M16 15h2"/>',
  robot: '<rect x="5" y="8" width="14" height="11" rx="2"/><path d="M12 8V4"/><circle cx="12" cy="3" r="1"/><path d="M9 13h.01"/><path d="M15 13h.01"/><path d="M2 12v3"/><path d="M22 12v3"/>',
  settings: '<circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.6 1.6 0 0 0 .3 1.8l.1.1a2 2 0 1 1-2.8 2.8l-.1-.1a1.6 1.6 0 0 0-2.7 1.1V21a2 2 0 1 1-4 0v-.1A1.6 1.6 0 0 0 7 19.4a1.6 1.6 0 0 0-1.8.3l-.1.1a2 2 0 1 1-2.8-2.8l.1-.1a1.6 1.6 0 0 0-1.1-2.7H1a2 2 0 1 1 0-4h.1A1.6 1.6 0 0 0 2.6 7a1.6 1.6 0 0 0-.3-1.8l-.1-.1a2 2 0 1 1 2.8-2.8l.1.1a1.6 1.6 0 0 0 1.8.3H7a1.6 1.6 0 0 0 1-1.5V1a2 2 0 1 1 4 0v.1a1.6 1.6 0 0 0 1 1.5 1.6 1.6 0 0 0 1.8-.3l.1-.1a2 2 0 1 1 2.8 2.8l-.1.1a1.6 1.6 0 0 0-.3 1.8V7a1.6 1.6 0 0 0 1.5 1H23a2 2 0 1 1 0 4h-.1a1.6 1.6 0 0 0-1.5 1z"/>',
  info: '<circle cx="12" cy="12" r="9"/><path d="M12 11v5"/><path d="M12 8h.01"/>',

  // —— 通用操作 / 箭头 ——
  close: '<path d="M6 6l12 12"/><path d="M18 6 6 18"/>',
  'arrow-right': '<path d="M5 12h14"/><path d="M13 6l6 6-6 6"/>',
  'arrow-left': '<path d="M19 12H5"/><path d="M11 18l-6-6 6-6"/>',
  'arrow-up': '<path d="M12 19V5"/><path d="M6 11l6-6 6 6"/>',
  'arrow-down': '<path d="M12 5v14"/><path d="M6 13l6 6 6-6"/>',
  'arrow-up-right': '<path d="M7 17 17 7"/><path d="M8 7h9v9"/>',
  plus: '<path d="M12 5v14"/><path d="M5 12h14"/>',
  check: '<path d="M5 12l5 5L20 7"/>',
  'check-circle': '<circle cx="12" cy="12" r="9"/><path d="M8 12l3 3 5-6"/>',
  search: '<circle cx="11" cy="11" r="7"/><path d="M21 21l-4.3-4.3"/>',
  trash: '<path d="M4 7h16"/><path d="M10 11v6"/><path d="M14 11v6"/><path d="M5 7l1 13h12l1-13"/><path d="M9 7V4h6v3"/>',
  edit: '<path d="M4 20h4L18.5 9.5a2.12 2.12 0 0 0-3-3L5 17v3z"/><path d="M13.5 6.5l3 3"/>',
  eye: '<path d="M2 12s3.5-7 10-7 10 7 10 7-3.5 7-10 7-10-7-10-7z"/><circle cx="12" cy="12" r="3"/>',
  refresh: '<path d="M21 12a9 9 0 1 1-3-6.7"/><path d="M21 3v5h-5"/>',
  question: '<circle cx="12" cy="12" r="9"/><path d="M9.5 9.5a2.5 2.5 0 0 1 4.5 1.5c0 1.7-2.5 2-2.5 3.5"/><path d="M12 18h.01"/>',
  clipboard: '<rect x="6" y="4" width="12" height="17" rx="2"/><path d="M9 4V3h6v1"/><path d="M9 10h6"/><path d="M9 14h6"/>',

  // —— 语义可切换（filled 控制实心/空心）——
  heart: '<path d="M12 21C7 17 3 13 3 8.5 3 6 5 4 7.5 4c1.6 0 3 .8 3.9 2 .9-1.2 2.3-2 3.9-2C20 4 22 6 22 8.5 22 13 17 17 12 21z"/>',
  dot: '<circle cx="12" cy="12" r="6"/>',

  // —— 行情 / 财经 ——
  'trend-up': '<path d="M3 17l6-6 4 4 8-8"/><path d="M17 7h4v4"/>',
  'trend-down': '<path d="M3 7l6 6 4-4 8 8"/><path d="M17 17h4v-4"/>',
  fire: '<path d="M12 3c.5 3 3 4 3 7.5A3 3 0 0 1 9 11c0-1 .4-1.7.8-2.2C8 10 7 12 7 14a5 5 0 0 0 10 0c0-4-3-6-5-11z"/>',
  money: '<circle cx="12" cy="12" r="9"/><path d="M12 7v10"/><path d="M14.5 9.3c0-1-1.1-1.8-2.5-1.8s-2.5.8-2.5 1.9c0 2.6 5 1.4 5 4 0 1.1-1.1 1.9-2.5 1.9s-2.5-.8-2.5-1.8"/>',
  globe: '<circle cx="12" cy="12" r="9"/><path d="M3 12h18"/><path d="M12 3c2.6 3 2.6 15 0 18"/><path d="M12 3c-2.6 3-2.6 15 0 18"/>',
  target: '<circle cx="12" cy="12" r="9"/><circle cx="12" cy="12" r="5"/><circle cx="12" cy="12" r="1.4"/>',
  trophy: '<path d="M8 4h8v5a4 4 0 0 1-8 0z"/><path d="M8 6H5a2 2 0 0 0 2 3.5"/><path d="M16 6h3a2 2 0 0 1-2 3.5"/><path d="M10 20h4"/><path d="M12 13v4"/>',
  rocket: '<path d="M9 12c1-4.5 4-8.5 10-9 0 6-3.5 9-8 10z"/><path d="M9 12l3 3"/><path d="M5.5 15c-1 1-1.3 4-1.3 4s3-.3 4-1.3a1.8 1.8 0 1 0-2.7-2.7z"/><circle cx="14.5" cy="9.5" r="1.3"/>',
  calendar: '<rect x="3" y="5" width="18" height="16" rx="2"/><path d="M3 9h18"/><path d="M8 3v4"/><path d="M16 3v4"/>',
  clock: '<circle cx="12" cy="12" r="8"/><path d="M12 8v4l3 2"/>',
  news: '<rect x="3" y="5" width="18" height="14" rx="2"/><path d="M7 9h6"/><path d="M7 13h6"/><path d="M16 9h2"/><path d="M16 13h2"/>',
  comment: '<path d="M21 12a8 8 0 0 1-11.5 7.2L4 20.5l1.3-5A8 8 0 1 1 21 12z"/>',

  // —— 分析 / 思考 ——
  brain: '<path d="M12 5a3 3 0 0 0-5.9-.6A3 3 0 0 0 4 9.5 3 3 0 0 0 5.5 15 3 3 0 0 0 12 19z"/><path d="M12 5a3 3 0 0 1 5.9-.6A3 3 0 0 1 20 9.5 3 3 0 0 1 18.5 15 3 3 0 0 1 12 19z"/><path d="M12 5v14"/>',
  thinking: '<path d="M8 15a4.5 4.5 0 0 1-.5-8.97A5 5 0 0 1 17 6.5a3.5 3.5 0 0 1-.5 7.5H8z"/><circle cx="6" cy="18.5" r="1"/><circle cx="9" cy="20.5" r="1"/>',
  sparkle: '<path d="M12 3l1.8 5.2L19 10l-5.2 1.8L12 17l-1.8-5.2L5 10l5.2-1.8z"/><path d="M18 15l.7 2 2 .7-2 .7-.7 2-.7-2-2-.7 2-.7z"/>',
  bolt: '<path d="M13 2 4 14h7l-1 8 9-12h-7z"/>',
  warning: '<path d="M12 3 2 20h20z"/><path d="M12 10v4"/><path d="M12 17h.01"/>',
  alarm: '<circle cx="12" cy="13" r="8"/><path d="M12 9v4l3 2"/><path d="M5 3 2 6"/><path d="M19 3l3 3"/>',
  volatility: '<path d="M3 12h3l3-8 4 16 3-8h5"/>',
  ruler: '<path d="M4 16 16 4l4 4L8 20z"/><path d="M8 8l2 2"/><path d="M11 5l2 2"/><path d="M5 11l2 2"/>',
  burst: '<path d="M12 2l2.2 6 5.8-2-3.5 5 3.5 5-5.8-2L12 22l-2.2-6-5.8 2 3.5-5-3.5-5 5.8 2z"/>',

  // —— 系统 / 设置 ——
  palette: '<path d="M12 3a9 9 0 0 0 0 18c1.4 0 2-1 2-2s-.5-1.3-.5-2 .6-1 1.5-1H18a3 3 0 0 0 3-3c0-4.4-4-7-9-7z"/><circle cx="7.5" cy="11" r="1"/><circle cx="12" cy="7.5" r="1"/><circle cx="16.5" cy="11" r="1"/>',
  bell: '<path d="M6 9a6 6 0 0 1 12 0c0 5 2 6 2 6H4s2-1 2-6z"/><path d="M10 20a2 2 0 0 0 4 0"/>',
  save: '<path d="M5 3h11l3 3v15H5z"/><path d="M8 3v5h7V3"/><rect x="8" y="13" width="8" height="5"/>',
  folder: '<path d="M3 7a2 2 0 0 1 2-2h4l2 2h8a2 2 0 0 1 2 2v8a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/>',
  'folder-open': '<path d="M3 7a2 2 0 0 1 2-2h4l2 2h8a2 2 0 0 1 2 2H3z"/><path d="M3 9h18l-2.2 8.3a2 2 0 0 1-1.9 1.7H5a2 2 0 0 1-2-2z"/>',
  cloud: '<path d="M7 18a4 4 0 0 1 .5-7.97A5.5 5.5 0 0 1 18 9.5a3.5 3.5 0 0 1-.5 8.5H7z"/>',
  device: '<rect x="7" y="3" width="10" height="18" rx="2"/><path d="M11 18h2"/>',
  github: '<path d="M9 19c-4 1.4-4-2-6-2m12 4v-3.4c0-1 .3-2-.5-2.7 2.4-.3 4.5-1.2 4.5-5a3.9 3.9 0 0 0-1-2.7 3.6 3.6 0 0 0-.1-2.7s-1-.3-3.3 1.3a11 11 0 0 0-6 0C6.8 2.5 5.8 2.8 5.8 2.8a3.6 3.6 0 0 0-.1 2.7A3.9 3.9 0 0 0 4.7 8.2c0 3.8 2.1 4.7 4.5 5-.6.6-.6 1.2-.5 2V21"/>',
  signal: '<circle cx="6" cy="18" r="2"/><path d="M5 12a7 7 0 0 1 7 7"/><path d="M5 6a13 13 0 0 1 13 13"/>',
  package: '<path d="M12 3l8 4.5v9L12 21l-8-4.5v-9z"/><path d="M4 7.5l8 4.5 8-4.5"/><path d="M12 12v9"/>',
  lock: '<rect x="5" y="11" width="14" height="10" rx="2"/><path d="M8 11V8a4 4 0 0 1 8 0v3"/>',
  bug: '<rect x="8" y="7" width="8" height="12" rx="4"/><path d="M12 7V4"/><path d="M9 5 8 3"/><path d="M15 5l1-2"/><path d="M8 11H4"/><path d="M8 15H4"/><path d="M16 11h4"/><path d="M16 15h4"/>',
  dragon: '<path d="M4 14c2-1 3-3 3-5a4 4 0 0 1 8 0c2 0 4 1 4 3s-2 3-4 3h-3c-1 0-2 1-2 2v2"/><path d="M11 7h.01"/>',
}

const inner = computed(() => PATHS[props.name] || '')
const px = computed(() => (typeof props.size === 'number' ? `${props.size}px` : props.size))
const fillColor = computed(() => (props.filled ? 'currentColor' : 'none'))
</script>

<template>
  <svg
    class="m-icon"
    :width="px"
    :height="px"
    viewBox="0 0 24 24"
    :fill="fillColor"
    :stroke-width="strokeWidth"
    stroke="currentColor"
    stroke-linecap="round"
    stroke-linejoin="round"
    aria-hidden="true"
    v-html="inner"
  />
</template>

<style scoped>
.m-icon {
  display: inline-block;
  flex-shrink: 0;
  vertical-align: middle;
}
</style>
