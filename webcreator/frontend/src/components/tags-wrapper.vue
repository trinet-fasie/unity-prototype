<template>
  <div class="tag-wrapper">
    <vue-tags-input
      :class="{loading: loading}"
      v-model="tag"
      :tags="value"
      :allow-edit-tags="true"
      :maxlength="TAG_MAX_LENGTH"
      :placeholder="$t('tags.addTag')"
      :autocomplete-min-length="2"
      :autocomplete-items="autocompleteTags"
      class="tags-input"
      @before-adding-tag="beforeAddingTag"
      @before-saving-tag="beforeSavingTag"
      @tags-changed="updateTags"/>

    <div
      v-if="loading"
      class="loading-wrap">
      <i class="fas fa-spin fa-spinner"/>
    </div>
  </div>
</template>

<script>
import debounce from 'lodash.debounce'
import VueTagsInput from '@johmun/vue-tags-input'

import {TAG_MAX_LENGTH, DEBOUNCE_DURATION} from '@/utils/constants'

export default {
  name: 'TagsWrapper',
  components: {
    VueTagsInput
  },
  props: {
    value: {
      type: Array,
      default: () => []
    },
    itemId: {
      type: Number,
      required: true
    },
    setTagsToItemApi: {
      type: Function,
      required: true
    },
    getTagsApi: {
      type: Function,
      required: true
    },
    createTagApi: {
      type: Function,
      required: true
    }
  },
  data () {
    return {
      TAG_MAX_LENGTH,
      tag: '',
      autocompleteTags: [],
      debounce: null,
      loading: false
    }
  },
  watch: {
    tag: debounce(function () {
      this.fetchTags()
    }, DEBOUNCE_DURATION)
  },
  methods: {
    updateTags (updatedTags) {
      this.loading = true
      this.setTagsToItemApi(this.itemId, updatedTags)
        .then(() => {
          this.$emit('input', updatedTags)
        })
        .finally(() => {
          this.loading = false
        })
    },

    fetchTags () {
      if (this.tag.length === 0) return
      this.loading = true
      this.getTagsApi(this.tag)
        .then(({data}) => {
          this.autocompleteTags = data.data
        })
        .finally(() => {
          this.loading = false
        })
    },

    beforeAddingTag (obj) {
      if (obj.tag.tiClasses.includes('duplicate')) return
      if (!obj.tag.id) {
        this.createNewTag(obj.tag)
          .then(({data}) => {
            obj.tag.id = data.data.id
            obj.addTag()
          })
      } else {
        obj.addTag()
      }
    },

    beforeSavingTag (obj) {
      this.createNewTag(obj.tag)
        .then(({data}) => {
          obj.tag.id = data.data.id
          obj.saveTag()
        })
    },

    createNewTag (tag) {
      return this.createTagApi(tag)
    }
  }
}
</script>

<style lang="scss">
  .tag-wrapper {
    position:relative;
  }

  .tags-input {
    .tags .tag {
      background-color: #dbdfe2;
      color: #000;
    }
  }

  .loading-wrap {
    position: absolute;
    bottom: 2px;
    right: 2px;
    padding: 5px;
    background-color: #fff;
  }
</style>
