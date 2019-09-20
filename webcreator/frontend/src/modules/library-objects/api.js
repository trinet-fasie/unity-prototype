import http from '@/utils/http'

export default {
  getList () {
    return http.get(`/v1/objects`)
  },

  delete (objectId) {
    if (!objectId) {
      return Promise.reject(new Error(`objectId is not passed`))
    }
    return http.delete(`/v1/delete-object/${objectId}`)
  },

  getTags (searchString = '') {
    return http.get(`/v1/object-tags${searchString ? `?search=${searchString}` : ''}`)
  },

  setTagsToObject (objectId, data) {
    if (!objectId) {
      return Promise.reject(new Error(`objectId is not passed`))
    }
    return http.post(`/v1/update-object-tags/${objectId}`, data)
  },

  createTag (tag) {
    if (!tag) {
      return Promise.reject(new Error(`tag is not passed`))
    }
    return http.post(`/v1/add-object-tag`, tag)
  }
}
