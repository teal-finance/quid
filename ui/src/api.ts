import conf from "@/conf";
import { notify, user } from './state';
import { useApi } from "restmix";
import { UserStatusContract } from "./interface";
import Namespace from "./models/namespace";

const api = useApi({ serverUrl: conf.quidUrl });

async function checkStatus(): Promise<{ ok: boolean, status: UserStatusContract }> {
  let _data: UserStatusContract = {} as UserStatusContract;
  try {
    _data = await api.get<UserStatusContract>("/status")
  } catch (e: any) {
    if (e.name == "ResponseError") {
      if (e.status == 401) {
        return { ok: false, status: {} as UserStatusContract }
      }
      throw new Error(e?.toString())
    } else {
      throw e
    }
  }
  return { ok: true, status: _data }
}

async function adminLogin(nsName: string, username: string, password: string): Promise<void> {
  const payload = {
    namespace: nsName,
    username: username,
    password: password,
  }
  const opts: RequestInit = {
    method: 'POST',
    credentials: 'include',
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
  const resp = await response.json();
  const ns = Namespace.empty();
  ns.name = resp.name;
  ns.id = resp.id;
  if (nsName != 'quid') {
    user.changeNs(ns.toTableRow());
  } else {
    user.type.value = "serverAdmin";
    user.adminUrl = "/admin";
    user.resetNs()
  }
}

export { api, adminLogin, checkStatus }
