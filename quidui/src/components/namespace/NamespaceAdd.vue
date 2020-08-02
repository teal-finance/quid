<template>
  <b-form>
    <b-card title="Add a namespace" style="max-width: 28rem;" class="mb-2">
      <b-card-text>
        <b-form-group>
          <b-form-input :state="isValidName" v-model="form.name" placeholder="name"></b-form-input>
          <b-form-invalid-feedback :state="isValidName">The name must be at least 4 characters long</b-form-invalid-feedback>
        </b-form-group>
        <b-form-group label="Access tokens max time to live">
          <b-form-input :state="isNotEmpty" v-model="form.maxTtl"></b-form-input>
        </b-form-group>
        <b-form-group label="Refresh tokens max time to live">
          <b-form-input :state="isNotEmpty" v-model="form.refreshMaxTtl"></b-form-input>
        </b-form-group>
        <b-form-group>
          <b-form-checkbox v-model="form.enableEndpoint" switch>Enable public enableEndpoint</b-form-checkbox>
        </b-form-group>
      </b-card-text>
      <b-button variant="success" v-if="isValidForm" @click="postForm">Save</b-button>
      <b-button variant="success" v-else disabled>Save</b-button>
      <b-button variant="warning" class="ml-2" @click="$store.commit('endAction')">Cancel</b-button>
    </b-card>
  </b-form>
</template>

<script>
export default {
  data() {
    return {
      form: {
        name: "",
        maxTtl: "20m",
        refreshMaxTtl: "24h",
        enableEndpoint: false,
      },
    };
  },
  methods: {
    async postForm() {
      let { error } = await this.$api.post("/admin/namespaces/add", {
        name: this.form.name,
        max_ttl: this.form.maxTtl,
        refresh_max_ttl: this.form.refreshMaxTtl,
        enable_endpoint: this.form.enableEndpoint,
      });
      if (error == null) {
        this.form.name = "";
        this.form.maxTtl = "20m";
        (this.form.refreshMaxTtl = "24h"), (this.form.enableEndpoint = false);
        this.$emit("refresh");
        this.$store.commit("endAction");
      }
    },
  },
  computed: {
    isValidForm() {
      if (this.isValidName && this.isNotEmpty) {
        return true;
      }
      return false;
    },
    isValidName() {
      if (this.form.name.length === 0) {
        return null;
      }
      return this.form.name.length >= 4;
    },
    isNotEmpty() {
      if (this.form.maxTtl.length === 0) {
        return false;
      }
      return true;
    },
  },
};
</script>