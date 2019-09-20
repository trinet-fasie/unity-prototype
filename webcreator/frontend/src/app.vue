<template>
  <router-view/>
</template>

<script>
export default {}
</script>

<style lang="scss">
  /* Import Font Awesome Icons Set */
  $fa-font-path: '~@fortawesome/fontawesome-free/webfonts';
  @import '~@fortawesome/fontawesome-free/scss/fontawesome.scss';
  @import '~@fortawesome/fontawesome-free/scss/solid.scss';
  @import '~@fortawesome/fontawesome-free/scss/regular.scss';
  @import '~@fortawesome/fontawesome-free/scss/brands.scss';
  /* Import Bootstrap Vue Styles */
  @import '~bootstrap-vue/dist/bootstrap-vue.css';
  // Import Main styles for this application
  @import 'assets/scss/style';
</style>
