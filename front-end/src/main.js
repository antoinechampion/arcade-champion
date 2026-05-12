import './assets/main.css'

import { createApp } from 'vue'
import App from './App.vue'
import { startGamepadPolling } from './gamepad'

createApp(App).mount('#app')
startGamepadPolling()
