import axios from 'axios'
import { createApp } from 'vue'
import App from './App.vue'
import router from './router'

// Configure axios to always send cookies
axios.defaults.withCredentials = true

const app = createApp(App)
app.use(router)
app.mount('#app')
