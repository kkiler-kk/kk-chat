import { createApp } from 'vue'
import Antd from 'ant-design-vue'
import 'ant-design-vue/dist/reset.css'
import App from './App.vue'
import router from './router'
import { setupStore } from './stores'

const app = createApp(App)
setupStore(app)
app.use(router)
app.use(Antd)
app.mount('#app')
