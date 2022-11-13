import conf from "@/conf";
import { notify, user } from './state';
import { useApi, ApiResponse } from "@/packages/restmix/main";
import { UserStatusContract } from "./interface";
import Namespace from "./models/namespace";

const api = useApi({ serverUrl: conf.quidUrl });
api.onResponse(async <T>(res: ApiResponse<T>): Promise<ApiResponse<T>> => {
  console.log("On resp", JSON.stringify(res, null, "  "));
  if (!res.ok) {
    if ([401, 403].includes(res.status)) {
      console.warn("Unauthorized request", res.status, "from", res.url)
    } else if (res.status == 500) {
      console.warn("Server error", res.status, "from", res.url)
    } else {
      console.warn("Error", res.status, "from", res.url)
    }
  }
  return res
});

async function checkStatus(): Promise<{ ok: boolean, status: UserStatusContract }> {
  let _data: UserStatusContract = {} as UserStatusContract;
  const res = await api.get<UserStatusContract>("/status");
  if (res.ok) {
    _data = res.data;
  } else {
    if (res.status == 401) {
      return { ok: false, status: {} as UserStatusContract }
    }
    throw new Error(res.data.toString())
  }
  return { ok: true, status: _data }
}

async function adminLogin(namespace: string, username: string, password: string): Promise<boolean> {
  const payload = {
    namespace: namespace,
    username: username,
    password: password,
  }
  const uri = "/admin_login";
  console.log("Login", uri, JSON.stringify(payload))

  const res = await api.post<Record<string, any>>(uri, payload);
  if (!res.ok) {
    console.log("RESP NOT OK", res);
    if (res.status === 401) {
      notify.warning("Login refused", "The server refused the login, please try again")
      return false
    }
  }
  // process response
  const ns = Namespace.empty();
  ns.name = res.data.name;
  ns.id = res.data.id;
  if (namespace != 'quid') {
    user.changeNs(ns.toTableRow());
  } else {
    user.type.value = "serverAdmin";
    user.adminUrl = "/admin";
    user.resetNs()
  }
  return true
}

export { api, adminLogin, checkStatus }
