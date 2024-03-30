import './assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'
import store from "./store";
// 导入antv组件库
import antDesign from 'ant-design-vue';
import "ant-design-vue/dist/antd.css";
const app = createApp(App)

app.use(antDesign)
app.use(createPinia())
app.use(router)


app.mount('#app')
