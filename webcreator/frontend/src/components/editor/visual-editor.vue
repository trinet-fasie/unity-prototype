<template>
  <div>
    <div
      class="blockly-area"
      style="width:100%;height:650px;"/>
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
        class="mr-1"
        variant="success"
        tabindex="2"
        @click.prevent="applyAndOpenEditor">
        {{ $t('actions.applyAndOpenEditor') }}
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
        tabindex="no1">
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
const Blockly = require('node-blockly/lib/blockly_compressed_browser')
const defaultMessages = {
  en: require('node-blockly/lib/i18n/en')(),
  ru: require('node-blockly/lib/i18n/ru')()
}

import BButton from 'bootstrap-vue/es/components/button/button'
import CommonBlocks from './сommon'
import CSharp from './csharp'

let messages = CommonBlocks.i18n

const locale = document.querySelector('html').getAttribute('lang')
Blockly.Msg = Object.assign(Blockly.Msg, defaultMessages[locale], messages[locale])

Blockly.utils.getMessageArray_ = function () {
  return Blockly.Msg
}

Blockly.Blocks = Object.assign(Blockly.Blocks, require('node-blockly/lib/blocks_compressed_browser')(Blockly))
Blockly.CSharp = CSharp.init(Blockly)
CommonBlocks.init(Blockly)

