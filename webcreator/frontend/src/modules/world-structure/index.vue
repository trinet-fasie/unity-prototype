<template>
  <div class="wrapper">
    <div class="animated fadeIn">
      <b-row>
        <b-col cols="12">
          <h1
            v-if="worldName"
            class="mb-4">{{ worldName }}</h1>
          <b-alert
            v-if="error"
            show
            variant="danger">{{ error }}</b-alert>

          <b-tabs
            v-if="worldId"
            v-model="tabIndex">
            <b-tab
              :href="'/worlds/' + worldId"
              active>
              <template slot="title">
                <i class="fa fa-sitemap"/> {{ $t('world.structure.header') }}
              </template>
              <tree
                ref="tree"
                :data="structure"
                :options="treeOptions"
                class="world-structure"
                @node:collapsed="onCollapsedNode"
                @node:expanded="onExpandedNode">
                <span
                  slot-scope="{ node }"
                  class="tree-text">
                  <template v-if="node.data.type === nodeType.worldLocations">
                    <span class="node-text">
                      <i class="fas fa-images"/> {{ $t('world.structure.locations') }}
                    </span>
                  </template>
                  <template v-if="node.data.type === nodeType.worldConfigurations">
                    <span class="node-text">
                      <i class="fas fa-cogs"/> {{ $t('world.structure.configurations') }}
                    </span>
                  </template>
                  <template v-if="node.data.type === nodeType.worldLocation">
                    <span class="node-text">
                      <i class="fas fa-image"/> {{ node.data.name }}
                    </span>
                    <location-form-popover
                      :id="node.id"
                      :location="node.data"
                      @submit="data => updateWorldLocation(node, data)" />
                    <delete-with-confirm
                      :confirm="$t('world.structure.location.deleteConfirm')"
                      :id="node.id"
                      @submit="deleteWorldLocation(node)"/>
                  </template>
                  <template v-else-if="node.data.type === nodeType.addWorldLocation">
                    <location-form-popover
                      :id="node.id"
                      :location="node.data"
                      @submit="data => addWorldLocation(node, data)"/>
                  </template>
                  <template v-else-if="node.data.type === nodeType.group">
                    <span class="node-text">
                      <i class="far fa-object-group"/> {{ node.data.name }}
                    </span>
                    <group-form-popover
                      :id="node.id"
                      :group="node.data"
                      @submit="data => updateGroup(node, data)" />
                    <b-button
                      v-b-tooltip.hover="$t('actions.editInVR')"
                      :href="node.data.editLink"
                      variant="ghost-secondary">
                      <i class="fab fa-simplybuilt"/>
                    </b-button>
                    <b-button
                      v-b-tooltip.hover="$t('actions.showVisualEditor')"
                      variant="ghost-secondary"
                      @click.stop="openVisualEditorTab(node.parent, node)">
                      <i class="fa fa-puzzle-piece"/>
                    </b-button>
                    <b-button
                      v-b-tooltip.hover="$t('actions.showCodeEditor')"
                      variant="ghost-secondary"
                      @click.stop="openCodeEditorTab(node.parent, node)">
                      <i class="fa fa-code"/>
                    </b-button>
                    <delete-with-confirm
                      :confirm="$t('world.structure.group.deleteConfirm')"
                      :id="node.id"
                      variant="ghost-secondary"
                      @submit="deleteGroup(node)">
                      <i class="fas fa-trash"/>
                    </delete-with-confirm>
                  </template>
                  <template v-else-if="node.data.type === nodeType.addGroup">
                    <group-form-popover
                      :id="node.id"
                      :group="node.data"
                      @submit="data => addGroup(node, data)" />
                  </template>
                  <template v-else-if="node.data.type === nodeType.groupObject">
                    <span class="node-text">
                      <b-img
                        :src="node.data.objectIcon"
                        :alt="node.data.name"
                        :width="30"
                        fluid />
                      {{ node.data.name }}
                    </span>
                    <object-form-popover
                      :id="node.id"
                      :group-object="node.data"
                      @submit="data => updateGroupObject(node, data)" />
                  </template>
                  <template v-else-if="node.data.type === nodeType.worldConfiguration">
                    <span class="node-text">
                      <i class="fas fa-cog"/> {{ node.data.name }}
                    </span>
                    <b-button
                      v-b-tooltip.hover="$t('actions.edit')"
                      :id="'mf-btn-' + node.id"
                      variant="ghost-secondary">
                      <i class="fas fa-edit"/>
                    </b-button>
                    <b-button
                      v-b-tooltip.hover="$t('actions.viewInVR')"
                      :href="node.data.viewLink"
                      variant="ghost-secondary">
                      <i class="fab fa-simplybuilt"/>
                    </b-button>
                    <delete-with-confirm
                      :confirm="$t('world.structure.configuration.deleteConfirm')"
                      :id="node.id"
                      variant="ghost-secondary"
                      @submit="deleteWorldConfiguration(node)">
                      <i class="fas fa-trash"/>
                    </delete-with-confirm>
                    <configuration-form-popover
                      :id="node.id"
                      :configuration="node.data"
                      :world-locations-node="$refs.tree.find({id: nodeType.worldLocations})[0]"
                      :node-type="nodeType"
                      @submit="data => updateWorldConfiguration(node, data)" />
                  </template>
                  <template v-else-if="node.data.type === nodeType.addWorldConfiguration">
                    <span class="node-text">
                      <a
                        :id="'mf-btn-' + node.id"
                        href="javascript:void(0)"
                        @click.prevent="">
                        <i class="fa fa-plus"/> {{ $t('world.structure.configuration.addButton') }}
                      </a>
                    </span>
                    <configuration-form-popover
                      :id="node.id"
                      :configuration="node.data"
                      :world-locations-node="$refs.tree.find({id: nodeType.worldLocations})[0]"
                      :node-type="nodeType"
                      @submit="data => addWorldConfiguration(node, data)" />
                  </template>
                  <template v-else>
                    <span class="node-text">
                      {{ node.text }}
                    </span>
                  </template>
                </span>
              </tree>
            </b-tab>
            <b-tab
              v-for="(tab, i) in tabs"
              :key="tab.id"
              :id="tab.id">
              <template slot="title">
                <b-button
                  :aria-label="$t('actions.close')"
                  class="close"
                  @click.prevent="tryCloseTab(tab.id)">
                  <span
                    class="d-inline-block"
                    aria-hidden="true">&times;</span>
                </b-button>
                <template v-if="tab.type === tabType.visualEditor">
                  <span class="mr-2">
                    <i class="fa fa-puzzle-piece"/> {{ tab.locationNode.data.name }}: {{ tab.groupNode.data.name }}
                  </span>
                </template>
                <template v-if="tab.type === tabType.codeEditor">
                  <span class="mr-2">
                    <i class="fa fa-code"/> {{ tab.locationNode.data.name }}: {{ tab.groupNode.data.name }}
                  </span>
                </template>
              </template>
              <template v-if="tab.type === tabType.visualEditor">
                <visual-editor
                  :ref="tab.id"
                  :unique-class-postfix="`Group${tab.groupNode.data.id}`"
                  :vr-link="tab.groupNode.data.editLink"
                  :editor-config="tab.editorConfig"
                  :world-locations="tab.worldLocations"
                  :world-configurations="tab.worldConfigurations"
                  :editor-data="tab.groupNode.data.editorData"
                  @open-editor="openCodeEditorTab(tab.groupNode.parent, tab.groupNode)"
                  @change="(code, editorData, lockedInstances) => onVisualEditorChange(tab, code, editorData, lockedInstances)"
                  @close="closeTab(i)"/>
              </template>
              <template v-if="tab.type === tabType.codeEditor">
                <code-editor
                  :ref="tab.id"
                  :unique-class-postfix="`Group${tab.groupNode.data.id}`"
                  :vr-link="tab.groupNode.data.editLink"
                  :code="tab.groupNode.data.code"
                  @change="(code) => onCodeEditorChange(tab, code)"
                  @close="closeTab(i)"/>
              </template>
            </b-tab>
          </b-tabs>
        </b-col>
      </b-row>
    </div>
  </div>
