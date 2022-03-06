<template>
  <div class="flex p-5 text-xl font-bold bg-primary dark:bg-sgray-800">
    <div class="flex-grow text-white">
      <router-link to="/">Quid</router-link>
    </div>
    <div class="text-gray-400 cursor-pointer dark:text-gray-400" @click="logout()">
      <i class="mr-2 align-middle fa-lg fas fa-sign-out-alt"></i>
      <span v-html="user.name.value"></span>
    </div>
    <button
      class="text-gray-400 border-none cursor-pointer dark:text-gray-400 focus:outline-none btn"
      @click="user.toggleDarkMode()"
    >
      <i :class="modeIcon"></i>
    </button>
  </div>
</template>

<script lang="ts">
import { computed, defineComponent } from "vue";
import { user } from "@/state";

export default defineComponent({
  setup() {
    const modeIcon = computed<string>(() => {
      if (user.isDarkMode.value === true) {
        return "fas fa-sun text-xl";
      }
      return "fas fa-moon text-xl";
    });

    function logout() {
      console.log("Logout");
    }

    return {
      user,
      logout,
      modeIcon
    };
  },
});
</script>