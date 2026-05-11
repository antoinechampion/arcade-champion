import './assets/main.css'

import { createApp } from 'vue'
import App from './App.vue'
import { GamepadPlugin } from './plugins/gamepad.ts'

createApp(App)
    .use(GamepadPlugin)
    .mount('#app')
