<template>
  <div>
    <div class="text-xl">Add administrators</div>
    <div class="px-8 mt-5 mb-16">
      <sw-progress-stepper
        class="stepper-secondary dark:stepper-success"
        :steps="steps"
        :active-index="activeStep"
      >
        <template #content="slotProps">
          <i-ant-design:save-outlined class="text-xl" v-if="slotProps.index == 2"></i-ant-design:save-outlined>
        </template>
      </sw-progress-stepper>
    </div>
    <component
      :is="subviews.component"
      :users="users"
      @users-found="step2($event)"
      @users-selected="setHasUsersSelected()"
      @users-unselected="setHasNoUsersSelected()"
      @submit-selection="submitSelection($event)"
      @cancel="onCancel()"
      @back="onBack()"
    ></component>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from "vue";
import User from "@/models/user/user";
import { SubViews } from "@snowind/subviews"
import SearchForUsers from "./subviews/SearchForUsers.vue";
import SelectUser from "./subviews/SelectUser.vue";
import SwProgressStepper from "@/packages/stepper/SwProgressStepper.vue";

const emit = defineEmits(["end"]);
const users = ref<Array<User>>([]);
const subviews = new SubViews({
  views: {
    "search": { component: SearchForUsers, props: { label: "Search for users" } },
    "select": { component: SelectUser, props: { label: "Select users to add to admin" } },
    "end": { component: SearchForUsers, props: { label: "Finish" } },
  }
});
const steps = reactive([
  { label: "Search for users" },
  { label: "Select users to add to admin" },
  { label: "Save" }
]);
const activeStep = ref(0);

function step2(u: Array<User>) {
  users.value = u;
  subviews.activate("select");
  activeStep.value = 1;
}

function setHasUsersSelected() {
  activeStep.value = 2;
}

function setHasNoUsersSelected() {
  activeStep.value = 1;
}

function submitSelection(users: Set<User>) {
  console.log("SUB", users)
  for (const user of users) {
    console.log("USER", user.id, user.name)
  }
}

function onCancel() {
  emit("end");
}

function onBack() {
  users.value = [];
  subviews.activate("search")
  activeStep.value = 0;
}
</script>
