// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import 'core-js/es6/promise'
import 'core-js/es6/string'
import 'core-js/es7/array'
// import cssVars from 'css-vars-ponyfill'
import Vue from 'vue'
import BootstrapVue from 'bootstrap-vue'
import App from './app'
import router from './router'
import axios from 'axios'
import i18n from './i18n/index'

Vue.use(BootstrapVue)

import VueSessionStorage from 'vue-sessionstorage'
Vue.use(VueSessionStorage)

const config = {
  rabbitmq: {
    stompUrl: process.env.VUE_APP_RABBITMQ_STOMP_URL,
    host: process.env.VUE_APP_RABBITMQ_HOST,
    login: process.env.VUE_APP_RABBITMQ_LOGIN,
    password: process.env.VUE_APP_RABBITMQ_PASSWORD
  },
  photon: {
    host: process.env.VUE_APP_PHOTON_HOST
  },
  api: {
    baseUrl: process.env.VUE_APP_ROOT_API
  },
  web: {
    baseUrl: window.location.protocol + '//' + window.location.hostname + ':' + window.location.port
  }
}

// Stomp
import MessageBus from './plugins/messageBus'
Vue.use(MessageBus, config.rabbitmq)

// Error handlers
import ErrorHandler from './plugins/errorHandler'
Vue.use(ErrorHandler)

// Shared config
Object.defineProperty(Vue.prototype, '$config', {
  get: function () {
    return config
  }
})

// Shared http
axios.defaults.baseURL = config.api.baseUrl
Object.defineProperty(Vue.prototype, '$http', {
  get: function () {
    return axios
  }
})

/* eslint-disable no-new */
new Vue({
  el: '#app',
  router,
  i18n,
  components: {
    App
  },
  template: '<App/>'
})

