<template>
  <div>
    <DataTable :value="users" class="main-table" v-model:expandedRows="expandedRows" data-key="id">
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
import DataTable from "primevue/datatable";
import Column from "primevue/column";
import { ref } from 'vue';
import { notify } from '@/state';
import { UserTable } from '@/models/user/interface';
import User from "@/models/user/user";

const expandedRows = ref<any>([]);

defineProps({
  users: {
    type: Array as () => Array<UserTable>,
    required: true,
  }
});
const emit = defineEmits(["reload"]);

function confirmDelete(row: UserTable) {
  notify.confirmDelete(
    `Delete the ${row.name} admin user?`,
    () => {
      User.delete(row.id).then(() => {
        notify.done("User deleted");
        emit("reload");
      })
    }
  )
}
</script>