<template>
  <div class="container">
    <b-alert
      v-if="error"
      show
      variant="danger">{{ error }}</b-alert>
    <b-row>
      <b-col>
        <a
          v-for="item in items"
          :key="item.id"
          :onclick="'spawnObject(' + item.id +')'"
          href="javascript:"
          class="menu-item">
          <b-img
            :src="item.icon"
            :alt="item.name"
            thumbnail
            fluid
            class="m-2"
            width="136"
            height="136" />
        </a>
      </b-col>
    </b-row>
  </div>
</template>

<script>
import BRow from 'bootstrap-vue/es/components/layout/row'
import BCol from 'bootstrap-vue/es/components/layout/col'
import BImg from 'bootstrap-vue/es/components/image/img'
import BAlert from 'bootstrap-vue/es/components/alert/alert'

export default {
  components: {
    BRow,
    BCol,
    BImg,
    BAlert
  },
  data () {
    return {
      items: [],
      error: null,
      loading: false
    }
  },
  created () {
    this.fetchData()
  },
  methods: {
    fetchData: function () {
      this.error = null
      this.loading = true
      this.items = []
      this.$http.get('/v1/objects')
        .then(response => {
          this.loading = false
          if (response.data['Status'] === 'success') {
            response.data['Data'].forEach(object => {
              this.items.push(this.objectToItem(object))
            })
          }
        })
        .catch(err => {
          this.loading = false
          this.error = this.$getResponseErrorMessage(err)
        })
    },
    objectToItem: function (object) {
      return {
        id: object['Id'],
        icon: this.$config.api.baseUrl + object['Resources']['Icon'],
        name: object['Name']
      }
    }
  }
}
</script>

<style>
  a.menu-item:hover .img-thumbnail {
    background-color: #ffc107;
  }
</style>
