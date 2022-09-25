<template>
  <the-header :lib-name="libName" :links="links"></the-header>
  <div class="container h-full min-h-screen pt-16 pb-8 background sm:pl-48">
    <router-view></router-view>
  </div>
  <the-sidebar class="fixed left-0 hidden w-48 h-screen p-3 top-16 sm:block secondary"></the-sidebar>
</template>

<script setup lang="ts">
import { onBeforeMount } from "vue";
import TheHeader from "./components/TheHeader.vue";
import TheSidebar from "./components/TheSidebar.vue";
import { libName, links } from "@/conf";
import { initState } from "./state";

onBeforeMount(() => initState())
</script>

