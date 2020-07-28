<template>
  <b-card title="Login" tag="article" style="max-width: 20rem;" class="mb-2">
    <b-form @submit.prevent="postForm">
      <b-card-text>
        <b-form-group>
          <b-form-input v-model="form.username" placeholder="name"></b-form-input>
        </b-form-group>
        <b-form-group>
          <b-form-input type="password" v-model="form.password" placeholder="password"></b-form-input>
        </b-form-group>
      </b-card-text>
      <b-button variant="primary" type="submit">Submit</b-button>
    </b-form>
  </b-card>
</template>

<script>
export default {
  data() {
    return {
      form: {
        username: "",
        password: "",
      },
    };
  },
  methods: {
    async postForm() {
      let { response, error } = await this.$api.post("/admin_login", {
        namespace: "quid",
        username: this.form.username,
        password: this.form.password,
      });
      if (error === null) {
        this.$store.commit("authenticate", {
          key: response.data.key,
          username: this.form.username,
        });
        this.$router.push({ path: "/namespaces" });
      } else {
        if (error.response.status === 401) {
          this.$bvToast.toast("Login refused", {
            title: "Authentication failed",
            variant: "danger",
          });
        }
      }
    },
  },
};
</script>
