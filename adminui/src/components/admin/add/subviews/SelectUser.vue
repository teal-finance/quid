<template>
  <div>
    <div class="flex flex-row ml-3">
      <simple-badge
        class="cursor-pointer primary"
        :class="isSelected(user) ? 'secondary' : ['block-lighter', 'txt-light']"
        v-for="user in users"
        @click="toggleSelectUser(user)"
      >
        <div class="flex flex-row items-center justify-center">
          <i-ep:close-bold class="mr-1" v-if="!isSelected(user)"></i-ep:close-bold>
          <i-ci:check-bold v-else class="mr-1"></i-ci:check-bold>
          <div>{{ user.name }}</div>
        </div>
      </simple-badge>
    </div>
    <div class="flex flex-row mt-12 mb-5">
      <button
        class="mr-3 btn success"
        :disabled="selectedUserIds.size == 0"
        @click="submitSelection()"
      >Save</button>
      <button class="mr-3 btn warning" @click="onCancel()">Cancel</button>
      <button class="btn lighter" @click="onBack()">Back</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import User from '@/models/user/user';
import SimpleBadge from "@/components/widgets/SimpleBadge.vue";
import { reactive } from 'vue';

defineProps({
  users: {
    type: Array as () => Array<User>,
    required: true,
  }
});
const emit = defineEmits(["users-selected", "users-unselected", "submit-selection", "back", "cancel"]);
const selectedUsers = new Set<User>();
const selectedUserIds = reactive(new Set<number>());

function isSelected(user: User) {
  return selectedUserIds.has(user.id)
}

function toggleSelectUser(user: User) {
  if (selectedUserIds.has(user.id)) {
    selectedUsers.delete(user);
    selectedUserIds.delete(user.id);
    if (selectedUserIds.size == 0) {
      emit("users-unselected")
    }
  } else {
    selectedUsers.add(user);
    selectedUserIds.add(user.id);
    emit("users-selected");
  }
}

function submitSelection() {
  emit("submit-selection", [...selectedUsers]);
}

function onBack() {
  emit("back")
}

function onCancel() {
  emit("cancel")
}
</script>