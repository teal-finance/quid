<template>
  <DataTable :value="orgs" class="p-datatable main-table" data-key="id">
    <Column field="id" header="Id"></Column>
    <Column field="name" header="Name"></Column>
    <Column field="actions">
      <template #body="slotProps">
        <action-button type="delete" class="ml-2" @click="confirmDelete(slotProps.data)">Delete</action-button>
      </template>
    </Column>
  </DataTable>
</template>

<script setup lang="ts">
import DataTable from "primevue/datatable";
import Column from "primevue/column";
import { OrgTable } from "@/models/org/interface";
import ActionButton from "@/components/widgets/ActionButton.vue";
import { notify } from "@/state";
import Org from "@/models/org";

defineProps({
  orgs: {
    type: Array as () => Array<OrgTable>,
    required: true,
  }
});

const emit = defineEmits(["reload"]);

function confirmDelete(ns: OrgTable) {
  notify.confirmDelete(
    `Delete the ${ns.name} org?`,
    () => {
      Org.delete(ns.id).then(() => {
        notify.done("Org deleted");
        emit("reload");
      })
    }
  )
}
</script>