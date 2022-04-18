import { QuidRequests, QuidRequestError } from '@/packages/quidjs'
import conf from "@/conf";
import { EnvType } from './env';
import { notify, user } from './state';

const requests = new QuidRequests({
  namespace: "quid",
  timeouts: {
    accessToken: "5m",
    refreshToken: "24h"
  },
  quidUri: conf.quidUrl,
  serverUri: conf.serverUri,
  accessTokenUri: conf.serverUri + "/admin_token/access/",
  verbose: conf.env === EnvType.local,
  onHasToLogin: async () => {
    notify.warning("Connection expired", "Please login again")
    user.isLoggedIn.value = false;
  }
});

async function adminLogin(namespace: string, username: string, password: string): Promise<void> {
  const payload = {
    namespace: namespace,
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
  const t = await response.json();
  console.log("T", t)
  requests.refreshToken = t.token;
}

export { requests, adminLogin }