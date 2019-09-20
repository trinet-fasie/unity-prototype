<template>
  <span>
    <b-button
      v-b-tooltip.hover="$t('actions.delete')"
      :id="buttonId"
      :size="size"
      :variant="variant"
      :disabled="disabled">
      <slot><i class="fas fa-trash"/></slot>
    </b-button>

    <b-popover
      ref="popover"
      :target="buttonId"
      placement="bottom"
      @show="onShow"
      @shown="onShown"
      @hiden="onHidden"
    >
      <template slot="title">
        <b-button
          :aria-label="$t('actions.close')"
          class="close"
          @click="onCancel">
          <span
            class="d-inline-block"
            aria-hidden="true">&times;</span>
        </b-button>
        <span class="title-body">
          {{ $t('deleteConfirm.title') }}
        </span>
      </template>
      <b-form
        ref="form"
        @submit.prevent="onSubmit"
        @keyup.esc="onCancel">
        <p>
          {{ confirm }}
        </p>
        <div>
          <b-button
            ref="submit"
            type="submit"
            size="sm"
            variant="danger"
            tabindex="1">{{ $t('actions.delete') }}</b-button>
          <b-button
            size="sm"
            variant="secondary"
            class="float-right"
            tabindex="2"
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
    confirm: {
      type: String,
      default: ''
    },
    disabled: {
      type: Boolean,
      default: false
    },
    variant: {
      type: String,
      default: 'ghost-secondary'
    },
    size: {
      type: String,
      default: null
    }
  },
  data () {
    return {
      buttonId: 'delete-button-' + this.id,
      popoverRef: 'delete-popover-' + this.id
    }
  },
  methods: {
    onSubmit: function () {
      this.onCancel()
      this.$emit('submit')
    },
    onCancel: function () {
      this.$refs.popover.$emit('close')
    },
    onShow () {
      this.$root.$emit('bv::hide::popover')
    },
    onShown () {
      this.setFocus(this.$refs.submit)
    },
    onHidden () {
      /* Called just after the popover has finished hiding */
      /* Bring focus back to the button */
      this.setFocus(this.$refs[this.buttonId])
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
</style>
