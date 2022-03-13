<template>
  <div class="text-3xl txt-primary">
    Admins
    <button
      class="ml-3 text-2xl border-none btn focus:outline-none txt-neutral"
      @click="collapse = !collapse"
      v-if="!mustSelectNamespace"
    >
      <icon icon="fa6-solid:plus" v-if="collapse === true"></icon>
      <icon icon="fa6-solid:minus" v-else></icon>
    </button>
  </div>
  <div
    :class="{
      'slide-y': true,
      'slideup': collapse === true,
      'slidedown': collapse === false
    }"
    class="my-4"
    v-if="!mustSelectNamespace"
  >
    <div class="p-5 mt-3 border border-light dark:border-light-dark w-96">
      <div class="text-xl">Add an admin</div>
      <add-admin class="mt-5" @end="endAdd()"></add-admin>
    </div>
  </div>
  <div class="w-full" v-else>
    <div class="mt-3 text-2xl">Select a namespace</div>
    <namespace-selector class="mt-5"></namespace-selector>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { Icon } from '@iconify/vue';
import AddAdmin from "@/components/admin/AddAdmin.vue";
import { mustSelectNamespace } from "@/state";
import NamespaceSelector from "@/components/namespace/NamespaceSelector.vue";

const collapse = ref(true);

function endAdd() {
  collapse.value = true;
}
</script>