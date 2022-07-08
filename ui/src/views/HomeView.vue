<template>
  <div class="p-5">
    <cards-grid :cards="categories"></cards-grid>
  </div>
</template>

<script setup lang="ts">
import { user } from "@/state";
import CardsGrid from "@/components/widgets/cards/CardsGrid.vue";
import { serverAdminCategories, nsAdminCategories } from "@/const/categories";

const categories = user.type.value == "serverAdmin" ? serverAdminCategories : nsAdminCategories;
</script>
