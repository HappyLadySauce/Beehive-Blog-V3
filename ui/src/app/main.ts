import 'virtual:uno.css';
import '@/app/styles/app.css';

import { createPinia } from 'pinia';
import { createApp } from 'vue';

import App from '@/app/App.vue';
import { router } from '@/app/router';

const app = createApp(App);

app.use(createPinia());
app.use(router);
app.mount('#app');
