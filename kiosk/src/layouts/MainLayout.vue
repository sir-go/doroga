<template>
  <q-layout view="lHr LPR lFr">
    <q-header elevated class="bg-primary text-white">
      <q-toolbar class="bg-red-10" style="padding: 2rem 7rem 0 7rem">
        <q-toolbar-title>Проект «Дорога Памяти»</q-toolbar-title>
        <q-space />
        <q-tabs v-model="tab" align="justify" class="col-5">
          <q-tab name="list" label="Список Героев" class="col-5" />
          <q-tab name="index" label="О проекте" class="col-5" />
<!--          <q-tab name="help" label="Помощь" />-->
        </q-tabs>
      </q-toolbar>
    </q-header>

    <q-tab-panels
      v-model="tab"
      swipeable
      style="top: 5rem"
    >

      <q-tab-panel name="list" class="q-pa-none tab-list">
        <q-dialog transition-show="none" transition-hide="none" auto-close v-model="infoPad" full-width full-height v-if="this.activeRecord != null">
          <div v-if="this.activeRecord != null" class="bg-white row p-dialog">
              <div class="col-6 bg-black q-pr-sm" style="box-shadow: #6f6f6f 2px 0 8px 1px; border-radius: 0 10px 10px 0;">
                <div style="height: 100%; display: flex; align-items: center;">
                  <q-img
                    contain
                    class="rounded-borders"
                    style="max-height: 976px"
                    :src="'/photos/o/' + this.activeRecord.id + '.jpg'"
                  />
                </div>
              </div>
              <div class="col-6">
                <div class="row justify-end">
                  <q-btn
                    flat
                    size="xl"
                    icon="close"
                    v-close-popup
                  />
                </div>
                <div class="p-descr">
                  <div class="p-name">
                    <div v-if="this.activeRecord.name">{{ this.activeRecord.name }}</div>
                    <div class="text-grey-6" v-else>- имя не указано -</div>
                  </div>
                  <div v-if="this.activeRecord['years']">
                    <q-separator spaced />
                    <div class="p-label">годы жизни</div>
                    <div class="p-value">{{ this.activeRecord['years'] }}</div>
                  </div>
                  <div v-if="this.activeRecord['bplace']">
                    <q-separator spaced />
                    <div class="p-label">Родина</div>
                    <div class="p-value">{{ this.activeRecord['bplace'] }}</div>
                  </div>
                  <div v-if="this.activeRecord['vdate']">
                    <q-separator spaced />
                    <div class="p-label">призыв</div>
                    <div class="p-value">{{ this.activeRecord['vdate'] }}</div>
                    <div v-if="this.activeRecord['vplace']">
                      <div class="p-value">{{ this.activeRecord['vplace'] }}</div>
                    </div>
                  </div>
                  <div v-if="this.activeRecord['rang']">
                    <q-separator spaced />
                    <div class="p-label">военское звание</div>
                    <div class="p-value">{{ this.activeRecord['rang'] }}</div>
                  </div>
                  <div v-if="this.activeRecord['awards']">
                    <q-separator spaced />
                    <div class="p-label">награды</div>
                    <div class="p-value">{{ this.activeRecord['awards'] }}</div>
                  </div>
                  <div v-if="this.activeRecord['info']">
                    <q-separator spaced />
                    <div class="p-label">биография</div>
                    <div class="p-value">{{ this.activeRecord['info'] }}</div>
                  </div>
                </div>
              </div>
          </div>
        </q-dialog>

        <q-dialog v-model="notFound" @hide="showSearchPan">
          <q-card>
            <q-card-section style="font-size: 26pt; padding: 2rem 3rem 1rem 3rem; color: #616161">
              По Вашему запросу ничего не найдено
            </q-card-section>

            <q-card-actions align="right">
              <q-btn flat label="OK" color="primary" v-close-popup />
            </q-card-actions>
          </q-card>
        </q-dialog>

        <div class="row">
          <q-scroll-area
            visible
            :thumb-style="{
              borderRadius: '5px',
              backgroundColor: 'grey',
              width: '1rem',
              opacity: 1
            }"
            ref="scrollArea"
            class="q-pl-lg q-pr-md col-11"
            :style="{
              height: searchPan.active ? '34rem' : '50rem',
              visibility: !loading
            }"
          >
            <q-card
              flat
              bordered
              v-for="r in records" :key="r.id"
              class="ds-card q-ma-sm"
              style="border-radius: 6px; border-bottom: 2px solid grey;"
              v-ripple
              @click="() => { showInfo(r) }"
              v-touch-hold.mouse="() => { showInfo(r) }"
            >
              <q-card-section horizontal>
                <q-card-section class="q-pa-none">
                  <q-avatar rounded size="100px">
                    <q-img :src="'/photos/t/' + r.id + '.jpg'" />
                  </q-avatar>
                </q-card-section>
                <q-card-section class="q-pt-xs q-pb-xs col-5">
                  <div>
                    <div class="card-name" v-if="r.name">
                      <text-highlight v-if="/name/.test(filter.search_fields)" :queries="filter.search">
                        {{ r.name }}
                      </text-highlight>
                      <div v-else>{{ r.name }}</div>
                    </div>
                    <div class="card-name text-grey-6" v-else>- имя не указано -</div>
                  </div>
                  <div class="row" v-if="r.years">
                    <div class="col-3 text-grey-8">годы жизни</div>
                    <div class="col text-no-wrap ellipsis">
                      <text-highlight v-if="/years/.test(filter.search_fields)" :queries="filter.search">
                        {{ r.years }}
                      </text-highlight>
                      <div v-else>{{ r.years }}</div>
                    </div>
                  </div>
                  <div class="row" v-if="r.bplace">
                    <div class="col-3 text-grey-8">Родина</div>
                    <div class="col text-no-wrap ellipsis">
                      <text-highlight v-if="/bplace/.test(filter.search_fields)" :queries="filter.search">
                      {{ r.bplace }}
                      </text-highlight>
                      <div v-else>{{ r.bplace }}</div>
                    </div>
                  </div>
                </q-card-section>
                <q-card-section class="q-pa-xs" v-if="r.info">
                  <div class="col-3 text-grey-8">военский путь</div>
                  <div class="col-5 ellipsis-3-lines">
                    <text-highlight v-if="/info/.test(filter.search_fields)" :queries="filter.search">
                      {{ r.info }}
                    </text-highlight>
                    <div v-else>{{ r.info }}</div>
                  </div>
                </q-card-section>
              </q-card-section>
            </q-card>
          </q-scroll-area>
          <q-card class="col-1 shadow-3 q-ml-sm" style="height: 30rem; width: 5rem">
            <q-pagination
              icon-prev="keyboard_arrow_up"
              icon-next="keyboard_arrow_down"
              v-model="pagination.current"
              :max="pagination.max"
              :max-pages="6"
              :boundary-numbers="true"
              :direction-links="true"
              color="grey-8"
              size="1.2rem"
              @input="fetchRecords"
            />
          </q-card>
        </div>

        <q-page-sticky position="bottom-right" :offset="[0, -2]" v-if="searchPan.active"
                       style="z-index: 5000 !important">
          <q-btn
            flat
            icon="keyboard_hide"
            text-color="grey-9"
            label="скрыть"
            size="lg"
            align="left"
            class="bg-white"
            style="width: 15rem; box-shadow: -7px -7px 8px 0 #00000029;"
            @click="hideSearchPan"
          />
        </q-page-sticky>

        <q-footer
          elevated
          class="text-black q-pb-xl bg-white"
          :style="searchPan.active ? {
              bottom: 0,
            } : {
              bottom: '-20rem',
              width: '70rem',
              left: '2rem',
              paddingTop: '1.5rem'
            }"
        >
          <div class="row flex-center full-width">
            <q-card id="searchPan" flat square class="bg-transparent">
              <!--  search field         -->
              <div
                class="row q-mb-md justify-between bg-white q-pl-md q-pr-sm"
                :style="!searchPan.active ? {
                  borderRadius: '0.5rem', boxShadow: '3px 3px 10px #a7a7a7'
                } : {}"
              >
                <q-input
                  :borderless="!searchPan.active"
                  v-model="filter.search"
                  placeholder="поиск"
                  type="search"
                  class="col-9"
                  style="font-size: 16pt;"
                  ref="searchInput"
                  @focus="showSearchPan"
                >
                  <template v-slot:prepend>
                    <q-icon class="gt-xs" name="search"/>
                  </template>
                  <template v-slot:after>
                    <q-btn
                      v-if="!searchIsEmpty"
                      text-color="red-10"
                      :ripple="false"
                      flat
                      icon="clear"
                      @click="searchIsEmpty = true"
                      size="lg"
                    />
                  </template>
                </q-input>

                <q-select
                  menu-anchor="top left"
                  menu-self="bottom left"
                  :borderless="!searchPan.active"
                  v-model="filter.search_fields"
                  :options="searchPan.fields.options"
                  label="где искать"
                  map-options
                  emit-value
                  class="col-2"
                >
                </q-select>
              </div>
              <!--  search field end        -->
                <!--            keyboard -->
                <div :class="['column', searchPan.active ? 'q-pt-none' : 'q-pt-xl']">
                  <div class="row">
                    <q-btn push @click="KLclick" class="kb-b kb-l sml col" :label="i" v-bind:key="i"
                           v-for="i in ['Ё', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0']"/>
                    <q-btn push color="grey-4" text-color="red-8"
                           :disable="searchIsEmpty"
                           v-touch-repeat:300:100.mouse.enter.space="BSclick"
                           @click="BSclick" class="kb-b kb-bs col-2" icon="backspace"/>
                  </div>
                  <div class="row q-ml-sm">
                    <q-btn push @click="KLclick" class="kb-b kb-l col" :label="i" v-bind:key="i"
                           v-for="i in ['Й', 'Ц', 'У', 'К', 'Е', 'Н', 'Г', 'Ш', 'Щ', 'З', 'Х', 'Ъ']"/>
                  </div>
                  <div class="row q-ml-lg q-mr-md">
                    <q-btn push @click="KLclick" class="kb-b kb-l col" :label="i" v-bind:key="i"
                           v-for="i in ['Ф', 'Ы', 'В', 'А', 'П', 'Р', 'О', 'Л', 'Д', 'Ж', 'Э', '-']"/>
                  </div>
                  <div class="row">

                    <div class="col">
                      <div class="row">
                          <q-btn push @click="KLclick" class="kb-b kb-l col" :label="i" v-bind:key="i"
                                 v-for="i in ['Я', 'Ч', 'С', 'М', 'И', 'Т', 'Ь', 'Б', 'Ю', '.']"/>
                      </div>
                      <div class="row">
                        <q-btn :disable="searchIsEmpty" push @click="SPCclick" class="kb-b kb-spc col" icon="space_bar"/>
                      </div>
                    </div>

                    <div class="col-2">
                      <q-btn
                        :disable="searchIsEmpty"
                        color="grey-4"
                        text-color="green-8"
                        push
                        @click="fetchRecords"
                        class="kb-sr q-mt-xs q-ml-sm"
                        icon="search"
                      />
                    </div>
                  </div>
                </div>
              <!--            keyboard end -->
            </q-card>
          </div>
        </q-footer>
      </q-tab-panel>

      <q-tab-panel name="index" class="tab-index no-padding">
        <IndexTab />
      </q-tab-panel>
