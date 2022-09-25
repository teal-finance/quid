<template>
  <div class="container p-3 mx-auto">
    <render-md-file file-url="/README.md" :hljs="hljs"></render-md-file>
  </div>
</template>

<script setup lang="ts">
import RenderMdFile from '@/widgets/RenderMdFile.vue';
import { hljs } from "@/conf";
</script>
