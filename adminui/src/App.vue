<template>
  <div class="w-screen h-screen" :class="{ 'dark': user.isDarkMode.value === true }">
    <div v-if="user.isLoggedIn.value === true" class="h-full overflow-hidden background">
      <the-topbar class="w-full"></the-topbar>
      <div class="absolute flex flex-row w-full h-full">
        <the-sidebar
          class="fixed pt-16 sidebar border-b"
          :sidebar="isSidebarOpened"
          @toggle="toggleSidebar()"
          v-if="!isMobile"
        ></the-sidebar>
        <div
          class="w-full px-5 pt-16 pb-8 overflow-auto slide-main container mx-auto"
          :class="mainCls"
        >
          <div class="w-full p-3">
            <router-view />
          </div>
        </div>
      </div>
    </div>
    <div v-else>
      <the-login></the-login>
    </div>
    <Toast position="top-right" group="main"></Toast>
    <Toast position="bottom-right" group="bottom-right"></Toast>
    <ConfirmDialog></ConfirmDialog>
  </div>
</template>

<script setup lang="ts">
import { onBeforeMount, ref } from "vue";
import ConfirmDialog from 'primevue/confirmdialog';
import TheSidebar from "@/components/TheSidebar.vue";
import { initState, user, isMobile } from "@/state";
import TheLogin from "./components/TheLogin.vue";
import TheTopbar from "./components/TheTopbar.vue";
import Toast from 'primevue/toast';
import { useToast } from "primevue/usetoast";
import { useConfirm } from "primevue/useconfirm";
import { computed } from "@vue/reactivity";

const toast = useToast();
const confirm = useConfirm();

const isSidebarOpened = ref(false);

const mainCls = computed<Array<string>>(() => {
  const cls = new Array<string>();
  if (!isMobile.value) {
    if (isSidebarOpened.value) {
      cls.push("main-opened")
    } else {
      cls.push("main-closed")
    }
  }
  return cls
})

function toggleSidebar() {
  isSidebarOpened.value = !isSidebarOpened.value;
}

onBeforeMount(() => initState(toast, confirm));
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
        border-color: #404040
        border-bottom-width: 1px
    & tbody
      & tr:not(:last-child)
        & [role="cell"]
          border-color: #404040
          border-bottom-width: 1px
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
