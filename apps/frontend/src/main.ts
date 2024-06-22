import { createApp } from 'vue'
import './assets/index.css'
import App from './App.vue'
import { router } from "@/plugins/vue-router.ts"
import { urqlCLient } from "@/plugins/urql.ts"
import urql from '@urql/vue'
const app = createApp(App)

app
	.use(router)
	.use(urql, urqlCLient)
	.mount('#app')
