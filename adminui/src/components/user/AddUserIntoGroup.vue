<template>
  <div v-if="step == 1" class="mt-5">
    <simple-badge
      v-for="group in state.nsGroups"
      :key="group.id"
      class="success cursor-pointer"
      @click="addInGroup(group)"
    >{{ group.name }}</simple-badge>
    <div class="mt-5 flex flex-row pl-5 text-sm">
      <div class="mt-3 txt-lighter">Select a group</div>
    </div>
  </div>
  <div class="mt-5 flex flex-row pl-5 text-sm" v-else-if="step == 0">
    <button class="btn hover:secondary" @click="fetchGroups()">Add user in group</button>
  </div>
</template>

<script setup lang="ts">
import SimpleBadge from "@/components/widgets/SimpleBadge.vue";
import Group from '@/models/group';
import { GroupTable } from '@/models/group/interface';
import { UserTable } from '@/models/user/interface';
import { user } from '@/state';
import { reactive, ref } from 'vue';

const props = defineProps({
  user: {
    type: Object as () => UserTable,
    required: true,
  },
  userGroups: {
    type: Array as () => Array<GroupTable>,
    required: true,
  }
});
const emit = defineEmits<{
  (event: 'addingroup', g: GroupTable): void
}>();

const state = reactive({
  nsGroups: new Array<GroupTable>()
})
const step = ref(0);

async function fetchGroups() {
  const gr = await Group.fetchAll(user.namespace.value.id);
  const grIds = props.userGroups.map<number>((item) => item.id);
  state.nsGroups = [];
  gr.forEach((g) => {
    if (!grIds.includes(g.id)) {
      state.nsGroups.push(g)
    }
  })
  step.value = 1;
}

function addInGroup(g: GroupTable) {
  emit("addingroup", g);
  step.value = 0
}
</script>