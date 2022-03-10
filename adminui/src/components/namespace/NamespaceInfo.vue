<template>
  <div class="flex flex-col p-3 ml-5">
    <div>
      <span class="pr-2 font-bold">Users:</span>
      {{ numUsers }}
    </div>
    <div class="mt-2 space-x-2 font-bold">
      <span class="pr-2">Groups:</span>
      <simple-badge
        v-for="group in groups"
        :key="group.id"
        class="inline-block"
        :text="group.name"
        color="secondary"
      ></simple-badge>
    </div>
  </div>
</template>

<script setup lang="ts">
import Group from '@/models/group';
import SimpleBadge from "@/components/widgets/SimpleBadge.vue";

defineProps({
  numUsers: {
    type: Number,
    required: true
  },
  groups: {
    type: Array as () => Array<Group>,
    required: true
  }
})
</script>