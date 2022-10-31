<!--suppress HtmlUnknownTag -->
<template>
  <div id="app">
    <section  id="s0" ref="s0">
      <div id="s_title" ref="s_title">
        <a href="#s1" v-smooth-scroll="{ duration: 2000}">
          <img id="l0Title" ref="l0Title" src="./assets/img/l0_title.png" alt=""
               :style="{top: l0titleTop, display: l0titleDisplay ? 'block' : 'none'}">
          <div :class="[showArrow ? '' : 'hidden', 'arrow']">
              <span></span>
          </div>
        </a>
      </div>
    </section>

    <div id="sky_bg">
      <section id="s1" ref="s1">
        <div id="s_text">
          <img id="planks" src="./assets/img/planks.png" alt="">
          <a href="#patriot" v-smooth-scroll="{ duration: 2000}">
            <img id="pglink" src="./assets/img/pg_logo.png" alt="">
          </a>

          <div class="tb" ref="tb0" id="tb0">
            <p>
              В 2020 г., к 75-летнему юбилею Победы в Великой Отечественной войне в
              подмосковном парке "Патриот" собираются открыть собор Воскресения Христова -
              главный храм Вооружённых Сил России.
            </p>
            <p>
              Ведущая к храму дорога, длинною в 1270 метров, из расчётов среднего шага,
              соответствует 1418 дням войны.
            </p>
            <p>
              Она будет размечена на каждый день, и вдоль неё на стене с использованием
              современных технологий будут фотографии всех, кто воевал
              в годы Великой Отечественной войны.
            </p>
          </div>

          <div class="tb" id="tb1">
            <p>
              Вы можете помочь в пополнении базы участников Великой Отечественной войны.
            </p>
            <p>
              Фотография вашего родственника появится на интерактивной панели в храме
              Вооружённых Сил Российской Федерации в честь Вознесения Христова.
            </p>
          </div>
        </div>
      </section>

      <div id="s_belt">
        <div id="phPlace" @click="$refs.phFile.click()" :class="busy ? 'disabled' : ''">
          <span v-if="phData === null && !busy">Загрузить фото *</span>
          <div id="photoImg" v-else :style="{ backgroundImage: 'url(' + phData + ')' }"></div>
        </div>
      </div>
    </div>

    <section id="s2" ref="s2">
      <div id="s_form">
        <form ref="form" id="af" @submit.prevent="" accept-charset="UTF-8">
          <div id="f_container" :class="busy ? 'disabled' : ''">
            <label>
              <input
                name="photo"
                ref="phFile"
                type="file"
                hidden="hidden"
                @change="fileSelected"
                :accept="allowedPhotoFormats"
                :disabled="busy"
              ></label>
            <div class="cpt">Укажите известные вам данные</div>
            <label v-for="field in fields" :key="field.name"><input
              class="field"
              type="text"
              :name="field.name"
              :placeholder="isEdge ? field.placeholder + (field.req ? ' *' : '') : ' '"
              maxlength="4000"
              @keyup="ev => {field.req ? checkValid(ev):null}"
              @change="ev => {field.req ? checkValid(ev):null}"
              :required="field.req"
              :disabled="busy"
            ><span v-if="!isEdge">{{field.placeholder}}{{field.req ? ' *' : ''}}</span></label>
            <label class="textarea">
              <textarea
                  class="field"
                  name="info"
                  :placeholder="isEdge ? 'доп. сведения' : ' '"
                  maxlength="100000"
                  rows="4"
                  v-bind:disabled="busy"
                ></textarea>
              <span v-if="!isEdge">доп. сведения</span>
            </label>
            <div class="cpt">Контактные данные</div>
            <label>
              <input
                class="field"
                type="text"
                name="sender_name"
                :placeholder="isEdge ? 'Ваши Ф.И.О. *' : ' '"
                maxlength="250"
                @keyup="checkValid"
                @change="checkValid"
                required="required"
                :disabled="busy"
              >
              <span v-if="!isEdge">Ваши Ф.И.О. *</span>
            </label>
            <label><input
              class="field"
              type="tel"
              name="phone"
              :placeholder="isEdge ? 'Ваш телефон *' : ' '"
              maxlength="50"
              @keyup="checkValid"
              @change="checkValid"
              required="required"
              :disabled="busy"
            >
              <span v-if="!isEdge">Ваш телефон *</span>
            </label>
            <div style="text-align: right; padding-right: 3rem">* обязательные для заполнения поля</div>
          </div>
          <label id="chBox"><input type="checkbox" checked="checked" disabled="disabled">
            соглашаюсь на обработку персональных данных
          </label>
          <div id="submitControl">
            <vue-simple-spinner v-if="busy"></vue-simple-spinner>
            <button v-else @click="submit" :disabled="!isFilled || busy">
              ОТПРАВИТЬ
            </button>
          </div>
        </form>
      </div>

      <div id="crt">
        <div class="snlinks">
          <span>Мы в соцсетях: </span>
          <div>
            <a href="https://vk.com/theroadofmemory" target="_blank">
              <img style="height: 4.5rem; max-width: 83%;" src="./assets/img/vk-logo.svg" alt="ТелеТайм в VK">
            </a>
            <a href="https://ok.ru/group/60370790645814" target="_blank">
              <img style="height: 4.5rem; max-width: 83%;" src="./assets/img/ok-logo.svg" alt="ТелеТайм в OK">
            </a>
          </div>
        </div>
        <a href="https://ttnet.ru/patriot" target="_blank">
          <div class="lgs">
            <img style="height: 4.5rem; max-width: 83%;" src="./assets/img/tt_logo.svg" alt="ТелеТайм">
            <img style="height: 5rem" src="./assets/img/pg_logo.svg" alt="Патриот Города">
          </div>
          <div class="txt">
            Портал создан компанией <span style="white-space: nowrap">ООО «ТелеТайм»</span>
            <br>
            по программе <span style="white-space: nowrap">«Патриот Города»</span>
          </div>
        </a>
      </div>
      <div id="patriot"></div>
    </section>

  </div>
</template>

<script src="./app.js"></script>
<style src="./app.css"></style>
