import Stomp from '@stomp/stompjs'

export default {
  install (Vue, options) {
    options = Object.assign({
      stompUrl: 'ws://localhost:15674/ws',
      login: 'guest',
      password: 'guest',
      reconnectDelay: 5000,
      debug: (str) => {
        // console.log(str)
      }
    }, options)

    let stomp = null
    let connected = false
    let idleSubscriptions = []
    let subscriptionCounter = 0

    function connect () {
      if (stomp !== null) {
        return
      }
      stomp = Stomp.over(new WebSocket(options.stompUrl))
      stomp.debug = options.debug
      stomp.reconnect_delay = options.reconnectDelay
      stomp.connect(options.login, options.password, () => {
        connected = true
        idleSubscriptions.forEach((subscription) => {
          stomp.subscribe(...subscription)
        })
        idleSubscriptions = []
      }, () => {
        connected = false
      })
    }

    Vue.prototype.$messageBus = {
      subscribe (destination, callback, headers) {
        if (headers == null) {
          headers = {}
        }
        if (!headers.id) {
          headers.id = 'sub-' + subscriptionCounter++
        }

        if (connected) {
          stomp.subscribe(destination, callback, headers)
        } else {
          connect()
          idleSubscriptions.push([destination, callback, headers])
        }

        return headers.id
      },
      unsubscribe (id, headers) {
        if (connected) {
          stomp.unsubscribe(id, headers)
        } else {
          idleSubscriptions = idleSubscriptions.filter(subsciption => subsciption[2].id !== id)
        }
      }
    }
  }
}
