<template>
  <div>
    <h1 class="text-muted mt-3">
      Users
      &nbsp;
      <b-icon-plus
        v-if="action!=='addUser'"
        class="mr-1"
        style="color:lightgrey"
        @click="$store.commit('action', 'addUser')"
      />
    </h1>
    <loading-indicator v-if="state.isLoading"></loading-indicator>
    <div>
      <b-collapse id="collapse-4" v-model="showActionBar" class="mt-2">
        <users-add v-if="action === 'addUser'" @refresh="refresh"></users-add>
      </b-collapse>
    </div>
    <b-table hover bordeless :items="data" :fields="fields" class="mt-4" style="max-width:650px">
      <template v-slot:cell(action)="row">
        <b-button
          class="mr-2"
          variant="outline-secondary"
          @click="toggleDetails(row)"
        >{{ row.detailsShowing ? 'Hide' : 'Show'}} groups</b-button>
        <b-button
          variant="outline-danger"
          v-if="row.item.username !== username"
          @click="confirmDeleteItem(row.item.id, row.item.username)"
        >Delete</b-button>
      </template>
      <template v-slot:row-details="row">
        <b-card v-if="users[row.index] !== undefined">
          <users-manage-groups
            :user="users[row.index]"
            @user-added-in-group="userAddedToGroup(row, $event)"
          ></users-manage-groups>
        </b-card>
      </template>
    </b-table>
    <b-modal title="Delete user" ref="delete-modal">
      Delete {{ itemToDelete.username }}?
      <template v-slot:modal-footer="{ ok, cancel }">
        <b-button variant="danger" @click="deleteItem(itemToDelete)">Delete</b-button>
        <b-button variant="warning" @click="cancel()">Cancel</b-button>
      </template>
    </b-modal>
  </div>
</template>

<script>
import { mapState, mapGetters } from "vuex";
import LoadingIndicator from "@/components/LoadingIndicator";
import UsersAdd from "@/components/users/UsersAdd";
import UsersManageGroups from "@/components/users/UsersManageGroups";

export default {
  components: {
    LoadingIndicator,
    UsersAdd,
    UsersManageGroups,
  },
  data() {
    return {
      data: [],
      state: {
        isLoading: false,
      },
      fields: [
        { key: "id", sortable: true },
        { key: "username", sortable: true },
        { key: "namespace", sortable: true },
        { key: "action", sortable: false },
      ],
      itemToDelete: {},
      users: {},
    };
  },
  methods: {
    async userAddedToGroup(row, group) {
      this.users[row.index].groups.add(group);
      row.toggleDetails();
      await new Promise((r) => setTimeout(r, 10));
      row.toggleDetails();
    },
    async fetchUserDetails(row) {
      let { response } = await this.$api.post("/admin/users/info", {
        id: row.item.id,
      });
      let user = row.item;
      user.index = row.index;
      user.groups = new Set(response.data.groups);
      this.users[row.index] = user;
    },
    async toggleDetails(row) {
      if (!row.detailsShowing) {
        await this.fetchUserDetails(row);
      }
      row.toggleDetails();
    },
    refresh() {
      this.fetchUsers();
      this.$bvToast.toast("ok", {
        title: "User saved",
        variant: "success",
        autoHideDelay: 1500,
      });
    },
    async fetchUsers() {
      let { response, error } = await this.$api.get("/admin/users/all");
      if (!error) {
        this.data = response.data;
      }
    },
    confirmDeleteItem(id, username) {
      this.itemToDelete = {
        username: username,
        id: id,
      };
      this.$refs["delete-modal"].show();
    },
    async deleteItem(ns) {
      this.$refs["delete-modal"].hide();
      let { error } = await this.$api.post("/admin/users/delete", {
        id: ns.id,
      });
      if (error === null) {
        this.$bvToast.toast("Ok", {
          title: "User deleted",
          autoHideDelay: 1000,
          variant: "success",
        });
        this.fetchUsers();
      }
    },
  },
  mounted: function () {
    this.fetchUsers();
  },
  computed: {
    ...mapState(["action", "username"]),
    ...mapGetters({
      s: "showActionBar",
    }),
    showActionBar: {
      get() {
        return this.s;
      },
      set(newName) {
        return newName;
      },
    },
  },
  //computed: mapState(["action", "showActionBar"])
};
</script>