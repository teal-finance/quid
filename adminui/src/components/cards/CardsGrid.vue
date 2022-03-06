<template>
  <div
    class="grid grid-cols-1 gap-6 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 2xl:grid-cols-6"
  >
    <div v-for="(card, i) in cards" :key="i">
      <router-link v-if="card.url" :to="card.url">
        <simple-card :card="card"></simple-card>
      </router-link>
      <div v-else-if="onCtrlClickCard !== null" v-on:click.ctrl="onCtrlClickCard(card)">
        <simple-card :card="card"></simple-card>
      </div>
      <simple-card v-else :card="card"></simple-card>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import Category from "@/models/category";
import SimpleCard from "./SimpleCard.vue";

export default defineComponent({
  components: { SimpleCard },
  props: {
    cards: {
      type: Object as () => Set<Category>,
      required: true,
    },
    onCtrlClickCard: {
      type: Function,
      default: () => null,
    },
  },
});
</script>