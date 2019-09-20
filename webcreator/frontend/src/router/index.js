import Vue from 'vue'
import Router from 'vue-router'

// Containers
const DefaultContainer = () => import('@/containers/DefaultContainer')
const MyPadContainer = () => import('@/containers/MyPadContainer')

// Views - Pages
const Page404 = () => import('@/views/pages/Page404')
const Page500 = () => import('@/views/pages/Page500')
const Login = () => import('@/views/pages/Login')
const Register = () => import('@/views/pages/Register')

// MyPad
const SpawnMenu = () => import('../views/myPad/Spawn')

import libraryObjectsRoutes from '@/modules/library-objects/routes'
import libraryLocationsRoutes from '@/modules/library-locations/routes'
import worldStructureRoutes from '@/modules/world-structure/routes'
import worldListRoutes from '@/modules/world-list/routes'

Vue.use(Router)

export default new Router({
  mode: 'history',
  linkActiveClass: 'open active',
  scrollBehavior: () => ({ y: 0 }),
  routes: [
    {
      path: '/',
      redirect: '/worlds',
      component: DefaultContainer,
      children: [
        {
          path: 'library',
          redirect: '/library/objects',
          name: 'Library',
          meta: {label: 'nav.library'},
          component: {
            render (c) { return c('router-view') }
          },
          children: [
            ...libraryObjectsRoutes,
            ...libraryLocationsRoutes
          ]
        },
        {
          path: 'worlds',
          meta: {label: 'nav.worlds'},
          component: {
            render (c) { return c('router-view') }
          },
          children: [
            ...worldStructureRoutes,
            ...worldListRoutes
          ]
        }
      ]
    },
    {
      path: '/pages',
      redirect: '/pages/404',
      name: 'Pages',
      component: {
        render (c) { return c('router-view') }
      },
      children: [
        {
          path: '404',
          name: 'Page404',
          component: Page404
        },
        {
          path: '500',
          name: 'Page500',
          component: Page500
        },
        {
          path: 'login',
          name: 'Login',
          component: Login
        },
        {
          path: 'register',
          name: 'Register',
          component: Register
        }
      ]
    },
    {
      path: '/mypad',
      name: 'MyPad',
      component: MyPadContainer,
      children: [
        {
          path: 'spawn',
          name: 'SpawnMenu',
          component: SpawnMenu
        }
      ]
    },
    {
      path: '*', component: Page404
    }
  ]
})
