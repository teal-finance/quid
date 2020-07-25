<template>
  <div class="mt-3">
    <div>
      <h4 class="text-muted" @click="$router.push('/namespaces')">
        <b-icon-box class="mr-1 text-muted" />Namespaces
      </h4>
      <a @click="action('/namespaces', 'addNamespace')">
        <b-icon-plus class="mr-1 text-muted" />Add
      </a>
    </div>
    <div class="mt-3">
      <h4 class="text-muted" @click="$router.push('/users')">
        <b-icon-person class="mr-1 text-muted" />Users
      </h4>
      <a @click="action('/users', 'addUser')">
        <b-icon-plus class="mr-1 text-muted" />Add
      </a>
    </div>
    <div class="mt-3">
      <h4 class="text-muted" @click="$router.push('/groups')">
        <b-icon-people class="mr-1" />Groups
      </h4>
      <a @click="action('/groups', 'addGroup')">
        <b-icon-plus class="mr-1 text-muted" />Add
      </a>
    </div>
    <!-- div class="mt-3">
      <h4 class="text-muted" @click="$router.push('/tokens')">
        <b-icon-puzzle class="mr-1" />Tokens
      </h4>
      <a @click="a">
        <b-icon-plus class="mr-1 text-muted" />Add
      </a>
    </div-->
  </div>
</template>

<script>
// v-if="$route.path.startsWith('/users')"
export default {
  methods: {
    action: async function (path, name) {
      if (!this.$route.path.startsWith(path)) {
        await this.$router.push(path);
      }
      this.$store.commit("action", name);
    },
    a: function () {},
  },
};
</script>

<style scoped>
h4 router-link,
h4 router-link:active,
h4 router-link:visited {
  color: lime;
}
</style>

