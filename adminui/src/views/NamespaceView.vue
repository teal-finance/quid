<template>
  <div class="text-3xl txt-primary">
    Namespaces
    <button
      class="ml-3 text-2xl border-none btn focus:outline-none txt-neutral"
      @click="collapse = !collapse"
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
    class="my-4"
  >
    <div class="p-5 mt-3 border border-light dark:border-light-dark w-96">
      <div class="text-xl">Add a namespace</div>
      <add-namespace class="mt-5" @end="endAddNamespace()"></add-namespace>
    </div>
  </div>
  <namespace-datatable :namespaces="namespaces" @reload="fetchNamespaces()"></namespace-datatable>
</template>
 
<script setup lang="ts">
import { onMounted, ref } from "vue";
import { Icon } from '@iconify/vue';
import AddNamespace from "@/components/namespace/AddNamespace.vue";
import NamespaceDatatable from "@/components/namespace/NamespaceDatatable.vue";
import NamespaceTable from "@/models/namespace/table";
import Namespace from "@/models/namespace";

const collapse = ref(true);
const namespaces = ref(new Array<NamespaceTable>());

function endAddNamespace() {
  collapse.value = !collapse.value;
  fetchNamespaces()
}

async function fetchNamespaces() {
  const ns = await Namespace.fetchAll();
  namespaces.value = Array.from(ns);
  //console.log("DATA", namespaces.value)
}

onMounted(() => fetchNamespaces())
</script>

<style lang="sass">
table
  & th
    @apply text-primary border-b p-3
  & td
    @apply px-5 py-1 border-b
</style>