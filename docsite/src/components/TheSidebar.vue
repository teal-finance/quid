<template>
  <div>
    <apidocs-list class="mt-3"></apidocs-list>
    <div class="mt-3 text-xl">Server</div>
    <server-list class="mt-3"></server-list>
    <div class="mt-3 text-xl">Client</div>
    <client-list class="mt-3"></client-list>
  </div>
</template>

<script setup lang="ts">
import ApidocsList from "@/components/ApidocsList.vue";
import ServerList from "./ServerList.vue";
import ClientList from "./ClientList.vue";
</script>