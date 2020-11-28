<template>
  <b-form-group>
    <b-form-input
      v-model="value"
      debounce="500"
      @change="onChange"
      @update="fetchOrgs"
      placeholder="org"
      style="width: 12em"
    ></b-form-input>
    <div class="mt-2">
      <b-icon
        v-if="state.isLoading"
        icon="three-dots"
        animation="cylon"
      ></b-icon>
      <b-badge
        variant="primary"
        class="mr-2"
        v-for="item in items"
        :key="item.id"
        @click="select(item)"
        role="button"
        >{{ item.name }}</b-badge
      >
    </div>
  </b-form-group>
</template>

<script>
export default {
  data() {
    return {
      value: "",
      items: [],
      selected: null,
      state: {
        isLoading: false,
      },
    };
  },
  props: {
    user: {
      type: Object,
      required: true,
    },
  },
  methods: {
    async fetchOrgs() {
      if (this.value === "") {
        this.items = [];
        return;
      }
      this.state.isLoading = true;
      let { response, error } = await this.$api.post("/admin/orgs/find", {
        name: this.value,
      });
      this.state.isLoading = false;
      if (error === null) {
        let orgs = new Set([]);
        let orgNames = new Set(
          Array.from(this.user.orgs).map((item) => item.name)
        );
        response.data.forEach((item) => {
          console.log(item, orgNames.has(item.name));
          if (!orgNames.has(item.name)) {
            orgs.add(item);
          }
        });
        this.items = orgs;
      }
    },
    select(item) {
      this.$emit("org-selected", item);
    },
    onChange() {
      this.selected = null;
    },
  },
};
</script>