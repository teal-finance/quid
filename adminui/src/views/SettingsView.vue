<template>
  <div class="text-3xl txt-primary">Settings</div>
  <sw-switch
    @change="user.toggleDarkMode()"
    :checked="user.isDarkMode.value"
    class="mt-8 ml-2 secondary w-max"
  >
    <div class="ml-2">Dark mode</div>
  </sw-switch>
  <div class="mt-8">
    <div
      v-if="user.isLoggedIn.value === true"
      class="align-middle cursor-pointer"
      @click="logout()"
    >
      <i-si-glyph:sign-out class="inline-block ml-3 text-2xl"></i-si-glyph:sign-out>
      &nbsp;&nbsp;Logout {{ user.name.value }}
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