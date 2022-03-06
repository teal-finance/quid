import { QuidRequests, QuidError } from '@/quidjs'
import conf from "@/conf";
import { EnvType } from './env';

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
  /*onHasToLogin: () => {
    user.isLoggedIn.value = false;
  }*/
});

async function adminLogin(username: string, password: string): Promise<void> {
  const payload = {
    namespace: "quid",
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
  const response = await fetch(uri, opts);
  if (!response.ok) {
    console.log("RESP NOT OK", response);
    if (response.status === 401) {
      throw new QuidError("Admin login refused");
    }
    throw new Error(response.statusText)
  }
  const t = await response.json();
  console.log("T", t)
  requests.refreshToken = t.token;
}

export { requests, adminLogin }