<template>
  <span>
    <slot
      :button-id="buttonId"
      name="popover-button">
      <b-button
        :id="buttonId"
        variant="primary">
        {{ $t('actions.apply') }}
      </b-button>
    </slot>

    <b-popover
      ref="popover"
      :target="buttonId"
      :title="title"
      :placement="placement"
      @show="onShow"
      @shown="onShown"
      @hide="onHide"
      @hidden="onHidden">
      <template slot="title">
        <b-button
          size="sm"
          variant="ghost-secondary"
          class="float-right"
          @click="onCancel">
          <i class="fa fa-remove"/>
        </b-button>
        <span class="title-body">
          {{ title }}
        </span>
      </template>
      <b-form
        ref="form"
        novalidate
        @submit.stop.prevent="onSubmit"
        @keyup.esc="onCancel">
        <div ref="form">
          <slot/>
        </div>
        <div class="popover-controls">
          <b-button
            ref="submit"
            :disabled="!canSubmit"
            type="submit"
            size="sm"
            variant="primary"
            tabindex="20">{{ submitLabel }}</b-button>
          <b-button
            ref="cancel"
            size="sm"
            variant="secondary"
            class="float-right"
            tabindex="21"
            @click="onCancel">{{ $t('actions.cancel') }}</b-button>
        </div>
      </b-form>
    </b-popover>
  </span>
</template>

<script>
import BButton from 'bootstrap-vue/es/components/button/button'
import BPopover from 'bootstrap-vue/es/components/popover/popover'
import BForm from 'bootstrap-vue/es/components/form/form'

export default {
  components: {
    BButton,
    BPopover,
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
    placement: {
      type: String,
      default: 'auto'
    },
    canSubmit: {
      type: Boolean,
      default: true
    }
  },
  data () {
    return {
      buttonId: 'pf-button-' + this.id
    }
  },
  methods: {
    onSubmit: function () {
      if (this.canSubmit) {
        this.onCancel()
        this.$emit('submit')
      }
    },
    onCancel: function () {
      this.$refs.popover.$emit('close')
      this.$emit('close')
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
      /* Called just after the popover has finished hiding */
      /* Bring focus back to the button */
      this.setFocus(this.$el.querySelector('#' + this.buttonId))
      this.$emit('hidden')
    },
    setFocus (ref) {
      /* Some references may be a component, functional component, or plain element */
      /* This handles that check before focusing, assuming a focus() method exists */
      /* We do this in a double nextTick to ensure components have updated & popover positioned first */
      this.$nextTick(() => {
        this.$nextTick(() => { (ref.$el || ref).focus() })
      })
    }
  }
}
</script>

<style>
  .popover .title-body {
    padding: 0.25rem 0.5rem 0.25rem  0;
    display: inline-block;
    vertical-align: middle;
  }
  .popover .popover-controls {
    padding: 0.25rem 0;
  }
</style>
