<template>
  <div v-if="isReady">
    <render-md :source="code" :hljs="hljs"></render-md>
  </div>
</template>

<script setup lang="ts">
import "highlight.js/styles/stackoverflow-dark.css";
import RenderMd from '@/widgets/RenderMd.vue';
import { api } from "@/state";
import { ref, watchEffect } from "vue";
import _hljs from 'highlight.js/lib/core';

const props = defineProps({
  fileUrl: {
    type: String,
    required: true
  },
  hljs: {
    type: Object as () => typeof _hljs,
    required: true
  }
})

const code = ref("");
const isReady = ref(false);

async function load() {
  isReady.value = false;
  code.value = await api.get<string>(props.fileUrl);
  isReady.value = true;
}

watchEffect(() => load())
</script>