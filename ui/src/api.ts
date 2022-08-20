import conf from "@/conf";
import { notify } from './state';
//import { useApi } from "@snowind/api";
import { useApi } from "@/packages/api";
import { ResponseError } from "./packages/errors";

const api = useApi(conf.quidUrl);

async function checkStatus(): Promise<{ ok: boolean, status: Record<string, any> }> {
  let _data = {};
  try {
    _data = await api.get("/status")
  } catch (e) {
    if (e instanceof ResponseError) {
      console.log("Response error", e);
      if (e.response.status == 401) {
        return { ok: false, status: _data }
      }
      throw new Error(e.toString())
    }
  }
  return { ok: true, status: _data }
}

async function adminLogin(namespaceName: string, username: string, password: string): Promise<void> {
  const payload = {
    namespace: namespaceName,
    username: username,
    password: password,
  }
  const opts: RequestInit = {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    } as HeadersInit,
    body: JSON.stringify(payload)
  };
  const uri = conf.quidUrl + "/admin_login";
  console.log("Login", uri, JSON.stringify(payload))
  const response = await fetch(uri, opts);
  if (!response.ok) {
    console.log("RESP NOT OK", response);
    if (response.status === 401) {
      notify.warning("Login refused", "The server refused the login, please try again")
    }
    throw new Error(response.statusText)
  }
  //const resp = await response.json();
  //console.log("RESP", resp)
  /*namespace.name = namespaceName;
  if (namespaceName != 'quid') {
    user.changeNs(namespace.toTableRow());
  } else {
    user.type.value = "serverAdmin";
    user.adminUrl = "/admin";
    user.resetNs()
  }*/
}

export { api, adminLogin, checkStatus }