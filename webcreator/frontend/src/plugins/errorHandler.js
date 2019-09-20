export default {
  install (Vue) {
    Vue.prototype.$getResponseErrorMessage = err => {
      if (err.response) {
        if (err.response.data.message) {
          return err.response.data.message
        }

        if (err.response.data.status === 'fail') {
          return 'Bad request. Invalid data.'
        }

        return 'Error response. Status: ' + err.response.status
      }

      if (err.request) {
        return 'Bad request'
      }

      if (err.message) {
        return err.message
      }

      this.error = 'Unknown error'
    }
  }
}
