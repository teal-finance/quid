function getQuidUrl() {
    let url = "";
    if (process.env.NODE_ENV === 'development') {
        url = "http://localhost:8082";
    }
    if (process.env.VUE_APP_SERVER_URL !== undefined) {
        return process.env.VUE_APP_SERVER_URL;
    }
    return url;
}

const Conf = { quidUrl: getQuidUrl(), isProduction: process.env.NODE_ENV === 'production' }

export default Conf;