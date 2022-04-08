<template>
  <button class="btn" :class="type" :disabled="disabled">
    <slot></slot>
  </button>
</template>

<script setup lang="ts">
defineProps({
  type: {
    type: String as () => "delete" | "action",
    default: "action",
  },
  disabled: {
    type: Boolean,
    default: false,
  }
})
</script>

<style lang="sass" scoped>
.btn
  @apply text-sm
  &.delete
    @apply bord-danger txt-danger hover:danger
  &.action
    @apply bord-primary txt-primary hover:primary dark:txt-lighter dark:bord-light
</style>