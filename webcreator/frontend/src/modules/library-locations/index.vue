<template>
  <div class="wrapper">
    <div class="animated fadeIn">
      <b-row>
        <b-col cols="12">
          <h1 class="mb-4">{{ $t('locations.list') }}</h1>
          <b-alert
            v-if="error"
            show
            variant="danger">{{ error }}</b-alert>
          <b-card>
            <!-- TODO: change URL-based upload to axios. CYB-477 -->
            <dropzone-wrapper
              id="install"
              :url="this.$config.api.baseUrl + '/v1/install-location'"
              :language="$t('installDropzone')"
              :use-font-awesome="true"
              :use-custom-dropzone-options="true"
              :max-file-size-in-mb="maxFileSizeInMB"
              :dropzone-options="dropzoneOptions"
              accepted-file-types=".zip"
              @vdropzone-success="installSuccess" />
          </b-card>
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
                slot="name"
                slot-scope="row">
                <b-img
                  :src="row.item.icon"
                  :alt="row.item.name"
                  fluid
                  width="40" />
                {{ row.item.name }}
              </template>
              <template
                slot="tags"
                slot-scope="row">
                <tags-wrapper
                  :item-id="row.item.id"
                  v-model="row.item.tags"
                  :set-tags-to-item-api="locationsApi.setTagsToLocation"
                  :get-tags-api="locationsApi.getTags"
                  :create-tag-api="locationsApi.createTag"/>
              </template>
              <template
                slot="actions"
                slot-scope="row">
                <span
                  v-b-tooltip.hover="row.item.usages > 0 ? $t('locations.deleteRestriction') : $t('actions.delete')"
                  class="d-inline-block">
                  <b-button
                    :disabled="row.item.usages > 0 || loading"
                    class="mr-1"
                    variant="ghost-secondary"
                    @click.stop="deleteItem(row.item)">
                    <i class="fas fa-trash"/>
                  </b-button>
                </span>
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
import DropzoneWrapper from '@/components/dropzone-wrapper'
import TagsWrapper from '@/components/tags-wrapper'

import locationsApi from './api'
import { convertToISODate, convertToLocalFullDate, sortObjectsByKey } from '@/utils/index'
import {LOCATION_MAX_FILE_SIZE} from '@/utils/constants'

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
    DropzoneWrapper,
    TagsWrapper
  },
  data () {
    return {
      items: [],
      error: null,
      loading: false,
      fields: [
        {key: 'id', label: this.$t('table.columns.id'), sortable: true, sortDirection: 'desc', class: 'align-middle'},
        {key: 'name', label: this.$t('table.columns.name'), sortable: true, sortDirection: 'desc', class: 'align-middle'},
        {key: 'usages', label: this.$t('table.columns.usages'), sortable: true, sortDirection: 'desc', class: 'align-middle'},
        {
          key: 'createdAt',
          label: this.$t('table.columns.createdAt'),
          sortable: true,
          sortDirection: 'desc',
          formatter: 'formatDateTime',
          class: 'align-middle'
        },
        {
          key: 'updatedAt',
          label: this.$t('table.columns.updatedAt'),
          sortable: true,
          sortDirection: 'desc',
          formatter: 'formatDateTime',
          class: 'align-middle'
        },
        {key: 'tags', label: this.$t('table.columns.tags'), sortable: false, sortDirection: 'desc', class: 'align-middle'},
        {key: 'actions', label: this.$t('table.columns.actions'), class: 'align-middle'}
      ],
      currentPage: 1,
      perPage: 15,
      totalRows: 0,
      pageOptions: [5, 10, 15],
      sortBy: 'id',
      sortDesc: false,
      sortDirection: 'asc',
      filter: null,
      maxFileSizeInMB: LOCATION_MAX_FILE_SIZE,
      dropzoneOptions: {
        showRemoveLink: false,
        parallelUploads: 1000,
        maxNumberOfFiles: 1000
      },
      sortObjectsByKey,
      locationsApi
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
      locationsApi.getList()
        .then(({data}) => {
          if (data.status === 'success') {
            this.items = data.data.map(location => {
              return {
                ...location,
                icon: this.$config.api.baseUrl + location.resources.icon,
                createdAt: convertToISODate(location.createdAt),
                updatedAt: convertToISODate(location.updatedAt)
              }
            })
            this.totalRows = this.items.length
          }
        })
        .catch(err => {
          this.error = this.$getResponseErrorMessage(err)
        })
        .finally(() => {
          this.loading = false
        })
    },
    locationToItem (location) {
      return {
        id: location['Id'],
        icon: this.$config.api.baseUrl + location['Resources']['Icon'],
        name: location['Name'],
        usages: location['Usages'],
        createdAt: convertToISODate(location['CreatedAt']),
        updatedAt: convertToISODate(location['UpdatedAt'])
      }
    },
    formatDateTime (value) {
      return convertToLocalFullDate(value)
    },
    installSuccess (file, response) {
      const statusEl = file.previewElement.querySelector('.dz-status')
      const installedItem = this.locationToItem(response['Data']['Location'])
      if (response['Data']['Created']) {
        statusEl.innerHTML = this.$t('locations.installCreated', {name: installedItem.name})
        this.items.push(installedItem)
        this.itemsCount++
      } else {
        statusEl.innerHTML = this.$t('locations.installUpdated', {name: installedItem.name})
        this.items.forEach((item, index) => {
          if (item.id === installedItem.id) {
            this.items.splice(index, 1, installedItem)
          }
        })
      }
    },
    deleteItem (item) {
      this.error = null
      this.loading = true
      locationsApi.delete(item.id)
        .then(() => {
          this.items.splice(this.items.indexOf(item), 1)
          this.totalRows--
        })
        .catch(err => {
          this.error = this.$getResponseErrorMessage(err)
        })
        .finally(() => {
          this.loading = false
        })
    }
  }
}
</script>
