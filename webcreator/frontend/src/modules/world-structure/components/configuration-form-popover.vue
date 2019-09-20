<template>
  <modal-form
    :id="id"
    :title="fields.id ? $t('world.structure.configuration.editTitle') : $t('world.structure.configuration.addTitle')"
    :submit-label="fields.id ? $t('actions.apply') : $t('actions.add')"
    :can-submit="canSubmit"
    :show-submit="hasGroups"
    @submit="onSubmit"
    @hidden="onHidden"
    @show="onShow">
    <div v-if="hasGroups">
      <b-form-group
        :label="$t('world.structure.configuration.fields.name')"
        :label-for="nameInputId"
        :state="nameState"
        :invalid-feedback="$t('validation.required')"
        :valid-feedback="$t('validation.success')"
        class="mb-2">
        <b-form-input
          :id="nameInputId"
          :state="nameState"
          v-model="fields.name"
          required />
      </b-form-group>
      <b-form-group :label="$t('world.structure.configuration.fields.selectedGroups')">
        <tree
          v-if="visible"
          ref="selectGroupsTree"
          :data="selectGroupsTreeData"
          :options="selectGroupsTreeOptions"
          class="tree-select"
          @node:checked="onChangeSelectedGroups"
          @node:unchecked="onChangeSelectedGroups" />
      </b-form-group>
      <b-form-group
        :label="$t('world.structure.configuration.fields.startWorldLocationId')"
        :label-for="startLocationInputId"
        :state="startWorldLocationState"
        :invalid-feedback="$t('world.structure.configuration.validation.startWorldLocationRequired')"
        :valid-feedback="$t('validation.success')">
        <b-form-select
          v-model="fields.startWorldLocationId"
          :state="startWorldLocationState"
          :options="chosenLocations"
          required/>
      </b-form-group>
    </div>
    <div v-else>
      {{ $t('world.structure.configuration.emptyGroupsAlert') }}
    </div>
  </modal-form>
</template>

<script>
import ModalForm from '../../../components/modal-form'
import BFormGroup from 'bootstrap-vue/es/components/form-group/form-group'
import BFormInput from 'bootstrap-vue/es/components/form-input/form-input'
import BButton from 'bootstrap-vue/es/components/button/button'
import BFormSelect from 'bootstrap-vue/es/components/form-select/form-select'
import Tree from 'liquor-tree'
import BAlert from 'bootstrap-vue/es/components/alert/alert'

export default {
  components: {
    ModalForm,
    BFormGroup,
    BFormInput,
    BButton,
    BFormSelect,
    Tree,
    BAlert
  },
  props: {
    id: {
      type: [String, Number],
      required: true
    },
    configuration: {
      type: Object,
      required: true
    },
    worldLocationsNode: {
      type: Object,
      required: true
    },
    nodeType: {
      type: Object,
      required: true
    },
    defaultValues: {
      type: Object,
      default: () => {
        return {
          name: '',
          groupIds: [],
          startWorldLocationId: null
        }
      }
    }
  },
  data () {
    return {
      visible: false,
      hasGroups: false,
      nameInputId: 'cf-name-' + this.id,
      startLocationInputId: 'cf-start-location-' + this.id,
      fields: Object.assign({}, this.defaultValues, this.configuration),
      selectGroupsTreeData: [],
      selectGroupsTreeOptions: {
        checkbox: true,
        checkOnSelect: true,
        parentSelect: true
      },
      chosenLocations: []
    }
  },
  computed: {
    nameState () {
      return this.fields.name.length > 0
    },
    startWorldLocationState () {
      return this.fields.startWorldLocationId > 0
    },
    canSubmit () {
      return this.hasGroups && this.nameState && this.startWorldLocationState
    }
  },
  methods: {
    onChangeSelectedGroups () {
      this.$nextTick(() => {
        this.chosenLocations = []
        this.fields.groupIds = []
        let resetStartLocation = true
        this.$refs.selectGroupsTree.findAll({}).forEach(node => {
          if (!node.states.checked && !node.states.indeterminate) {
            return
          }
          if (node.data.type === this.nodeType.worldLocation) {
            if (node.data.id === this.fields.startWorldLocationId) {
              resetStartLocation = false
            }
            this.chosenLocations.push({
              value: node.data.id,
              text: node.text
            })
          }
          if (node.data.type === this.nodeType.group) {
            this.fields.groupIds.push(node.data.id)
          }
        })
        if (resetStartLocation) {
          this.fields.startWorldLocationId = null
        }
      })
    },
    onSubmit () {
      this.$emit('submit', this.fields)
    },
    onHidden () {
      this.visible = false
    },
    onShow () {
      this.fields = Object.assign({}, this.defaultValues, this.configuration)
      this.selectGroupsTreeData = []
      const selectedLocations = new Map()
      this.chosenLocations = []
      this.hasGroups = false
      this.worldLocationsNode.children.forEach(worldLocationNode => {
        if (worldLocationNode.data.type !== this.nodeType.worldLocation) {
          return
        }
        const groupNodes = []
        worldLocationNode.children.forEach(groupNode => {
          if (groupNode.data.type !== this.nodeType.group) {
            return
          }

          if (this.fields.groupIds.includes(groupNode.data.id) && !selectedLocations.has(groupNode.parent.data.id)) {
            selectedLocations.set(groupNode.parent.data.id, true)
            this.chosenLocations.push({
              value: groupNode.parent.data.id,
              text: groupNode.parent.text
            })
          }

          groupNodes.push({
            text: groupNode.data.name,
            state: {
              expanded: true,
              selectable: false,
              checked: this.fields.groupIds.includes(groupNode.data.id)
            },
            data: {
              id: groupNode.data.id,
              type: groupNode.data.type
            }
          })
        })

        if (groupNodes.length > 0) {
          this.selectGroupsTreeData.push({
            text: worldLocationNode.data.name,
            children: groupNodes,
            state: {
              expanded: true,
              selectable: false
            },
            data: {
              id: worldLocationNode.data.id,
              type: worldLocationNode.data.type
            }
          })
          this.hasGroups = true
        }
      })
      this.visible = true
    }
  }
}
</script>

<style>
  .tree-select .tree-arrow {
    display: none;
  }

  .tree-select .tree-checkbox {
    width: 16px;
    height: 16px;
  }

  .tree-select .tree-checkbox.checked:after {
    left: 4px;
    top: 1px;
    height: 7px;
    width: 4px;
  }

  .tree-select .tree-anchor {
    line-height: 14px;
    padding-left: 5px !important;
  }

  .tree-select .tree-node:not(.selected) > .tree-content:hover {
    background: inherit;
  }
</style>
