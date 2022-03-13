
import { ToastServiceMethods } from "primevue/toastservice";
import { computed, reactive, ref } from "vue";
import { requests } from "./api";
import conf from "./conf";
import { EnvType } from "./env";
import { ConfirmOptions, NotifyService } from "./interface";
import Namespace from "./models/namespace";
import User from "./models/user";
import useNotify from "./notify";
import { PopToast } from "./type";

const user = new User();
let notify: NotifyService;
const state = reactive({
  namespace: Namespace.empty(),
});

const namespaceMutations = {
  change: (ns: Namespace) => {
    state.namespace = ns;
  },
  reset: () => {
    state.namespace = Namespace.empty();
  }
}

const mustSelectNamespace = computed<boolean>(() => {
  return state.namespace.id == 0
  //&& user.type == "serverAdmin";
});

function initState(toast: ToastServiceMethods, confirm: ConfirmOptions, popToast: PopToast): void {
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
  notify = useNotify(toast, confirm, popToast)
}

export { user, initState, notify, state, namespaceMutations, mustSelectNamespace }