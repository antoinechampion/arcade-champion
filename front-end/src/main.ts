import './assets/main.css'

import { createApp } from 'vue'
import { createRouter, createWebHistory, RouterView } from 'vue-router'
import HomePage from '@/pages/HomePage.vue'
import { startGamepadPolling } from './gamepad'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', component: HomePage },
    { path: '/backoffice', component: () => import('@/pages/BackOfficePage.vue') },
    { path: '/backoffice/add', component: () => import('@/pages/GameFormPage.vue') },
    { path: '/backoffice/edit/:id', component: () => import('@/pages/GameFormPage.vue') },
    { path: '/backoffice/settings', component: () => import('@/pages/SettingsPage.vue') },
  ],
})

createApp(RouterView).use(router).mount('#app')
startGamepadPolling()
