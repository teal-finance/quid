<template>
  <div class="h-screen" :class="{ 'dark': user.isDarkMode.value === true }">
    <div
      v-if="user.isLoggedIn.value === true"
      class="flex flex-col h-full bg-background dark:bg-background-dark text-foreground dark:text-foreground-dark"
    >
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
  </div>
</template>

<script lang="ts">
import { defineComponent, onBeforeMount, ref } from "vue";
import TheSidebar from "@/components/TheSidebar.vue";
import { initState, user } from "@/state";
import TheLogin from "./components/TheLogin.vue";
import TheTopbar from "./components/TheTopbar.vue";

export default defineComponent({
  components: {
    TheSidebar,
    TheLogin,
    TheTopbar,
  },
  setup() {
    const sidebar = ref(true);

    onBeforeMount(() => {
      initState();
    });

    return {
      sidebar,
      user,
    };
  },
});
</script>

<style lang="sass">
body,
html
  margin: 0
  font-family: Arial, Helvetica, sans-serif
</style>
