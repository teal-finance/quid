<template>
  <div>
    <h1 class="text-muted mt-3">
      Orgs&nbsp;
      <b-icon-plus
        v-if="action !== 'addOrg'"
        class="mr-1"
        style="color: lightgrey"
        @click="$store.commit('action', 'addOrg')"
      />
    </h1>
    <loading-indicator v-if="state.isLoading"></loading-indicator>
    <div>
      <b-collapse id="collapse-4" v-model="showActionBar" class="mt-2">
        <orgs-add v-if="action === 'addOrg'" @refresh="refresh"></orgs-add>
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
          variant="outline-danger"
          @click="confirmDeleteItem(row.item.id, row.item.name)"
          >Delete</b-button
        >
      </template>
    </b-table>
    <b-modal title="Delete org" ref="delete-modal">
      Delete {{ itemToDelete.name }}?
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
import OrgsAdd from "@/components/orgs/OrgsAdd";

export default {
  components: {
    LoadingIndicator,
    OrgsAdd,
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
        { key: "action", sortable: false },
      ],
      itemToDelete: {},
    };
  },
  methods: {
    async fetchOrgs() {
      this.state.isLoading = true;
      let { response } = await this.$api.get("/admin/orgs/all");
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
      let { error } = await this.$api.post("/admin/orgs/delete", {
        id: ns.id,
      });
      if (error === null) {
        this.$bvToast.toast("Ok", {
          title: "Org deleted",
          autoHideDelay: 1000,
          variant: "success",
        });
        this.fetchOrgs();
      }
    },
    refresh() {
      this.fetchOrgs();
      this.$bvToast.toast("ok", {
        title: "Org saved",
        variant: "success",
        autoHideDelay: 1500,
      });
    },
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
  mounted() {
    this.fetchOrgs();
  },
};
</script>