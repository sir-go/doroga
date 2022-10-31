<template>
  <div class="container">
    <q-img src="../assets/img/l0_title.png" contain height="200px" class="i-title"></q-img>
    <div class="i-text">
      <p>
        К 75-летнему юбилею Победы в Великой Отечественной войне в подмосковном парке "Патриот" откроют собор
        Воскресения Христова - главный храм Вооружённых Сил России.
      </p>
      <p>
        Ведущая к храму дорога, длинною в 1270 метров, из расчётов среднего шага, соответствует 1418 дням войны.
      </p>
      <p>
        Она будет размечена на каждый день, и вдоль неё на стене с использованием современных
        технологий будут фотографии всех, кто воевал в годы Великой Отечественной войны.
      </p>
      <div class="row">
        <div class="col-6">
          Вы можете помочь в пополнении базы.<br/>
          Для этого перейдите по ссылке справа и заполните форму анкеты.
        </div>
        <div class="col-6 row qr q-pl-xl items-stretch justify-center">
            <img src="../assets/img/qr_doroga.png" height="120px" />
            <div class="q-pa-lg grad">
              <div style="color: blue; font-size: 18pt">doroga.ttnet.ru</div>
              <div>Проект "Дорога Памяти"</div>
            </div>
        </div>
      </div>
    </div>
    <div class="row justify-around q-mt-lg">
      <div class="col-5 row qr q-ml-sm q-mt-xl grad items-center">
        <img src="../assets/img/qr_patriot.png" height="100" />
        <div class="row grad full-height items-center">
          <img src="../assets/img/pg_logo.svg" height="90" />
        </div>
      </div>
      <div class="col-5 row qr q-ml-sm q-mt-xl grad items-center">
        <img src="../assets/img/qr_ttnet.png" height="100" />
        <div class="row grad full-height items-center">
          <img src="../assets/img/tt_logo.svg" height="80"/>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'IndexTab'
}
</script>

<style scoped lang="stylus">
  .container
    height 59rem
    background url('../assets/img/l0_bg.jpg') center
    padding 1rem
    overflow no-content
    font-family JetBrainsMono

    .i-title
      margin .5rem 0

    .i-text
      text-align justify
      border-radius 10px
      margin 1rem
      padding 2rem
      font-size 16pt
      background rgba(255, 255, 255, 0.6)
      box-shadow 0 0 20px 3px white

    .qr
      img
        border-radius 20px 0 0 20px

      .grad
        background linear-gradient(90deg, white 10%, rgba(255,255,255,0) 100%)
        border-radius 0 20px 20px 0
</style>
