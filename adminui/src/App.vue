<template>
  <div class="h-screen" :class="{ 'dark': user.isDarkMode.value === true }">
    <div v-if="user.isLoggedIn.value === true" class="flex flex-col h-full background">
      <the-topbar class="w-full"></the-topbar>
      <div class="flex flex-row h-full">
        <the-sidebar :sidebar="sidebar" @toggle="sidebar = !sidebar"></the-sidebar>
        <div class="px-5 pt-3 pb-8">
          <router-view />
        </div>
      </div>
    </div>
    <div v-else>
      <the-login></the-login>
    </div>
    <Toast position="top-right" group="main"></Toast>
  </div>
</template>

<script setup lang="ts">
import { onBeforeMount, ref } from "vue";
import TheSidebar from "@/components/TheSidebar.vue";
import { initState, user } from "@/state";
import TheLogin from "./components/TheLogin.vue";
import TheTopbar from "./components/TheTopbar.vue";
import Toast from 'primevue/toast';
import useNotify from "./notify";
import { useToast } from "primevue/usetoast";

const sidebar = ref(true);
const toast = useToast();

onBeforeMount(() => {
  initState(toast);
});
</script>

<style lang="sass">
body,
html
  margin: 0
  font-family: Arial, Helvetica, sans-serif
</style>
