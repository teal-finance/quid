<template>
  <div>
    <DataTable :value="groups" class="main-table p-datatable-sm" v-model:expandedRows="expandedRows" data-key="id">
      <Column body-class="col-id" field="id" header="Id"></Column>
      <Column body-class="col-name" field="name" header="Name"></Column>
      <Column body-class="col-actions" field="actions">
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
import ActionButton from '../widgets/ActionButton.vue';

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