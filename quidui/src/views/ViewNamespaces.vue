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
        <namespace-add v-if="action === 'addNamespace'" @refresh="refresh"></namespace-add>
      </b-collapse>
    </div>
    <b-table hover bordeless :items="data" :fields="fields" class="mt-4" style="max-width:950px">
      <template v-slot:cell(public_endpoint_enabled)="row">
        <b-form-group class="text-center" v-if="row.item.name !== 'quid'">
          <b-form-checkbox
            v-model="row.item.public_endpoint_enabled"
            switch
            @change="toggleEndpoint(row)"
          ></b-form-checkbox>
        </b-form-group>
      </template>
      <template v-slot:cell(max_token_ttl)="row">
        <b-icon-stopwatch
          class="mr-1 text-superlight"
          @click="toggleEditTtl(row.item.id)"
          v-if="row.item.name!=='quid'"
        />&nbsp;
        <b-icon-stopwatch style="color:transparent" v-else />&nbsp;
        <span v-if="editTtl !== row.item.id ">{{ row.item.max_token_ttl }}</span>
        <namespace-edit-max-ttl
          v-else
          :namespaceId="row.item.id"
          :value="row.item.max_token_ttl"
          @end-edit="undeditTtl(row, $event)"
        ></namespace-edit-max-ttl>
      </template>
      <template v-slot:cell(max_refresh_token_ttl)="row">
        <b-icon-stopwatch
          class="mr-1 text-superlight"
          @click="toggleEditRefreshTtl(row.item.id)"
          v-if="row.item.name!=='quid'"
        />&nbsp;
        <b-icon-stopwatch style="color:transparent" v-else />&nbsp;
        <span
          v-if="editRefreshTtl !== row.item.id || row.item.name=='quid'"
        >{{ row.item.max_refresh_token_ttl }}</span>
        <namespace-edit-max-refresh-ttl
          v-else
          :namespaceId="row.item.id"
          :value="row.item.max_refresh_token_ttl"
          @end-edit="undeditRefreshTtl(row, $event)"
        ></namespace-edit-max-refresh-ttl>
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
        <namespace-info v-if="rowDetails[row.index] != undefined" :details="rowDetails[row.index]"></namespace-info>
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
import NamespaceAdd from "@/components/namespace/NamespaceAdd";
import NamespaceInfo from "@/components/namespace/NamespaceInfo";
import NamespaceEditMaxTtl from "@/components/namespace/NamespaceEditMaxTtl";
import NamespaceEditMaxRefreshTtl from "@/components/namespace/NamespaceEditMaxRefreshTtl";

export default {
  components: {
    NamespaceAdd,
    NamespaceInfo,
    NamespaceEditMaxTtl,
    NamespaceEditMaxRefreshTtl,
  },
  data() {
    return {
      data: [],
      fields: [
        { key: "id", sortable: true },
        { key: "name", sortable: true },
        {
          key: "public_endpoint_enabled",
          label: "Public endpoint",
          sortable: false,
        },
        { key: "max_token_ttl", label: "Access tokens ttl", sortable: true },
        {
          key: "max_refresh_token_ttl",
          label: "Refresh tokens ttl",
          sortable: true,
        },
        { key: "action", label: "", sortable: false },
      ],
      namespaceToDelete: {},
      rowDetails: {},
      selectedNs: { title: "" },
      isEditTtl: false,
      editTtl: null,
      isEditRefreshTtl: false,
      editRefreshTtl: null,
    };
  },
  methods: {
    undeditRefreshTtl(row, value) {
      this.isEditRefreshTtl = false;
      this.editRefreshTtl = null;
      if (value != null) {
        this.data[row.index].max_refresh_token_ttl = value;
      }
    },
    toggleEditRefreshTtl(id) {
      if (this.isEditRefreshTtl) {
        this.editRefreshTtl = null;
        this.isEditRefreshTtl = false;
      } else {
        this.editRefreshTtl = id;
        this.isEditRefreshTtl = true;
      }
    },
    undeditTtl(row, value) {
      this.isEditTtl = false;
      this.editTtl = null;
      if (value != null) {
        this.data[row.index].max_token_ttl = value;
      }
    },
    toggleEditTtl(id) {
      if (this.isEditTtl) {
        this.editTtl = null;
        this.isEditTtl = false;
      } else {
        this.editTtl = id;
        this.isEditTtl = true;
      }
    },
    async toggleEndpoint(row) {
      let enabled = !row.item.public_endpoint_enabled;
      let { error } = await this.$api.post("/admin/namespaces/endpoint", {
        id: row.item.id,
        enable: enabled,
      });
      if (error === null) {
        let msg = row.item.public_endpoint_enabled ? "enabled" : "disabled";
        this.$notify.done(`Namespace endpoint ${msg}`);
      }
    },
    async showKey(id, title) {
      this.selectedNs = { id: id, title: title, key: null };
      let { response, error } = await this.$api.post("/admin/namespaces/key", {
        id: id,
      });
      if (!error) {
        this.selectedNs.key = response.data.key;
        this.$refs["nskey-modal"].show();
      }
    },
    async toggleDetails(row) {
      if (!row.detailsShowing) {
        let { response } = await this.$api.post("/admin/namespaces/info", {
          id: row.item.id,
        });
        this.rowDetails[row.index] = response.data;
      }
      row.toggleDetails();
    },
    refresh() {
      this.fetchNamespaces();
      this.$notify.done("Namespace saved");
    },
    async fetchNamespaces() {
      let { response } = await this.$api.get("/admin/namespaces/all");
      this.data = response.data;
    },
    async deleteNamespace(ns) {
      this.$refs["delete-modal"].hide();
      let { error } = await this.$api.post("/admin/namespaces/delete", {
        id: ns.id,
      });
      if (error === null) {
        this.$notify.done("Namespace deleted");
        this.fetchNamespaces();
      }
    },
    getRowDetails(row) {
      if (!row.detailsShowing) {
        if (this.rowDetails[row.index] !== undefined) {
          return this.rowDetails[row.index];
        }
      }
      return { num_users: 0, groups: [] };
    },
    confirmDeleteNamespace(id, name) {
      this.namespaceToDelete = {
        name: name,
        id: id,
      };
      this.$refs["delete-modal"].show();
    },
    nsKeyModalTitle() {
      if (this.selectedNs === null) return;
      return `Key for namespace ${this.selectedNs.title}`;
    },
  },
  mounted: function () {
    this.fetchNamespaces();
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