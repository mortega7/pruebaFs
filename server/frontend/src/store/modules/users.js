import axios from 'axios'

export default {
  namespaced: true,
  state: {
    users: []
  },
  mutations: {
    setUsers (state, payload) {
      state.users = payload
    }
  },
  actions: {
    async getUsers ({ commit }) {
      try {
        const url = process.env.VUE_APP_BACKEND_SERVER + '/user'
        const data = await axios.get(url)
        commit('setUsers', data.data)
      } catch (error) {
        commit('setUsers', [])
        console.log('Error in getUsers', error)
      }
    }
  },
  getters: {
    getUsersLength: (state) => {
      return (state.users != null) ? state.users.length : 0
    },
    getUserChannel: () => (user) => {
      return (user.channel != null) ? user.channel.name : ''
    }
  }
}
