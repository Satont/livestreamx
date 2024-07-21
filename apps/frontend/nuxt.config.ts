import path from 'node:path'
import { vite as vidstack } from 'vidstack/plugins'
import { watch } from 'vite-plugin-watch'

// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2024-04-03',
  devtools: {
    enabled: true,

    timeline: {
      enabled: true
    }
  },

  modules: [
    '@nuxtjs/tailwindcss',
    '@nuxtjs/color-mode',
    '@vueuse/nuxt',
    '@nuxt/icon',
    '@bicou/nuxt-urql',
    '@nuxt-alt/proxy',
    '@nuxt/image',
    '@nuxt/fonts'
  ],

  tailwindcss: { exposeConfig: true, editorSupport: true },
  build: { transpile: ['vue-sonner'] },
  colorMode: { classSuffix: '' },
  css: [
    '~/assets/css/global.css',
    'notivue/notifications.css',
    'notivue/animations.css'
  ],
  imports: {
    imports: [
      {
        from: 'tailwind-variants',
        name: 'tv'
      },
      {
        from: 'tailwind-variants',
        name: 'VariantProps',
        type: true
      },
      {
        from: 'vue-sonner',
        name: 'toast',
        as: 'useSonner'
      }
    ]
  },
  vue: {
    compilerOptions: {
      isCustomElement: (tag) => tag.startsWith('media-')
    }
  },
  vite: {
    plugins: [
      watch({
        onInit: true,
        pattern: 'api/**/*.ts',
        command: 'graphql-codegen'
      }),
      vidstack()

    ],
    resolve: {
      alias: {
        '@': path.resolve(__dirname)
      }
    }
  },

  proxy: {
    debug: true,
    experimental: {
      listener: true
    },
    proxies: {
      '/api': {
        target: 'http://127.0.0.1:1337',
        changeHost: false,
        rewrite: (path) => path.replace(/^\/api/, ''),
        ws: true
      }
    }
  },

  devServer: {
    port: 5173
  },

  nitro: {
    devProxy: {
      '/api': {
        target: 'http://127.0.0.1:1337',
        changeOrigin: true,
        ws: true
      }
    }
  },

  notivue: {
    enqueue: true,
    pauseOnHover: true,
    pauseOnTabChange: true,
    position: 'top-right',
    teleportTo: 'body'
  },

  app: {
    head: {
      title: 'Streamx'
    }
  },

  urql: {
    endpoint:
      process.env.NODE_ENV === 'production'
        ? 'https://streamx.satont.dev/api/query'
        : 'http://localhost:5173/api/query',
    client: './configs/urql.ts'
  },

  icon: {
    localApiEndpoint: '/localapi/icons'
  }
})