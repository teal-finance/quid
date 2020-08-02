<template>
  <b-form-input
    :value="value"
    v-on:keyup.enter="submit"
    v-on:keydown.escape="$emit('end-edit',null)"
    v-model="newValue"
    style="display:inline-block;width:52px"
    size="sm"
    autofocus
  ></b-form-input>
</template>

<script>
export default {
  data() {
    return {
      newValue: this.value,
    };
  },
  props: {
    namespaceId: {
      type: Number,
      required: true,
    },
    value: {
      type: String,
      required: true,
    },
  },
  methods: {
    async submit() {
      let { error } = await this.$api.post("/admin/namespaces/max-ttl", {
        id: this.namespaceId,
        max_ttl: this.newValue,
      });
      if (error == null) {
        this.$notify.done("Namespace max access token ttl set");
      }
      this.$emit("end-edit", this.value);
    },
  },
};
</script>