<template>
  <div>
    <b-row class="mb-2">
      <b-col sm="3" class="text-sm-right">
        <b>Orgs</b>
      </b-col>
      <b-col v-if="user.orgs.size > 0">
        <b-badge
          variant="primary"
          class="mr-2"
          v-for="org in user.orgs"
          :key="org.id"
        >
          {{ org.name }}
          <b-icon
            icon="x-circle-fill"
            class="ml-1"
            @click="confirmRemoveUserFromOrg(org)"
          ></b-icon>
        </b-badge>
      </b-col>
      <b-col v-else class="text-superlight">the user has no orgs</b-col>
    </b-row>
    <b-button
      variant="link"
      size="sm"
      class="float-right"
      v-if="!addingToOrg"
      @click="addingToOrg = true"
    >
      <b-icon-plus class="mr-1" />&nbsp;Add user to org
    </b-button>
    <div v-else>
      <users-add-to-org
        :user="user"
        @add-user-to-org="addUserToOrg($event)"
        @cancel="addingToOrg = false"
      ></users-add-to-org>
    </div>
    <b-modal title="Remove user from org" ref="remove-from-org-modal">
      <span v-if="orgToRemove !== null"
        >Remove user from org {{ orgToRemove.name }}?</span
      >
      <template v-slot:modal-footer="mod">
        <b-button variant="danger" @click="removeUserFromOrg(orgToRemove)"
          >Remove</b-button
        >
        <b-button variant="warning" @click="mod.cancel()">Cancel</b-button>
      </template>
    </b-modal>
  </div>
</template>

<script>
import UsersAddToOrg from "./UsersAddToOrg";

export default {
  components: {
    UsersAddToOrg,
  },
  data: function () {
    return {
      addingToOrg: false,
      orgToRemove: null,
    };
  },
  props: {
    user: {
      type: Object,
      required: true,
    },
  },
  methods: {
    async removeUserFromOrg(o) {
      this.$refs["remove-from-org-modal"].hide();
      if (this.addingToOrg) {
        this.addingToOrg = false;
      }
      let { error } = await this.$api.post("/admin/orgs/remove_user", {
        user_id: this.user.id,
        org_id: this.orgToRemove.id,
      });
      if (error === null) {
        this.user.orgs.delete(o);
        this.$bvToast.toast(`Ok`, {
          title: "User removed from org",
          variant: "success",
          autoHideDelay: 1500,
        });
      }
      this.orgToRemove = null;
    },
    confirmRemoveUserFromOrg(o) {
      this.orgToRemove = o;
      this.$refs["remove-from-org-modal"].show();
    },
    async addUserToOrg(o) {
      this.addingToOrg = false;
      let { error } = await this.$api.post("/admin/orgs/add_user", {
        user_id: this.user.id,
        org_id: o.id,
      });
      if (error === null) {
        this.user.orgs.add(o);
        this.$emit("user-added-in-org", o);
        this.$bvToast.toast(`Ok`, {
          title: "User added in org",
          variant: "success",
          autoHideDelay: 1500,
        });
      }
    },
  },
};
</script>