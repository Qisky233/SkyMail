import { createApp } from 'vue'
import './style.css'
import App from './App.vue'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'

// 导入element plus
import 'element-plus/dist/index.css'
import ElementPlus from 'element-plus'
import zhCn from 'element-plus/es/locale/lang/zh-cn'

// 导入vue router
import router from './route/index'

import 'quill/dist/quill.snow.css';

const app = createApp(App)

Object.keys(ElementPlusIconsVue).forEach((key) => {
    app.component(key, ElementPlusIconsVue[key])
})
app.use(router)
app.use(ElementPlus, {
    locale: zhCn,
})
app.mount('#app')