export default {
  components: {
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
    editorConfig: {
      type: Object,
      required: true
    },
    worldLocations: {
      type: Array,
      required: true
    },
    worldConfigurations: {
      type: Array,
      required: true
    },
    code: {
      type: String,
      default: ''
    },
    editorData: {
      type: Object,
      required: true
    }
  },
  data () {
    return {
      config: this.editorConfig,
      locations: this.worldLocations,
      configurations: this.worldConfigurations,
      modified: false,
      workspace: null,
      ignoreChangeEvents: false,
      blocksScope: new Map()
    }
  },
  watch: {
    config () {
      this.refresh()
    }
  },
  mounted () {
    window.setTimeout(() => {
      const blocklyArea = this.$el.querySelector('div.blockly-area')
      this.workspace = Blockly.inject(blocklyArea, {
        toolbox: this.getToolbox(),
        collapse: true,
        comments: true,
        disable: true,
        maxBlocks: Infinity,
        trashcan: true,
        horizontalLayout: false,
        toolboxPosition: 'start',
        css: true,
        // media: '/static/media/',
        rtl: false,
        scrollbars: true,
        sounds: true,
        oneBasedIndex: true,
        grid: {
          spacing: 20,
          length: 3,
          colour: '#CCC',
          snap: true
        },
        zoom: {
          controls: true,
          wheel: false,
          startScale: 1,
          maxScale: 3,
          minScale: 0.3,
          scaleSpeed: 1.2
        }
      })

      this.updateWorkspace(this.editorData)

      this.workspace.addChangeListener((event) => {
        if (event.type === Blockly.Events.UI) {
          return
        }

        if (this.ignoreChangeEvents) {
          return
        }
        this.modified = true
      })
    })
  },
  methods: {
    refresh () {
      if (!this.workspace) {
        return
      }
      this.workspace.updateToolbox(this.getToolbox())
      // trick for refresh workspace
      const editorData = this.editorData
      editorData['Blockly'] = Blockly.Xml.domToText(Blockly.Xml.workspaceToDom(this.workspace))
      this.updateWorkspace(editorData)
      this.workspace.getToolbox().clearSelection()
    },
    apply () {
      this.modified = false
      const editorData = this.editorData
      editorData['Blockly'] = Blockly.Xml.domToText(Blockly.Xml.workspaceToDom(this.workspace))

      const code = this.generateCode()
      const lockedInstances = []
      for (let id in Blockly.CSharp.lockedInstancesForDelete) {
        if (Blockly.CSharp.lockedInstancesForDelete.hasOwnProperty(id)) {
          lockedInstances.push(Number(id))
        }
      }

      this.$emit('change', code, editorData, lockedInstances)
    },
    applyAndOpenEditor () {
      this.apply()
      this.$emit('open-editor')
    },
    reset () {
      this.modified = false
      this.updateWorkspace(this.editorData)
    },
    updateWorkspace (editorData) {
      this.ignoreChangeEvents = true
      this.workspace.clear()

      if (editorData['Blockly']) {
        Blockly.Xml.domToWorkspace(Blockly.Xml.textToDom(editorData['Blockly']), this.workspace)
      }
      setTimeout(() => { this.ignoreChangeEvents = false })
    },
    generateCode () {
      let update = []
      let initialize = []
      let events = []

      const blocks = this.workspace.getTopBlocks(true)
      Blockly.CSharp.init(this.workspace)

      for (let x = 0; x < blocks.length; x++) {
        let block = blocks[x]
        let line = Blockly.CSharp.blockToCode(block)
        if (line instanceof Array) {
          // Value blocks return tuples of code and operator order.
          // Top-level blocks don't care about operator order.
          line = line[0]
        }
        if (line) {
          if (block.outputConnection && this.scrubNakedValue) {
            // This block is a naked value.  Ask the language's code generator if
            // it wants to append a semicolon, or something.
            line = Blockly.CSharp.scrubNakedValue(line)
          }
          // Events logic
          if (block.type === 'event_on_init') {
            initialize.push(line)
          } else if (this.blocksScope.get(block.type) === 'Events') {
            events.push(line)
          } else {
            update.push(line)
          }
        }
      }

      let using = ['using System;', 'using System.Collections.Generic;']
      for (let namespace in Blockly.CSharp.using) {
        if (Blockly.CSharp.using.hasOwnProperty(namespace)) {
          using.push('using ' + namespace + ';')
        }
      }
      let definitions = []
      for (let name in Blockly.CSharp.definitions_) {
        if (Blockly.CSharp.definitions_.hasOwnProperty(name)) {
          definitions.push(Blockly.CSharp.definitions_[name])
        }
      }
      const usingCode = using.join('\n')
      const definitionsCode = this.indentCode(definitions.join('\n'), '        ')
      const initializeCode = this.indentCode(initialize.join('\n'), '            ')
      const eventsCode = this.indentCode(events.join('\n'), '            ')
      const updateCode = this.indentCode(update.join('\n'), '            ')

      return `${usingCode}

namespace OpenWorld
{
    public class LogicOf${this.uniqueClassPostfix} : ILogic
    {
        ${definitionsCode}

        public void Initialize(WrappersCollection collection)
        {
            ${initializeCode}
        }

        public void Update(WrappersCollection collection)
        {
            ${updateCode}
        }

        public void Events(WrappersCollection collection)
        {
            ${eventsCode}
        }
    }
}
`
    },
    indentCode (code, leftSpaces) {
      code = code.split('\n').join('\n' + leftSpaces)
      code = code.replace(/^\s+\n/, '')
      code = code.replace(/\n\s+$/, '\n')
      code = code.replace(/[ \t]+\n/g, '\n')

      return code
    },
    getToolbox () {
      let logic = CommonBlocks.toolbox.logic
      let actions = CommonBlocks.toolbox.actions
      let lists = CommonBlocks.toolbox.lists
      let loops = CommonBlocks.toolbox.loops
      let math = CommonBlocks.toolbox.math
      let text = CommonBlocks.toolbox.text
      let events = CommonBlocks.toolbox.events

      actions += CommonBlocks.getSwitchLocationBlock(Blockly, this.locations)
      actions += CommonBlocks.getSwitchConfigurationBlock(Blockly, this.configurations)

      const objects = {}

      const anyObjectInstanceBlock = this.formAnyObjectInstanceBlock(this.config.allInstances)
      const objectTypes = []
      const usedTypes = {}
      for (const objectId in this.config.objects) {
        if (!this.config.objects.hasOwnProperty(objectId)) {
          continue
        }
        const object = this.config.objects[objectId]
        if (!object.config.type) {
          console.error('Undefined type of object:', object)
          continue
        }

        const objectCategoryName = this.getObjectInstanceUniqueId(object)
        if (!objects[objectCategoryName]) {
          objects[objectCategoryName] = {
            name: this.getLocalizedName(object.config, object.config.type.toLowerCase()),
            blocks: ''
          }
        }

        if (!usedTypes[object.config.type]) {
          objectTypes.push([object.config.type, object.config.type])
          usedTypes[object.config.type] = true
        }

        objects[objectCategoryName].blocks += this.formObjectInstanceBlock(object)
        objects[objectCategoryName].blocks += this.formAllObjectInstancesListBlock(object)

        if (!object.config.blocks) {
          continue
        }

        let objectLogic = ''
        let objectActions = ''
        let objectState = ''
        let objectEvents = ''
        object.config.blocks.forEach(blockConfig => {
          if (!blockConfig.type) {
            console.error('Undefined type of block:', blockConfig)
            return
          }
          if (!blockConfig.name) {
            console.error('Undefined name of block:', blockConfig)
            return
          }
          if (!blockConfig.items) {
            console.error('Undefined items of block:', blockConfig)
            return
          }
          switch (blockConfig.type) {
            case 'checker':
              objectLogic += this.formCheckerBlock(blockConfig, object)
              break
            case 'action':
              objectActions += this.formActionBlock(blockConfig, object)
              break
            case 'values':
              objectState += this.formValuesBlock(blockConfig, object)
              break
            case 'getter':
              objectState += this.formGetterBlock(blockConfig, object)
              break
            case 'setter':
              objectState += this.formSetterBlock(blockConfig, object)
              break
            case 'event':
              objectEvents += this.formEventBlock(blockConfig, object)
              break
            default:
              console.error(`Unknown block type "${blockConfig.type}" for block:`, blockConfig)
          }
        })

        if (objectLogic !== '') {
          objects[objectCategoryName].blocks += `
          <label text="${this.$t('editor.category.logic')}" web-class="logicLabel"></label>
          ${objectLogic}
`
        }
        if (objectActions !== '') {
          objects[objectCategoryName].blocks += `
          <label text="${this.$t('editor.category.actions')}" web-class="actionsLabel"></label>
          ${objectActions}
`
        }
        if (objectState !== '') {
          objects[objectCategoryName].blocks += `
          <label text="${this.$t('editor.category.state')}" web-class="stateLabel"></label>
          ${objectState}
`
        }
        if (objectEvents !== '') {
          objects[objectCategoryName].blocks += `
          <label text="${this.$t('editor.category.events')}" web-class="eventsLabel"></label>
          ${objectEvents}
`
        }
      }

      if (Object.keys(objects).length) {
        objects.any = {
          name: window.$t('editor.anyObject'),
          blocks: `
        ${anyObjectInstanceBlock}
        <label text="${this.$t('editor.category.logic')}" web-class="logicLabel"></label>
        ${this.formAnyObjectCheckTypeBlock(objectTypes)}
        ${this.formAnyObjectCheckBlock()}
        <label text="${this.$t('editor.category.actions')}" web-class="actionsLabel"></label>
        ${this.formAnyObjectActionBlock()}
        <label text="${this.$t('editor.category.state')}" web-class="stateLabel"></label>
        ${this.formAnyObjectHierarchicPropertiesBlock()}
        ${this.formAnyObjectHierarchicListsBlock()}
        ${this.formAnyObjectGetterBlock()}
        ${this.formAnyObjectSetterBlock()}
`
        }
      }

      let objectCategories = ''
      for (const categoryId in objects) {
        if (!objects.hasOwnProperty(categoryId)) {
          continue
        }

        const category = objects[categoryId]
        objectCategories += `
           <category name="${category.name}">
              ${category.blocks}
           </category>
`
      }

      let categoryObjects = ``
      if (objectCategories) {
        categoryObjects = `
        <category name="${this.$t('editor.category.objects')}" colour="#A65C81">
            ${objectCategories}
        </category>
        `
      }

      return `
        <xml id="toolbox" ref="toolbox" style="display: none">
          <category name="${this.$t('editor.category.logic')}" colour="#5C81A6">
            ${logic}
          </category>
          <category name="${this.$t('editor.category.actions')}" colour="#6da55b">
            ${actions}
          </category>
          <category name="${this.$t('editor.category.variables')}" colour="#A65C81" custom="VARIABLE"></category>
          <category name="${this.$t('editor.category.events')}" colour="#985aa5">
            ${events}
          </category>
          ${categoryObjects}
          <category name="${this.$t('editor.category.lists')}" colour="#745CA6">
            ${lists}
          </category>
          <category name="${this.$t('editor.category.loops')}" colour="#5CA65C">
            ${loops}
          </category>
          <category name="${this.$t('editor.category.math')}" colour="#5C68A6">
            ${math}
          </category>
          <category name="${this.$t('editor.category.text')}" colour="#5CA68D">
            ${text}
          </category>
          <category name="${this.$t('editor.category.functions')}" colour="#9A5CA6" custom="PROCEDURE"></category>
        </xml>
`
    },
    formAnyObjectCheckTypeBlock (objectTypes) {
      if (!objectTypes.length) {
        return ''
      }
      Blockly.Blocks['check_object_type'] = {
        init () {
          this.jsonInit({
            'type': 'check_object_type',
            'message0': '%{BKY_CHECK_OBJECT_TYPE_TITLE}',
            'args0': [
              {
                'type': 'input_value',
                'name': 'instance',
                'check': 'Object'
              },
              {
                'type': 'field_dropdown',
                'name': 'type',
                'options': objectTypes
              }
            ],
            'output': 'Boolean',
            'colour': '%{BKY_LOGIC_HUE}',
            'tooltip': '%{BKY_CHECK_OBJECT_TYPE_TOOLTIP}',
            'helpUrl': ''
          })
        }
      }
      Blockly.CSharp['check_object_type'] = block => {
        const instance = Blockly.CSharp.valueToCode(block, 'instance', Blockly.CSharp.ORDER_ATOMIC) || 'null'
        const type = block.getFieldValue('type')
        const code = `${instance} is ${type}Wrapper`

        Blockly.CSharp.use(`OpenWorld.Types.${type}`)

        return [code, Blockly.CSharp.ORDER_TYPEOF]
      }

      return `
        <block type="check_object_type">
            <value name="instance">
                <shadow type="object_any"></shadow>
            </value>
        </block>
    `
    },
    formObjectInstanceBlock (object) {
      if (!object.config) {
        return ''
      }
      const blockType = this.getObjectInstanceUniqueId(object)
      const name = this.getLocalizedName(object.config, object.config.type.toLowerCase())

      Blockly.Blocks[blockType] = {
        init () {
          this.jsonInit({
            'type': blockType,
            'message0': `${name} %1`,
            'args0': [
              {
                'type': 'field_dropdown',
                'name': 'instance',
                'options': object.instances
              }
            ],
            'output': 'Object',
            'colour': '%{BKY_VARIABLES_HUE}',
            'tooltip': '',
            'helpUrl': ''
          })
        }
      }

      Blockly.CSharp[blockType] = block => {
        const id = block.getFieldValue('instance')

        Blockly.CSharp.use(`OpenWorld.Types.${object.config.type}`)
        Blockly.CSharp.lockInstance(id)

        return [`collection.Get<${object.config.type}Wrapper>(${id})`, Blockly.CSharp.ORDER_ATOMIC]
      }

      return `
        <block type="${blockType}"></block>
    `
    },
    formAllObjectInstancesListBlock (object) {
      if (!object.config) {
        return ''
      }
      const blockType = this.getObjectInstanceUniqueId(object) + '_instances_list'
      const name = this.getLocalizedName(object.config, object.config.type.toLowerCase())

      Blockly.Blocks[blockType] = {
        init () {
          this.jsonInit({
            'type': blockType,
            'message0': window.$t('editor.allObjectInstancesTitle', {name}),
            'output': 'Array',
            'colour': '%{BKY_LISTS_HUE}',
            'tooltip': '',
            'helpUrl': ''
          })
        }
      }

      Blockly.CSharp[blockType] = () => {
        Blockly.CSharp.use(`OpenWorld.Types.${object.config.type}`)

        return [`collection.GetWrappersOfType<${object.config.type}Wrapper>()`, Blockly.CSharp.ORDER_ATOMIC]
      }

      return `
        <block type="${blockType}"></block>
    `
    },
    formAnyObjectInstanceBlock (allInstances) {
      Blockly.Blocks['object_any'] = {
        init () {
          this.jsonInit({
            'type': 'object_any',
            'message0': '%{BKY_OBJECT_ANY_TITLE}',
            'args0': [
              {
                'type': 'field_dropdown',
                'name': 'instance',
                'options': allInstances
              }
            ],
            'output': 'Object',
            'colour': '%{BKY_VARIABLES_HUE}',
            'tooltip': '%{BKY_OBJECT_ANY_TOOLTIP}',
            'helpUrl': ''
          })
        }
      }

      Blockly.CSharp['object_any'] = block => {
        const id = block.getFieldValue('instance')
        Blockly.CSharp.lockInstance(id)

        return [`collection.Get(${id})`, Blockly.CSharp.ORDER_ATOMIC]
      }

      return `
      <block type="object_any"></block>
`
    },
    formCheckerBlock (blockConfig, object) {
      const methods = []
      blockConfig.items.forEach(item => {
        if (!item.method) {
          console.error('Undefined item method for block:', blockConfig)
          return
        }
        const methodTitle = this.getLocalizedName(item, '')
        methods.push([methodTitle, item.method])
      })
      if (!methods.length) {
        console.error('Empty methods for block:', blockConfig)
        return
      }

      const args = [
        {
          'type': 'input_value',
          'name': 'instance',
          'check': 'Object'
        },
        {
          'type': 'field_dropdown',
          'name': 'method',
          'options': methods
        }
      ]
      let argsMessage = ''
      let argsShadows = ''
      if (!blockConfig.args) {
        blockConfig.args = []
      }

      for (let i = 0; i < blockConfig.args.length; i++) {
        const arg = blockConfig.args[i]
        const name = 'arg' + i
        const placeholder = '%' + (args.length + 1)
        const title = this.getLocalizedName(arg, '')

        if (!arg.valueType) {
          console.error('Undefined valueType for arg:', arg)
          return
        }

        argsMessage += title + ' ' + placeholder

        let shadowBlock = this.getShadowBlockForValueType(arg.valueType)
        if (arg.values) {
          const valuesBlockType = this.getObjectInstanceUniqueId(object) + '_' + arg.values
          if (Blockly.Blocks[valuesBlockType]) {
            shadowBlock = `<shadow type="${valuesBlockType}"></shadow>`
          } else {
            console.error('Undefined values block:', valuesBlockType)
          }
        }
        argsShadows += `
          <value name="${name}">
            ${shadowBlock}
          </value>
`
        args.push({
          'type': 'input_value',
          'name': name,
          'check': arg.valueType
        })
      }

      const objectInstanceBlockType = this.getObjectInstanceUniqueId(object)
      const blockType = objectInstanceBlockType + '_' + blockConfig.name

      Blockly.Blocks[blockType] = {
        init () {
          this.jsonInit({
            'type': blockType,
            'message0': '%1 %2 ' + argsMessage,
            'args0': args,
            'inputsInline': true,
            'output': 'Boolean',
            'colour': '%{BKY_LOGIC_HUE}',
            'tooltip': '',
            'helpUrl': ''
          })
        }
      }
      Blockly.CSharp[blockType] = block => {
        const instance = Blockly.CSharp.valueToCode(block, 'instance', Blockly.CSharp.ORDER_ATOMIC) || 'null'
        const method = block.getFieldValue('method')

        const args = []
        for (let i = 0; i < blockConfig.args.length; i++) {
          args.push(Blockly.CSharp.valueToCode(block, 'arg' + i, Blockly.CSharp.ORDER_ATOMIC) || 'null')
        }

        const code = `((${object.config.type}Wrapper) ${instance}).${method}(${args.join(', ')})`

        return [code, Blockly.CSharp.ORDER_FUNCTION_CALL]
      }

      return `
        <block type="${blockType}">
            <value name="instance">
                <shadow type="${objectInstanceBlockType}"></shadow>
            </value>
            ${argsShadows}
        </block>
    `
    },
    formGetterBlock (blockConfig, object) {
      if (!blockConfig.valueType) {
        console.error('Undefined valueType for block:', blockConfig)
        return ''
      }

      const states = []
      blockConfig.items.forEach(item => {
        if (!item.property) {
          console.error('Undefined item property for block:', blockConfig)
          return ''
        }
        const methodTitle = this.getLocalizedName(item, item.property)
        states.push([methodTitle, item.property])
      })
      if (!states.length) {
        console.error('Empty states for block:', blockConfig)
        return ''
      }

      const objectInstanceBlockType = this.getObjectInstanceUniqueId(object)
      const blockType = objectInstanceBlockType + '_' + blockConfig.name
      Blockly.Blocks[blockType] = {
        init () {
          this.jsonInit({
            'type': blockType,
            'message0': '%1 %2',
            'args0': [
              {
                'type': 'input_value',
                'name': 'instance',
                'check': 'Object'
              },
              {
                'type': 'field_dropdown',
                'name': 'state',
                'options': states
              }
            ],
            'output': [
              blockConfig.valueType
            ],
            'colour': '%{BKY_VARIABLES_HUE}',
            'tooltip': '',
            'helpUrl': ''
          })
        }
      }
      Blockly.CSharp[blockType] = block => {
        const instance = Blockly.CSharp.valueToCode(block, 'instance', Blockly.CSharp.ORDER_ATOMIC) || 'null'
        const state = block.getFieldValue('state')
        const code = `((${object.config.type}Wrapper) ${instance}).${state}`

        Blockly.CSharp.use(`OpenWorld.Types.${object.config.type}`)

        return [code, Blockly.CSharp.ORDER_MEMBER]
      }

      return `
        <block type="${blockType}">
            <value name="instance">
                <shadow type="${objectInstanceBlockType}"></shadow>
            </value>
        </block>
    `
    },
    formSetterBlock (blockConfig, object) {
      if (!blockConfig.valueType) {
        console.error('Undefined valueType for block:', blockConfig)
        return ''
      }

      const states = []
      blockConfig.items.forEach(item => {
        if (!item.property) {
          console.error('Undefined item property for block:', blockConfig)
          return ''
        }
        const methodTitle = this.getLocalizedName(item, item.property)
        states.push([methodTitle, item.property])
      })
      if (!states.length) {
        console.error('Empty states for block:', blockConfig)
        return ''
      }

      const objectInstanceBlockType = this.getObjectInstanceUniqueId(object)
      const blockType = objectInstanceBlockType + '_' + blockConfig.name
      Blockly.Blocks[blockType] = {
        init () {
          this.jsonInit({
            'type': blockType,
            'message0': window.$t('editor.setterBlockPrefix') + '%1 %2 = %3',
            'args0': [
              {
                'type': 'input_value',
                'name': 'instance',
                'check': 'Object'
              },
              {
                'type': 'field_dropdown',
                'name': 'state',
                'options': states
              },
              {
                'type': 'input_value',
                'name': 'value',
                'check': blockConfig.valueType
              }
            ],
            'inputsInline': true,
            'nextStatement': null,
            'previousStatement': null,
            'colour': '%{BKY_VARIABLES_HUE}',
            'tooltip': '',
            'helpUrl': ''
          })
        }
      }
      Blockly.CSharp[blockType] = block => {
        const instance = Blockly.CSharp.valueToCode(block, 'instance', Blockly.CSharp.ORDER_ATOMIC) || 'null'
        const state = block.getFieldValue('state')
        const value = Blockly.CSharp.valueToCode(block, 'value', Blockly.CSharp.ORDER_ATOMIC) || null

        Blockly.CSharp.use(`OpenWorld.Types.${object.config.type}`)

        return `((${object.config.type}Wrapper) ${instance}).${state} = ${value};\n`
      }

      let shadowBlock = this.getShadowBlockForValueType(blockConfig.valueType)
      if (blockConfig.values) {
        const valuesBlockType = this.getObjectInstanceUniqueId(object) + '_' + blockConfig.values
        if (Blockly.Blocks[valuesBlockType]) {
          shadowBlock = `<shadow type="${valuesBlockType}"></shadow>`
        } else {
          console.error('Undefined values block:', valuesBlockType)
        }
      }

      return `
        <block type="${blockType}">
            <value name="instance">
                <shadow type="${objectInstanceBlockType}"></shadow>
            </value>
            <value name="value">
                ${shadowBlock}
            </value>
        </block>
    `
    },
    formValuesBlock (blockConfig, object) {
      const options = []
      blockConfig.items.forEach(item => {
        if (!item.name) {
          console.error('Undefined item name for block:', blockConfig)
          return ''
        }
        const title = this.getLocalizedName(item, item.name)
        options.push([title, item.name])
      })
      if (!options.length) {
        console.error('Empty values for block:', blockConfig)
        return ''
      }

      const objectInstanceBlockType = this.getObjectInstanceUniqueId(object)
      const blockType = objectInstanceBlockType + '_' + blockConfig.name
      Blockly.Blocks[blockType] = {
        init () {
          this.jsonInit({
            'type': blockType,
            'message0': '%1',
            'args0': [
              {
                'type': 'field_dropdown',
                'name': 'value',
                'options': options
              }
            ],
            'output': [
              'String'
            ],
            'colour': '%{BKY_VARIABLES_HUE}',
            'tooltip': '',
            'helpUrl': ''
          })
        }
      }
      Blockly.CSharp[blockType] = block => {
        const code = Blockly.CSharp.quote_(block.getFieldValue('value'))
        return [code, Blockly.CSharp.ORDER_ATOMIC]
      }

      return `
        <block type="${blockType}">
        </block>
    `
    },
    formActionBlock (blockConfig, object) {
      const methods = []
      blockConfig.items.forEach(item => {
        if (!item.method) {
          console.error('Undefined item method for block:', blockConfig)
          return
        }
        const methodTitle = this.getLocalizedName(item, '')
        methods.push([methodTitle, item.method])
      })
      if (!methods.length) {
        console.error('Empty methods for block:', blockConfig)
        return
      }

      const args = [
        {
          'type': 'input_value',
          'name': 'instance',
          'check': 'Object'
        },
        {
          'type': 'field_dropdown',
          'name': 'method',
          'options': methods
        }
      ]
      let argsMessage = ''
      let argsShadows = ''
      if (!blockConfig.args) {
        blockConfig.args = []
      }

      for (let i = 0; i < blockConfig.args.length; i++) {
        const arg = blockConfig.args[i]
        const name = 'arg' + i
        const placeholder = '%' + (args.length + 1)
        const title = this.getLocalizedName(arg, '')

        if (!arg.valueType) {
          console.error('Undefined valueType for arg:', arg)
          return
        }

        argsMessage += title + ' ' + placeholder

        let shadowBlock = this.getShadowBlockForValueType(arg.valueType)
        if (arg.values) {
          const valuesBlockType = this.getObjectInstanceUniqueId(object) + '_' + arg.values
          if (Blockly.Blocks[valuesBlockType]) {
            shadowBlock = `<shadow type="${valuesBlockType}"></shadow>`
          } else {
            console.error('Undefined values block:', valuesBlockType)
          }
        }
        argsShadows += `
          <value name="${name}">
            ${shadowBlock}
          </value>
`
        args.push({
          'type': 'input_value',
          'name': name,
          'check': arg.valueType
        })
      }

      const objectInstanceBlockType = this.getObjectInstanceUniqueId(object)
      const blockType = objectInstanceBlockType + '_' + blockConfig.name

      Blockly.Blocks[blockType] = {
        init () {
          this.jsonInit({
            'type': blockType,
            'message0': window.$t('editor.actionBlockPrefix') + '%1 %2 ' + argsMessage,
            'args0': args,
            'inputsInline': true,
            'nextStatement': null,
            'previousStatement': null,
            'colour': '%{BKY_ACTIONS_HUE}',
            'tooltip': '',
            'helpUrl': ''
          })
        }
      }
      Blockly.CSharp[blockType] = block => {
        const instance = Blockly.CSharp.valueToCode(block, 'instance', Blockly.CSharp.ORDER_ATOMIC) || 'null'
        const method = block.getFieldValue('method')

        const args = []
        for (let i = 0; i < blockConfig.args.length; i++) {
          args.push(Blockly.CSharp.valueToCode(block, 'arg' + i, Blockly.CSharp.ORDER_ATOMIC) || 'null')
        }

        return `((${object.config.type}Wrapper) ${instance}).${method}(${args.join(', ')});\n`
      }

      return `
        <block type="${blockType}">
            <value name="instance">
                <shadow type="${objectInstanceBlockType}"></shadow>
            </value>
            ${argsShadows}
        </block>
    `
    },
    formAnyObjectCheckBlock () {
      const blockType = 'check_object_state'
      Blockly.Blocks[blockType] = {
        init () {
          this.jsonInit({
            'type': blockType,
            'message0': '%{BKY_CHECK_OBJECT_STATE_TITLE}',
            'args0': [
              {
                'type': 'input_value',
                'name': 'instance',
                'check': 'Object'
              },
              {
                'type': 'field_dropdown',
                'name': 'method',
                'options': [
                  ['%{BKY_CHECK_OBJECT_STATE_ACTIVE}', 'IsActive'],
                  ['%{BKY_CHECK_OBJECT_STATE_INACTIVE}', 'IsInactive']
                ]
              }
            ],
            'output': 'Boolean',
            'colour': '%{BKY_LOGIC_HUE}',
            'tooltip': '%{BKY_CHECK_OBJECT_STATE_TOOLTIP}',
            'helpUrl': ''
          })
        }
      }
      Blockly.CSharp[blockType] = block => {
        const instance = Blockly.CSharp.valueToCode(block, 'instance', Blockly.CSharp.ORDER_ATOMIC) || 'null'
        const method = block.getFieldValue('method')
        const code = `${instance}.${method}()`

        return [code, Blockly.CSharp.ORDER_FUNCTION_CALL]
      }

      return `
        <block type="${blockType}">
            <value name="instance">
                <shadow type="object_any"></shadow>
            </value>
        </block>
    `
    },
    formAnyObjectActionBlock () {
      const blockType = 'object_actions'

      Blockly.Blocks[blockType] = {
        init () {
          this.jsonInit({
            'type': blockType,
            'message0': '%{BKY_OBJECT_ACTIONS_TITLE}',
            'args0': [
              {
                'type': 'input_value',
                'name': 'instance',
                'check': 'Object'
              },
              {
                'type': 'field_dropdown',
                'name': 'method',
                'options': [
                  ['%{BKY_OBJECT_ACTIONS_OPTION_ACTIVATE}', 'Activate'],
                  ['%{BKY_OBJECT_ACTIONS_OPTION_DEACTIVATE}', 'Deactivate']
                ]
              }],
            'inputsInline': true,
            'nextStatement': null,
            'previousStatement': null,
            'colour': '%{BKY_ACTIONS_HUE}',
            'tooltip': '',
            'helpUrl': ''
          })
        }
      }
      Blockly.CSharp[blockType] = block => {
        const instance = Blockly.CSharp.valueToCode(block, 'instance', Blockly.CSharp.ORDER_ATOMIC) || 'null'
        const method = block.getFieldValue('method')

        return `${instance}.${method}();\n`
      }

      return `
        <block type="${blockType}">
            <value name="instance">
                <shadow type="object_any"></shadow>
            </value>
        </block>
    `
    },
    formAnyObjectHierarchicPropertiesBlock () {
      const blockType = 'get_object_hierarchic_properties'
      Blockly.Blocks[blockType] = {
        init () {
          this.jsonInit({
            'type': blockType,
            'message0': '%{BKY_GET_OBJECT_HIERARCHIC_PROPERTIES_TITLE}',
            'args0': [
              {
                'type': 'input_value',
                'name': 'instance',
                'check': 'Object'
              },
              {
                'type': 'field_dropdown',
                'name': 'method',
                'options': [
                  ['%{BKY_GET_OBJECT_HIERARCHIC_PROPERTIES_PARENT}', 'GetParent']
                ]
              }
            ],
            'output': [
              'Object'
            ],
            'colour': '%{BKY_VARIABLES_HUE}',
            'tooltip': '%{BKY_GET_OBJECT_HIERARCHIC_PROPERTIES_TOOLTIP}',
            'helpUrl': ''
          })
        }
      }
      Blockly.CSharp[blockType] = block => {
        const instance = Blockly.CSharp.valueToCode(block, 'instance', Blockly.CSharp.ORDER_ATOMIC) || 'null'
        const method = block.getFieldValue('method')
        const code = `collection.${method}(${instance})`

        return [code, Blockly.CSharp.ORDER_MEMBER]
      }

      return `
        <block type="${blockType}">
            <value name="instance">
                <shadow type="object_any"></shadow>
            </value>
        </block>
    `
    },
    formAnyObjectHierarchicListsBlock () {
      const blockType = 'get_object_hierarchic_lists'
      Blockly.Blocks[blockType] = {
        init () {
          this.jsonInit({
            'type': blockType,
            'message0': '%{BKY_GET_OBJECT_HIERARCHIC_LISTS_TITLE}',
            'args0': [
              {
                'type': 'input_value',
                'name': 'instance',
                'check': 'Object'
              },
              {
                'type': 'field_dropdown',
                'name': 'method',
                'options': [
                  ['%{BKY_GET_OBJECT_HIERARCHIC_LISTS_CHILDREN}', 'GetChildren'],
                  ['%{BKY_GET_OBJECT_HIERARCHIC_LISTS_DESCENDANTS}', 'GetDescendants'],
                  ['%{BKY_GET_OBJECT_HIERARCHIC_LISTS_ANCESTRY}', 'GetAncestry'],
                ]
              }
            ],
            'output': [
              'Array'
            ],
            'colour': '%{BKY_LISTS_HUE}',
            'tooltip': '%{BKY_GET_OBJECT_HIERARCHIC_LISTS_TOOLTIP}',
            'helpUrl': ''
          })
        }
      }
      Blockly.CSharp[blockType] = block => {
        const instance = Blockly.CSharp.valueToCode(block, 'instance', Blockly.CSharp.ORDER_ATOMIC) || 'null'
        const method = block.getFieldValue('method')
        const code = `collection.${method}(${instance})`

        return [code, Blockly.CSharp.ORDER_MEMBER]
      }

      return `
        <block type="${blockType}">
            <value name="instance">
                <shadow type="object_any"></shadow>
            </value>
        </block>
    `
    },
    formAnyObjectGetterBlock () {
      const blockType = 'get_object_state'
      Blockly.Blocks[blockType] = {
        init () {
          this.jsonInit({
            'type': blockType,
            'message0': '%{BKY_GET_OBJECT_STATE_TITLE}',
            'args0': [
              {
                'type': 'input_value',
                'name': 'instance',
                'check': 'Object'
              },
              {
                'type': 'field_dropdown',
                'name': 'state',
                'options': [
                  ['%{BKY_GET_OBJECT_STATE_ACTIVITY}', 'Activity']
                ]
              }
            ],
            'output': [
              'Boolean'
            ],
            'colour': '%{BKY_VARIABLES_HUE}',
            'tooltip': '%{BKY_GET_OBJECT_STATE_TOOLTIP}',
            'helpUrl': ''
          })
        }
      }
      Blockly.CSharp[blockType] = block => {
        const instance = Blockly.CSharp.valueToCode(block, 'instance', Blockly.CSharp.ORDER_ATOMIC) || 'null'
        const state = block.getFieldValue('state')
        const code = `${instance}.${state}`

        return [code, Blockly.CSharp.ORDER_MEMBER]
      }

      return `
        <block type="${blockType}">
            <value name="instance">
                <shadow type="object_any"></shadow>
            </value>
        </block>
    `
    },
    formAnyObjectSetterBlock () {
      const blockType = 'set_object_state'
      Blockly.Blocks[blockType] = {
        init () {
          this.jsonInit({
            'type': blockType,
            'message0': window.$t('editor.setterBlockPrefix') + '%1 %2 = %3',
            'args0': [
              {
                'type': 'input_value',
                'name': 'instance',
                'check': 'Object'
              },
              {
                'type': 'field_dropdown',
                'name': 'state',
                'options': [
                  ['%{BKY_GET_OBJECT_STATE_ACTIVITY}', 'Activity']
                ]
              },
              {
                'type': 'input_value',
                'name': 'value',
                'check': 'Boolean'
              }
            ],
            'inputsInline': true,
            'nextStatement': null,
            'previousStatement': null,
            'colour': '%{BKY_VARIABLES_HUE}',
            'tooltip': '',
            'helpUrl': ''
          })
        }
      }
      Blockly.CSharp[blockType] = block => {
        const instance = Blockly.CSharp.valueToCode(block, 'instance', Blockly.CSharp.ORDER_ATOMIC) || 'null'
        const state = block.getFieldValue('state')
        const value = Blockly.CSharp.valueToCode(block, 'value', Blockly.CSharp.ORDER_ATOMIC) || null

        return `${instance}.${state} = ${value};\n`
      }

      return `
        <block type="${blockType}">
            <value name="instance">
                <shadow type="object_any"></shadow>
            </value>
            <value name="value">
                <shadow type="logic_boolean"></shadow>
            </value>
        </block>
    `
    },
    formEventBlock (blockConfig, object) {
      const objectName = this.getLocalizedName(object.config, object.config.type.toLowerCase())
      const methods = []
      blockConfig.items.forEach(item => {
        if (!item.method) {
          console.error('Undefined item method for block:', blockConfig)
          return
        }
        const methodTitle = this.getLocalizedName(item, item.method)
        methods.push([methodTitle, item.method])
      })
      if (!methods.length) {
        console.error('Empty methods for block:', blockConfig)
        return
      }

      const args = [
        {
          'type': 'field_dropdown',
          'name': 'method',
          'options': methods
        },
        {
          'type': 'input_dummy'
        },
        {
          'type': 'field_dropdown',
          'name': 'instance',
          'options': object.instances
        }
      ]
      let paramsMessage = ''
      if (!blockConfig.params) {
        blockConfig.params = []
      }

      if (blockConfig.params.length > 0) {
        args.push({
          'type': 'input_dummy'
        })
      }
      for (let i = 0; i < blockConfig.params.length; i++) {
        const param = blockConfig.params[i]
        paramsMessage += ' %' + (args.length + 1)

        if (!param.valueType) {
          console.error('Undefined valueType for param:', param)
          return
        }

        if (!param.name) {
          console.error('Undefined name for param:', param)
          return
        }

        args.push({
          'type': 'field_variable',
          'name': param.name,
          'variable': param.name
        })
      }

      const objectInstanceBlockType = this.getObjectInstanceUniqueId(object)
      const blockType = objectInstanceBlockType + '_' + blockConfig.name
      this.blocksScope.set(blockType, 'Events')

      args.push({
        'type': 'input_dummy'
      })

      args.push({
        'type': 'input_statement',
        'name': 'stack'
      })

      Blockly.Blocks[blockType] = {
        init () {
          this.jsonInit({
            'type': blockType,
            'message0': `${window.$t('editor.eventsBlockPrefix')} %1 %2 ${objectName} %3 %4 ${window.$t('editor.eventsBlockParams')}: ` + paramsMessage + ' %' + (args.length - 1) + ' %' + args.length,
            'args0': args,
            'colour': '%{BKY_PROCEDURES_HUE}',
            'tooltip': '',
            'helpUrl': ''
          })
        }
      }
      Blockly.CSharp[blockType] = block => {
        const id = block.getFieldValue('instance')
        Blockly.CSharp.lockInstance(id)
        const method = block.getFieldValue('method')
        Blockly.CSharp.use(`OpenWorld.Types.${object.config.type}`)
        let args = blockConfig.params.map(param => {
          return Blockly.CSharp.variableDB_.getName(block.getFieldValue(param.name), Blockly.Variables.NAME_TYPE)
        })
        let branch = Blockly.CSharp.statementToCode(block, 'stack')

        return `
collection.Get<${object.config.type}Wrapper>(${id}).${method} += (${args.join(', ')}) => {
  ${branch}
};`
      }

      return `
        <block type="${blockType}">
        </block>
    `
    },
    getShadowBlockForValueType (valueType) {
      switch (valueType) {
        case 'Boolean':
          return `<shadow type="logic_boolean"></shadow>`
        case 'Number':
          return `<shadow type="math_number"></shadow>`
        case 'String':
          return `<shadow type="text"></shadow>`
        case 'Object':
          return `<shadow type="object_any"></shadow>`
      }

      return ''
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
    getObjectInstanceUniqueId (object) {
      return 'o_' + object.guid.replace('-', '_')
    },
    close () {
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
  .logicLabel>.blocklyFlyoutLabelText {
    font-size: 20px;
    fill: #5a7fa5;
  }
  .actionsLabel>.blocklyFlyoutLabelText {
    font-size: 20px;
    fill: #73a45b;
  }
  .stateLabel>.blocklyFlyoutLabelText {
    font-size: 20px;
    fill: #a45a80;
  }
  .eventsLabel>.blocklyFlyoutLabelText {
    font-size: 20px;
    fill: #985aa5;
  }
</style>
