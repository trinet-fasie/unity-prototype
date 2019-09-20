<template>
  <popover-form
    :id="'lf-' + id"
    :title="fields.id ? $t('world.structure.location.editTitle') : $t('world.structure.location.addTitle')"
    :submit-label="fields.id ? $t('actions.apply') : $t('actions.add')"
    :can-submit="canSubmit"
    @submit="onSubmit"
    @show="onShow">
    <template
      slot="popover-button"
      slot-scope="button">
      <b-button
        v-b-tooltip.hover="$t('actions.edit')"
        v-if="fields.id"
        :id="button.buttonId"
        variant="ghost-secondary">
        <i class="fas fa-edit"/>
      </b-button>
      <span
        v-else
        class="node-text">
        <a
          :id="button.buttonId"
          href="javascript:void(0)"
          @click.prevent=""><i class="fa fa-plus"/> {{ $t('world.structure.location.addButton') }}</a>
      </span>
    </template>
    <b-form-group
      :label="$t('world.structure.location.fields.name')"
      :label-for="nameInputId"
      :state="nameState"
      :invalid-feedback="$t('validation.required')"
      :valid-feedback="$t('validation.success')"
      class="mb-2">
      <b-form-input
        :id="nameInputId"
        :state="nameState"
        v-model="fields.name"
        size="sm"
        required />
    </b-form-group>
    <b-form-group
      :label="$t('world.structure.location.fields.locationId')"
      :label-for="locationInputId"
      :state="locationState"
      :invalid-feedback="$t('validation.required')"
      :valid-feedback="$t('validation.success')">
      <b-form-select
        v-model="fields.locationId"
        :state="locationState"
        :options="locations"
        size="sm"
        required/>
    </b-form-group>
  </popover-form>
</template>

<script>

import PopoverForm from '../../../components/popover-form'
import BFormGroup from 'bootstrap-vue/es/components/form-group/form-group'
import BFormInput from 'bootstrap-vue/es/components/form-input/form-input'
import BButton from 'bootstrap-vue/es/components/button/button'
import BFormSelect from 'bootstrap-vue/es/components/form-select/form-select'

export default {
  components: {
    PopoverForm,
    BFormGroup,
    BFormInput,
    BButton,
    BFormSelect
  },
  props: {
    id: {
      type: [String, Number],
      required: true
    },
    location: {
      type: Object,
      required: true
    },
    defaultValues: {
      type: Object,
      default: () => {
        return {
          id: null,
          name: '',
          locationId: null
        }
      }
    }
  },
  data () {
    return {
      nameInputId: 'lf-name-' + this.id,
      locationInputId: 'lf-location-' + this.id,
      fields: Object.assign({}, this.defaultValues, this.location),
      locations: []
    }
  },
  computed: {
    nameState () {
      return this.fields.name.length > 0
    },
    locationState () {
      return this.fields.locationId > 0
    },
    canSubmit () {
      return this.nameState && this.locationState
    }
  },
  methods: {
    onSubmit () {
      this.$emit('submit', this.fields)
    },
    onShow () {
      this.fields = Object.assign({}, this.defaultValues, this.location)
      this.locations = []
      this.$http.get('/v1/locations')
        .then(response => {
          if (response.data['Status'] === 'success') {
            response.data['Data'].forEach(location => {
              this.locations.push(this.locationToOption(location))
            })
          }
        })
    },
    locationToOption: function (location) {
      return {
        value: location['Id'],
        text: location['Name']
      }
    }
  }
}
</script>
