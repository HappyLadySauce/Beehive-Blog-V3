import 'virtual:uno.css'
import './app/styles/app.css'

import { createPinia } from 'pinia'
import { createApp } from 'vue'

import App from './app/App.vue'
import { router } from './app/router'
import { i18n, syncDocumentLocale } from './shared/i18n'

const app = createApp(App)

app.use(createPinia())
app.use(i18n)
app.use(router)
syncDocumentLocale()
app.mount('#app')
