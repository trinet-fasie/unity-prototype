<template>
  <div class="wrapper">
    <div class="animated fadeIn">
      <b-row>
        <b-col cols="12">
          <h1 class="mb-4">{{ $t('objects.list') }}</h1>
          <b-alert
            v-if="error"
            show
            variant="danger">{{ error }}</b-alert>
          <b-card>
            <!-- TODO: change URL-based upload to axios. CYB-477 -->
            <dropzone-wrapper
              id="install"
              :url="this.$config.api.baseUrl + '/v1/install-object'"
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
                  :set-tags-to-item-api="objectsApi.setTagsToObject"
                  :get-tags-api="objectsApi.getTags"
                  :create-tag-api="objectsApi.createTag"/>
              </template>
              <template
                slot="actions"
                slot-scope="row">
                <span
                  v-b-tooltip.hover="row.item.usages > 0 ? $t('objects.deleteRestriction') : $t('actions.delete')"
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

import objectsApi from './api'
import { convertToISODate, convertToLocalFullDate, sortObjectsByKey } from '@/utils/index'
import {OBJECT_MAX_FILE_SIZE} from '@/utils/constants'

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
      maxFileSizeInMB: OBJECT_MAX_FILE_SIZE,
      dropzoneOptions: {
        showRemoveLink: false,
        parallelUploads: 1000,
        maxNumberOfFiles: 1000
      },
      sortObjectsByKey,
      objectsApi
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

      objectsApi.getList()
        .then(({data}) => {
          if (data.status === 'success') {
            this.items = data.data.map(object => {
              return {
                ...object,
                name: this.getLocalizedName(object.config, object.config.type),
                icon: this.$config.api.baseUrl + object.resources.icon,
                createdAt: convertToISODate(object.createdAt),
                updatedAt: convertToISODate(object.updatedAt)
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
    objectToItem (object) {
      // TODO: remove this mapping after CYB-477 (dropzone module removing)
      return {
        id: object['Id'],
        icon: this.$config.api.baseUrl + object['Resources']['Icon'],
        name: this.getLocalizedName(object['Config'], object['Config']['type']),
        usages: object['Usages'],
        createdAt: convertToISODate(object['CreatedAt']),
        updatedAt: convertToISODate(object['UpdatedAt'])
      }
    },
    getLocalizedName (config, fallback) {
      let name = fallback
      if (config.i18n) {
        if (config.i18n[this.$i18n.locale]) {
          name = config.i18n[this.$i18n.locale]
        } else if (config.i18n[this.$i18n.fallbackLocale]) {
          name = config.i18n[this.$i18n.fallbackLocale]
        }
      }

      return name
    },
    formatDateTime (value) {
      return convertToLocalFullDate(value)
    },
    installSuccess (file, response) {
      const statusEl = file.previewElement.querySelector('.dz-status')
      const installedItem = this.objectToItem(response['Data']['Object'])
      if (response['Data']['Created']) {
        statusEl.innerHTML = this.$t('objects.installCreated', {name: installedItem.name})
        this.items.push(installedItem)
        this.itemsCount++
      } else {
        statusEl.innerHTML = this.$t('objects.installUpdated', {name: installedItem.name})
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
      objectsApi.delete(item.id)
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

<style lang="scss" scoped>
  .popover-trigger {
    cursor: pointer;
    &:hover {
      opacity: .7;
    }
  }

  .tag {
    display: inline-block;
    padding: 2px 5px;
    margin: 2px 3px;
    background-color: #dbdfe2;
    border-radius: 3px;
    font-size: 13px;

    &:last-of-type {
      margin-right: 0;
    }
  }
</style>
