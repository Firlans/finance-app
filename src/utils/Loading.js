import { ref, createApp } from 'vue'
import LoadingFeature from '../components/features/LoadingFeaures.vue'

export class Loading {
  static TYPES = {
    INDETERMINATE: 'indeterminate',
    DETERMINATE: 'determinate',
  }

  constructor({ label = 'Loading...', type = Loading.TYPES.INDETERMINATE, duration = null } = {}) {
    this.show = ref(false)
    this.label = ref(label)
    this._baseLabel = ref(label)
    this.type = type
    this.duration = Number.isFinite(duration) ? duration : null
    this.remaining = ref(this.duration)
    this._timer = null
    this._container = null
    this._app = null

    if (typeof document !== 'undefined') {
      this._container = document.createElement('div')
      document.body.appendChild(this._container)
      this._app = createApp(LoadingFeature, { loading: this })
      this._app.mount(this._container)
    }
  }

  start({ label, type, duration } = {}) {
    if (label !== undefined) {
      this._baseLabel.value = label
    }

    if (type !== undefined) {
      this.type = type
    }

    if (duration !== undefined) {
      this.duration = Number.isFinite(duration) ? duration : null
    }

    this.show.value = true
    this.remaining.value = this.duration
    this._updateLabel()

    if (this.type === Loading.TYPES.DETERMINATE && this.duration > 0) {
      this._startCountdown()
    } else {
      this._clearTimer()
    }
  }

  stop() {
    this.show.value = false
    this._clearTimer()
  }

  update({ label, type, duration } = {}) {
    if (label !== undefined) {
      this._baseLabel.value = label
    }

    if (type !== undefined) {
      this.type = type
    }

    if (duration !== undefined) {
      this.duration = Number.isFinite(duration) ? duration : null
      this.remaining.value = this.duration
    }

    if (this.show.value) {
      this._updateLabel()

      if (this.type === Loading.TYPES.DETERMINATE && this.duration > 0) {
        this._clearTimer()
        this._startCountdown()
      } else {
        this._clearTimer()
      }
    }
  }

  _updateLabel() {
    if (this.type === Loading.TYPES.DETERMINATE && Number.isFinite(this.remaining.value)) {
      this.label.value = `${this._baseLabel.value} (${this.remaining.value}s remaining)`
    } else {
      this.label.value = this._baseLabel.value
    }
  }

  _startCountdown() {
    this._clearTimer()

    if (!this.duration || this.duration <= 0) {
      return
    }

    this.remaining.value = this.duration
    this._updateLabel()

    this._timer = setInterval(() => {
      this.remaining.value -= 1

      if (this.remaining.value <= 0) {
        this.stop()
      } else {
        this._updateLabel()
      }
    }, 1000)
  }

  _clearTimer() {
    if (this._timer) {
      clearInterval(this._timer)
      this._timer = null
    }
  }
}
