<template>
  <div>
    <b-card title="Add user" style="max-width: 28rem;" class="mb-2">
      <b-card-text>
        <namespace-select @validate="validateNamespace" @unvalidate="unvalidateNamespace"></namespace-select>
        <b-form-group>
          <b-form-input :state="isValidName" v-model="form.username" placeholder="name"></b-form-input>
        </b-form-group>
        <b-form-group>
          <b-form-input
            :state="isValidPassword"
            type="password"
            v-model="form.password"
            placeholder="password"
          ></b-form-input>
        </b-form-group>
      </b-card-text>
      <b-button variant="success" v-if="isValidForm" @click="postForm">Save</b-button>
      <b-button variant="success" v-else disabled>Save</b-button>
      <b-button variant="warning" class="ml-2" @click="$store.commit('endAction')">Cancel</b-button>
    </b-card>
  </div>
</template>

<script>
import NamespaceSelect from "@/components/namespace/NamespaceSelect.vue";

export default {
  components: {
    NamespaceSelect,
  },
  data() {
    return {
      form: {
        username: "",
        password: "",
        namespace: null,
      },
      isNamespaceValid: false,
    };
  },
  methods: {
    unvalidateNamespace() {
      this.isNamespaceValid = false;
    },
    validateNamespace(name) {
      this.form.namespace = name;
      this.isNamespaceValid = true;
    },
    async postForm() {
      let { error } = await this.$api.post("/admin/users/add", {
        name: this.form.username,
        password: this.form.password,
        namespace_id: this.form.namespace.id,
      });
      if (error == null) {
        this.form.name = "";
        this.form.password = "";
        this.form.namespace = null;
        this.$emit("refresh");
        this.$store.commit("endAction");
      } else {
        //console.log("ERR RESP", error.response);
        if (error.reponse.status === 409) {
          this.form.name = "";
          this.form.password = "";
          this.form.namespace = null;
        }
      }
    },
  },
  computed: {
    isValidForm() {
      return this.isValidName && this.isValidPassword && this.isNamespaceValid;
    },
    isValidName() {
      if (this.form.username.length === 0) {
        return null;
      }
      return this.form.username.length >= 3;
    },
    isValidPassword() {
      if (this.form.password.length === 0) {
        return null;
      }
      return this.form.password.length >= 7;
    },
  },
};
</script>