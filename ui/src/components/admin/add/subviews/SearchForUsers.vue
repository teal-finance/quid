<template>
  <div class="flex flex-col space-y-5">
    <sw-input
      v-model:value="form.name.val"
      v-model:isvalid="form.name.isValid"
      :validator="form.name.validator"
      inline-label="Username"
      v-on:keyup.enter="search()"
      required
      autofocus
    ></sw-input>
    <div class="flex flex-row">
      <button class="mr-3 btn success" :disabled="!isFormValid === true" @click="search()">Search</button>
      <button class="btn warning" @click="cancel()">Cancel</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, reactive } from "vue";
import SwInput from "@snowind/input";
import { notify, user } from "@/state";
import User from "@/models/user/user";

const emit = defineEmits(["users-found", "cancel"]);

const form = reactive({
  name: {
    val: "",
    isValid: null,
    validator: (v: string) => {
      return (v.length >= 2);
    },
  },
});

const isFormValid = computed<boolean>(() => {
  return form.name.isValid === true
});

function cancel(): void {
  emit("cancel");
}

async function search() {
  if (!isFormValid.value == true) {
    return
  }
  try {
    const users = await User.search(user.namespace.value.id, form.name.val)
    if (users.length == 0) {
      notify.warning("No users found", `No usernames starting with ${form.name.val} found`);
      return
    }
    emit("users-found", users)
    console.log("Users", users)
  } catch (e) {
    console.log(e)
    notify.error("Error creating org")
  }
}
</script>
