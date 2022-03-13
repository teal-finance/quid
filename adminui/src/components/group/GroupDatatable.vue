<template>
  <div>
    <DataTable :value="groups" class="main-table" v-model:expandedRows="expandedRows" data-key="id">
      <Column field="id" header="Id"></Column>
      <Column field="name" header="Name"></Column>
      <Column field="actions">
        <template #body="slotProps">
          <action-button type="delete" class="ml-2" @click="confirmDelete(slotProps.data)">Delete</action-button>
        </template>
      </Column>
    </DataTable>
  </div>
</template>

<script setup lang="ts">
import { GroupTable } from '@/models/group/interface';
import DataTable from "primevue/datatable";
import Column from "primevue/column";
import { ref } from 'vue';
import { notify } from '@/state';
import Group from '@/models/group';

const expandedRows = ref<any>([]);

const props = defineProps({
  groups: {
    type: Array as () => Array<GroupTable>,
    required: true,
  }
});
const emit = defineEmits(["reload"]);

function confirmDelete(g: GroupTable) {
  notify.confirmDelete(
    `Delete the ${g.name} group?`,
    () => {
      Group.delete(g.id).then(() => {
        notify.done("Group deleted");
        emit("reload");
      })
    }
  )
}
</script>