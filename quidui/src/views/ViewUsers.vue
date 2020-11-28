<template>
  <div>
    <h1 class="text-muted mt-3">
      Users &nbsp;
      <b-icon-plus
        v-if="action !== 'addUser'"
        class="mr-1"
        style="color: lightgrey"
        @click="$store.commit('action', 'addUser')"
      />
    </h1>
    <loading-indicator v-if="state.isLoading"></loading-indicator>
    <div>
      <b-collapse id="collapse-4" v-model="showActionBar" class="mt-2">
        <users-add v-if="action === 'addUser'" @refresh="refresh"></users-add>
      </b-collapse>
    </div>
    <b-table
      hover
      bordeless
      :items="data"
      :fields="fields"
      class="mt-4"
      style="max-width: 650px"
    >
      <template v-slot:cell(action)="row">
        <b-button
          class="mr-2"
          variant="outline-secondary"
          @click="toggleUserGroups(row)"
          >{{
            row.detailsShowing && actionType === "groups" ? "Hide" : "Show"
          }}
          groups</b-button
        >
        <b-button
          class="mr-2"
          variant="outline-secondary"
          @click="toggleUserOrgs(row)"
          >{{
            row.detailsShowing && actionType === "orgs" ? "Hide" : "Show"
          }}
          orgs</b-button
        >
        <b-button
          variant="outline-danger"
          v-if="row.item.username !== username"
          @click="confirmDeleteItem(row.item.id, row.item.username)"
          >Delete</b-button
        >
      </template>
      <template v-slot:row-details="row">
        <div v-if="users[row.index] !== undefined">
          <b-card v-if="actionType === 'groups'">
            <users-manage-groups
              :user="users[row.index]"
              @user-added-in-group="userAddedToGroup(row, $event)"
            ></users-manage-groups>
          </b-card>
          <b-card v-else-if="actionType === 'orgs'">
            <users-manage-orgs
              :user="users[row.index]"
              @user-added-in-org="userAddedToOrg(row, $event)"
            ></users-manage-orgs>
          </b-card>
        </div>
      </template>
    </b-table>
    <b-modal title="Delete user" ref="delete-modal">
      Delete {{ itemToDelete.username }}?
      <template v-slot:modal-footer="mod">
        <b-button variant="danger" @click="deleteItem(itemToDelete)"
          >Delete</b-button
        >
        <b-button variant="warning" @click="mod.cancel()">Cancel</b-button>
      </template>
    </b-modal>
  </div>
</template>

<script>
import { mapState, mapGetters } from "vuex";
import LoadingIndicator from "@/components/LoadingIndicator";
import UsersAdd from "@/components/users/UsersAdd";
import UsersManageGroups from "@/components/users/UsersManageGroups";
import UsersManageOrgs from "@/components/users/UsersManageOrgs";

export default {
  components: {
    LoadingIndicator,
    UsersAdd,
    UsersManageGroups,
    UsersManageOrgs,
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
      actionType: null,
    };
  },
  methods: {
    async userAddedToGroup(row, group) {
      this.users[row.index].groups.add(group);
      row.toggleDetails();
      await new Promise((r) => setTimeout(r, 10));
      row.toggleDetails();
    },
    async userAddedToOrg(row, org) {
      this.users[row.index].orgs.add(org);
      row.toggleDetails();
      await new Promise((r) => setTimeout(r, 10));
      row.toggleDetails();
    },
    async fetchUserGroups(row) {
      let { response } = await this.$api.post("/admin/users/groups", {
        id: row.item.id,
      });
      let user = row.item;
      user.index = row.index;
      user.groups = new Set(response.data.groups);
      this.users[row.index] = user;
    },
    async fetchUserOrgs(row) {
      let { response } = await this.$api.post("/admin/users/orgs", {
        id: row.item.id,
      });
      let user = row.item;
      user.index = row.index;
      user.orgs = new Set(response.data.orgs);
      this.users[row.index] = user;
    },
    async toggleUserOrgs(row) {
      if (!row.detailsShowing) {
        this.actionType = "orgs";
        await this.fetchUserOrgs(row);
      } else {
        this.actionType = null;
      }
      row.toggleDetails();
    },
    async toggleUserGroups(row) {
      if (!row.detailsShowing) {
        this.actionType = "groups";
        await this.fetchUserGroups(row);
      } else {
        this.actionType = null;
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