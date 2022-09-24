import { EnvType, getEnv } from "./env";

const env = getEnv();

function getServerUrl() {
    let url = "";
    if (env == EnvType.local) {
        url = "http://localhost:8090";
    }
    /*if (process.env.VUE_APP_SERVER_URL !== undefined) {
        return process.env.VUE_APP_SERVER_URL;
    }*/
    return url;
}

const conf = {
    env: env,
    quidUrl: getServerUrl(),
    serverUri: getServerUrl(),
    isProduction: import.meta.env.PROD
}

export default conf;