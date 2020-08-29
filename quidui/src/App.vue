<template>
  <div class="main">
    <div v-if="isAuthenticated">
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
import conf from "@/conf";
import TheNavbar from "@/components/TheNavbar.vue";
import TheSidebar from "@/components/TheSidebar.vue";
import TheLogin from "@/components/TheLogin";

export default {
  components: {
    TheNavbar,
    TheSidebar,
    TheLogin,
  },
  computed: {
    ...mapState(["isAuthenticated"]),
  },
  mounted() {
    if (!conf.isProduction) {
      let rt = localStorage.getItem("refreshToken");
      if (rt) {
        this.$api.requests.refreshToken = rt;
        let username = localStorage.getItem("username");
        this.$store.commit("authenticate", username);
      }
    }
    //if (process.env.VUE_APP_DEBUG === "true") {
    //console.log("IS AUTHENTICATED", this.isAuthenticated);
    console.log("ENV", process.env.NODE_ENV);
    //}
  },
};
</script>

<style lang="scss">
@import "./scss/main.scss";
</style>
