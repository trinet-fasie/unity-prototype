export default {
  nav: {
    library: 'Library',
    objects: 'Objects',
    locations: 'Locations',
    worlds: 'Worlds',
    worldStructure: 'World structure',
    worldConfigurations: 'World configurations',
    groupLogicEditor: 'Logic editor'
  },
  lang: {
    ru: {
      title: 'Русский',
      icon: 'flag-icon flag-icon-ru'
    },
    en: {
      title: 'English',
      icon: 'flag-icon flag-icon-us'
    }
  },
  error404: {
    message: 'Page not found.',
    details: 'We could not find the page you were looking for.'
  },
  objects: {
    list: 'Objects list',
    installUpdated: 'Object "{name}" updated successfully',
    installCreated: 'Object "{name}" created successfully',
    deleteRestriction: 'It is not possible to delete an object while it is being used in worlds'
  },
  locations: {
    list: 'Locations list',
    installUpdated: 'Location "{name}" updated successfully',
    installCreated: 'Location "{name}" created successfully',
    deleteRestriction: 'While the location is used in worlds, it is not possible to remove it'
  },
  worlds: {
    list: 'Worlds list',
    configurations: 'Configurations',
    deleteConfirm: 'When you delete a world, it will erase all settings and logic. You cannot cancel this operation.',
    addButton: 'Add world',
    addTitle: 'New world',
    editTitle: 'Edit world',
    importButton: 'Import world',
    importTitle: 'Choose file for import',
    fields: {
      name: 'Name',
      file: 'File'
    }
  },
  world: {
    structure: {
      header: 'World structure',
      locations: 'Locations',
      location: {
        deleteConfirm: 'When you delete a location, all groups in it will be deleted. You cannot cancel the operation.',
        addButton: 'Add location',
        addTitle: 'New location',
        editTitle: 'Edit location',
        fields: {
          name: 'Name',
          locationId: 'Location type'
        }
      },
      group: {
        deleteConfirm: 'When you delete a group, all objects and their logic are deleted. You cannot cancel the operation.',
        addButton: 'Add group',
        addTitle: 'New group',
        editTitle: 'Edit group',
        fields: {
          name: 'Name'
        }
      },
      groupObject: {
        editTitle: 'Edit object',
        alreadyNamed: 'An object with the same name already exists',
        fields: {
          name: 'Name'
        }
      },
      configurations: 'Configurations',
      configuration: {
        deleteConfirm: 'When you delete a configuration, all settings are deleted. You will not be able to cancel the operation.',
        addButton: 'Add configuration',
        addTitle: 'New configuration',
        editTitle: 'Edit configuration',
        emptyGroupsAlert: 'At least one location and group must be created to add a configuration.',
        fields: {
          name: 'Name',
          selectedGroups: 'Used groups',
          startWorldLocationId: 'Starting location'
        },
        validation: {
          startWorldLocationRequired: 'Select the groups to use and then the starting location'
        }
      }
    }
  },
  deleteConfirm: {
    title: 'Are you sure?'
  },
  actions: {
    delete: 'Remove',
    cancel: 'Cancel',
    add: 'Add',
    edit: 'Edit',
    apply: 'Apply',
    applyAndOpenEditor: 'Apply and show code',
    upload: 'Upload',
    close: 'Close',
    reset: 'Reset changes',
    vr: 'Run VR',
    viewInVR: 'Run in VR in view mode',
    editInVR: 'Run in VR in edit mode',
    showVisualEditor: 'Open the visual logic editor',
    showCodeEditor: 'Open the code editor',
    exportToFile: 'Export to file',
    buildPackage: 'Build package',
    importFromFile: 'Import from file'
  },
  validation: {
    required: 'Required field',
    success: '<i class="fa fa-check"></i> Success'
  },
  installDropzone: {
    filenamePrefix: 'Installing ',
    filesizePrefix: 'Size: ',
    close: 'Close',
    dictDefaultMessage: '<br>Drop files here to install',
    dictCancelUpload: 'Cancel upload',
    dictCancelUploadConfirmation: 'Are you sure you want to cancel this upload?',
    dictFallbackMessage: 'Your browser does not support drag and drop file uploads.',
    dictFallbackText: 'Please use the fallback form below to upload your files like in the olden days.',
    dictFileTooBig: 'File is too big ({{filesize}}MiB). Max filesize: {{maxFilesize}}MiB.',
    dictInvalidFileType: `You can't upload files of this type.`,
    dictMaxFilesExceeded: 'You can not upload any more files. (max: {{maxFiles}})',
    dictRemoveFile: 'Remove',
    dictRemoveFileConfirmation: null,
    dictResponseError: 'Server responded with {{statusCode}} code.'
  },
  table: {
    filter: 'Filter:',
    filterClear: 'Reset',
    emptyText: 'There are no records to show',
    emptyFilteredText: 'There are no records matching your request',
    columns: {
      id: 'ID',
      name: 'Name',
      type: 'Type',
      usages: 'Usages',
      createdAt: 'Created',
      updatedAt: 'Updated',
      tags: 'Tags',
      actions: 'Actions'
    }
  },
  editor: {
    category: {
      logic: 'Logic',
      actions: 'Actions',
      objects: 'Objects',
      state: 'Variables',
      lists: 'Lists',
      loops: 'Loops',
      math: 'Math',
      text: 'Text',
      variables: 'Variables',
      functions: 'Procedures',
      events: 'Events',
      other: 'Other'
    },
    anyObject: 'any',
    actionBlockPrefix: 'do',
    setterBlockPrefix: 'set',
    eventsBlockPrefix: 'on',
    eventsBlockParams: 'params',
    allObjectInstancesTitle: '{name} list',
    unsavedChangesConfirm: 'Do you really want to leave? You have unsaved changes!'
  },
  tags: {
    addTag: 'Add tag'
  }
}
