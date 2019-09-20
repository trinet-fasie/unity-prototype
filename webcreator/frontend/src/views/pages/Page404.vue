<template>
  <div class="app flex-row align-items-center">
    <div class="container">
      <b-row class="justify-content-center">
        <b-col md="6">
          <div class="clearfix">
            <h1 class="float-left display-3 mr-4">404</h1>
            <h4 class="pt-3">{{ $t('error404.message') }}</h4>
            <p class="text-muted">{{ $t('error404.details') }}</p>
          </div>
        </b-col>
      </b-row>
    </div>
  </div>
</template>

<script>
import BRow from 'bootstrap-vue/es/components/layout/row'
import BCol from 'bootstrap-vue/es/components/layout/col'

export default {
  components: {
    BRow,
    BCol
  }
}
</script>
