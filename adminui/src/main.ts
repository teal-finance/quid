import { createApp } from 'vue'
import App from './App.vue';
import router from "./router";
import PrimeVue from 'primevue/config';
import 'primevue/resources/themes/tailwind-light/theme.css'
import 'primevue/resources/primevue.min.css'
import 'primeicons/primeicons.css'
import './assets/index.css';
import ToastService from 'primevue/toastservice';

createApp(App).use(router).use(PrimeVue).use(ToastService).mount('#app')
