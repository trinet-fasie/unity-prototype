<template>
  <div class="wrapper">
    <div class="animated fadeIn">
      <b-row>
        <b-col cols="12">
          <h1 class="mb-4">{{ $t('worlds.list') }}</h1>
          <b-alert
            v-if="error"
            show
            variant="danger">{{ error }}</b-alert>
          <b-card>
            <b-row>
              <b-col md="3">
                <b-form-group
                  :label="$t('table.filter')"
                  horizontal>
                  <b-input-group>
                    <b-form-input v-model="filter" />
                    <b-input-group-append>
                      <b-button
                        :disabled="!filter"
                        @click="filter = ''">{{ $t('table.filterClear') }}</b-button>
                    </b-input-group-append>
                  </b-input-group>
                </b-form-group>
              </b-col>
            </b-row>

            <b-table
              :items="items"
              :fields="fields"
              :current-page="currentPage"
              :per-page="perPage"
              :filter="filter"
              :sort-by.sync="sortBy"
              :sort-desc.sync="sortDesc"
              :sort-direction="sortDirection"
              :empty-text="$t('table.emptyText')"
              :empty-filtered-text="$t('table.emptyFilteredText')"
              :sort-compare="sortObjectsByKey"
              show-empty
              stacked="md"
              @filtered="onFiltered">
              <template
                slot="actions"
                slot-scope="row">
                <world-form-popover
                  :id="row.item.id"
                  :world="row.item"
                  @submit="updateItem" />
                <delete-with-confirm
                  :id="row.item.id"
                  :confirm="$t('worlds.deleteConfirm')"
                  :disabled="loading"
                  @submit="deleteItem(row.item)"/>
                <b-button
                  v-b-tooltip.hover="$t('actions.exportToFile')"
                  :href="row.item.exportLink"
                  variant="ghost-secondary">
                  <slot><i class="fa fa-download"/></slot>
                </b-button>
                <b-button
                  v-b-tooltip.hover="$t('actions.buildPackage')"
                  :href="row.item.buildLink"
                  variant="ghost-secondary">
                  <slot><i class="fas fa-hammer"/></slot>
                </b-button>
              </template>
              <template
                slot="name"
                slot-scope="row">
                <b-link :to="{ name: 'WorldStructure', params: { worldId: row.item.id }}">{{ row.item.name }}</b-link>
              </template>
            </b-table>

            <b-row v-if="totalRows > perPage">
              <b-col
                md="6"
                class="my-1">
                <b-pagination
                  :total-rows="totalRows"
                  :per-page="perPage"
                  v-model="currentPage"
                  class="my-0" />
              </b-col>
            </b-row>

            <footer>
              <world-form-popover
                id="new-world"
                :world="{name: ''}"
                @submit="addItem" />
              <world-import-popover @submit="importItem" />
            </footer>
          </b-card>
        </b-col>
      </b-row>
    </div>
  </div>
</template>

<script>
import BRow from 'bootstrap-vue/es/components/layout/row'
import BCol from 'bootstrap-vue/es/components/layout/col'
import BCard from 'bootstrap-vue/es/components/card/card'
import BTable from 'bootstrap-vue/es/components/table/table'
import BPagination from 'bootstrap-vue/es/components/pagination/pagination'
import BButton from 'bootstrap-vue/es/components/button/button'
import BFormGroup from 'bootstrap-vue/es/components/form-group/form-group'
import BInputGroup from 'bootstrap-vue/es/components/input-group/input-group'
import BFormInput from 'bootstrap-vue/es/components/form-input/form-input'
import BInputGroupAppend from 'bootstrap-vue/es/components/input-group/input-group-append'
import BAlert from 'bootstrap-vue/es/components/alert/alert'
import BImg from 'bootstrap-vue/es/components/image/img'
import BLink from 'bootstrap-vue/es/components/link/link'
import DeleteWithConfirm from '../../components/delete-with-confirm'
import WorldFormPopover from './components/world-form-popover'
import WorldImportPopover from './components/world-import-popover'

import { convertToISODate, convertToLocalFullDate, sortObjectsByKey } from '@/utils/index'

