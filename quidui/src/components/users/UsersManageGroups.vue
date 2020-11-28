<template>
  <div>
    <b-row class="mb-2">
      <b-col sm="3" class="text-sm-right">
        <b>Groups</b>
      </b-col>
      <b-col v-if="user.groups.size > 0">
        <b-badge
          variant="primary"
          class="mr-2"
          v-for="group in user.groups"
          :key="group.id"
        >
          {{ group.name }}
          <b-icon
            icon="x-circle-fill"
            class="ml-1"
            @click="confirmRemoveUserFromGroup(group)"
          ></b-icon>
        </b-badge>
      </b-col>
      <b-col v-else class="text-superlight">the user has no groups</b-col>
    </b-row>
    <b-button
      variant="link"
      size="sm"
      class="float-right"
      v-if="!addingToGroup"
      @click="addingToGroup = true"
    >
      <b-icon-plus class="mr-1" />&nbsp;Add user to group
    </b-button>
    <div v-else @click="addingToGroup = false">
      <users-add-to-group
        :user="user"
        @add-user-to-group="addUserToGroup($event)"
      ></users-add-to-group>
    </div>
    <b-modal title="Remove user from group" ref="remove-from-group-modal">
      <span v-if="groupToRemove !== null"
        >Remove user from group {{ groupToRemove.name }}?</span
      >
      <template v-slot:modal-footer="mod">
        <b-button variant="danger" @click="removeUserFromGroup(groupToRemove)"
          >Remove</b-button
        >
        <b-button variant="warning" @click="mod.cancel()">Cancel</b-button>
      </template>
    </b-modal>
  </div>
</template>

<script>
import UsersAddToGroup from "./UsersAddToGroup";
export default {
  components: {
    UsersAddToGroup,
  },
  data: function () {
    return {
      addingToGroup: false,
      groupToRemove: null,
    };
  },
  props: {
    user: {
      type: Object,
      required: true,
    },
  },
  methods: {
    async removeUserFromGroup(g) {
      this.$refs["remove-from-group-modal"].hide();
      if (this.addingToGroup) {
        this.addingToGroup = false;
      }
      let { error } = await this.$api.post("/admin/groups/remove_user", {
        user_id: this.user.id,
        group_id: this.groupToRemove.id,
      });
      if (error === null) {
        this.user.groups.delete(g);
        this.$bvToast.toast(`Ok`, {
          title: "User removed from group",
          variant: "success",
          autoHideDelay: 1500,
        });
      }
      this.groupToRemove = null;
    },
    confirmRemoveUserFromGroup(g) {
      this.groupToRemove = g;
      this.$refs["remove-from-group-modal"].show();
    },
    async addUserToGroup(g) {
      this.addingToGroup = false;
      let { error } = await this.$api.post("/admin/groups/add_user", {
        user_id: this.user.id,
        group_id: g.id,
      });
      if (error === null) {
        this.user.groups.add(g);
        this.$emit("user-added-in-group", g);
        this.$bvToast.toast(`Ok`, {
          title: "User added in group",
          variant: "success",
          autoHideDelay: 1500,
        });
      }
    },
  },
  /*mounted() {
    console.log("USER", this.user);
  },*/
};
</script>