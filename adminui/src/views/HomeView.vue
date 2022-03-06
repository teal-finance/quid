<template>
  <div class="p-5">
    <cards-grid :cards="categories"></cards-grid>
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import categories from "@/const/categories";

import CardsGrid from "@/components/cards/CardsGrid.vue";

export default defineComponent({
  components: {
    CardsGrid,
  },
  setup() {
    return {
      categories,
    };
  },
});
</script>
