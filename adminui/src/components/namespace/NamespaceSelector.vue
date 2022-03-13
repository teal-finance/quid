<template>
  <div>
    <loading-indicator v-if="isLoading"></loading-indicator>
    <div v-else class="flex flex-wrap space-x-1">
      <div v-for="ns in namespaces">
        <SimpleBadge
          v-if="ns.name != 'quid'"
          :text="ns.name"
          class="mr-2 cursor-pointer"
          @click="selectNamespace(ns)"
        ></SimpleBadge>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue";
import Namespace from "@/models/namespace";
import NamespaceTable from "@/models/namespace/interface";
import SimpleBadge from "../widgets/SimpleBadge.vue";
import LoadingIndicator from "@/components/widgets/LoadingIndicator.vue";
import { namespaceMutations, notify } from "@/state";

const isLoading = ref(false);
const namespaces = ref(new Array<NamespaceTable>());

async function fetchNamespaces() {
  const ns = await Namespace.fetchAll();
  namespaces.value = Array.from(ns);
  //console.log("DATA", namespaces.value)
}

function selectNamespace(ns: NamespaceTable) {
  namespaceMutations.change(Namespace.fromNamespaceTable(ns));
  notify.toastInfo("Namespace selected")
}

onMounted(() => fetchNamespaces())
</script>