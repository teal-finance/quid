<template>
  <div class="flex flex-col space-y-2">
    <div v-for="doc in state.apidocs" @click="$router.push(`/apidoc/${encodeURIComponent(doc)}`)"
      class="capitalize cursor-pointer">
      {{doc}}
    </div>
  </div>
</template>

<script setup lang="ts">
import { state } from '@/state';
</script>