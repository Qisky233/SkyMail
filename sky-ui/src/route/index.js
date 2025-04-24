import { createWebHistory, createRouter } from 'vue-router'

import HomeView from '@/views/home/index.vue'
import NavOne from '@/views/navOne/index.vue'
import NavTwo from '@/views/navTwo/index.vue'
import NavThree from '@/views/navThree/index.vue'
import NavFour from '@/views/navFour/index.vue'

const routes = [
  { path: '/', component: NavOne },
  { path: '/send', component: NavOne },
  { path: '/inbox', component: NavTwo },
  { path: '/spam', component: NavThree },
  { path: '/data', component: NavFour },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router