
const notify = function (vue) {
    return {
        error(content) {
            vue.$bvToast.toast(content, {
                title: "Error",
                variant: "danger",
                noAutoHide: true,
                appendToast: true
            });
        },
        warning({ title, content }) {
            vue.$bvToast.toast(content, {
                title: title,
                variant: "danger",
                autoHideDelay: 5000,
                appendToast: true
            });
        },
        success({ title, content, timeOnScreen = 1500 }) {
            vue.$bvToast.toast(content, {
                title: title,
                variant: "success",
                autoHideDelay: timeOnScreen,
                appendToast: true
            });
        },
        done(content) {
            vue.$bvToast.toast(content, {
                title: "Done",
                variant: "success",
                autoHideDelay: 1500,
                appendToast: true
            });
        }
    }
}

export default notify;