<template>
  <div>
    <h1 class="text-muted mt-3">Users</h1>
    <loading v-if="state.isLoading"></loading>
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
        >{{ row.detailsShowing ? 'Hide' : 'Show'}} info</b-button>
        <b-button
          variant="outline-danger"
          v-if="row.item.name !== username"
          @click="confirmDeleteItem(row.item.id, row.item.name)"
        >Delete</b-button>
      </template>
      <template v-slot:row-details="row">
        <users-info v-if="rowDetails[row.index] != undefined" :details="rowDetails[row.index]"></users-info>
      </template>
    </b-table>
    <b-modal title="Delete user" ref="delete-modal">
      Delete {{ itemToDelete.name }}?
      <template v-slot:modal-footer="{ ok, cancel }">
        <b-button variant="danger" @click="deleteItem(itemToDelete)">Delete</b-button>
        <b-button variant="warning" @click="cancel()">Cancel</b-button>
      </template>
    </b-modal>
  </div>
</template>

<script>
import { mapState, mapGetters } from "vuex";
import UsersAdd from "@/components/users/UsersAdd";
import UsersInfo from "@/components/users/UsersInfo";

export default {
  components: {
    UsersAdd,
    UsersInfo,
  },
  data() {
    return {
      data: [],
      state: {
        isLoading: false,
      },
      fields: [
        { key: "id", sortable: true },
        { key: "name", sortable: true },
        { key: "namespace", sortable: true },
        { key: "action", sortable: false },
      ],
      itemToDelete: {},
      rowDetails: {},
    };
  },
  methods: {
    getRowDetails(row) {
      if (!row.detailsShowing) {
        if (this.rowDetails[row.index] !== undefined) {
          return this.rowDetails[row.index];
        }
      }
      return { groups: [] };
    },
    async toggleDetails(row) {
      if (!row.detailsShowing) {
        let { response } = await this.$api.post("/admin/users/info", {
          id: row.item.id,
        });
        this.rowDetails[row.index] = response.data;
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
    confirmDeleteItem(id, name) {
      this.itemToDelete = {
        name: name,
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