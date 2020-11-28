<template>
  <b-form>
    <b-card title="Add an org" style="max-width: 28rem" class="mb-2">
      <b-card-text>
        <b-form-group>
          <b-form-input
            v-model="form.name"
            :state="isValidName"
            placeholder="name"
          ></b-form-input>
          <b-form-invalid-feedback :state="isValidName"
            >The name must be at least 4 characters
            long</b-form-invalid-feedback
          >
        </b-form-group>
      </b-card-text>
      <b-button variant="success" v-if="isValidName" @click="postForm"
        >Save</b-button
      >
      <b-button variant="success" v-else disabled>Save</b-button>
      <b-button
        variant="warning"
        class="ml-2"
        @click="$store.commit('endAction')"
        >Cancel</b-button
      >
    </b-card>
  </b-form>
</template>

<script>
export default {
  data() {
    return {
      form: {
        name: "",
      },
    };
  },
  methods: {
    async postForm() {
      let { error } = await this.$api.post("/admin/orgs/add", {
        name: this.form.name,
      });
      if (error == null) {
        this.form.name = "";
        this.$emit("refresh");
        this.$store.commit("endAction");
      } else {
        //console.log("ERR RESP", error.response);
        if (error.reponse.status === 409) {
          this.form.name = "";
        }
      }
    },
  },
  computed: {
    isValidName() {
      if (this.form.name.length === 0) {
        return null;
      }
      return this.form.name.length >= 4;
    },
  },
};
</script>