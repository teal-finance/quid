<template>
  <b-form>
    <b-card title="Add a namespace" style="max-width: 28rem;" class="mb-2">
      <b-card-text>
        <b-form-group>
          <b-form-input :state="isValidName" v-model="form.name" placeholder="name"></b-form-input>
          <b-form-invalid-feedback :state="isValidName">The name must be at least 4 characters long</b-form-invalid-feedback>
        </b-form-group>
        <b-form-group label="Tokens max time to live">
          <b-form-input :state="isValidTtl" v-model="form.ttl"></b-form-input>
        </b-form-group>
        <b-form-group>
          <b-form-checkbox v-model="form.endpoint" switch>Enable public endpoint</b-form-checkbox>
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
  name: "NamespaceAdd",
  data() {
    return {
      form: {
        name: "",
        ttl: "1h",
        endpoint: false
      }
    };
  },
  methods: {
    async postForm() {
      let { error } = await this.$api.post("/admin/namespaces/add", {
        name: this.form.name,
        ttl: this.form.ttl,
        endpoint: this.form.endpoint
      });
      if (error == null) {
        this.form.name = "";
        this.form.ttl = "1h";
        this.form.endpoint = false;
        this.$emit("refresh");
        this.$store.commit("endAction");
      }
    }
  },
  computed: {
    isValidForm() {
      if (this.isValidName && this.isValidTtl) {
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
    isValidTtl() {
      if (this.form.ttl.length === 0) {
        return false;
      }
      return true;
    }
  }
};
</script>