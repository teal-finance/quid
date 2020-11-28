<template>
  <div>
    <b-navbar toggleable="lg" type="dark" variant="dark">
      <b-navbar-brand tag="router-link" :to="{ path: '/' }"
        >Quid</b-navbar-brand
      >
      <b-navbar-toggle target="nav-collapse"></b-navbar-toggle>
      <b-collapse id="nav-collapse" is-nav>
        <b-navbar-nav class="ml-auto">
          <b-nav-item tag="router-link" :to="{ path: '/namespaces' }"
            >Namespaces</b-nav-item
          >
          <b-nav-item tag="router-link" :to="{ path: '/users' }"
            >Users</b-nav-item
          >
          <b-nav-item tag="router-link" :to="{ path: '/groups' }"
            >Groups</b-nav-item
          >
          <b-nav-item tag="router-link" :to="{ path: '/orgs' }"
            >Orgs</b-nav-item
          >
          <!-- b-nav-item tag="router-link" :to="{ path: '/tokens' }">Tokens</b-nav-item -->
          <b-button size="sm" class="ml-3 my-2 my-sm-0" @click="logout"
            >Logout</b-button
          >
        </b-navbar-nav>
      </b-collapse>
    </b-navbar>
  </div>
</template>

<script>
export default {
  methods: {
    async logout() {
      this.$store.commit("unauthenticate");
      if (!this.$router.currentRoute.path !== "/") {
        this.$router.push({ path: "/" });
      }
      await this.$api.get("/admin/logout");
    },
  },
};
</script>
