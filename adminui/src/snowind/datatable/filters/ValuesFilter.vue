<template>
  <component
    :is="renderer"
    :col="col"
    :values="vals"
    @include="includeRow($event)"
    @exclude="excludeRow($event)"
  ></component>
</template>

<script lang="ts">
import { defineComponent, toRefs, onMounted, reactive } from 'vue'
import SwDatatableModel from '../models/datatable'
import ValuesFilterSwitchRender from './ValuesFilterSwitchRender.vue';

export default defineComponent({
  components: {
    ValuesFilterSwitchRender
  },
  props: {
    model: {
      type: Object as () => SwDatatableModel,
      required: true,
    },
    col: {
      type: String,
      required: true
    },
    renderer: {
      type: Object,
      default: () => ValuesFilterSwitchRender
    },
  },
  setup(props) {
    const { model, col } = toRefs(props);
    const vals = reactive(new Set());

    function distinctValues() {
      model.value.state.rows.forEach((row) => {
        vals.add(row[col.value]);
      });
    }

    function excludeRow(evt: any) {
      model.value.addExcludeFilter(col.value, evt)
    }

    function includeRow(evt: any) {
      model.value.removeExcludeFilter(col.value, evt)
    }

    onMounted(() => distinctValues());

    return {
      vals,
      includeRow,
      excludeRow,
    }
  }
});
</script>