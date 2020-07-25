<template>
  <div>
    <h1 class="text-muted mt-3">Groups</h1>
    <loading-indicator v-if="state.isLoading"></loading-indicator>
    <div>
      <b-collapse id="collapse-4" v-model="showActionBar" class="mt-2">
        <groups-add v-if="action === 'addGroup'" @refresh="refresh"></groups-add>
      </b-collapse>
    </div>
    <b-table hover bordeless :items="data" :fields="fields" class="mt-4" style="max-width:650px">
      <template v-slot:cell(action)="row">
        <b-button
          class="mr-2"
          variant="outline-secondary"
          @click="toggleDetails(row)"
        >{{ row.detailsShowing ? 'Hide' : 'Show'}} users</b-button>
        <b-button
          variant="outline-danger"
          v-if="row.item.name !== 'quid_admin'"
          @click="confirmDeleteItem(row.item.id, row.item.name)"
        >Delete</b-button>
      </template>
      <template v-slot:row-details="row">
        <groups-info v-if="rowDetails[row.index] != undefined" :details="rowDetails[row.index]"></groups-info>
      </template>
    </b-table>
    <b-modal title="Delete group" ref="delete-modal">
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
import GroupsAdd from "@/components/groups/GroupsAdd";
import GroupsInfo from "@/components/groups/GroupsInfo";
import LoadingIndicator from "@/components/LoadingIndicator";

export default {
  components: {
    GroupsAdd,
    GroupsInfo,
    LoadingIndicator,
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
      return { users: null };
    },
    async toggleDetails(row) {
      if (!row.detailsShowing) {
        let { response } = await this.$api.post("/admin/groups/info", {
          id: row.item.id,
        });
        this.rowDetails[row.index] = response.data;
      }
      row.toggleDetails();
    },
    refresh() {
      this.fetchGroups();
      this.$bvToast.toast("ok", {
        title: "Group saved",
        variant: "success",
        autoHideDelay: 1500,
      });
    },
    async fetchGroups() {
      this.state.isLoading = true;
      let { response } = await this.$api.get("/admin/groups/all");
      this.data = response.data;
      this.state.isLoading = false;
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
      let { error } = await this.$api.post("/admin/groups/delete", {
        id: ns.id,
      });
      if (error === null) {
        this.$bvToast.toast("Ok", {
          title: "Group deleted",
          autoHideDelay: 1000,
          variant: "success",
        });
        this.fetchGroups();
      }
    },
  },
  mounted: function () {
    this.fetchGroups();
  },
  computed: {
    ...mapState(["action"]),
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
};
</script>