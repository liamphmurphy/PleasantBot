import axios from 'axios';

// prepare initial state for a command (see /bot/commands.go for the back-end structure)
const state = () => ({
    commands: [],
})

// getters
const getters = {
    getCommands: (state) => {
        return state.commands;
    }
}

// mutations
const mutations = {
    setCommands: (state, commands) => {
        state.commands = commands;
    }
}
// actions
const actions = {
    loadCommands: ({commit}) => {
        axios.get("http://localhost:8080/getcoms")
        .then(res => commit('setCommands', res && res.data))
    }
}

export default {
    namespaced: true,
    state,
    getters,
    actions,
    mutations
  }