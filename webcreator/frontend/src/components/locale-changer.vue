<template>
  <b-nav-item-dropdown
    :text="currentLang"
    right>
    <b-dropdown-item
      v-for="(lang, i) in languages"
      :key="`lang-${i}`"
      @click.prevent="changeLocale(lang)">
      <i :class="$t('lang.' + lang + '.icon')"/> {{ $t('lang.' + lang + '.title') }}
    </b-dropdown-item>
  </b-nav-item-dropdown>
</template>

<script>
import BNavItemDropdown from 'bootstrap-vue/es/components/nav/nav-item-dropdown'
import BDropdownItem from 'bootstrap-vue/es/components/dropdown/dropdown-item'

export default {
  components: {
    BNavItemDropdown,
    BDropdownItem
  },
  data () {
    return {
      languages: [
        'ru',
        'en'
      ]
    }
  },
  computed: {
    currentLang: function () {
      return '<i class="' + this.$t('lang.' + this.$i18n.locale + '.icon') + '"></i> ' + this.$t('lang.' + this.$i18n.locale + '.title')
    }
  },
  methods: {
    changeLocale: function (lang) {
      localStorage.lang = lang
      window.location.reload()
    }
  }
}
</script>
