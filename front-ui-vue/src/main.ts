import './assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'
import { setupStore } from '@/store';
// 导入antv组件库
import antDesign from 'ant-design-vue';
import "ant-design-vue/dist/antd.css";
// 倒入 emoji
import Vue3EmojiPicker from 'vue3-emoji-picker' 
//挂载websocket
const app = createApp(App)

app.use(antDesign)
app.use(createPinia())
app.use(router)
app.component('Vue3EmojiPicker', Vue3EmojiPicker)
// setupStore(app);

// 挂载状态管理
app.mount('#app')
