<template>
  <div class="flex flex-col space-y-5">
    <sw-input v-model:value="form.name.val" v-model:isvalid="form.name.isValid" :validator="form.name.validator"
      inline-label="Name" required autofocus></sw-input>
    <div class="flex flex-row">
      <button class="w-20 mr-3 btn success" :disabled="!isFormValid === true" @click="postForm()">Save</button>
      <button class="w-20 btn warning" @click="onCancel()">Cancel</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, reactive } from "vue";
import SwInput from "@snowind/input";
import { api } from "@/api";
import { notify, user } from "@/state";

const emit = defineEmits(["end"]);

const form = reactive({
  name: {
    val: "",
    isValid: null,
    validator: (v: string) => {
      return (v.length >= 3);
    },
  },
});

const isFormValid = computed<boolean>(() => {
  return form.name.isValid === true
});

function onCancel(): void {
  emit("end");
}

async function postForm() {
  try {
    await api.post(user.adminUrl + "/groups/add", {
      name: form.name.val,
      namespace_id: user.namespace.value.id,
    });
    emit("end");
    notify.done("Group added")
  } catch (e) {
    console.log(e)
    notify.error("Error creating group")
  }
}

</script>
