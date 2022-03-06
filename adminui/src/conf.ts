import { getEnv } from "./env";

function getServerUrl() {
    let url = "";
    if (import.meta.env.DEV) {
        url = "http://localhost:8082";
    }
    /*if (process.env.VUE_APP_SERVER_URL !== undefined) {
        return process.env.VUE_APP_SERVER_URL;
    }*/
    return url;
}

const conf = {
    env: getEnv(),
    quidUrl: getServerUrl(),
    serverUri: getServerUrl(),
    isProduction: import.meta.env.PROD
}

export default conf;