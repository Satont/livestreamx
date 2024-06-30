import { createRouter, createWebHistory } from 'vue-router'

export const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      component: () => import('../components/layout/layout.vue'),
      children: [
        {
          path: '/',
          component: () => import('../pages/index.vue')
        },
        {
          path: '/dashboard',
          component: () => import('../pages/dashboard.vue')
        },
        {
          path: '/:channelName',
          name: 'Channel',
          component: () => import('../pages/stream.vue')
        }
      ]
    }
  ]
})
