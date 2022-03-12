<template>
  <span
    class="inline-block px-2 py-1 mr-3 text-xs font-bold rounded-full"
    :class="color"
    v-html="text"
  ></span>
</template>

<script setup lang="ts">
defineProps({
  text: {
    type: String,
    required: true,
  },
  color: {
    type: String,
    default: "primary"
  }
})
</script>