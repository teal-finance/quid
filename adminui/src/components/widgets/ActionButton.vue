<template>
  <button class="btn" :class="type">
    <slot></slot>
  </button>
</template>

<script setup lang="ts">
defineProps({
  type: {
    type: String as () => "delete" | "action",
    default: "action",
  },
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