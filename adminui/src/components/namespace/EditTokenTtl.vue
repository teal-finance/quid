<template>
  <div v-if="!isEditing">
    <div class="cursor-pointer group" @click="startEdit()">
      <i-mdi:clock-edit-outline class="hidden group-hover:inline-block txt-neutral"></i-mdi:clock-edit-outline>
      <i-mdi:clock-outline class="group-hover:hidden txt-light"></i-mdi:clock-outline>
      <span v-html="ttl" class="ml-2"></span>
    </div>
  </div>
  <div v-else>
    <i-line-md:cancel class="inline-block txt-neutral" @click="endEdit()"></i-line-md:cancel>
    <span class="ml-2">
      <sw-input
        class="inline-block ttl-inline"
        v-model:value="form.ttlval.val"
        v-model:isvalid="form.ttlval.isValid"
        :validator="form.ttlval.validator"
        :autofocus="true"
        required
      ></sw-input>
      <i-ant-design:save-outlined
        class="ml-2 text-xl"
        :class="form.ttlval.isValid ? ['txt-success', 'cursor-pointer'] : ['txt-light']"
        @click="save()"
      ></i-ant-design:save-outlined>
    </span>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue';
import SwInput from "@snowind/input";
import Namespace from '@/models/namespace';
import { notify } from '@/state';

const props = defineProps({
  id: {
    type: Number,
    required: true,
  },
  ttl: {
    type: String,
    required: true,
  },
  tokenType: {
    type: String as () => "refresh" | "access",
    required: true,
  }
});

const emit = defineEmits(["end"]);

const form = reactive({
  ttlval: {
    val: props.ttl,
    isValid: null,
    validator: (v: string) => v.length >= 2
  }
});
const isEditing = ref(false);

function startEdit() {
  isEditing.value = true;
}

function endEdit() {
  isEditing.value = false;
}

async function save() {
  if (!form.ttlval.isValid) {
    return
  }
  isEditing.value = false;
  if (props.tokenType == "refresh") {
    await Namespace.saveMaxRefreshTokenTtl(props.id, form.ttlval.val);
  } else {
    await Namespace.saveMaxAccessTokenTtl(props.id, form.ttlval.val);
  }
  notify.done("Ttl changed");
  emit("end", form.ttlval.val);
}
</script>

<style lang="sass">
.ttl-inline
  & input
    @apply w-16 border-none
</style>