<template>
  <div class="flex items-center sw-stepper">
    <template v-for="(step, i) in steps">
      <div class="relative flex items-center" :class="getRowClass(i)">
        <div class="stepper-step">
          <slot name="content" :index="i">
            <div v-html="step?.content ?? i + 1"></div>
          </slot>
        </div>
        <div class="stepper-label" :index="i">
          <slot name="label">
            <div v-html="step.label"></div>
          </slot>
        </div>
      </div>
      <div class="stepper-line" v-if="!isLastLoop(i)" :class="getRowClass(i)"></div>
    </template>
  </div>
</template>

<script setup lang="ts">import { toRefs } from 'vue';

const props = defineProps({
  steps: {
    type: Array as () => Array<{ label: string, content?: string }>,
    required: true,
  },
  activeIndex: {
    type: Number,
    required: true,
  }
});

const { steps, activeIndex } = toRefs(props);

function getRowClass(i: number): string {
  let cls = "";
  if (activeIndex.value == i) {
    cls = "active"
  } else if (i < activeIndex.value) {
    cls = "done"
  }
  return cls
}

function isLastLoop(i: number): boolean {
  return (i + 1) == props.steps.length;
}
</script>