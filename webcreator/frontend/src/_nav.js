export default {
  items: [
    {
      name: window.$t('nav.worlds'),
      url: '/worlds',
      icon: 'fas fa-globe'
    },
    {
      title: true,
      name: window.$t('nav.library')
    },
    {
      name: window.$t('nav.objects'),
      url: '/library/objects',
      icon: 'fas fa-puzzle-piece'
    },
    {
      name: window.$t('nav.locations'),
      url: '/library/locations',
      icon: 'fas fa-image'
    }
  ]
}
