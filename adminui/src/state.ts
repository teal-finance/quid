import { requests } from "./api";
import conf from "./conf";
import { EnvType } from "./env";
import User from "./models/user";


const user = new User();

function initState(): void {
  console.log("Running in env", conf.env);
  if (conf.env == EnvType.local) {
    let t = import.meta.env.VITE_DEV_TOKEN;
    console.log("ENVT", t)
    if (t) {
      user.devRefreshToken = t.toString();
      user.name.value = "devuser";
      console.log("Logging in user from dev token")
      requests.refreshToken = user.devRefreshToken;
      user.isLoggedIn.value = true;
    }
  }
}

export { user, initState }