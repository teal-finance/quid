<template>
  <b-form>
    <b-card title="Add a group" style="max-width: 28rem;" class="mb-2">
      <b-card-text>
        <namespace-select @validate="validateNamespace" @unvalidate="unvalidateNamespace"></namespace-select>
        <b-form-group>
          <b-form-input v-model="form.name" :state="isValidName" placeholder="name"></b-form-input>
          <b-form-invalid-feedback :state="isValidName">The name must be at least 4 characters long</b-form-invalid-feedback>
        </b-form-group>
      </b-card-text>
      <b-button variant="success" v-if="isValidForm" @click="postForm">Save</b-button>
      <b-button variant="success" v-else disabled>Save</b-button>
      <b-button variant="warning" class="ml-2" @click="$store.commit('endAction')">Cancel</b-button>
    </b-card>
  </b-form>
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
        name: "",
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
      let { error } = await this.$api.post("/admin/groups/add", {
        name: this.form.name,
        namespace_id: this.form.namespace.id,
      });
      if (error == null) {
        this.form.name = "";
        this.form.namespace = null;
        this.$emit("refresh");
        this.$store.commit("endAction");
      } else {
        //console.log("ERR RESP", error.response);
        if (error.reponse.status === 409) {
          this.form.name = "";
          this.form.namespace = null;
        }
      }
    },
  },
  computed: {
    isValidForm() {
      return this.isValidName && this.isNamespaceValid;
    },
    isValidName() {
      if (this.form.name.length === 0) {
        return null;
      }
      return this.form.name.length >= 4;
    },
  },
};
</script>