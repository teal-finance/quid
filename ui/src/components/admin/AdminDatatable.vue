<template>
  <div>
    <DataTable
      :value="users"
      class="main-table p-datatable-sm"
      v-model:expandedRows="expandedRows"
      data-key="id"
    >
      <Column field="id" header="Id"></Column>
      <Column field="userName" header="Name"></Column>
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
import { notify, user } from '@/state';
import { UserTable } from '@/models/user/interface';
import AdminUser from "@/models/adminuser";
import { AdminUserTable } from "@/models/adminuser/interface";

const expandedRows = ref<any>([]);

defineProps({
  users: {
    type: Array as () => Array<AdminUserTable>,
    required: true,
  }
});
const emit = defineEmits(["reload"]);

function confirmDelete(row: AdminUserTable) {
  notify.confirmDelete(
    `Delete the ${row.userName} admin user?`,
    () => {
      AdminUser.delete(row.userId, user.namespace.value.id).then(() => {
        notify.done("User deleted");
        emit("reload");
      })
    }
  )
}
</script>