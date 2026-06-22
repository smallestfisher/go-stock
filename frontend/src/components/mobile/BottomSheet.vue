<script setup>
import {computed} from 'vue'

const props = defineProps({
    show: {type: Boolean, default: false},
    title: {type: String, default: ''},
    // 抽屉高度，默认 90vh
    height: {type: String, default: '90vh'},
    closable: {type: Boolean, default: true},
    // 内容是否可滚动（图表等固定高度内容建议关闭，由内部布局自行滚动）
    bodyScrollable: {type: Boolean, default: true},
})

const emit = defineEmits(['update:show'])

const drawerStyle = computed(() => ({
    height: props.height,
    paddingBottom: 'var(--safe-bottom)',
}))
</script>

<template>
    <n-drawer
        :show="show"
        @update:show="(v) => emit('update:show', v)"
        :placement="'bottom'"
        :height="height"
        class="mobile-bottom-sheet"
    >
        <n-drawer-content
            :title="title"
            :closable="closable"
            :native-scrollbar="!bodyScrollable"
            body-content-style="padding: 0;"
        >
            <slot></slot>
            <template v-if="$slots.footer" #footer>
                <slot name="footer"></slot>
            </template>
        </n-drawer-content>
    </n-drawer>
</template>

<style scoped>
.mobile-bottom-sheet {
    border-radius: 12px 12px 0 0;
    overflow: hidden;
}

.mobile-bottom-sheet :deep(.n-drawer-content) {
    border-radius: 12px 12px 0 0;
}
</style>
