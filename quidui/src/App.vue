<template>
  <div class="main">
    <div v-if="isAuthorized()">
      <the-navbar></the-navbar>
      <b-container fluid class="h-100">
        <b-row class="full-height-content-zone">
          <b-col class="bg-light">
            <the-sidebar></the-sidebar>
          </b-col>
          <b-col cols="10">
            <router-view></router-view>
          </b-col>
        </b-row>
      </b-container>
    </div>
    <div class="vertical-center" v-else>
      <div class="inner-block">
        <the-login></the-login>
      </div>
    </div>
  </div>
</template>

<script>
import { mapState } from "vuex";
import TheNavbar from "@/components/TheNavbar.vue";
import TheSidebar from "@/components/TheSidebar.vue";
import TheLogin from "@/components/TheLogin";

export default {
  components: {
    TheNavbar,
    TheSidebar,
    TheLogin,
  },
  methods: {
    isAuthorized() {
      if (!this.isProduction) {
        if (this.isDevModeEnabled === true) {
          console.log("Dev mode is on");
          return true;
        }
      }
      if (this.isAuthenticated) {
        return true;
      }
      return false;
    },
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
    },
  },
  mounted() {
    console.log("IS AUTHENTICATED", this.isAuthenticated);
    console.log("IS PRODUCTION", this.isProduction);
    console.log("IS DEV_MODE", this.isDevModeEnabled);
  },
};
</script>

<style lang="scss">
@import "./scss/main.scss";
</style>
