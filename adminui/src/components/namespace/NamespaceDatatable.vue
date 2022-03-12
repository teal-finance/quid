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
            @change="togglePublicEndpoint(slotProps.data.id, Boolean($event))"
            v-if="slotProps.data.name != 'quid'"
          ></sw-switch>
        </template>
      </Column>
      <Column field="maxTokenTtl" header="Access token ttl">
        <template #body="slotProps">
          <edit-token-ttl
            v-if="slotProps.data.name != 'quid'"
            :id="slotProps.data.id"
            :ttl="slotProps.data.maxTokenTtl"
            token-type="access"
            @end="slotProps.data.maxTokenTtl = $event"
          ></edit-token-ttl>
          <span v-else class="ml-6" v-html="slotProps.data.maxTokenTtl"></span>
        </template>
      </Column>
      <Column field="maxRefreshTokenTtl" header="Refresh token ttl">
        <template #body="slotProps">
          <edit-token-ttl
            v-if="slotProps.data.name != 'quid'"
            :id="slotProps.data.id"
            :ttl="slotProps.data.maxRefreshTokenTtl"
            token-type="refresh"
            @end="slotProps.data.maxRefreshTokenTtl = $event"
          ></edit-token-ttl>
          <span v-else class="ml-6" v-html="slotProps.data.maxRefreshTokenTtl"></span>
        </template>
      </Column>
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
          <action-button
            type="delete"
            class="ml-2"
            v-if="slotProps.data.name != 'quid'"
            @click="confirmDelete(slotProps.data)"
          >Delete</action-button>
        </template>
      </Column>
      <template #expansion>
        <namespace-info :num-users="nsInfo.numUsers" :groups="nsInfo.groups"></namespace-info>
      </template>
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
    <ConfirmDialog></ConfirmDialog>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from "vue";
import SwSwitch from "@snowind/switch";
import ConfirmDialog from 'primevue/confirmdialog';
import { useToast } from "primevue/usetoast";
import Toast from "primevue/toast";
import DataTable from "primevue/datatable";
import Column from "primevue/column";
import Namespace from "@/models/namespace";
import NamespaceTable from "@/models/namespace/table";
import NamespaceInfo from "./NamespaceInfo.vue";
import { useConfirm } from "primevue/useconfirm";
import { notify } from "@/state";
import Group from "@/models/group";
import ActionButton from "../widgets/ActionButton.vue";
import EditTokenTtl from "./EditTokenTtl.vue";

const props = defineProps({
  namespaces: {
    type: Array as () => Array<NamespaceTable>,
    required: true,
  }
});
const emit = defineEmits(["reload"])

const expandedRows = ref<any>([]);
const expandedKey = ref(0);
const toast = useToast();
const nsInfo = reactive({
  numUsers: 0,
  groups: new Array<Group>(),
});
const confirm = useConfirm();

function confirmDelete(ns: NamespaceTable) {
  confirm.require({
    message: `Delete the ${ns.name} namespace?`,
    header: 'Delete Confirmation',
    icon: 'pi pi-info-circle',
    acceptClass: 'p-button-danger',
    accept: () => {
      console.log("ACCEPT")
      Namespace.delete(ns.id).then(() => {
        notify.done("Namespace deleted");
        emit("reload");
      })
    },
    reject: () => { }
  });
}

async function expand(id: number) {
  const info = await Namespace.fetchRowInfo(id);
  nsInfo.numUsers = info.numUsers;
  nsInfo.groups = info.groups;
  expandedRows.value = props.namespaces.filter((p) => p.id == id);
  expandedKey.value = id;
}

function unexpand() {
  expandedRows.value = [];
  expandedKey.value = 0;
  nsInfo.numUsers = 0;
  nsInfo.groups = new Array<Group>()
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
  console.log("K", key);
  toast.add({ severity: "info", summary: name, detail: key });
}

async function togglePublicEndpoint(id: number, enabled: boolean) {
  await Namespace.togglePublicEndpoint(id, enabled);
  notify.done("Endpoint toggled");
}
</script>

<style lang="sass">
#nstable
  & td
    @apply p-3
.ui-toast-message
  width: 96em
</style>