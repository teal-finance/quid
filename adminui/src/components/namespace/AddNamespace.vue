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
      v-model:value="form.accessTokenTtl.val"
      v-model:isvalid="form.accessTokenTtl.isValid"
      :validator="form.accessTokenTtl.validator"
      inline-label="Access tokens max time to live"
      required
    ></sw-input>
    <sw-input
      v-model:value="form.refreshTokenTtl.val"
      v-model:isvalid="form.refreshTokenTtl.isValid"
      :validator="form.refreshTokenTtl.validator"
      inline-label="Refresh tokens max time to live"
      required
    ></sw-input>
    <sw-switch
      id="ns-switch"
      class="switch-success"
      :v-model:value="enablePublicEndpoint"
    >&nbsp;Enable public endpoint</sw-switch>
    <div class="flex flex-row">
      <button
        class="w-20 mr-3 btn success"
        :disabled="!isFormValid === true"
        @click="postForm()"
      >Save</button>
      <button class="w-20 btn warning" @click="onCancel()">Cancel</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, reactive, ref } from "vue";
import SwInput from "@snowind/input";
import SwSwitch from "@snowind/switch";
import { requests } from "@/api";
import { notify } from "@/state";

const emit = defineEmits(["end"]);

const enablePublicEndpoint = ref(false);
const form = reactive({
  name: {
    val: "",
    isValid: null,
    // eslint-disable-next-line
    validator: (v: string) => {
      if (v.length >= 3) {
        return true;
      }
      return false;
    },
  },
  accessTokenTtl: {
    val: "20m",
    isValid: true,
    validator: (v: string) => {
      if (v.length >= 2) {
        return true;
      }
      return false;
    },
  },
  refreshTokenTtl: {
    val: "24h",
    isValid: true,
    validator: (v: string) => {
      if (v.length >= 2) {
        return true;
      }
      return false;
    },
  },
});

const isFormValid = computed<boolean>(() => {
  return form.name.isValid === true
    && form.accessTokenTtl.isValid === true
    && form.refreshTokenTtl.isValid === true
});

function onCancel(): void {
  emit("end");
}

async function postForm() {
  try {
    await requests.post("/admin/namespaces/add", {
      name: form.name.val,
      max_ttl: form.accessTokenTtl.val,
      refresh_max_ttl: form.refreshTokenTtl.val,
      enable_endpoint: enablePublicEndpoint.value,
    });
    emit("end");
    notify.done("Namespace added")
  } catch (e) {
    console.log(e)
    notify.error("Error creating namespace")
  }
}

</script>
