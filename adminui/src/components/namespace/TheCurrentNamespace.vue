<template>
  <span class="group" v-if="user.namespace.value.id != 0" @click="unselectNamespace()">
    <span
      class="px-3 py-2 text-sm font-bold rounded-full group-hover:light group-hover:cursor-pointer"
    >
      <i-eos-icons:namespace class="inline-block mr-2 text-xl"></i-eos-icons:namespace>
      <span>{{ user.namespace.value.name }}</span>
      <span class="invisible ml-2 group-hover:visible">
        <i-fontisto:close></i-fontisto:close>
      </span>
    </span>
  </span>
</template>

<script setup lang="ts">
import { notify, user } from '@/state';

function unselectNamespace() {
  user.resetNs();
  notify.done("Namespace unselected")
}
</script>