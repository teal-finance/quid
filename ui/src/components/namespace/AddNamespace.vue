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
      <div class="mr-2">Signing algorithm:</div>
      <select v-model="algo">
        <option title="HMAC using SHA-256 (SHA-2)">HS256</option>
        <option title="HMAC using SHA-384 (SHA-2)">HS384</option>
        <option title="HMAC using SHA-512 (SHA-2)">HS512</option>
        <option title="Ed25519 signature scheme using SHA-512 and Curve25519">EdDSA</option>
        <option title="ECDSA using P-256 and SHA-256. Deprecated: Filippo Valsorda recommends EdDSA instead.">ES256</option>
        <option title="ECDSA using P-384 and SHA-384. Deprecated: Filippo Valsorda recommends EdDSA instead.">ES384</option>
        <option title="ECDSA using P-521 and SHA-512. Deprecated: Filippo Valsorda recommends EdDSA instead.">ES512</option>
        <option title="RSASSA-PKCS-v1.5 using SHA-256. Deprecated: use EdDSA or PS256 instead (see RFC 8017, section 8)">RS256</option>
        <option title="RSASSA-PKCS-v1.5 using SHA-384. Deprecated: use EdDSA or PS384 instead (see RFC 8017, section 8)">RS384</option>
        <option title="RSASSA-PKCS-v1.5 using SHA-512. Deprecated: use EdDSA or PS512 instead (see RFC 8017, section 8)">RS512</option>
        <option title="RSASSA-PSS using SHA-256 and MGF1 with SHA-256">PS256</option>
        <option title="RSASSA-PSS using SHA-384 and MGF1 with SHA-384">PS384</option>
        <option title="RSASSA-PSS using SHA-512 and MGF1 with SHA-512">PS512</option>
      </select>
    </div>
    <p>
    The Quid team recommends HMAC (HS256, HS384, HS512) when the secret key can be securely shared/stored.
    Else prefer EdDSA: produces small signatures and faster/safer than RSA/ECDSA based algorithms.
    Quid enforces secure use by design: unsigned tokens are rejected.
    No support for encrypted tokens either (use wire encryption instead).
    </p>
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
