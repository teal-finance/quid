export default function quidException({ error, hasToLogin = false, unauthorized = false }) {
    return {
        hasToLogin: hasToLogin,
        error: error,
        unauthorized: unauthorized
    }
}