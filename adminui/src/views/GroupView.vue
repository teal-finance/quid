<template>
  <div class="text-3xl txt-primary">
    Groups
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
      <div class="text-xl">Add a group</div>
      <add-group class="mt-5" @end="endAdd()"></add-group>
    </div>
  </div>
  <div class="w-full" v-else>
    <div class="mt-3 text-2xl">Select a namespace</div>
    <namespace-selector class="mt-5" @selectns="fetchGroups()"></namespace-selector>
  </div>
  <group-datatable :groups="groups" v-if="!mustSelectNamespace" @reload="fetchGroups()"></group-datatable>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue";
import { Icon } from '@iconify/vue';

import { mustSelectNamespace, state } from "@/state";
import NamespaceSelector from "@/components/namespace/NamespaceSelector.vue";
import Group from "@/models/group";
import { GroupTable } from "@/models/group/interface";
import GroupDatatable from "@/components/group/GroupDatatable.vue";
import AddGroup from "@/components/group/AddGroup.vue";

const collapse = ref(true);
const groups = ref<Array<GroupTable>>([]);

function endAdd() {
  collapse.value = true;
  fetchGroups()
}

async function fetchGroups() {
  const g = await Group.fetchAll(state.namespace.id);
  groups.value = Array.from(g);
}

onMounted(() => {
  if (!mustSelectNamespace.value == true) {
    fetchGroups();
  }
})
</script>