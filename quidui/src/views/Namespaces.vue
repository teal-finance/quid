<template>
  <div>
    <h1 class="text-muted mt-3">
      Namespaces&nbsp;
      <b-icon-plus
        v-if="action!=='addNamespace'"
        class="mr-1"
        style="color:lightgrey"
        @click="$store.commit('action', 'addNamespace')"
      />
    </h1>
    <div>
      <b-collapse id="collapse-4" v-model="showActionBar" class="mt-2">
        <add v-if="action === 'addNamespace'" @refresh="refresh"></add>
      </b-collapse>
    </div>
    <b-table hover bordeless :items="data" :fields="fields" class="mt-4" style="max-width:850px">
      <template v-slot:cell(public_endpoint_enabled)="row">
        <b-form-group class="text-center" v-if="row.item.name !== 'quid'">
          <b-form-checkbox
            v-model="row.item.public_endpoint_enabled"
            switch
            @change="toggleEndpoint(row)"
          ></b-form-checkbox>
        </b-form-group>
      </template>
      <template v-slot:cell(action)="row">
        <b-button
          class="mr-2"
          variant="outline-secondary"
          @click="toggleDetails(row)"
        >{{ row.detailsShowing ? 'Hide' : 'Show'}} info</b-button>
        <b-button
          class="mr-2"
          variant="outline-secondary"
          v-if="row.item.name !== 'quid'"
          @click="showKey(row.item.id, row.item.name)"
        >Show key</b-button>
        <b-button
          variant="outline-danger"
          v-if="row.item.name !== 'quid'"
          @click="confirmDeleteNamespace(row.item.id, row.item.name)"
        >Delete</b-button>
      </template>
      <template v-slot:row-details="row">
        <info v-if="rowDetails[row.index] != undefined" :details="rowDetails[row.index]"></info>
      </template>
    </b-table>
    <b-modal title="Delete namespace" ref="delete-modal">
      Delete {{ namespaceToDelete.name }}?
      <template v-slot:modal-footer="{ ok, cancel }">
        <b-button variant="danger" @click="deleteNamespace(namespaceToDelete)">Delete</b-button>
        <b-button variant="warning" @click="cancel()">Cancel</b-button>
      </template>
    </b-modal>
    <b-modal :title="'Key for namespace '.concat(selectedNs.title)" ref="nskey-modal">
      {{ selectedNs.key }}
      <template v-slot:modal-footer="{ cancel }">
        <b-button variant="warning" @click="cancel()">Ok</b-button>
      </template>
    </b-modal>
  </div>
</template>

<script>
import { mapState, mapGetters } from "vuex";
import Add from "@/components/namespace/Add";
import Info from "@/components/namespace/Info";

export default {
  components: {
    Add,
    Info
  },
  data() {
    return {
      data: [],
      fields: [
        { key: "id", sortable: true },
        { key: "name", sortable: true },
        { key: "public_endpoint_enabled", sortable: false },
        { key: "max_token_ttl", sortable: true },
        { key: "action", sortable: false }
      ],
      namespaceToDelete: {},
      rowDetails: {},
      selectedNs: { title: "" }
    };
  },
  methods: {
    async toggleEndpoint(row) {
      let enabled = !row.item.public_endpoint_enabled;
      let { error } = await this.$api.post("/admin/namespaces/endpoint", {
        id: row.item.id,
        enable: enabled
      });
      if (error === null) {
        this.$bvToast.toast("ok", {
          title: "Namespace endpoint availability saved",
          variant: "success",
          autoHideDelay: 1500
        });
      }
    },
    async showKey(id, title) {
      this.selectedNs = { id: id, title: title, key: null };
      let { response } = await this.$api.post("/admin/namespaces/key", {
        id: id
      });
      this.selectedNs.key = response.data.key;
      this.$refs["nskey-modal"].show();
    },
    getRowDetails(row) {
      if (!row.detailsShowing) {
        if (this.rowDetails[row.index] !== undefined) {
          return this.rowDetails[row.index];
        }
      }
      return { num_users: 0, groups: [] };
    },
    async toggleDetails(row) {
      if (!row.detailsShowing) {
        let { response } = await this.$api.post("/admin/namespaces/info", {
          id: row.item.id
        });
        this.rowDetails[row.index] = response.data;
      }
      row.toggleDetails();
    },
    refresh() {
      this.fetchNamespaces();
      this.$bvToast.toast("ok", {
        title: "Namespace saved",
        variant: "success",
        autoHideDelay: 1500
      });
    },
    async fetchNamespaces() {
      let { response } = await this.$api.get("/admin/namespaces/all");
      this.data = response.data;
    },
    confirmDeleteNamespace(id, name) {
      this.namespaceToDelete = {
        name: name,
        id: id
      };
      this.$refs["delete-modal"].show();
    },
    async deleteNamespace(ns) {
      this.$refs["delete-modal"].hide();
      let { error } = await this.$api.post("/admin/namespaces/delete", {
        id: ns.id
      });
      if (error === null) {
        this.$bvToast.toast("Ok", {
          title: "Namespace deleted",
          autoHideDelay: 1000,
          variant: "success"
        });
      }
      this.fetchNamespaces();
    },
    nsKeyModalTitle() {
      if (this.selectedNs === null) return;
      return `Key for namespace ${this.selectedNs.title}`;
    }
  },
  mounted: function() {
    this.fetchNamespaces();
  },
  computed: {
    ...mapState(["action"]),
    ...mapGetters({
      s: "showActionBar"
    }),
    showActionBar: {
      get() {
        return this.s;
      },
      set(newName) {
        return newName;
      }
    }
  }
};
</script>