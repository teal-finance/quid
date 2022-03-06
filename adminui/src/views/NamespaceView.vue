<template>
  <div class="text-3xl txt-primary">
    Namespaces
    <button
      class="ml-3 text-2xl border-none btn focus:outline-none txt-neutral"
      @click="collapse = !collapse"
    >
      <i class="fas fa-plus" v-if="collapse === true"></i>
      <i class="fas fa-minus" v-else></i>
    </button>
  </div>
  <div
    :class="{
      'slide-y': true,
      'slideup': collapse === true,
      'slidedown': collapse === false
    }"
    class="mt-4"
  >
    <div class="p-5 mt-3 border border-light dark:border-light-dark w-96">
      <div class="text-xl">Add a namespace</div>
      <add-namespace class="mt-5" @end="collapse = !collapse"></add-namespace>
    </div>
  </div>

  <sw-datatable
    :model="datatable"
    :renderers="{ publicEndpointEnabled: SwitchCellRenderer }"
    :sortable-cols="['name', 'publicEndpointEnabled']"
    class="table mt-3 table-zebra"
  ></sw-datatable>
  <!-- div v-for="ns in namespaces" :key="ns.id" class="mt-10">
    {{ ns }}
  </div-->
  <!-- button
    @click="mut()"
    v-if="datatable.state.rows.length > 0"
  >Mutate: {{ datatable.state.rows[0].publicEndpointEnabled }}</button-->
</template>
 
<script setup lang="ts">
import Namespace from "@/models/namespace";
import { onMounted, ref } from "vue";
import SwDatatable from "@/snowind/datatable/SwDatatable.vue";
import SwDatatableModel from "@/snowind/datatable/models/datatable";
import SwitchCellRenderer from "@/snowind/datatable/renderers/SwitchCellRenderer.vue";
import NamespaceContract from "@/models/namespace/contract";
import AddNamespace from "@/components/namespace/AddNamespace.vue";

const namespaces = ref(new Set<Namespace>([]));
const datatable = ref(new SwDatatableModel<Namespace>());
const collapse = ref(true);

async function fetchData() {
  namespaces.value = await Namespace.fetchAll();
  datatable.value = new SwDatatableModel<Namespace>({ idCol: "id", rows: Array.from(namespaces.value) });
  datatable.value.setColumnsFromData();
  console.log("DATA", namespaces.value)
}

function mut() {
  //datatable.rows.value[0].publicEndpointEnabled = !datatable.rows.value[0].publicEndpointEnabled;
}

onMounted(() => {
  fetchData();
});
</script>

<style lang="sass">
table
  & th
    @apply text-primary border-b p-3
  & td
    @apply px-5 py-1 border-b
</style>