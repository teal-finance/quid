<template>
  <b-card>
    <b-row class="mb-2">
      <b-col sm="3" class="text-sm-right">
        <b>Users</b>
      </b-col>
      <b-col>{{ details.num_users }}</b-col>
    </b-row>
    <b-row class="mb-2">
      <b-col sm="3" class="text-sm-right">
        <b>Groups</b>
      </b-col>
      <b-col>
        <b-badge
          variant="primary"
          class="mr-2"
          v-for="group in details.groups"
          :key="group.id"
        >{{ group.name }}</b-badge>
      </b-col>
    </b-row>
  </b-card>
</template>

<script>
export default {
  props: {
    details: {
      type: Object,
      default: () => {}
    }
  }
};
</script>