<template>
  <div>
    <loading-indicator v-if="isLoading"></loading-indicator>
    <div v-else class="flex flex-wrap space-x-1">
      <div v-for="ns in namespaces">
        <SimpleBadge
          v-if="ns.name != 'quid'"
          :text="ns.name"
          class="mr-2 cursor-pointer primary"
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
import { user, notify } from "@/state";

const isLoading = ref(false);
const namespaces = ref(new Array<NamespaceTable>());

const emit = defineEmits(["selectns"]);

async function fetchNamespaces() {
  const ns = await Namespace.fetchAll();
  namespaces.value = Array.from(ns);
  //console.log("DATA", namespaces.value)
}

function selectNamespace(nst: NamespaceTable) {
  user.changeNs(nst);
  emit("selectns");
  notify.done("Namespace selected")
}

onMounted(() => fetchNamespaces())
</script>