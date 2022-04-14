<template>
  <span class="inline-block px-2 py-1 mr-3 text-xs font-bold rounded-full" @click="onClick()">
    <slot>{{ text }}</slot>
  </span>
</template>

<script setup lang="ts">
defineProps({
  text: {
    type: String,
    default: ""
  },
});

const emit = defineEmits(["select"]);

function onClick() {
  emit("select")
}
</script>