<template>
  <b-form-group class="mt-4">
    <h4 class="text-muted">Add user to org</h4>
    <div class="mt-3">
      <orgs-select
        :user="user"
        @org-selected="addUserToOrg($event)"
      ></orgs-select>
    </div>
    <div class="mt-4">
      <b-button variant="warning" @click="$emit('cancel')">Cancel</b-button>
    </div>
  </b-form-group>
</template>

<script>
import OrgsSelect from "@/components/orgs/OrgsSelect";

export default {
  components: {
    OrgsSelect,
  },
  data() {
    return {
      state: {
        isLoading: false,
      },
      newOrgs: new Set([]),
    };
  },
  props: {
    user: {
      type: Object,
      required: true,
    },
  },
  methods: {
    addUserToOrg(o) {
      this.$emit("add-user-to-org", o);
    },
  },
};
</script>