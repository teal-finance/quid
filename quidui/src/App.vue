<template>
  <div class="main">
    <navbar v-if="isAuthenticated"></navbar>
    <b-container fluid v-if="isAuthenticated" class="h-100">
      <b-row style="align-items: stretch;height:100%">
        <b-col class="bg-light">
          <sidebar></sidebar>
        </b-col>
        <b-col cols="10">
          <router-view></router-view>
        </b-col>
      </b-row>
    </b-container>
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
  created: function() {
    /* let exists = this.$cookies.isKey("quidkey");
    if (exists) {
      let key = this.$cookies.get("quidkey");
      this.$store.commit("authenticate", key);
    }*/
  },
  computed: mapState(["isAuthenticated"])
};
</script>

<style lang="scss">
* {
  box-sizing: border-box;
}

body {
  min-height: 100vh;
  display: flex;
  font-weight: 400;
}

body,
html,
.main,
.vertical-center {
  width: 100%;
  height: 100%;
}

.vertical-center {
  display: flex;
  justify-content: center;
  flex-direction: column;
}

.inner-block {
  width: 450px;
  margin: auto;
  background: #ffffff;
  box-shadow: 0px 14px 80px rgba(34, 35, 58, 0.2);
  padding: 40px 55px 45px 55px;
  border-radius: 15px;
  transition: all 0.3s;
}
.sidebar {
  background-color: lightskyblue;
}
</style>
