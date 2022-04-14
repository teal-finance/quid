<template>
  <div class="w-full text-center" :class="small ? 'py-8' : 'pt-24'">
    <i-eos-icons:loading class="text-5xl txt-light"></i-eos-icons:loading>
  </div>
</template>

<script setup lang="ts">
defineProps({
  small: {
    type: Boolean,
    default: false,
  }
})
</script>