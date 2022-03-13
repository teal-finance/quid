<template>
  <div>
    <DataTable :value="groups" class="main-table" v-model:expandedRows="expandedRows" data-key="id">
      <Column field="id" header="Id"></Column>
      <Column field="name" header="Name"></Column>
    </DataTable>
  </div>
</template>

<script setup lang="ts">
import { GroupTable } from '@/models/group/interface';
import DataTable from "primevue/datatable";
import Column from "primevue/column";
import { ref } from 'vue';

const expandedRows = ref<any>([]);

const props = defineProps({
  groups: {
    type: Array as () => Array<GroupTable>,
    required: true,
  }
});
const emit = defineEmits(["reload"])
</script>