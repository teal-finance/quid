<template>
  <sw-sidebar id="sidebar" :opened="sidebar" class="h-full sm:flex-col sm:flex">
    <div class="flex-grow mt-5 ml-6 space-y-6" v-if="user.type.value == 'serverAdmin'">
      <div v-for="(category, i) in new Set([...serverAdminCategories, ...nsAdminCategories])" :key="i"
        :class="{ 'ml-1': i < 2, 'cursor-pointer': true }">
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
    <div class="flex-grow mt-5 ml-6 space-y-6" v-else-if="user.type.value == 'nsAdmin'">
      <div v-for="(category, i) in nsAdminCategories" :key="i" :class="{ 'ml-1': i < 2, 'cursor-pointer': true }">
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
    <div class="flex-none w-full h-12 pl-5 mb-3 text-2xl cursor-pointer" @click="toggle()">
      <icon icon="fa-solid:angle-double-left" v-if="sidebar === true"></icon>
      <icon icon="fa-solid:angle-double-right" v-else></icon>
    </div>
  </sw-sidebar>
</template>

<script setup lang="ts">
import { Icon } from '@iconify/vue';
import { serverAdminCategories, nsAdminCategories } from "@/const/categories";
import SwSidebar from "@snowind/sidebar";
import { user } from "@/state";

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
