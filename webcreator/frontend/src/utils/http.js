import axios from 'axios'
import camelcaseKeys from 'camelcase-keys'

export default axios.create({
  baseUrl: `http://${process.env.VUE_APP_ROOT_API}:3000`,
  transformResponse: [(data) => {
    if (data) {
      return camelcaseKeys(JSON.parse(data), {deep: true})
    }
  }]
})
