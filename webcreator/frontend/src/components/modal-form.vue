<template>
  <b-modal
    ref="modal"
    :id="modalId"
    :title="title"
    :size="size"
    centered
    no-close-on-backdrop
    lazy
    @show="onShow"
    @shown="onShown"
    @hide="onHide"
    @hidden="onHidden">
    <b-form
      novalidate
      @submit.stop.prevent="onSubmit">
      <div ref="form">
        <slot/>
      </div>
    </b-form>

    <div
      slot="modal-footer"
      class="w-100">
      <b-button
        v-if="showSubmit"
        ref="submit"
        :disabled="!canSubmit"
        type="submit"
        variant="primary"
        class="float-left"
        tabindex="20"
        @click="onSubmit">{{ submitLabel }}</b-button>
      <b-button
        ref="cancel"
        variant="secondary"
        class="float-right"
        tabindex="21"
        @click="onCancel">{{ $t('actions.cancel') }}</b-button>
    </div>
  </b-modal>
</template>

<script>
import BButton from 'bootstrap-vue/es/components/button/button'
import BModal from 'bootstrap-vue/es/components/modal/modal'
import BForm from 'bootstrap-vue/es/components/form/form'

export default {
  components: {
    BButton,
    BModal,
    BForm
  },
  props: {
    id: {
      type: [String, Number],
      required: true
    },
    title: {
      type: String,
      required: true
    },
    submitLabel: {
      type: String,
      required: true
    },
    size: {
      type: String,
      default: 'lg'
    },
    canSubmit: {
      type: Boolean,
      default: true
    },
    canShow: {
      type: Boolean,
      default: true
    },
    showSubmit: {
      type: Boolean,
      default: true
    }
  },
  data () {
    return {
      buttonId: 'mf-btn-' + this.id,
      modalId: 'mf-' + this.id
    }
  },
  mounted () {
    const button = document.querySelector('#' + this.buttonId)
    if (!button) {
      throw new Error('Button #' + this.buttonId + ' is not found.')
    }
    button.addEventListener('click', () => {
      if (this.canShow) {
        this.$refs.modal.show()
      }
    })
  },
  methods: {
    onSubmit: function () {
      if (this.canSubmit) {
        this.onCancel()
        this.$emit('submit')
      }
    },
    onCancel: function () {
      this.$refs.modal.hide()
    },
    onShow () {
      this.$root.$emit('bv::hide::popover')
      this.$emit('show')
    },
    onShown () {
      const formElements = this.$refs.form.querySelectorAll('button,input,textarea,checkbox,select')
      formElements.forEach((el, index) => {
        el.setAttribute('tabindex', index + 1)
      })
      if (formElements.length > 0) {
        formElements[0].focus()
      } else if (this.$refs.submit) {
        this.setFocus(this.$refs.submit)
      } else if (this.$refs.cancel) {
        this.setFocus(this.$refs.cancel)
      }
      this.$emit('shown')
    },
    onHide () {
      this.$emit('hide')
    },
    onHidden () {
      const button = document.querySelector('#' + this.buttonId)
      if (button) {
        this.setFocus(button)
      }
      this.$emit('hidden')
    },
    setFocus (ref) {
      this.$nextTick(() => {
        this.$nextTick(() => { (ref.$el || ref).focus() })
      })
    }
  }
}
</script>
