<template>
  <div class="p-4 text-lg text-center bg-transparent border shadow dark:bg-gray-500">
    <div class="flex justify-center w-full">
      <icon :icon="card.icon" class="mt-2 text-3xl txt-neutral"></icon>
    </div>
    <div class="mt-3 mb-3 dark:font-bold">{{ card.title }}</div>
  </div>
</template>

<script setup lang="ts">
import { defineComponent } from "vue";
import Category from "@/models/category";
import { Icon } from '@iconify/vue';

defineProps({
  card: {
    type: Category,
    required: true,
  },
});
</script>