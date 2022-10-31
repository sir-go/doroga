import axios from 'axios'
import Spinner from 'vue-simple-spinner'

const MSG_NO_FILE_API = 'Ваш браузер не поддерживает работу с локальными файлами, ' +
  'установите, например, последний google-chrome или firefox'

// noinspection ES6ModulesDependencies, JSUnresolvedVariable
const serverTimeout = process.env['NODE_ENV'] === 'production' ? 300 : 300 // Sec
const maxFileSize = 5 // MBytes
// noinspection ES6ModulesDependencies, JSUnresolvedVariable
const inactiveTimeoutLimit = process.env['NODE_ENV'] === 'production' ? 300 : 300 // Sec
const tightBreakPoint = 600 // px

// noinspection ES6ModulesDependencies, JSUnresolvedVariable
const BACKEND_URL = '/'
const MSG_TOO_BIG_FILE = 'Файл не должен быть больше ' + maxFileSize + ' мегабайт'
const MSG_NO_PHOTO = 'Вы не выбрали фото'
const MSG_SEND_SUCCESS = 'Данные успешно отправлены'
const MSG_SEND_FAILED = 'Сервер не отвечает, попробуйте позже'
const MSG_FILE_FORMAT = 'Недопустимый формат файла! Поддерживаются: jpeg, png, gif, tiff, bmp'
const ALLOWED_PHOTO_FORMATS = ['image/jpeg', 'image/gif', 'image/png', 'image/bmp', 'image/tiff']

const html = document.documentElement

if (window.File && window.FileReader && window.FileList && window.Blob) {
} else {
  alert(MSG_NO_FILE_API)
  throw new Error("browser doesn't support File Api")
}

// noinspection JSUnusedGlobalSymbols
export default {
  name: 'App',
  components: {
    VueSimpleSpinner: Spinner
  },
  data () {
    return {
      isEdge: window.navigator.userAgent.indexOf('Edge') > -1,
      // scrolled: {prev: 0, now: 0, direction: 0, delta: 0},
      scrolled: {now: 0, offset: 0},
      allowedPhotoFormats: ALLOWED_PHOTO_FORMATS,
      phWidth: null,
      showArrow: false,
      arrowShowen: false,
      inactiveTimeout: null,
      isMounted: false,
      busy: false,
      valids: {
        photo: false,
        name: false,
        sender_name: false,
        phone: false
      },
      l0titleTop: '20vh',
      l0titleDisplay: true,
      fields: [
        {
          name: 'name',
          placeholder: 'Ф.И.О.',
          req: true,
          check_cb: this.checkValid
        }, {
          name: 'bplace',
          placeholder: 'место рождения'
        }, {
          name: 'years',
          placeholder: 'годы жизни'
        }, {
          name: 'vdate',
          placeholder: 'дата призыва'
        }, {
          name: 'vplace',
          placeholder: 'пункт призыва'
        }, {
          name: 'rang',
          placeholder: 'воинское звание'
        }, {
          name: 'awards',
          placeholder: 'награды'
        }, {
          name: 'phDate',
          placeholder: 'когда сделано фото'
        }
      ],
      phData: null
    }
  },
  computed: {
    isFilled () {
      return Object.values(this.valids).every((val) => {
        return val
      })
    },
    // psL0b () {
    //   return {bottom: -35 + ((this.scrolled.now > 0.17 ? 0.17 : this.scrolled.now) * 200) + 'vh'}
    // },
    tight () {
      return html.clientWidth <= tightBreakPoint || html.clientHeight <= tightBreakPoint
    }
  },
  methods: {
    checkValid (ev) {
      let el = ev.target
      this.valids[el.name] = (
        typeof el.value !== 'undefined' && el.value !== null && el.value.trim().length > 0
      )
    },

    fileSelected (ev) {
      let fInput = ev.target
      if (fInput.files && fInput.files[0]) {
        if (fInput.files[0].size > maxFileSize * 1e6) {
          alert(MSG_TOO_BIG_FILE)
          ev.preventDefault()
          fInput.value = ''
          return
        }

        if (!ALLOWED_PHOTO_FORMATS.includes(fInput.files[0].type)) {
          alert(MSG_FILE_FORMAT)
          ev.preventDefault()
          fInput.value = ''
          return
        }

        let reader = new FileReader()

        reader.onload = (loadEv) => {
          this.phData = loadEv.target['result']
          this.valids.photo = true
        }
        reader.readAsDataURL(fInput.files[0])
      }
    },

    submit () {
      if (this.phData === null) {
        alert(MSG_NO_PHOTO)
        return
      }
      let form = this.$refs['form']
      let formData = new FormData(form)

      // noinspection JSUnusedGlobalSymbols
      this.busy = true
      axios.post(BACKEND_URL, formData, {
        headers: {'Content-Type': 'multipart/form-data'},
        timeout: serverTimeout * 1000
      })
        .then(() => {
          alert(MSG_SEND_SUCCESS)
          form.reset()
          this.phData = null
          this.valids.photo = false
        })
        .catch(() => {
          alert(MSG_SEND_FAILED)
        })
        .finally(() => {
          // noinspection JSUnusedGlobalSymbols
          this.busy = false
        })
    },

    scrollHandler () {
      window.requestAnimationFrame(() => {
        this.scrolled = {
          // prev: this.scrolled.now,
          now: Math.round(window.pageYOffset / (html.scrollHeight - html.clientHeight) * 1000) / 1000,
          offset: window.pageYOffset
        }

        if (this.inactiveTimeout != null && this.scrolled.offset > 50) {
          clearTimeout(this.inactiveTimeout)
          this.inactiveTimeout = null
          this.showArrow = false
          // noinspection JSUnusedGlobalSymbols
          this.arrowShowen = true
        }

        if (this.isMounted && !this.tight) {
          let s1 = this.$refs['s1']
          if (s1.getBoundingClientRect().top > 0) {
            this.l0titleTop = (20 - this.scrolled.now * 100).toFixed(2) + 'vh'
            this.l0titleDisplay = true
          } else {
            this.l0titleDisplay = false
          }
        }
      })
    }
  },
  created () {
    window.addEventListener('scroll', this.scrollHandler)
    // window.addEventListener('resize', this.scrollHandler)
  },
  mounted () {
    this.scrollHandler()
    // noinspection JSUnusedGlobalSymbols
    this.inactiveTimeout = setTimeout(() => {
      if (!this.arrowShowen) {
        this.showArrow = true
        this.arrowShowen = true
      }
    }, 1000 * inactiveTimeoutLimit)
    this.isMounted = true
  },
  beforeDestroy () {
    window.removeEventListener('scroll', this.scrollHandler)
    // window.removeEventListener('resize', this.scrollHandler)
  }
}
