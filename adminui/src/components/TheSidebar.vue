<template>
  <sw-sidebar :opened="sidebar" class="flex flex-col h-full secondary">
    <div class="flex-grow mt-5 ml-6 space-y-6">
      <div
        v-for="(category, i) in categories"
        :key="i"
        :class="{ 'ml-1': i < 2, 'cursor-pointer': true }"
      >
        <router-link :to="category.url">
          <Icon :icon="category.icon" class="inline-block text-2xl" />
          <span v-if="sidebar === true" class="ml-3 text-lg">
            {{
              category.title
            }}
          </span>
        </router-link>
      </div>
    </div>
    <div
      class="flex-none w-full h-12 pl-5 mb-3 text-2xl cursor-pointer txt-light"
      @click="toggle()"
    >
      <icon icon="fa-solid:angle-double-left" v-if="sidebar === true"></icon>
      <icon icon="fa-solid:angle-double-right" v-else></icon>
    </div>
  </sw-sidebar>
</template>

<script setup lang="ts">
import { Icon } from '@iconify/vue';
import categories from "@/const/categories";
import SwSidebar from "@snowind/sidebar";

defineProps({
  sidebar: {
    type: Boolean,
    required: true,
  },
});

const emit = defineEmits(["toggle"]);

function toggle() {
  emit("toggle")
}
</script>