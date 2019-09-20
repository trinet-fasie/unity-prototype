<template>
  <popover-form
    id="import-world-form"
    :title="$t('worlds.importTitle')"
    :submit-label="$t('actions.upload')"
    :can-submit="fileState"
    @submit="onSubmit">
    <template
      slot="popover-button"
      slot-scope="button">
      <b-button
        :id="button.buttonId"
        variant="primary"
        class="ml-2">
        <i class="fa fa-upload"/> {{ $t('worlds.importButton') }}
      </b-button>
    </template>

    <b-form-group
      :state="fileState"
      :invalid-feedback="$t('validation.required')"
      :valid-feedback="$t('validation.success')"
      class="mb-2">
      <b-form-file
        v-model="file"
        :state="fileState"
        accept=".owws"/>
    </b-form-group>
  </popover-form>
</template>

<script>
import PopoverForm from '../../../components/popover-form'
import BFormGroup from 'bootstrap-vue/es/components/form-group/form-group'
import BFormFile from 'bootstrap-vue/es/components/form-file/form-file'
import BButton from 'bootstrap-vue/es/components/button/button'

export default {
  components: {
    PopoverForm,
    BFormGroup,
    BFormFile,
    BButton
  },
  data () {
    return {
      file: null
    }
  },
  computed: {
    fileState () {
      return Boolean(this.file)
    }
  },
  methods: {
    onSubmit () {
      this.$emit('submit', this.file)
    }
  }
}
</script>
