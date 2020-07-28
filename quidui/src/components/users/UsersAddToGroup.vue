<template>
  <b-form-group class="mt-4">
    <h4 class="text-muted">Add user to group</h4>
    <div class="mt-2">
      <b-icon v-if="state.isLoading" icon="three-dots" animation="cylon"></b-icon>
      <div v-if="newGroups.size > 0">
        <b-badge
          variant="success"
          class="mr-2"
          role="button"
          v-for="group in newGroups"
          :key="group.id"
          @click="addUserToGroup(group)"
        >{{ group.name }}</b-badge>
      </div>
      <p class="text-superlight" v-else>No groups found</p>
    </div>
    <div class="mt-4">
      <b-button variant="warning">Cancel</b-button>
    </div>
  </b-form-group>
</template>

<script>
export default {
  data() {
    return {
      selected: null,
      state: {
        isLoading: false,
      },
      newGroups: new Set([]),
    };
  },
  props: {
    user: {
      type: Object,
      required: true,
    },
  },
  methods: {
    addUserToGroup(g) {
      this.$emit("add-user-to-group", g);
    },
    async fetchData() {
      this.state.isLoading = true;
      let { response, error } = await this.$api.post(
        "/admin/namespaces/groups",
        {
          namespace: this.user.namespace,
        }
      );
      this.state.isLoading = false;
      if (error === null) {
        let ng = new Set([]);
        for (let newGroup of response.data.groups) {
          let has = false;
          for (let g of this.user.groups) {
            if (g.name === newGroup.name) {
              has = true;
              break;
            }
          }
          if (!has) {
            ng.add(newGroup);
          }
        }
        this.newGroups = ng;
      }
    },
  },
  mounted() {
    this.fetchData();
  },
};
</script>