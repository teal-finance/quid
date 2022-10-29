<template>
  <div>
    <sw-header class="fixed top-0 left-0 z-20 w-full h-16 topbar" @togglemenu="isMenuVisible = !isMenuVisible">
      <template #mobile-back>
        <i-ion-arrow-back-outline class="inline-flex ml-2 text-3xl" v-if="!isHome"></i-ion-arrow-back-outline>
      </template>
      <template #mobile-branding>
        <div class="inline-flex flex-row items-center h-full pt-1 ml-2 text-2xl truncate" @click="$router.push('/')">
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
          <the-current-namespace class="pr-3" v-if="user.type.value == 'serverAdmin'"></the-current-namespace>
          <div class="pr-5 text-lg cursor-pointer txt-lighter dark:txt-light" @click="user.toggleDarkMode()">
            <i-fa-solid:moon v-if="!user.isDarkMode.value"></i-fa-solid:moon>
            <i-fa-solid:sun v-else></i-fa-solid:sun>
          </div>
          <div class="pr-8 text-lg cursor-pointer">
            <i-mdi:logout v-if="user.isLoggedIn" @click="logout()"></i-mdi:logout>
          </div>
        </div>
      </template>
    </sw-header>
    <sw-mobile-menu :is-visible="isMenuVisible">
      <div class="flex flex-col p-3 pt-16 space-y-5">
        <router-link to="/namespaces" @click="closeMenu()">Namespaces</router-link>
        <div>
          <i-mdi:logout v-if="user.isLoggedIn" @click="logout()"></i-mdi:logout>
          <span>&nbsp;&nbsp;Logout</span>
        </div>
      </div>
    </sw-mobile-menu>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { SwHeader, SwMobileMenu } from "@snowind/header";
import router from '@/router';
import { user, notify } from "@/state";
import TheCurrentNamespace from './namespace/TheCurrentNamespace.vue';
import { api } from '@/api';

const isMenuVisible = ref(false);

const isHome = computed<boolean>(() => router.currentRoute.value.path == "/");

function closeMenu() {
  isMenuVisible.value = false;
}

async function logout() {

  notify.confirmDelete("Logout from Quid admin?", async () => {
    user.isLoggedIn.value = false;
    await api.get("/logout");
  }, () => null, "Disconnect")
}
</script>