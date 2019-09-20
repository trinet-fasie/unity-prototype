import Vue from 'vue'
import VueI18n from 'vue-i18n'

import ru from './ru'
import en from './en'
import datetime from './datetime'

Vue.use(VueI18n)

if (localStorage.lang) {
  document.querySelector('html').setAttribute('lang', localStorage.lang)
}

const i18n = new VueI18n({
  locale: document.querySelector('html').getAttribute('lang'),
  dateTimeFormats: datetime,
  messages: {ru, en}
})

window.$t = (key, values) => i18n.t(key, values)

export default i18n
