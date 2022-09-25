<template>
  <div class="flex flex-col space-y-2">
    <div v-for="doc in state.server" @click="$router.push(`/server/${encodeURIComponent(doc)}`)"
      class="capitalize cursor-pointer">
      {{doc}}
    </div>
  </div>
</template>

<script setup lang="ts">
import { state } from '@/state';
</script>