<template>
  <div>
    <DataTable :value="users" class="main-table p-datatable-sm" v-model:expandedRows="expandedRows" data-key="id"
      removableSort>
      <Column body-class="col-id" field="id" header="Id"></Column>
      <Column body-class="col-name" field="name" header="Name" :sortable="true"></Column>
      <Column body-class="col-actions" field="actions">
        <template #body="slotProps">
          <action-button @click="expand(slotProps.data.id)" v-if="expandedKey != slotProps.data.id">Show groups
          </action-button>
          <action-button @click="unexpand()" v-else>Hide groups&nbsp;</action-button>
          <action-button type="delete" class="ml-2" @click="confirmDelete(slotProps.data)">Delete</action-button>
        </template>
      </Column>
      <template #expansion="slotProps">
        <div class="p-3 pb-8 ml-5">
          <user-groups-info :user="slotProps.data" :is-loading="expandedIsLoading" :groups="state.userGroups"
            @user-removed="userRemovedFromGroup($event, slotProps.data)"></user-groups-info>
          <add-user-into-group :user="slotProps.data" @addingroup="userAddedIntoGroup($event, slotProps.data)"
            :user-groups="state.userGroups"></add-user-into-group>
        </div>
      </template>
    </DataTable>
  </div>
</template>

<script setup lang="ts">
import DataTable from "primevue/datatable";
import ActionButton from "../widgets/ActionButton.vue";
import Column from "primevue/column";
import { reactive, ref } from 'vue';
import { notify } from '@/state';
import { UserTable } from '@/models/user/interface';
import User from "@/models/user/user";
import UserGroupsInfo from "./UserGroupsInfo.vue";
import AddUserIntoGroup from "./AddUserIntoGroup.vue";
import { GroupTable } from "@/models/group/interface";
import Group from "@/models/group";

const expandedRows = ref<any>([]);
const expandedKey = ref(0);

const props = defineProps({
  users: {
    type: Array as () => Array<UserTable>,
    required: true,
  }
});
const emit = defineEmits(["reload"]);
const state = reactive({
  userGroups: new Array<GroupTable>()
})
const expandedIsLoading = ref(false);

async function userAddedIntoGroup(g: GroupTable, u: UserTable) {
  await Group.addUserToGroup(u.id, g.id)
  state.userGroups.push(g);
}

async function userRemovedFromGroup(g: GroupTable, u: UserTable) {
  state.userGroups = state.userGroups.filter((v) => {
    if (g.id != v.id) {
      return v
    }
  });
}

function expand(id: number) {
  expandedIsLoading.value = true;
  Group.fetchUserGroups(id).then((g) => {
    expandedIsLoading.value = false;
    state.userGroups.push(...g);
  })
  expandedRows.value = props.users.filter((p) => p.id == id);
  expandedKey.value = id;
}

function unexpand() {
  state.userGroups = [];
  expandedKey.value = 0;
  expandedRows.value = [];
}

function confirmDelete(row: UserTable) {
  notify.confirmDelete(
    `Delete the ${row.name} user?`,
    () => {
      User.delete(row.id).then(() => {
        notify.done("User deleted");
        emit("reload");
      })
    }
  )
}
</script>