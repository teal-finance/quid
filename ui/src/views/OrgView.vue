<template>
  <div class="text-3xl txt-primary dark:txt-light">
    Orgs
    <button id="add-org" class="ml-3 text-2xl border-none btn focus:outline-none txt-neutral"
      @click="collapse = !collapse">
      <icon icon="fa6-solid:plus" v-if="collapse === true"></icon>
      <icon icon="fa6-solid:minus" v-else></icon>
    </button>
  </div>
  <div :class="{
    'slide-y': true,
    'slideup': collapse === true,
    'slidedown': collapse === false
  }" class="mb-8">
    <div class="p-5 mt-3 border bord-lighter w-96">
      <div class="text-xl">Add an org</div>
      <add-org class="mt-5" @end="endAddOrg()"></add-org>
    </div>
  </div>
  <org-datatable :orgs="orgs" @reload="fetchData()"></org-datatable>
</template>

<script setup lang="ts">
import Org from "@/models/org";
import { Icon } from '@iconify/vue';
import { onMounted, ref } from "vue";
import OrgDatatable from "@/components/org/OrgDatatable.vue";
import AddOrg from "@/components/org/AddOrg.vue";
import { OrgTable } from "@/models/org/interface";

const orgs = ref<Array<OrgTable>>([]);
const collapse = ref(true);

async function fetchData() {
  console.log("FETCH ORGS")
  orgs.value = await Org.fetchAll();
}

function endAddOrg() {
  collapse.value = true;
  fetchData()
}

onMounted(() => fetchData())
</script>
