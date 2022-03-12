<template>
  <div>
    <sw-header class="w-full h-16 primary" @togglemenu="isMenuVisible = !isMenuVisible">
      <template #mobile-back>
        <i-ion-arrow-back-outline class="inline-flex ml-2 text-3xl" v-if="!isHome"></i-ion-arrow-back-outline>
      </template>
      <template #mobile-branding>
        <div
          class="inline-flex flex-row items-center h-full pt-1 ml-2 text-2xl truncate"
          @click="$router.push('/')"
        >
          <img alt="Snowind logo" src="@/assets/logo.png" v-if="isHome" class="inline-block mx-3" />
          <span class="text-lg">Quid</span>
        </div>
      </template>
      <template #branding>
        <div class="flex flex-row items-center h-full cursor-pointer" @click="$router.push('/')">
          <img alt="Snowind logo" src="@/assets/logo.png" class="inline-block mx-3" />
          <span class="text-lg">Quid</span>
        </div>
      </template>
      <template #menu>
        <div class="flex flex-row items-center justify-end w-full h-full space-x-1 txt-light">
          <button class="border-none btn" @click="$router.push('/namespaces')">Namespaces</button>
          <div class="px-5 text-lg cursor-pointer" @click="$router.push('/settings')">
            <i-fluent-settings-32-regular class="text-light dark:text-light-dark"></i-fluent-settings-32-regular>
          </div>
        </div>
      </template>
    </sw-header>
    <sw-mobile-menu :is-visible="isMenuVisible">
      <div class="flex flex-col p-3 space-y-5">
        <router-link to="/namespaces" @click="closeMenu()">Namespaces</router-link>
        <router-link to="/settings" @click="closeMenu()">
          <i-clarity-settings-line class="inline-block"></i-clarity-settings-line>
          <span>&nbsp;&nbsp;Options</span>
        </router-link>
      </div>
    </sw-mobile-menu>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { SwHeader, SwMobileMenu } from "@snowind/header";
import router from '@/router';

const isMenuVisible = ref(false);
const isHome = computed<boolean>(() => router.currentRoute.value.path == "/");
function closeMenu() {
  isMenuVisible.value = false;
}
</script>