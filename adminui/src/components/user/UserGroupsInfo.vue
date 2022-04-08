<template>
  <div class="flex flex-col" v-if="!isLoading">
    <div class="mt-2 space-x-2">
      <span class="pr-2 font-bold">Groups:</span>
      <span v-if="groups.length > 0">
        <simple-badge
          v-for="group in groups"
          :key="group.id"
          class="inline-block secondary cursor-pointer"
          :text="group.name"
          @click="removeUserFromGroup(group)"
        ></simple-badge>
      </span>
      <span v-else class="txt-lighter">the user has no groups</span>
    </div>
  </div>
  <div v-else>
    <loading-indicator :small="true"></loading-indicator>
  </div>
</template>

<script setup lang="ts">
import LoadingIndicator from "@/components/widgets/LoadingIndicator.vue";
import SimpleBadge from "@/components/widgets/SimpleBadge.vue";
import Group from "@/models/group";
import { GroupTable } from '@/models/group/interface';
import { UserTable } from "@/models/user/interface";
import { notify } from "@/state";

const props = defineProps({
  isLoading: {
    type: Boolean,
    default: true,
  },
  groups: {
    type: Array as () => Array<GroupTable>,
    required: true
  },
  user: {
    type: Object as () => UserTable,
    required: true,
  },
});

const emit = defineEmits(["user-removed"])

function removeUserFromGroup(g: GroupTable) {
  notify.confirmDelete(
    `Remove the ${props.user.name} from group ${g.name}?`,
    () => {
      Group.removeUserFromGroup(props.user.id, g.id).then(() => {
        notify.done("User removed from group");
        emit("user-removed", g);
      })
    }
  )
}
</script>