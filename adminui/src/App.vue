<template>
  <div class="h-screen" :class="{ 'dark': user.isDarkMode.value === true }">
    <div v-if="user.isLoggedIn.value === true" class="h-full overflow-hidden background">
      <the-topbar class="w-full"></the-topbar>
      <div class="absolute flex flex-row h-full">
        <the-sidebar class="fixed pt-16" :sidebar="isSidebarOpened" @toggle="toggleSidebar()"></the-sidebar>
        <div
          class="w-full px-5 pt-16 pb-8 overflow-auto slide-main"
          :class="isSidebarOpened ? 'main-opened' : 'main-closed'"
        >
          <div class="p-3">
            <router-view />
          </div>
        </div>
      </div>
    </div>
    <div v-else>
      <the-login></the-login>
    </div>
    <Toast position="top-right" group="main"></Toast>
    <ConfirmDialog></ConfirmDialog>
    <sw-toast :class="toastType" :show="isToastVisible">{{ toastMsg }}</sw-toast>
  </div>
</template>

<script setup lang="ts">
import { onBeforeMount, ref, nextTick } from "vue";
import ConfirmDialog from 'primevue/confirmdialog';
import TheSidebar from "@/components/TheSidebar.vue";
import { initState, user } from "@/state";
import TheLogin from "./components/TheLogin.vue";
import TheTopbar from "./components/TheTopbar.vue";
import Toast from 'primevue/toast';
import { useToast } from "primevue/usetoast";
import { useConfirm } from "primevue/useconfirm";
import { SwToast, useToast as swToast } from "@snowind/toast";
import { PopToast, ColorVariant } from "./type";

const toast = useToast();
const confirm = useConfirm();
//const { popToast, isToastVisible } = swToast();
const toastMsg = ref("");
const toastType = ref<ColorVariant>("secondary");
const isToastVisible = ref(false);

const isSidebarOpened = ref(false);

function toggleSidebar() {
  isSidebarOpened.value = !isSidebarOpened.value;
}

function genPopSwToast(): PopToast {
  return async function (msg: string, type: ColorVariant, delay?: number | undefined) {
    console.log("Run pop")
    toastType.value = type;
    toastMsg.value = msg;
    console.log("End init pop")
    try {
      isToastVisible.value = true;
      setTimeout(() => { isToastVisible.value = false; }, delay);
      console.log("wait")
      console.log("endwait")
      isToastVisible.value = false;
      console.log("end")
    } catch (e) {
      console.log("Error poping toast", e)
    }
  }
}

onBeforeMount(() => initState(toast, confirm, genPopSwToast()));
</script>

<style lang="sass">

#sidebar
  @apply w-20
  &.opened
    width: 16rem
.main-opened
  padding-left: 16rem !important
.main-closed
  padding-left: 5rem !important
.dark
  .p-datatable
    & th
      background: transparent !important
    & thead
      & [role="cell"]
        @apply border-b border-neutral
    & tbody
      & tr:not(:last-child)
        & [role="cell"]
          @apply border-b border-neutral
    & tbody
      & tr:last-child
        & [role="cell"]
          border: 0
.main-table
  & td
    border-width: 0
    @apply background px-5 py-1
.slide-main
  @apply overflow-x-hidden transition-all duration-300
</style>
