import { ref, createApp } from 'vue'
import NotificationFeature from '../components/features/NotificationFeature.vue'

export class Notification {
  static TYPES = {
    SUCCESS: 'success',
    ERROR: 'error'
  }

  constructor({ duration = 2000 } = {}) {
    this.show = ref(false)
    this.message = ref('')
    this.type = Notification.TYPES.SUCCESS
    this.duration = Number.isFinite(duration) ? duration : 2000
    this._timer = null
    this._container = null
    this._app = null

    if (typeof document !== 'undefined') {
      this._container = document.createElement('div')
      document.body.appendChild(this._container)
      this._app = createApp(NotificationFeature, { notification: this })
      this._app.mount(this._container)
    }
  }

  _show({ message, type, duration } = {}) {
    if (message !== undefined) {
      this.message.value = message
    }

    if (type !== undefined) {
      this.type = type
    }

    if (duration !== undefined) {
      this.duration = Number.isFinite(duration) ? duration : 2000
    }

    this.show.value = true

    if (this._timer) {
      clearTimeout(this._timer)
    }

    this._timer = setTimeout(() => {
      this.show.value = false
      this._timer = null
    }, this.duration)
  }

  showSuccess(message, duration) {
    this._show({ message, type: Notification.TYPES.SUCCESS, duration })
  }

  showError(message, duration) {
    this._show({ message, type: Notification.TYPES.ERROR, duration })
  }
}
