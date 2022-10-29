<template>
  <div class="text-3xl txt-primary dark:txt-light">
    Admins
    <button id="add-admin" class="ml-3 text-2xl border-none btn focus:outline-none txt-light"
      @click="collapse = !collapse" v-if="!user.mustSelectNamespace">
      <icon icon="fa6-solid:plus" v-if="collapse === true"></icon>
      <icon icon="fa6-solid:minus" v-else></icon>
    </button>
  </div>
  <div :class="{
    'slide-y': true,
    'slideup': collapse === true,
    'slidedown': collapse === false
  }" class="w-full mb-8" v-if="!user.mustSelectNamespace">
    <div class="w-1/2 p-5 mt-3 border bord-lighter">
      <add-admin class="mt-3" @end="endAdd()"></add-admin>
    </div>
  </div>
  <admin-datatable v-if="!user.mustSelectNamespace" class="mt-5" :users="users" @reload="fetchData()"></admin-datatable>
  <div class="w-full" v-else>
    <div class="mt-3 text-2xl">Select a namespace</div>
    <namespace-selector class="mt-5" @selectns="fetchData()"></namespace-selector>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue";
import { Icon } from '@iconify/vue';
import AddAdmin from "@/components/admin/add/AddAdmin.vue";
import { user } from "@/state";
import NamespaceSelector from "@/components/namespace/NamespaceSelector.vue";
import AdminDatatable from "@/components/admin/AdminDatatable.vue";
import AdminUser from "@/models/adminuser";
import { AdminUserTable } from "@/models/adminuser/interface";

const users = ref<Array<AdminUserTable>>([]);
const collapse = ref(true);

async function fetchData() {
  console.log("Fetching data")
  users.value = await AdminUser.fetchAll(user.namespace.value.id);
}

function endAdd() {
  collapse.value = true;
  fetchData()
}

onMounted(() => {
  if (!user.mustSelectNamespace == true) {
    fetchData();
  }
})
</script>