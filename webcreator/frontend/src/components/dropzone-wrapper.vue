<template>
  <form
    :action="url"
    :id="id"
    class="vue-dropzone dropzone">
    <slot/>

    <div
      :id="'template-' + id"
      style="display:none">
      <b-card>
        <div class="h4 m-0">
          {{ language.filenamePrefix }}<span data-dz-name/>
        </div>
        <div>
          {{ language.filesizePrefix }}<span data-dz-size/>
        </div>
        <div class="progress">
          <div
            class="progress-bar"
            role="progressbar"
            aria-valuemin="0"
            aria-valuemax="100"
            data-dz-uploadprogress/>
        </div>
        <div
          class="dz-status mt-1"
          data-dz-errormessage/>
      </b-card>
    </div>
  </form>
</template>

<script>
import BCard from 'bootstrap-vue/es/components/card/card'
import BAlert from 'bootstrap-vue/es/components/alert/alert'
import BProgress from 'bootstrap-vue/es/components/progress/progress'

export default {
  components: {
    BCard,
    BAlert,
    BProgress
  },
  props: {
    id: {
      type: String,
      required: true
    },
    url: {
      type: String,
      required: true
    },
    clickable: {
      type: [Boolean, String],
      default: true
    },
    confirm: {
      type: Function,
      default: undefined
    },
    paramName: {
      type: String,
      default: 'file'
    },
    acceptedFileTypes: {
      type: String,
      default: ''
    },
    thumbnailHeight: {
      type: Number,
      default: 100
    },
    thumbnailWidth: {
      type: Number,
      default: 200
    },
    showRemoveLink: {
      type: Boolean,
      default: true
    },
    maxFileSizeInMb: {
      type: Number,
      default: 2
    },
    maxNumberOfFiles: {
      type: Number,
      default: 5
    },
    autoProcessQueue: {
      type: Boolean,
      default: true
    },
    useFontAwesome: {
      type: Boolean,
      default: false
    },
    headers: {
      type: Object,
      default: () => {}
    },
    language: {
      type: Object,
      default: () => {}
    },
    useCustomDropzoneOptions: {
      type: Boolean,
      default: false
    },
    dropzoneOptions: {
      type: Object,
      default: () => {}
    },
    resizeWidth: {
      type: Number,
      default: null
    },
    resizeHeight: {
      type: Number,
      default: null
    },
    resizeMimeType: {
      type: String,
      default: null
    },
    resizeQuality: {
      type: Number,
      default: 0.8
    },
    resizeMethod: {
      type: String,
      default: 'contain'
    },
    uploadMultiple: {
      type: Boolean,
      default: false
    },
    duplicateCheck: {
      type: Boolean,
      default: false
    },
    parallelUploads: {
      type: Number,
      default: 2
    },
    timeout: {
      type: Number,
      default: 1000000
    },
    method: {
      type: String,
      default: 'POST'
    },
    withCredentials: {
      type: Boolean,
      default: false
    },
    capture: {
      type: String,
      default: null
    },
    hiddenInputContainer: {
      type: String,
      default: 'body'
    }
  },
  computed: {
    languageSettings () {
      let defaultValues = {
        dictDefaultMessage: '<br>Drop files here to upload',
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
      }

      for (let attrname in this.language) {
        defaultValues[attrname] = this.language[attrname]
      }

      if (this.useCustomDropzoneOptions) {
        if (this.dropzoneOptions.language) {
          for (let attrname in this.dropzoneOptions.language) {
            defaultValues[attrname] = this.dropzoneOptions.language[attrname]
          }
        }
      }

      return defaultValues
    },
    cloudIcon: function () {
      if (this.useFontAwesome) {
        return '<i class="fa fa-cloud-upload"></i>'
      } else {
        return '<i class="material-icons">cloud_upload</i>'
      }
    },
    doneIcon: function () {
      if (this.useFontAwesome) {
        return '<i class="fa fa-check"></i>'
      } else {
        return ' <i class="material-icons">done</i>'
      }
    },
    errorIcon: function () {
      if (this.useFontAwesome) {
        return '<i class="fa fa-warning"></i>'
      } else {
        return ' <i class="material-icons">error</i>'
      }
    }
  },
  beforeDestroy () {
    this.dropzone.destroy()
  },
  mounted () {
    if (this.$isServer) {
      return
    }
    let Dropzone = require('dropzone')
    Dropzone.autoDiscover = false
    if (this.confirm) {
      Dropzone.confirm = this.getProp(this.confirm, this.dropzoneOptions.confirm)
    }
    let element = document.getElementById(this.id)
    this.dropzone = new Dropzone(element, {
      clickable: this.getProp(this.clickable, this.dropzoneOptions.clickable),
      paramName: this.getProp(this.paramName, this.dropzoneOptions.paramName),
      thumbnailWidth: this.getProp(this.thumbnailWidth, this.dropzoneOptions.thumbnailWidth),
      thumbnailHeight: this.getProp(this.thumbnailHeight, this.dropzoneOptions.thumbnailHeight),
      maxFiles: this.getProp(this.maxNumberOfFiles, this.dropzoneOptions.maxNumberOfFiles),
      maxFilesize: this.getProp(this.maxFileSizeInMB, this.dropzoneOptions.maxFileSizeInMB),
      addRemoveLinks: this.getProp(this.showRemoveLink, this.dropzoneOptions.showRemoveLink),
      acceptedFiles: this.getProp(this.acceptedFileTypes, this.dropzoneOptions.acceptedFileTypes),
      autoProcessQueue: this.getProp(this.autoProcessQueue, this.dropzoneOptions.autoProcessQueue),
      headers: this.getProp(this.headers, this.dropzoneOptions.headers),
      previewTemplate: document.querySelector('#template-' + this.id).innerHTML,
      dictDefaultMessage: this.cloudIcon + this.languageSettings.dictDefaultMessage,
      dictCancelUpload: this.languageSettings.dictCancelUpload,
      dictCancelUploadConfirmation: this.languageSettings.dictCancelUploadConfirmation,
      dictFallbackMessage: this.languageSettings.dictFallbackMessage,
      dictFallbackText: this.languageSettings.dictFallbackText,
      dictFileTooBig: this.languageSettings.dictFileTooBig,
      dictInvalidFileType: this.languageSettings.dictInvalidFileType,
      dictMaxFilesExceeded: this.languageSettings.dictMaxFilesExceeded,
      dictRemoveFile: this.languageSettings.dictRemoveFile,
      dictRemoveFileConfirmation: this.languageSettings.dictRemoveFileConfirmation,
      dictResponseError: this.languageSettings.dictResponseError,
      resizeWidth: this.getProp(this.resizeWidth, this.dropzoneOptions.resizeWidth),
      resizeHeight: this.getProp(this.resizeHeight, this.dropzoneOptions.resizeHeight),
      resizeMimeType: this.getProp(this.resizeMimeType, this.dropzoneOptions.resizeMimeType),
      resizeQuality: this.getProp(this.resizeQuality, this.dropzoneOptions.resizeQuality),
      resizeMethod: this.getProp(this.resizeMethod, this.dropzoneOptions.resizeMethod),
      uploadMultiple: this.getProp(this.uploadMultiple, this.dropzoneOptions.uploadMultiple),
      parallelUploads: this.getProp(this.parallelUploads, this.dropzoneOptions.parallelUploads),
      timeout: this.getProp(this.timeout, this.dropzoneOptions.timeout),
      method: this.getProp(this.method, this.dropzoneOptions.method),
      capture: this.getProp(this.capture, this.dropzoneOptions.capture),
      hiddenInputContainer: this.getProp(this.hiddenInputContainer, this.dropzoneOptions.hiddenInputContainer),
      withCredentials: this.getProp(this.withCredentials, this.dropzoneOptions.withCredentials)
    })

    // Handle the dropzone events
    let vm = this

    this.dropzone.on('thumbnail', function (file, dataUrl) {
      vm.$emit('vdropzone-thumbnail', file, dataUrl)
    })

    this.dropzone.on('addedfile', function (file) {
      /**
       * If Duplicate Check enabled remove duplicate file and emit the event.
       */
      if (vm.duplicateCheck) {
        if (this.files.length) {
          var _i, _len
          for (_i = 0, _len = this.files.length; _i < _len - 1; _i++) {
            if (this.files[_i].name === file.name) {
              this.removeFile(file)
              vm.$emit('duplicate-file', file)
            }
          }
        }
      }

      vm.$emit('vdropzone-file-added', file)
    })

    this.dropzone.on('addedfiles', function (files) {
      vm.$emit('vdropzone-files-added', files)
    })

    this.dropzone.on('removedfile', function (file) {
      vm.$emit('vdropzone-removed-file', file)
    })

    this.dropzone.on('success', function (file, response) {
      file.previewElement.className += 'text-white bg-success'
      vm.$emit('vdropzone-success', file, response)
    })

    this.dropzone.on('successmultiple', function (file, response) {
      vm.$emit('vdropzone-success-multiple', file, response)
    })

    this.dropzone.on('error', function (file, error, xhr) {
      file.previewElement.className += 'text-white bg-danger'

      const statusEl = file.previewElement.querySelector('.dz-status')
      if (xhr) {
        try {
          const errInfo = JSON.parse(xhr.responseText)
          statusEl.innerHTML = errInfo['Message']
        } catch (e) {
          statusEl.innerHTML = xhr.responseText
        }
      }
      if (error) {
        statusEl.innerHTML = error
      }
      vm.$emit('vdropzone-error', file, error, xhr)
    })

    this.dropzone.on('sending', function (file, xhr, formData) {
      vm.$emit('vdropzone-sending', file, xhr, formData)
    })

    this.dropzone.on('sendingmultiple', function (file, xhr, formData) {
      vm.$emit('vdropzone-sending-multiple', file, xhr, formData)
    })

    this.dropzone.on('queuecomplete', function (file, xhr, formData) {
      vm.$emit('vdropzone-queue-complete', file, xhr, formData)
    })

    this.dropzone.on('totaluploadprogress', function (totaluploadprogress, totalBytes, totalBytesSent) {
      vm.$emit('vdropzone-total-upload-progress', totaluploadprogress, totalBytes, totalBytesSent)
    })

    vm.$emit('vdropzone-mounted')
  },
  methods: {
    manuallyAddFile: function (file, fileUrl, callback, crossOrigin, options) {
      this.dropzone.emit('addedfile', file)
      this.dropzone.emit('thumbnail', file, fileUrl)
      this.dropzone.createThumbnailFromUrl(file, fileUrl, callback, crossOrigin)
      this.dropzone.emit('complete', file)
      if ((typeof options.dontSubstractMaxFiles === undefined) || !options.dontSubstractMaxFiles) {
        this.dropzone.options['maxFiles'] = this.dropzone.options['maxFiles'] - 1
      }
      if ((typeof options.addToFiles !== undefined) && options.addToFiles) {
        this.dropzone.files.push(file)
      }
      this.$emit('vdropzone-file-added-manually', file)
    },
    setOption: function (option, value) {
      this.dropzone.options[option] = value
    },
    removeAllFiles: function () {
      this.dropzone.removeAllFiles(true)
    },
    processQueue: function () {
      let dropzoneEle = this.dropzone
      this.dropzone.processQueue()
      this.dropzone.on('success', function () {
        dropzoneEle.options.autoProcessQueue = true
      })
      this.dropzone.on('queuecomplete', function () {
        dropzoneEle.options.autoProcessQueue = false
      })
    },
    removeFile: function (file) {
      this.dropzone.removeFile(file)
    },
    getAcceptedFiles: function () {
      return this.dropzone.getAcceptedFiles()
    },
    getRejectedFiles: function () {
      return this.dropzone.getRejectedFiles()
    },
    getUploadingFiles: function () {
      return this.dropzone.getUploadingFiles()
    },
    getQueuedFiles: function () {
      return this.dropzone.getQueuedFiles()
    },
    getProp: function (attributeProp, objectProp) {
      if (!this.useCustomDropzoneOptions) {
        return attributeProp
      }

      if (objectProp !== undefined && objectProp !== null && objectProp !== '') {
        return objectProp
      }
      return attributeProp
    }
  }
}
</script>

<style>
@import url('~dropzone/dist/dropzone.css');

.vue-dropzone {
  border: 2px solid #E5E5E5;
  color: #2f353a;
  padding: 0.5rem;
}

.vue-dropzone:hover {
  background-color: #F6F6F6;
}

.dz-complete .progress, .dz-max-files-reached .progress {
  display: none;
}

.vue-dropzone .dz-status {
  display: none;
}

.dz-complete .dz-status, .dz-max-files-reached .dz-status {
  display: block;
}
</style>
