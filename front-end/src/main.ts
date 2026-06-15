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
    { path: '/playing', component: () => import('@/pages/PlayingPage.vue') },
  ],
})

createApp(RouterView).use(router).mount('#app')
startGamepadPolling()

document.body.focus()
document.addEventListener('keydown', (e) => {
  if (['ArrowUp', 'ArrowDown', 'ArrowLeft', 'ArrowRight', ' '].includes(e.key)) {
    const tag = (e.target as HTMLElement).tagName
    if (tag !== 'INPUT' && tag !== 'TEXTAREA') e.preventDefault()
  }
})