</template>

<script>
import BRow from 'bootstrap-vue/es/components/layout/row'
import BCol from 'bootstrap-vue/es/components/layout/col'
import BCard from 'bootstrap-vue/es/components/card/card'
import BButton from 'bootstrap-vue/es/components/button/button'
import BAlert from 'bootstrap-vue/es/components/alert/alert'
import BTabs from 'bootstrap-vue/es/components/tabs/tabs'
import BTab from 'bootstrap-vue/es/components/tabs/tab'
import BImg from 'bootstrap-vue/es/components/image/img'
import Tree from 'liquor-tree'
import DeleteWithConfirm from '../../components/delete-with-confirm'
import LocationFormPopover from './components/location-form-popover'
import GroupFormPopover from './components/group-form-popover'
import ObjectFormPopover from './components/object-form-popover'
import ConfigurationFormPopover from './components/configuration-form-popover'
import VisualEditor from '../../components/editor/visual-editor'
import CodeEditor from '../../components/editor/code-editor'
import Vue from 'vue'

export default {
  components: {
    LocationFormPopover,
    GroupFormPopover,
    ObjectFormPopover,
    ConfigurationFormPopover,
    VisualEditor,
    CodeEditor,
    DeleteWithConfirm,
    BRow,
    BCol,
    BCard,
    BButton,
    BAlert,
    BTabs,
    BTab,
    BImg,
    Tree
  },
  data () {
    return {
      error: null,
      loading: false,
      tabIndex: 0,
      structure: null,
      objects: null,
      worldId: null,
      worldName: null,
      treeOptions: {
        multiple: false,
        parentSelect: true
      },
      nodeType: {
        worldLocations: 'wls',
        worldConfigurations: 'wcs',
        worldLocation: 'wl',
        addWorldLocation: 'awl',
        addWorldConfiguration: 'awc',
        group: 'g',
        addGroup: 'ag',
        groupObject: 'go',
        worldConfiguration: 'wc'
      },
      tabType: {
        visualEditor: 1,
        codeEditor: 2
      },
      tabs: [],
      lang: document.querySelector('html').getAttribute('lang'),
      groupChangedSubscription: null,
      expandedNodes: {}
    }
  },
  watch: {
    expandedNodes (newValue) {
      const worldId = this.$route.params['worldId']

      let expandedNodes = {}

      if (this.$session.exists('expandedNodes')) {
        expandedNodes = JSON.parse(this.$session.get('expandedNodes'))
      }

      expandedNodes[worldId] = newValue
      this.$session.set('expandedNodes', JSON.stringify(expandedNodes))
    }
  },
  mounted () {
    this.$root.$on('changed::tab', (tabs, tabIndex, tab) => {
      if (this.$refs[tab.id]) {
        window.setTimeout(() => {
          if (this.$refs[tab.id][0]) {
            this.$refs[tab.id][0].refresh()
          }
        })
      }
    })

    const worldId = this.$route.params['worldId']

    this.expandedNodes = {}
    this.expandedNodes[this.nodeType.worldLocations] = true
    this.expandedNodes[this.nodeType.worldConfigurations] = true

    if (this.$session.exists('expandedNodes')) {
      const expandedNodes = JSON.parse(this.$session.get('expandedNodes'))
      if (expandedNodes[worldId]) {
        this.expandedNodes = expandedNodes[worldId]
      }
    }
  },
  created () {
    this.fetchData()
    this.groupChangedSubscription = this.$messageBus.subscribe('/exchange/owd.group_objects.changed/world.' + this.$route.params['worldId'] + '.#', (message) => {
      this.refreshGroupObjects(JSON.parse(message.body))
    })
  },
  destroyed () {
    this.$messageBus.unsubscribe(this.groupChangedSubscription)
  },
  beforeRouteLeave (to, from, next) {
    for (let i = 0; i < this.tabs.length; i++) {
      if (this.$refs[this.tabs[i].id] && this.$refs[this.tabs[i].id][0]) {
        if (!this.$refs[this.tabs[i].id][0].close()) {
          next(false)
          return
        }
      }
    }
    next()
  },
  methods: {
    refreshGroupObjects (groupData) {
      const groupId = groupData['GroupId']
      const groupNodes = this.$refs.tree.find({id: 'gr-' + groupId})
      if (!groupNodes) {
        return
      }
      const groupNode = groupNodes[0]

      const nodesToRemove = []
      groupNode.children.forEach(node => {
        nodesToRemove.push(node)
      })
      nodesToRemove.forEach(node => {
        node.remove()
      })

      groupData['Objects'].forEach(objectData => {
        const object = this.formObject(objectData)
        this.objects[object.id] = object
      })

      this.$nextTick(() => {
        groupData['GroupObjects'].forEach(childObject => {
          groupNode.append(this.groupObjectToNode(childObject))
        })

        // replace objects in Visual Editor
        if (this.$refs['visual-editor-' + groupId] && this.$refs['visual-editor-' + groupId][0]) {
          this.$refs['visual-editor-' + groupId][0].config = this.formVisualEditorConfig(groupNode)
        }
      })
    },
    fetchData () {
      this.error = null
      this.loading = true
      this.structure = null
      this.worldId = null
      this.worldName = null
      this.objects = null
      this.$http.get('/v1/world-structure/' + this.$route.params['worldId'])
        .then(response => {
          this.loading = false
          if (response.data['Status'] === 'success') {
            this.worldId = response.data['Data']['WorldId']
            this.worldName = response.data['Data']['WorldName']

            this.objects = {}
            response.data['Data']['Objects'].forEach(objectData => {
              const object = this.formObject(objectData)
              this.objects[object.id] = object
            })

            this.structure = []

            const locations = {
              id: this.nodeType.worldLocations,
              data: {
                type: this.nodeType.worldLocations,
                name: ''
              },
              state: {
                selectable: false,
                expanded: this.expandedNodes[this.nodeType.worldLocations] || false
              },
              children: []
            }

            response.data['Data']['WorldLocations'].forEach(worldLocation => {
              locations.children.push(this.worldLocationToNode(worldLocation))
            })

            locations.children.push({
              data: {
                id: null,
                type: this.nodeType.addWorldLocation,
                locationId: null,
                name: ''
              },
              state: {
                selectable: false
              }
            })
            this.structure.push(locations)

            const configurations = {
              id: this.nodeType.worldConfigurations,
              data: {
                type: this.nodeType.worldConfigurations,
                name: ''
              },
              state: {
                selectable: false,
                expanded: this.expandedNodes[this.nodeType.worldConfigurations] || false
              },
              children: []
            }

            response.data['Data']['WorldConfigurations'].forEach(worldConfiguration => {
              configurations.children.push(this.worldConfigurationToNode(worldConfiguration))
            })

            configurations.children.push({
              data: {
                id: null,
                type: this.nodeType.addWorldConfiguration,
                name: ''
              },
              state: {
                selectable: false
              }
            })
            this.structure.push(configurations)
          }
        })
        .catch(err => {
          this.loading = false
          this.error = this.$getResponseErrorMessage(err)
        })
    },
    onCollapsedNode (node) {
      Vue.delete(this.expandedNodes, node.id)
    },
    onExpandedNode (node) {
      Vue.set(this.expandedNodes, node.id, true)
    },
    formObject (objectData) {
      return {
        id: objectData['Id'],
        guid: objectData['Guid'],
        type: objectData['Type'],
        icon: this.$config.api.baseUrl + objectData['Resources']['Icon'],
        config: objectData['Config']
      }
    },
    worldLocationToNode (worldLocation) {
      let groups = []
      if (worldLocation['Groups']) {
        worldLocation['Groups'].forEach(group => {
          groups.push(this.groupToNode(group, worldLocation['LocationId']))
        })
      }

      groups.push({
        text: this.$t('world.structure.group.addButton'),
        data: {
          type: this.nodeType.addGroup,
          name: ''
        },
        state: {
          selectable: false
        }
      })

      return {
        text: worldLocation['Name'],
        id: 'wl-' + worldLocation['Id'],
        data: {
          type: this.nodeType.worldLocation,
          id: worldLocation['Id'],
          sid: worldLocation['Sid'],
          locationId: worldLocation['LocationId'],
          name: worldLocation['Name']
        },
        children: groups,
        state: {
          selectable: false,
          expanded: this.expandedNodes['wl-' + worldLocation['Id']] || false
        }
      }
    },
    groupToNode (group, locationId) {
      const groupObjects = []
      if (group['GroupObjects']) {
        group['GroupObjects'].forEach(groupObject => {
          groupObjects.push(this.groupObjectToNode(groupObject))
        })
      }

      const editConfig = btoa(JSON.stringify({
        gm: 2,
        lang: this.lang,
        worldId: this.worldId,
        groupId: group['Id'],
        locationId: locationId,
        api: this.$config.api,
        web: this.$config.web,
        rabbitmq: this.$config.rabbitmq,
        photon: this.$config.photon
      }))

      return {
        text: group['Name'],
        id: 'gr-' + group['Id'],
        data: {
          type: this.nodeType.group,
          id: group['Id'],
          name: group['Name'],
          code: group['Code'],
          editLink: 'owlp://' + editConfig,
          editorData: group['EditorData']
        },
        children: groupObjects,
        state: {
          selectable: false,
          expanded: this.expandedNodes['gr-' + group['Id']] || false
        }
      }
    },
    groupObjectToNode (groupObject) {
      const childNodes = []
      groupObject['GroupObjects'].forEach(childObject => {
        childNodes.push(this.groupObjectToNode(childObject))
      })

      return {
        text: groupObject['Name'],
        id: 'go-' + groupObject['Id'],
        data: {
          type: this.nodeType.groupObject,
          id: groupObject['Id'],
          name: groupObject['Name'],
          instanceId: groupObject['InstanceId'],
          objectId: groupObject['ObjectId'],
          objectIcon: this.objects[groupObject['ObjectId']] ? this.objects[groupObject['ObjectId']].icon : ''
        },
        state: {
          selectable: false,
          expanded: childNodes.length > 0 && this.expandedNodes['go-' + groupObject['Id']] || false
        },
        children: childNodes
      }
    },
    worldConfigurationToNode (worldConfiguration) {
      const viewConfig = btoa(JSON.stringify({
        gm: 0,
        lang: this.lang,
        worldId: this.worldId,
        worldConfigurationId: worldConfiguration['Id'],
        api: this.$config.api,
        web: this.$config.web,
        rabbitmq: this.$config.rabbitmq,
        photon: this.$config.photon
      }))

      return {
        text: worldConfiguration['Name'],
        id: 'wc-' + worldConfiguration['Id'],
        data: {
          type: this.nodeType.worldConfiguration,
          id: worldConfiguration['Id'],
          sid: worldConfiguration['Sid'],
          name: worldConfiguration['Name'],
          groupIds: worldConfiguration['GroupIds'],
          startWorldLocationId: worldConfiguration['StartWorldLocationId'],
          viewLink: 'owlp://' + viewConfig
        },
        state: {
          selectable: false
        }
      }
    },
    addWorldLocation (addButtonNode, data) {
      this.error = null
      this.loading = true
      this.$http.put('/v1/add-world-location', {
        'WorldId': this.worldId,
        'LocationId': data.locationId,
        'Name': data.name
      })
        .then(response => {
          this.loading = false
          if (response.data['Status'] === 'success') {
            const node = addButtonNode.parent.insertAt(this.worldLocationToNode(response.data['Data']), addButtonNode.parent.children.length - 1)
            node.expand()
          }
        })
        .catch(err => {
          this.loading = false
          this.error = this.$getResponseErrorMessage(err)
        })
    },
    updateWorldLocation (node, data) {
      this.error = null
      this.loading = true
      this.$http.post('/v1/update-world-location/' + node.data.id, {
        'Name': data.name,
        'LocationId': data.locationId
      })
        .then(response => {
          this.loading = false
          if (response.data['Status'] === 'success') {
            node.data.name = response.data['Data']['Name']
            node.data.locationId = response.data['Data']['LocationId']
          }
        })
        .catch(err => {
          this.loading = false
          this.error = this.$getResponseErrorMessage(err)
        })
    },
    deleteWorldLocation (node) {
      this.error = null
      this.loading = true
      this.$http.delete('/v1/delete-world-location/' + node.data.id)
        .then(response => {
          this.loading = false
          if (response.data['Status'] === 'success') {
            node.remove()
          }
        })
        .catch(err => {
          this.loading = false
          this.error = this.$getResponseErrorMessage(err)
        })
    },
    addGroup (addButtonNode, data) {
      this.error = null
      this.loading = true
      this.$http.put('/v1/add-group', {
        'WorldLocationId': addButtonNode.parent.data.id,
        'Name': data.name
      })
        .then(response => {
          this.loading = false
          if (response.data['Status'] === 'success') {
            addButtonNode.parent.insertAt(this.groupToNode(response.data['Data'], addButtonNode.parent.data.locationId), addButtonNode.parent.children.length - 1)
          }
        })
        .catch(err => {
          this.loading = false
          this.error = this.$getResponseErrorMessage(err)
        })
    },
    updateGroup (node, data) {
      this.error = null
      this.loading = true
      this.$http.post('/v1/update-group/' + node.data.id, {
        'Name': data.name
      })
        .then(response => {
          this.loading = false
          if (response.data['Status'] === 'success') {
            node.data.name = response.data['Data']['Name']
          }
        })
        .catch(err => {
          this.loading = false
          this.error = this.$getResponseErrorMessage(err)
        })
    },
    deleteGroup (node) {
      this.error = null
      this.loading = true
      this.$http.delete('/v1/delete-group/' + node.data.id)
        .then(response => {
          this.loading = false
          if (response.data['Status'] === 'success') {
            node.remove()
          }
        })
        .catch(err => {
          this.loading = false
          this.error = this.$getResponseErrorMessage(err)
        })
    },

    findParentGroup (node) {
      if (node.parent === null) {
        return null
      }
      if (node.parent.data.type !== this.nodeType.group) {
        return this.findParentGroup(node.parent)
      } else {
        return node.parent
      }
    },

    isObjectNameUnique (node, newName) {
      if (node.data.name === newName) {
        return
      }

      const parentNode = this.findParentGroup(node)
      if (parentNode) {
         return !this.getFlattenGroupObjectsData(parentNode).some(object => object.name === newName)
      }
      return true
    },

    updateGroupObject (node, data) {
      if (!this.isObjectNameUnique(node, data.name)) {
        this.error = this.$t('world.structure.groupObject.alreadyNamed')
        return
      }

      this.error = null
      this.loading = true
      this.$http.post('/v1/update-group-object/' + node.data.id, {
        'Name': data.name
      })
        .then(response => {
          this.loading = false
          if (response.data['Status'] === 'success') {
            node.data.name = response.data['Data']['Name']

            // Update name in Visual Editor
            const groupId = node.parent.data.id
            if (this.$refs['visual-editor-' + groupId] && this.$refs['visual-editor-' + groupId][0]) {
              this.$refs['visual-editor-' + groupId][0].config = this.formVisualEditorConfig(node.parent)
            }
          }
        })
        .catch(err => {
          this.loading = false
          this.error = this.$getResponseErrorMessage(err)
        })
    },
    addWorldConfiguration (addButtonNode, data) {
      this.error = null
      this.loading = true

      this.$http.put('/v1/add-world-configuration', {
        'WorldId': this.worldId,
        'Name': data.name,
        'GroupIds': data.groupIds,
        'StartWorldLocationId': data.startWorldLocationId
      })
        .then(response => {
          this.loading = false
          if (response.data['Status'] === 'success') {
            addButtonNode.parent.insertAt(this.worldConfigurationToNode(response.data['Data']), addButtonNode.parent.children.length - 1)
          }
        })
        .catch(err => {
          this.loading = false
          this.error = this.$getResponseErrorMessage(err)
        })
    },
    updateWorldConfiguration (node, data) {
      this.error = null
      this.loading = true
      this.$http.post('/v1/update-world-configuration/' + node.data.id, {
        'Name': data.name,
        'GroupIds': data.groupIds,
        'StartWorldLocationId': data.startWorldLocationId
      })
        .then(response => {
          this.loading = false
          if (response.data['Status'] === 'success') {
            node.data.name = response.data['Data']['Name']
            node.data.groupIds = response.data['Data']['GroupIds']
            node.data.startWorldLocationId = response.data['Data']['StartWorldLocationId']
          }
        })
        .catch(err => {
          this.loading = false
          this.error = this.$getResponseErrorMessage(err)
        })
    },
    deleteWorldConfiguration (node) {
      this.error = null
      this.loading = true
      this.$http.delete('/v1/delete-world-configuration/' + node.data.id)
        .then(response => {
          this.loading = false
          if (response.data['Status'] === 'success') {
            node.remove()
          }
        })
        .catch(err => {
          this.loading = false
          this.error = this.$getResponseErrorMessage(err)
        })
    },
    openVisualEditorTab (locationNode, groupNode) {
      let existTabIndex = null
      const tabId = 'visual-editor-' + groupNode.data.id
      for (let i = 0; i < this.tabs.length; i++) {
        if (this.tabs[i].id === tabId) {
          existTabIndex = i
          break
        }
      }
      if (existTabIndex !== null) {
        this.tabIndex = existTabIndex + 1
        return
      }

      this.tabs.push({
        id: tabId,
        type: this.tabType.visualEditor,
        locationNode: locationNode,
        groupNode: groupNode,
        editorConfig: this.formVisualEditorConfig(groupNode),
        worldLocations: this.getWorldLocationsData(),
        worldConfigurations: this.getWorldConfigurationsData()
      })
      window.setTimeout(() => {
        this.tabIndex = this.tabs.length
      })
    },
    tryCloseTab (tabId) {
      if (this.$refs[tabId] && this.$refs[tabId][0]) {
        this.$refs[tabId][0].close()
      }
    },
    closeTab (index) {
      this.tabs.splice(index, 1)
      this.tabIndex = 0
    },
    onVisualEditorChange (tab, code, editorData, lockedInstances) {
      tab.groupNode.data.editorData = editorData
      tab.groupNode.data.code = code

      this.error = null
      this.loading = true
      this.$http.post('/v1/update-group/' + tab.groupNode.data.id, {
        'EditorData': editorData,
        'LockedInstances': lockedInstances,
        'Code': code
      })
        .then(() => {
          this.loading = false
        })
        .catch(err => {
          this.loading = false
          this.error = this.$getResponseErrorMessage(err)
        })
    },
    onCodeEditorChange (tab, code) {
      tab.groupNode.data.code = code

      this.error = null
      this.loading = true
      this.$http.post('/v1/update-group/' + tab.groupNode.data.id, {
        'Code': code
      })
        .then(() => {
          this.loading = false
        })
        .catch(err => {
          this.loading = false
          this.error = this.$getResponseErrorMessage(err)
        })
    },
    openCodeEditorTab (locationNode, groupNode) {
      let existTabIndex = null
      const tabId = 'code-editor-' + groupNode.data.id
      for (let i = 0; i < this.tabs.length; i++) {
        if (this.tabs[i].id === tabId) {
          existTabIndex = i
          break
        }
      }
      if (existTabIndex !== null) {
        this.tabIndex = existTabIndex + 1
        return
      }

      this.tabs.push({
        id: tabId,
        type: this.tabType.codeEditor,
        locationNode: locationNode,
        groupNode: groupNode
      })
      window.setTimeout(() => {
        this.tabIndex = this.tabs.length
      })
    },
    formVisualEditorConfig (groupNode) {
      const objects = {}
      const allInstances = []

      this.getFlattenGroupObjectsData(groupNode).forEach(groupObjectData => {
        const objectId = groupObjectData['objectId']
        if (!objects[objectId]) {
          const object = this.objects[objectId]
          objects[objectId] = {
            guid: object.guid,
            config: object.config,
            instances: []
          }
        }
        const instance = [groupObjectData.name, groupObjectData.instanceId + '']
        objects[objectId].instances.push(instance)
        allInstances.push(instance)
      })

      return {
        'objects': objects,
        'allInstances': allInstances
      }
    },
    getFlattenGroupObjectsData (groupNode) {
      const result = []
      if (groupNode.data.type === this.nodeType.groupObject) {
        result.push(groupNode.data)
      }

      groupNode.children.forEach(childNode => {
        result.push(...this.getFlattenGroupObjectsData(childNode))
      })

      return result
    },
    getWorldLocationsData () {
      const worldLocationsNodes = this.$refs.tree.find({id: this.nodeType.worldLocations})
      if (!worldLocationsNodes || !worldLocationsNodes[0]) {
        return []
      }

      const result = []
      worldLocationsNodes[0].children.forEach(childNode => {
        if (childNode.data.type === this.nodeType.worldLocation) {
          result.push(childNode.data)
        }
      })

      return result
    },
    getWorldConfigurationsData () {
      const worldConfigurationsNodes = this.$refs.tree.find({id: this.nodeType.worldConfigurations})
      if (!worldConfigurationsNodes || !worldConfigurationsNodes[0]) {
        return []
      }

      const result = []
      worldConfigurationsNodes[0].children.forEach(childNode => {
        if (childNode.data.type === this.nodeType.worldConfiguration) {
          result.push(childNode.data)
        }
      })

      return result
    }
  }
}
</script>

<style>
  .world-structure a {
    text-decoration: none
  }
  .world-structure .node-text {
    display: inline-block;
    vertical-align: middle;
    line-height: 1.5;
    padding: 0.25rem 0.5rem 0.25rem  0;
  }

  .world-structure .tree-anchor {
    padding: 0.25rem 0;
    margin: 0;
  }

  .world-structure .tree > .tree-root,
  .world-structure .tree > .tree-filter-empty {
    padding: 0;
  }
</style>
