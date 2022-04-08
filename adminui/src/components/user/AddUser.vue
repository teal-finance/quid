<template>
  <div class="flex flex-col space-y-5">
    <sw-input
      v-model:value="form.name.val"
      v-model:isvalid="form.name.isValid"
      :validator="form.name.validator"
      inline-label="Name"
      required
      autofocus
    ></sw-input>
    <sw-input
      class="mt-3"
      :type="showPasswordFields ? 'text' : 'password'"
      v-model:value="form.pwd.val"
      v-model:isvalid="form.pwd.isValid"
      :validator="form.pwd.validator"
      inline-label="Password"
      required
    ></sw-input>
    <sw-input
      class="mt-3"
      :type="showPasswordFields ? 'text' : 'password'"
      v-model:value="form.pwdVerif.val"
      v-model:isvalid="form.pwdVerif.isValid"
      :validator="form.pwdVerif.validator"
      inline-label="Password again"
      required
    ></sw-input>
    <div class="flex flex-row">
      <button
        class="w-20 mr-3 btn success"
        :disabled="!isFormValid === true"
        @click="submitForm()"
      >Save</button>
      <button class="w-20 btn warning" @click="onCancel()">Cancel</button>
      <div class="inline-block ml-2">
        <button
          class="btn lighter"
          v-if="showPasswordFields"
          @click="showPasswordFields = false"
        >Hide password</button>
        <button class="btn lighter" v-else @click="showPasswordFields = true">Show password</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, reactive, ref } from "vue";
import SwInput from "@snowind/input";
import { requests } from "@/api";
import { notify, user } from "@/state";

const emit = defineEmits(["end"]);
const showPasswordFields = ref(false)

const form = reactive({
  name: {
    val: "",
    isValid: null,
    validator: (v: string) => {
      if (v.length >= 3) {
        return true;
      }
      return false;
    },
  },
  pwd: {
    val: "",
    isValid: null,
    validator: (v: string) => {
      if (v.length >= 5) {
        return true;
      }
      return false;
    },
  },
  pwdVerif: {
    val: "",
    isValid: null,
    validator: (v: string) => {
      if (v.length >= 5) {
        return true;
      }
      return false;
    },
  },
});

const isFormValid = computed<boolean>(() => {
  return form.name.isValid === true && form.pwd.isValid === true && form.pwdVerif.isValid === true;
});

function submitForm() {
  if (form.pwd.val !== form.pwdVerif.val) {
    notify.warning("Password mismatch", "Please verify the typing");
    return
  }
  postForm();
}

function resetForm() {
  form.name.val = ""
}

async function postForm() {
  try {
    await requests.post("/admin/users/add", {
      name: form.name.val,
      password: form.pwd.val,
      namespace_id: user.namespace.value.id,
    });
    emit("end");
    notify.done("User added")
  } catch (e) {
    console.log(e)
    notify.error("Error creating user")
  }
}

function onCancel(): void {
  emit("end");
}
</script>
