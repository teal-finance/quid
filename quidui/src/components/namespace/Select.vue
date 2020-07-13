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
  name: "SelectNamespace",
  data() {
    return {
      value: "",
      items: [],
      selected: null,
      state: {
        isLoading: false,
        isFormValid: false
      }
    };
  },
  methods: {
    fetchNamespaces() {
      this.state.isFormValid = false;
      this.$emit("unvalidate");
      if (this.value === "") {
        this.items = [];
        return;
      }
      this.state.isLoading = true;
      let vue = this;
      vue.$axios
        .post("/admin/namespaces/find", {
          name: vue.value
        })
        .then(function(response) {
          vue.items = response.data;
        })
        .catch(e => {
          if (e.response.status != 200) {
            if (e.response.status === 401) {
              vue.$bvToast.toast(e.response.error, {
                title: "Error",
                variant: "danger"
              });
            } else {
              console.log(e);
            }
          }
        })
        .finally(function() {
          vue.state.isLoading = false;
        });
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
    }
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
    }
  }
};
</script>