export default {
  components: {
    BRow,
    BCol,
    BCard,
    BTable,
    BPagination,
    BButton,
    BFormGroup,
    BInputGroup,
    BFormInput,
    BInputGroupAppend,
    BAlert,
    BImg,
    BLink,
    DeleteWithConfirm,
    WorldFormPopover,
    WorldImportPopover
  },
  data () {
    return {
      items: [],
      error: null,
      loading: false,
      fields: [
        { key: 'id', label: this.$t('table.columns.id'), sortable: true, sortDirection: 'desc', class: 'align-middle' },
        { key: 'name', label: this.$t('table.columns.name'), sortable: true, sortDirection: 'desc', class: 'align-middle' },
        { key: 'configurations', label: this.$t('worlds.configurations'), sortable: true, sortDirection: 'desc', class: 'align-middle' },
        { key: 'createdAt', label: this.$t('table.columns.createdAt'), sortable: true, sortDirection: 'desc', formatter: 'formatDateTime', class: 'align-middle' },
        { key: 'updatedAt', label: this.$t('table.columns.updatedAt'), sortable: true, sortDirection: 'desc', formatter: 'formatDateTime', class: 'align-middle' },
        { key: 'actions', label: this.$t('table.columns.actions'), class: 'align-middle' }
      ],
      currentPage: 1,
      perPage: 5,
      totalRows: 0,
      pageOptions: [ 5, 10, 15 ],
      sortBy: 'id',
      sortDesc: false,
      sortDirection: 'asc',
      filter: null,
      sortObjectsByKey
    }
  },
  created () {
    this.fetchData()
  },
  methods: {
    onFiltered (filteredItems) {
      this.totalRows = filteredItems.length
      this.currentPage = 1
    },
    fetchData () {
      this.error = null
      this.loading = true
      this.items = []
      this.totalRows = 0
      this.$http.get('/v1/worlds')
        .then(response => {
          this.loading = false
          if (response.data['Status'] === 'success') {
            response.data['Data'].forEach(world => {
              this.items.push(this.worldToItem(world))
            })
            this.totalRows = this.items.length
          }
        })
        .catch(err => {
          this.loading = false
          this.error = this.$getResponseErrorMessage(err)
        })
    },
    worldToItem (world) {
      return {
        id: world['Id'],
        name: world['Name'],
        configurations: world['Configurations'],
        exportLink: this.$config.api.baseUrl + '/v1/world-structure/' + world['Id'] + '?export',
        buildLink: this.$config.api.baseUrl + '/v1/world-structure/' + world['Id'] + '?build',
        createdAt: convertToISODate(world['CreatedAt']),
        updatedAt: convertToISODate(world['UpdatedAt'])
      }
    },
    formatDateTime (value) {
      return convertToLocalFullDate(value)
    },
    addItem (data) {
      this.error = null
      this.loading = true
      this.$http.put('/v1/add-world', {
        'Name': data.name
      })
        .then(response => {
          this.loading = false
          if (response.data['Status'] === 'success') {
            this.items.push(this.worldToItem(response.data['Data']))
            this.totalRows++
          }
        })
        .catch(err => {
          this.loading = false
          this.error = this.$getResponseErrorMessage(err)
        })
    },
    importItem (file) {
      this.error = null
      this.loading = true

      let formData = new FormData()
      formData.append('file', file)

      this.$http.put('/v1/import-world',
        formData,
        {
          headers: {
            'Content-Type': 'multipart/form-data'
          }
        })
        .then(response => {
          this.loading = false
          if (response.data['Status'] === 'success') {
            this.items.push(this.worldToItem(response.data['Data']))
            this.totalRows++
          }
        })
        .catch(err => {
          this.loading = false
          this.error = this.$getResponseErrorMessage(err)
        })
    },
    updateItem (data) {
      this.error = null
      this.loading = true
      this.$http.post('/v1/update-world/' + data.id, {
        'Name': data.name
      })
        .then(response => {
          this.loading = false
          let updatedItem = this.worldToItem(response.data['Data'])
          this.items.forEach((item, index) => {
            if (item.id === updatedItem.id) {
              this.items.splice(index, 1, updatedItem)
            }
          })
        })
        .catch(err => {
          this.loading = false
          this.error = this.$getResponseErrorMessage(err)
        })
    },
    deleteItem (item) {
      this.error = null
      this.loading = true
      this.items.splice(this.items.indexOf(item), 1)
      this.totalRows--
      this.$http.delete('/v1/delete-world/' + item.id)
        .then(() => {
          this.loading = false
        })
        .catch(err => {
          this.loading = false
          this.error = this.$getResponseErrorMessage(err)
        })
    }
  }
}
</script>
