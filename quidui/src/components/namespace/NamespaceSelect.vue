<template>
  <b-form-group>
    <b-form-input
      v-model="value"
      debounce="500"
      :state="isValid"
      @change="onChange"
      @update="fetchNamespaces"
      placeholder="namespace"
    ></b-form-input>
    <div class="mt-2">
      <b-icon v-if="state.isLoading" icon="three-dots" animation="cylon"></b-icon>
      <b-badge
        variant="primary"
        class="mr-2"
        v-for="item in items"
        :key="item.id"
        @click="select(item)"
      >{{ item.name }}</b-badge>
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
        isFormValid: false,
      },
    };
  },
  methods: {
    async fetchNamespaces() {
      this.state.isFormValid = false;
      this.$emit("unvalidate");
      if (this.value === "") {
        this.items = [];
        return;
      }
      this.state.isLoading = true;
      let { response, error } = await this.$api.post("/admin/namespaces/find", {
        name: this.value,
      });
      this.state.isLoading = false;
      if (error === null) {
        this.items = response.data;
      }
    },
    select(item) {
      this.value = item.name;
      this.selected = item;
      this.state.isFormValid = true;
      this.$emit("validate", this.selected);
      this.items = [];
    },
    onChange() {
      this.selected = null;
      this.$emit("unvalidate");
    },
  },
  computed: {
    isValid() {
      if (this.value === "") {
        return null;
      }
      if (this.selected === null) {
        return false;
      }
      return this.state.isFormValid;
    },
  },
};
</script>