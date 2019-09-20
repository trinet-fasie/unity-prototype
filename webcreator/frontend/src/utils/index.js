import moment from 'moment'
import i18n from '@/i18n'

export const convertToISODate = (dateString) => {
  return moment(dateString, moment.ISO_8601)
}

/* istanbul ignore next */
export const convertToLocalFullDate = (date) => {
  return i18n.d(date, 'full')
}

export const sortObjectsByKey = (a, b, key) => {
  if ((typeof a[key] === 'number' && typeof b[key] === 'number') || (typeof a[key] === 'object' && typeof b[key] === 'object')) {
    return a[key] < b[key] ? -1 : (a[key] > b[key] ? 1 : 0)
  } else {
    return (a[key].toString()).localeCompare(b[key].toString(), undefined, {
      numeric: true
    })
  }
}
