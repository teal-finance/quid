<template>
  <div class="container p-3 mx-auto">
    <render-md-file :file-url="`/client/${file}.md`" :hljs="hljs"></render-md-file>
  </div>
</template>

<script setup lang="ts">
import { ref, watchEffect } from 'vue';
import RenderMdFile from '@/widgets/RenderMdFile.vue';
import router from '@/router';
import { hljs } from "@/conf";

const file = ref("")

watchEffect(() => {
  file.value = router.currentRoute.value.params?.file.toString() ?? ""
})
</script>

