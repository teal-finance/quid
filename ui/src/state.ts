import { ToastServiceMethods } from "primevue/toastservice";
import { api, checkStatus } from "./api";
import conf from "./conf";
import { EnvType } from "./env";
import { ConfirmOptions, NotifyService } from "./interface";
import SiteUser from "./models/siteuser";
import useNotify from "./notify";
import { useScreenSize } from "@snowind/state";
import Namespace from "./models/namespace";

const user = new SiteUser();
let notify: NotifyService;
const { isMobile, isTablet, isDesktop } = useScreenSize();

async function initState(toast: ToastServiceMethods, confirm: ConfirmOptions) {
  notify = useNotify(toast, confirm)
  await initUserState();
}

async function initUserState() {
  const { ok, status } = await checkStatus();
  if (!ok) {
    console.log("Status unauthorized");
    return
  }
  console.log("Status", status)
  user.isLoggedIn.value = true;
  user.name.value = status["username"];
  const ns = Namespace.empty();
  ns.id = status.ns.id;
  ns.name = status.ns.name;
  if (status.user.admin === true) {
    user.type.value = "serverAdmin";
    user.adminUrl = "/admin";
    user.resetNs()
  } else {
    user.changeNs(ns.toTableRow());
  }
}

export { user, initState, initUserState, notify, isMobile, isTablet, isDesktop }