<!--      <q-tab-panel name="help" class="tab-help">-->
<!--      </q-tab-panel>-->

    </q-tab-panels>
  </q-layout>
</template>

<script>
import IndexTab from '../components/IndexTab'

const SERVER_TIMEOUT = 5
const PAGE_SIZE = 25

// ScrollAnywhere   extension  for chrome

export default {
  components: { IndexTab },
  data () {
    return {
      notFound: false,
      tab: 'list',
      searchPan: {
        active: false,
        fields: {
          options: [
            { label: 'ФИО', value: 'name' },
            { label: 'БИОГРАФИЯ', value: 'info' },
            { label: 'ВЕЗДЕ', value: 'name,info,bplace,years,vdate,vplace,rang,awards' }
          ]
        }
      },
      filter: {
        applyed: false,
        search: '',
        search_fields: 'name'
      },
      activeRecord: null,
      records: [],
      pagination: {
        max: 500,
        current: 1
      },
      loading: true,
      infoPad: false
    }
  },
  computed: {
    searchIsEmpty: {
      get () {
        return this.filter.search === null || this.filter.search === ''
      },
      set () {
        this.filter.search = ''
        this.filter.search_fields = 'name'
        this.pagination.current = 1
        this.fetchRecords()
      }
    }
  },
  mounted () {
    this.fetchRecords()
  },
  methods: {
    showInfo (rec) {
      this.activeRecord = rec
      this.infoPad = true
    },
    showSearchPan () {
      this.searchPan.active = true
    },
    hideSearchPan () {
      this.searchPan.active = false
    },
    KLclick (ev) {
      if (this.filter.search === null) {
        this.filter.search = ev.target.outerText
      } else {
        this.filter.search += ev.target.outerText
      }
      this.$refs.searchInput.$el.focus()
    },
    BSclick () {
      if (this.filter.search !== '') this.filter.search = this.filter.search.slice(0, -1)
      if (this.searchIsEmpty) this.fetchRecords()
      this.$refs.searchInput.$el.focus()
    },
    SPCclick () {
      if (this.filter.search !== '') {
        this.filter.search += ' '
      }
      this.$refs.searchInput.$el.focus()
    },
    fetchRecords () {
      this.loading = true
      this.$axios.get('/ds/pub', {
        timeout: SERVER_TIMEOUT * 1000,
        params: {
          limit: PAGE_SIZE,
          sort: 'name',
          offset: PAGE_SIZE * (this.pagination.current - 1),
          search: this.filter.search,
          search_fields: this.filter.search_fields
        }
      })
        .then((d) => {
          if (d.data.records.length < 1) {
            this.notFound = true
            return
          }
          this.notFound = false
          this.records = []
          d.data.records.forEach((r) => {
            if (this.filter.applyed) {
              switch (this.filter.search_fields) {
                case 'name':
                  break
                case 'info':
                  break
                default:
                  break
              }
            }
            this.records.push(r)
          })
          this.pagination.max = Math.ceil(d.data.count / PAGE_SIZE)
          if (typeof this.$refs.scrollArea !== 'undefined') this.$refs.scrollArea.setScrollPosition(0)
          // this.searchPan.active = false
          this.filter.applyed = !this.searchIsEmpty
        })
        .catch((e) => {
          console.error(e)
          this.$q.notify({
            position: 'top',
            type: 'negative',
            message: 'Сервер не отвечает'
          })
        })
        .finally(() => {
          this.loading = false
        })
      this.loading = false
    }
  }
}
</script>
