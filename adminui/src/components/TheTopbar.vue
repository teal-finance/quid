<template>
  <div>
    <sw-header class="fixed z-20 w-full h-16 topbar" @togglemenu="isMenuVisible = !isMenuVisible">
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
        <div class="flex flex-row items-center justify-end w-full h-full">
          <the-current-namespace class="pr-3"></the-current-namespace>
          <div
            class="pr-5 text-lg cursor-pointer txt-lighter dark:txt-light"
            @click="user.toggleDarkMode()"
          >
            <i-fa-solid:moon v-if="!user.isDarkMode.value"></i-fa-solid:moon>
            <i-fa-solid:sun v-else></i-fa-solid:sun>
          </div>
          <div class="pr-8 text-lg cursor-pointer">
            <i-fluent-settings-32-regular class="txt-lighter dark:txt-light"></i-fluent-settings-32-regular>
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
import { user } from "@/state";
import TheCurrentNamespace from './namespace/TheCurrentNamespace.vue';

const isMenuVisible = ref(false);
const isHome = computed<boolean>(() => router.currentRoute.value.path == "/");
function closeMenu() {
  isMenuVisible.value = false;
}
</script>