import { fileURLToPath, URL } from 'node:url'

import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'
import tailwindcss from '@tailwindcss/vite'

// Proxies external images through localhost so the cropper's <canvas> can read
// their pixels without being tainted by cross-origin restrictions.
function imageProxyPlugin() {
  return {
    name: 'image-proxy',
    configureServer(server) {
      server.middlewares.use('/api/local-dev-image-proxy', async (req, res) => {
        const { searchParams } = new URL(req.url, 'http://localhost')
        const target = searchParams.get('url')
        if (!target) { res.writeHead(400); res.end(); return }
        try {
          const response = await fetch(target)
          res.writeHead(response.status, {
            'content-type': response.headers.get('content-type') || 'image/jpeg',
            'cache-control': 'public, max-age=86400',
          })
          const buffer = Buffer.from(await response.arrayBuffer())
          res.end(buffer)
        } catch {
          res.writeHead(502); res.end()
        }
      })
    },
  }
}

// https://vite.dev/config/
export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd())
  const backendUrl = env.VITE_BACKEND_URL

  return {
    plugins: [
      vue(),
      vueDevTools(),
      tailwindcss(),
      imageProxyPlugin(),
    ],
    resolve: {
      alias: {
        '@': fileURLToPath(new URL('./src', import.meta.url))
      },
    },
    server: {
      proxy: {
        '/api': {
          target: backendUrl,
        },
        '/images': {
          target: backendUrl,
        },
      },
    },
    test: {
      environment: 'jsdom',
    },
  }
})
