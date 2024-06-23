import { createApp } from 'vue'

import './assets/index.css'

import urql from '@urql/vue'

import { urqlCLient } from '@/plugins/urql.ts'
import { router } from '@/plugins/vue-router.ts'
import App from './App.vue'

const app = createApp(App)

app.use(router).use(urql, urqlCLient).mount('#app')
