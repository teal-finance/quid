<template>
  <div v-html="content" class="prose md-content dark:prose-invert"></div>
</template>

<script setup lang="ts">
import hljs from 'highlight.js/lib/core';
import { computed } from 'vue';
import MarkdownIt from 'markdown-it';

const props = defineProps({
  source: {
    type: String,
    required: true,
  },
  hljs: {
    type: Object as () => typeof hljs,
    required: true
  }
});
const md = new MarkdownIt({
  html: true,
  highlight: function (str, lang) {
    if (lang && props.hljs.getLanguage(lang)) {

      try {
        return hljs.highlight(str, { language: lang }).value;
      } catch (e) {
        throw new Error(`Code parse error ${e}`)
      }
    }
    return ''; // use external default escaping
  }
})

const content = computed(() => {
  const src = props.source
  let res = md?.render(src)
  return res
})
</script>

<style lang="sass">
  .md-content
    max-width: 120ch
</style>