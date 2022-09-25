<template>
  <div class="container mx-auto">
    <div class="p-5" v-if="isReady">
      <code-editor :code="initialCode" lang="javascript" @edit="codeChange($event)" :hljs="hljs"></code-editor>
      <button class="mt-3 btn secondary" @click="runCode()">Run</button>
      <div class="mt-5">{{ result }}</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watchEffect } from 'vue';
import { CodeEditor } from "vuecodit";
import "vuecodit/style.css";
import "highlight.js/styles/stackoverflow-light.css";
import { api } from '@/state';
import { hljs } from "@/conf"
import router from "@/router";
import { examplesExtension } from '@/conf';

const isReady = ref(false);
let initialCode = ref("");
let editedCode = ref("");
const result = ref("");

function codeChange(e: string) {
  // update the code
  editedCode.value = e;
}
function runCode() {
  // execute the code
  result.value = eval(editedCode.value)
}

async function load() {
  isReady.value = false;
  const file = router.currentRoute.value.params?.file.toString() ?? "";
  initialCode.value = await api.get<string>(`/examples/${file + examplesExtension}`);
  editedCode.value = initialCode.value;
  isReady.value = true;
}

watchEffect(() => load())
</script>