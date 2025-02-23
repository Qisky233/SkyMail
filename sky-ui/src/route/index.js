import { createWebHistory, createRouter } from 'vue-router'

import LoginView from '@/views/login/index.vue'

import HomeView from '@/views/home/index.vue'

const routes = [
  { path: '/', component: LoginView },
  { path: '/home', component: HomeView },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router