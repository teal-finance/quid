<template>
  <b-card>
    <b-row class="mb-2">
      <b-col sm="3" class="text-sm-right">
        <b>Users</b>
      </b-col>
      <b-col>{{ details.num_users }}</b-col>
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