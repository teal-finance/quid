import { ToastServiceMethods } from "primevue/toastservice";
import { requests } from "./api";
import conf from "./conf";
import { EnvType } from "./env";
import { ConfirmOptions, NotifyService } from "./interface";
import SiteUser from "./models/siteuser";
import useNotify from "./notify";
import { useScreenSize } from "@snowind/state";

const user = new SiteUser();
let notify: NotifyService;
const { isMobile, isTablet, isDesktop } = useScreenSize();

function initState(toast: ToastServiceMethods, confirm: ConfirmOptions): void {
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
  notify = useNotify(toast, confirm)
  console.log("NS", user.namespace.value)
}

export { user, initState, notify, isMobile, isTablet, isDesktop }