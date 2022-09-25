<template>
  <div class="flex flex-col space-y-2">
    <div v-for="doc in state.examples" @click="$router.push(`/example/${encodeURIComponent(doc)}`)"
      class="cursor-pointer">
      {{doc}}
    </div>
  </div>
</template>

<script setup lang="ts">
import { state } from '@/state';
</script>