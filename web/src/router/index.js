import { createRouter, createWebHistory } from 'vue-router'
import { qiankunWindow } from 'vite-plugin-qiankun/dist/helper'

const routes = [
  {
    path: '/',
    redirect: '/mail/inbox' // Default to inbox
  },
  {
    path: '/mail',
    redirect: '/mail/inbox'
  },
  {
    path: '/mail/inbox',
    name: 'inbox-list',
    component: () => import('../components/EmptyView.vue')
  },
  {
    path: '/mail/inbox/:id',
    name: 'inbox-detail',
    component: () => import('../components/MailDetail.vue')
  },
  {
    path: '/mail/sent',
    name: 'sent-list',
    component: () => import('../components/EmptyView.vue')
  },
  {
    path: '/mail/sent/:id',
    name: 'sent-detail',
    component: () => import('../components/MailDetail.vue')
  },
  {
    path: '/mail/compose',
    name: 'compose',
    component: () => import('../components/ComposeView.vue')
  }
]

const router = createRouter({
  // Always use root base to allow handling both /mail and /im in the same app
  history: createWebHistory('/'),
  routes
})

export default router
