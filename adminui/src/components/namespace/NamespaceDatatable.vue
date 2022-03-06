<template>
  <div>
    <DataTable
      :value="namespaces"
      class="p-datatable-lg"
      id="nstable"
      v-model:expandedRows="expandedRows"
      data-key="id"
    >
      <Column field="id" header="Id"></Column>
      <Column field="name" header="Name"></Column>
      <Column field="publicEndpointEnabled" header="Public endpoint">
        <template #body="slotProps">
          <sw-switch
            label="Switch"
            v-model:value="slotProps.data.publicEndpointEnabled"
            class="w-max secondary"
          ></sw-switch>
        </template>
      </Column>
      <Column field="maxTokenTtl" header="Access token ttl"></Column>
      <Column field="maxRefreshTokenTtl" header="Refresh token ttl"></Column>
      <Column field="actions">
        <template #body="slotProps">
          <action-button
            @click="expand(slotProps.data.id)"
            v-if="expandedKey != slotProps.data.id"
          >Show info</action-button>
          <action-button @click="unexpand()" v-else>Hide info</action-button>
          <action-button
            class="ml-2"
            @click="showKey(slotProps.data.id, slotProps.data.name)"
            v-if="slotProps.data.name != 'quid'"
          >Show key</action-button>
          <action-button type="delete" class="ml-2" v-if="slotProps.data.name != 'quid'">Delete</action-button>
        </template>
      </Column>
      <template #expansion="slotProps">INFO</template>
    </DataTable>
    <Toast position="top-center">
      <template #message="slotProps">
        <div class="flex flex-col">
          <div>
            <div class="text-xl">Namespace {{ slotProps.message.summary }}</div>
            <div class="mt-3">{{ slotProps.message.detail }}</div>
          </div>
          <div class="flex flex-row mt-5 space-x-2">
            <button class="btn primary" @click="copyKey(slotProps.message.detail)">Copy</button>
            <button class="btn" @click="closeKey()">Close</button>
          </div>
        </div>
      </template>
    </Toast>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue';
import SwSwitch from "@snowind/switch";
import { useToast } from "primevue/usetoast";
import Toast from 'primevue/toast';
import DataTable from 'primevue/datatable';
import Column from 'primevue/column';
import Namespace from "@/models/namespace";
import NamespaceTable from '@/models/namespace/table';
import ActionButton from '../widgets/ActionButton.vue';

const namespaces = ref(new Array<NamespaceTable>());
const expandedRows = ref<any>([]);
const expandedKey = ref(0);
const toast = useToast();

async function fetchData() {
  const ns = await Namespace.fetchAll();
  namespaces.value = Array.from(ns);
  //console.log("DATA", namespaces.value)
}

function expand(id: number) {
  expandedRows.value = namespaces.value.filter((p) => p.id == id);
  expandedKey.value = id;
}

function unexpand() {
  expandedRows.value = [];
  expandedKey.value = 0;
}

function closeKey() {
  toast.removeAllGroups();
}

function copyKey(k: string) {
  navigator.clipboard.writeText(k);
  toast.removeAllGroups();
}

async function showKey(id: number, name: string) {
  const key = await Namespace.getKey(id);
  console.log("K", key)
  toast.add({ severity: 'info', summary: name, detail: key });
}

onMounted(() => {
  fetchData();
});
</script>

<style lang="sass">
#nstable
  & td
    @apply p-3
.ui-toast-message
  width: 96em
</style>