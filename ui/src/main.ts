import 'virtual:uno.css'
import './app/styles/app.css'

import { QueryClient, VueQueryPlugin } from '@tanstack/vue-query'
import { createPinia } from 'pinia'
import { createApp } from 'vue'

import App from './app/App.vue'
import { router } from './app/router'
import { i18n, syncDocumentLocale } from './shared/i18n'

const app = createApp(App)
const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 30_000,
      gcTime: 5 * 60_000,
      refetchOnWindowFocus: false,
      retry: 1,
    },
  },
})

app.use(createPinia())
app.use(VueQueryPlugin, { queryClient })
app.use(i18n)
app.use(router)
syncDocumentLocale()
app.mount('#app')
