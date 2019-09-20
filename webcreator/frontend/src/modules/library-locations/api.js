import http from '@/utils/http'

export default {
  getList () {
    return http.get(`/v1/locations`)
  },

  delete (locationId) {
    if (!locationId) {
      return Promise.reject(new Error(`locationId is not passed`))
    }
    return http.delete(`/v1/delete-location/${locationId}`)
  },

  getTags (searchString = '') {
    return http.get(`/v1/location-tags${searchString ? `?search=${searchString}` : ''}`)
  },

  setTagsToLocation (locationId, data) {
    if (!locationId) {
      return Promise.reject(new Error(`locationId is not passed`))
    }
    return http.post(`/v1/update-location-tags/${locationId}`, data)
  },

  createTag (tag) {
    if (!tag) {
      return Promise.reject(new Error(`tag is not passed`))
    }
    return http.post(`/v1/add-location-tag`, tag)
  }
}
