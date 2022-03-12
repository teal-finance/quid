<template>
  <div class="text-xl text-gray-500 w-max dark:text-gray-300">User {{ user.name }}</div>
  <sw-switch
    @change="user.toggleDarkMode()"
    :checked="user.isDarkMode.value"
    class="mt-5 secondary w-max"
  >
    <div class="ml-2">Dark mode</div>
  </sw-switch>
  <div class="mt-5">
    <div
      v-if="user.isLoggedIn.value === true"
      class="align-middle cursor-pointer"
      @click="logout()"
    >
      <i-si-glyph:sign-out class="inline-block ml-3"></i-si-glyph:sign-out>&nbsp;&nbsp;Logout
    </div>
  </div>
</template>

<script setup lang="ts">
import SwSwitch from "@snowind/switch";
import { user } from "@/state";
import { requests } from "@/api";

async function logout() {
  //await requests.logout();
}
</script>