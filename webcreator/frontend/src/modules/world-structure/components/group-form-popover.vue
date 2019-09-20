<template>
  <popover-form
    :id="'gf-' + id"
    :title="fields.id ? $t('world.structure.group.editTitle') : $t('world.structure.group.addTitle')"
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
          @click.prevent=""><i class="fa fa-plus"/> {{ $t('world.structure.group.addButton') }}</a>
      </span>
    </template>
    <b-form-group
      :label="$t('world.structure.group.fields.name')"
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
  </popover-form>
</template>

<script>
import PopoverForm from '../../../components/popover-form'
import BFormGroup from 'bootstrap-vue/es/components/form-group/form-group'
import BFormInput from 'bootstrap-vue/es/components/form-input/form-input'
import BButton from 'bootstrap-vue/es/components/button/button'

export default {
  components: {
    PopoverForm,
    BFormGroup,
    BFormInput,
    BButton
  },
  props: {
    id: {
      type: [String, Number],
      required: true
    },
    group: {
      type: Object,
      required: true
    },
    defaultValues: {
      type: Object,
      default: () => {
        return {
          id: null,
          name: ''
        }
      }
    }
  },
  data () {
    return {
      nameInputId: 'gf-name-' + this.id,
      fields: Object.assign({}, this.defaultValues, this.group)
    }
  },
  computed: {
    nameState () {
      return this.fields.name.length > 0
    },
    canSubmit () {
      return this.nameState
    }
  },
  methods: {
    onSubmit () {
      this.$emit('submit', this.fields)
    },
    onShow () {
      this.fields = Object.assign({}, this.defaultValues, this.group)
    }
  }
}
</script>
