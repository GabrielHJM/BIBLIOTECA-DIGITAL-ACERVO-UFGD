// import Vue from 'vue';
// import Vuetify from 'vuetify/lib';

// Vue.use(Vuetify);

// export default new Vuetify({});

// import '@mdi/font/css/materialdesignicons.css'
// import 'vuetify/styles'

// const vuetify = createVuetify({
//   ssr: true,
// })

import { createVuetify } from 'vuetify'
import { aliases, mdi } from 'vuetify/iconsets/mdi'

export default createVuetify({
  icons: {
    defaultSet: 'mdi',
    aliases,
    sets: {
      mdi,
    },
  },
})
