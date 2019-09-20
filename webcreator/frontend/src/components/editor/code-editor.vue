<template>
  <div>
    <codemirror
      ref="editor"
      v-model="newCode"
      :options="options"
      @input="onChange"/>
    <div class="controls mt-3">
      <b-button
        :disabled="!modified"
        class="mr-1"
        variant="success"
        tabindex="2"
        @click.prevent="apply">
        {{ $t('actions.apply') }}
      </b-button>
      <b-button
        :disabled="!modified"
        class="mr-1"
        variant="danger"
        tabindex="2"
        @click.prevent="reset">
        {{ $t('actions.reset') }}
      </b-button>
      <b-button
        :href="vrLink"
        class="mr-1"
        variant="primary"
        tabindex="1">
        {{ $t('actions.vr') }}
      </b-button>
      <b-button
        variant="secondary"
        tabindex="2"
        @click.prevent="close">
        {{ $t('actions.close') }}
      </b-button>
    </div>
  </div>
</template>

<script>

import BButton from 'bootstrap-vue/es/components/button/button'
import { codemirror } from 'vue-codemirror'

import 'codemirror/lib/codemirror.css'
import 'codemirror/mode/clike/clike.js'
import 'codemirror/theme/idea.css'

export default {
  components: {
    codemirror,
    BButton
  },
  props: {
    uniqueClassPostfix: {
      type: String,
      required: true
    },
    vrLink: {
      type: String,
      required: true
    },
    code: {
      type: String,
      default: ''
    }
  },
  data () {
    return {
      newCode: this.code || `using System;

namespace OpenWorld
{
    public class LogicOf${this.uniqueClassPostfix} : ILogic
    {
        public void Initialize(WrappersCollection collection)
        {
            // Write on init group logic here
            // Don't forget using header for native type
            // Example:
            // using OpenWorld.Types.DisplayWrapper;
            // ...
            // collection.Get<DisplayWrapper>(4).Text = "Hello world!"
        }

        public void Update(WrappersCollection collection)
        {
            // Write on every fps update logic here
        }

        public void Events(WrappersCollection collection)
        {
            // Write events subscribers here
        }
    }
}
`,
      options: {
        tabSize: 4,
        mode: 'text/x-csharp',
        theme: 'idea',
        lineNumbers: true,
        line: true
      }
    }
  },
  computed: {
    modified () {
      return this.code !== this.newCode
    }
  },
  watch: {
    code: function () {
      this.reset()
    }
  },
  methods: {
    refresh: function () {
      window.setTimeout(() => {
        this.$refs.editor.refresh()
      })
    },
    onChange: function (newCode) {
      this.newCode = newCode
    },
    apply: function () {
      this.$emit('change', this.newCode)
    },
    reset: function () {
      this.newCode = this.code
    },
    close: function () {
      if (!this.modified) {
        this.$emit('close')
        return true
      }

      if (window.confirm(this.$t('editor.unsavedChangesConfirm'))) {
        this.$emit('close')
        return true
      }

      return false
    }
  }
}
</script>

<style>
  .CodeMirror {
    height: 600px;
  }
</style>
