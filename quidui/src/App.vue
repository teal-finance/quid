<template>
  <div class="main">
    <div v-if="isAuthorized()">
      <navbar></navbar>
      <b-container fluid class="h-100">
        <b-row style="align-items:stretch;height:100%">
          <b-col class="bg-light">
            <sidebar></sidebar>
          </b-col>
          <b-col cols="10">
            <router-view></router-view>
          </b-col>
        </b-row>
      </b-container>
    </div>
    <div class="vertical-center" v-else>
      <div class="inner-block">
        <login></login>
      </div>
    </div>
  </div>
</template>

<script>
import { mapState } from "vuex";
import Navbar from "@/components/Navbar.vue";
import Sidebar from "@/components/Sidebar.vue";
import Login from "@/views/Login";

export default {
  components: {
    Navbar,
    Sidebar,
    Login
  },
  methods: {
    isAuthorized() {
      if (!this.isProduction) {
        if (this.isDevModeEnabled) {
          return true;
        }
      }
      if (this.isAuthenticated) {
        return true;
      }
      return false;
    }
  },
  computed: {
    ...mapState(["isAuthenticated"]),
    isDevModeEnabled() {
      if (process.env.VUE_APP_ENABLE_DEV_MODE !== undefined) {
        return process.env.VUE_APP_ENABLE_DEV_MODE;
      }
      return false;
    },
    isProduction() {
      return process.env.NODE_ENV === "production";
    }
  },
  mounted() {
    console.log("IS AUTHENTICATED", this.isAuthenticated);
    console.log("IS PRODUCTION", this.isProduction);
    console.log("IS DEV_MODE", this.isDevModeEnabled);
  }
};
</script>

<style lang="scss">
@import "./scss/main.scss";
</style>
