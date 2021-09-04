import Vue from 'vue'
import Vuex from 'vuex'

import commands from './modules/commands';

Vue.use(Vuex)

export default new Vuex.Store({
  modules: {
    commands,
  }
})
