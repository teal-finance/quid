<template>
  <div class="flex flex-col space-y-5">
    <sw-input v-model:value="form.name.val" v-model:isvalid="form.name.isValid" :validator="form.name.validator"
      inline-label="Name" required autofocus></sw-input>
    <sw-input v-model:value="form.accessTokenTtl.val" v-model:isvalid="form.accessTokenTtl.isValid"
      :validator="form.accessTokenTtl.validator" inline-label="Access tokens max time to live" required></sw-input>
    <sw-input v-model:value="form.refreshTokenTtl.val" v-model:isvalid="form.refreshTokenTtl.isValid"
      :validator="form.refreshTokenTtl.validator" inline-label="Refresh tokens max time to live" required></sw-input>
    <sw-switch id="ns-switch" class="switch-success" v-model:value="enablePublicEndpoint">&nbsp;Enable public endpoint
    </sw-switch>
    <div class="flex flex-row items-center">
      <div class="mr-2">Algorithm:</div>
      <select v-model="algo">
        <option>HS256</option>
        <option>HS384</option>
        <option>HS512</option>
        <option>RS256</option>
        <option>RS384</option>
        <option>RS512</option>
        <option>PS256</option>
        <option>PS384</option>
        <option>PS512</option>
        <option>ES256</option>
        <option>ES384</option>
        <option>ES512</option>
        <option>EDDSA</option>
      </select>
    </div>

    <div class="flex flex-row">
      <button class="w-20 mr-3 btn success" :disabled="!isFormValid === true" @click="postForm()">Save</button>
      <button class="w-20 btn warning" @click="onCancel()">Cancel</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, reactive, ref } from "vue";
import SwInput from "@snowind/input";
import SwSwitch from "@snowind/switch";
import { api } from "@/api";
import { notify } from "@/state";
import { AlgoType } from "@/interface";

const emit = defineEmits(["end"]);

const enablePublicEndpoint = ref(false);
const algo = ref<AlgoType>("HS256");
const form = reactive({
  name: {
    val: "",
    isValid: null,
    // eslint-disable-next-line
    validator: (v: string) => {
      return (v.length >= 3);
    },
  },
  accessTokenTtl: {
    val: "20m",
    isValid: true,
    validator: (v: string) => {
      return (v.length >= 2);
    },
  },
  refreshTokenTtl: {
    val: "24h",
    isValid: true,
    validator: (v: string) => {
      return (v.length >= 2);
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
    await api.post("/admin/namespaces/add", {
      name: form.name.val,
      alg: algo.value,
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
