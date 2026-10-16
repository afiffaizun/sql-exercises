<script setup>
defineProps({
  toasts: { type: Array, default: () => [] }
})
</script>

<template>
  <div class="fixed right-4 top-4 z-[60] flex w-full max-w-sm flex-col gap-2">
    <transition-group name="toast">
      <div
        v-for="t in toasts"
        :key="t.id"
        class="flex items-start gap-3 rounded-xl border p-4 shadow-lg backdrop-blur animate-fade-in"
        :class="t.type === 'error'
          ? 'border-red-200 bg-red-50/95 text-red-800'
          : 'border-emerald-200 bg-emerald-50/95 text-emerald-800'"
      >
        <span class="mt-0.5 text-lg leading-none">{{ t.type === 'error' ? '⚠️' : '✅' }}</span>
        <p class="text-sm font-medium">{{ t.message }}</p>
      </div>
    </transition-group>
  </div>
</template>

<style scoped>
.toast-enter-active,
.toast-leave-active {
  transition: all 0.25s ease;
}
.toast-enter-from,
.toast-leave-to {
  opacity: 0;
  transform: translateX(20px);
}
</style>
