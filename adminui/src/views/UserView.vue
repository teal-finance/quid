<template>
  <div class="text-3xl txt-primary">
    Users
    <button
      class="ml-3 text-2xl border-none btn focus:outline-none txt-neutral"
      @click="collapse = !collapse"
      v-if="!mustSelectNamespace"
    >
      <icon icon="fa6-solid:plus" v-if="collapse === true"></icon>
      <icon icon="fa6-solid:minus" v-else></icon>
    </button>
  </div>
  <div
    :class="{
      'slide-y': true,
      'slideup': collapse === true,
      'slidedown': collapse === false
    }"
    class="mb-4"
    v-if="!mustSelectNamespace"
  >
    <div class="p-5 mt-3 border border-light dark:border-light-dark w-96">
      <div class="text-xl">Add a user</div>
      <add-user class="mt-5" @end="endAdd()" v-if="collapse === false"></add-user>
    </div>
  </div>
  <div class="w-full" v-else>
    <div class="mt-3 text-2xl">Select a namespace</div>
    <namespace-selector class="mt-5" @selectns="fetchData()"></namespace-selector>
  </div>
  <user-datatable :users="users" v-if="!mustSelectNamespace" @reload="fetchData()"></user-datatable>
</template>

<script setup lang="ts">
import { Icon } from '@iconify/vue';
import { onMounted, ref } from "vue";
import { UserTable } from "@/models/user/interface";
import User from "@/models/user/user";
import UserDatatable from '@/components/user/UserDatatable.vue';
import AddUser from '@/components/user/AddUser.vue';
import { state, mustSelectNamespace } from '@/state';

const users = ref<Array<UserTable>>([]);
const collapse = ref(true);

async function fetchData() {
  users.value = await User.fetchAll(state.namespace.id);
}

function endAdd() {
  collapse.value = true;
  fetchData()
}

onMounted(() => {
  if (!mustSelectNamespace.value == true) {
    fetchData();
  }
})
</script>
