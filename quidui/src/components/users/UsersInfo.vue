<template>
  <div>
    <b-row class="mb-2">
      <b-col sm="3" class="text-sm-right">
        <b>Groups</b>
      </b-col>
      <b-col v-if="details.groups.length > 0">
        <b-badge
          variant="primary"
          class="mr-2"
          v-for="group in details.groups"
          :key="group"
        >{{ group }}</b-badge>
      </b-col>
      <b-col v-else class="text-superlight">the user has no groups</b-col>
    </b-row>
  </div>
</template>

<script>
export default {
  props: {
    details: {
      type: Object,
      default: () => {},
    },
  },
};
</script>