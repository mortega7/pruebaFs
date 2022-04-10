import axios from 'axios'

export default {
  namespaced: true,
  state: {
    channels: []
  },
  mutations: {
    setChannels (state, payload) {
      state.channels = payload
    }
  },
  actions: {
    async getChannels ({ commit }) {
      try {
        const url = process.env.VUE_APP_BACKEND_SERVER + '/channel'
        const data = await axios.get(url)
        commit('setChannels', data.data)
      } catch (error) {
        commit('setChannels', [])
        console.log('Error in getChannels', error)
      }
    }
  },
  getters: {
    getChannelsLength: (state) => {
      return (state.channels != null) ? state.channels.length : 0
    }
  }
}
