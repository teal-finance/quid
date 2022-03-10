import { ToastServiceMethods } from "primevue/toastservice";
import { requests } from "./api";
import conf from "./conf";
import { EnvType } from "./env";
import User from "./models/user";
import useNotify from "./notify";


const user = new User();
let notify: {
  error: (content: string) => void;
  warning: (title: string, content: string, timeOnScreen?: number) => void;
  success: (title: string, content: string, timeOnScreen?: number) => void;
  done(content: string): void;
};

function initState(toast: ToastServiceMethods): void {
  console.log("Running in env", conf.env);
  if (conf.env == EnvType.local) {
    let t = import.meta.env.VITE_DEV_TOKEN;
    if (t) {
      user.devRefreshToken = t.toString();
      user.name.value = "devuser";
      console.log("Logging in user from dev token")
      requests.refreshToken = user.devRefreshToken;
      user.isLoggedIn.value = true;
    }
  }
  notify = useNotify(toast)
}

export { user, initState, notify }