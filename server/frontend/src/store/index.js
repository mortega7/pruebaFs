import { createStore } from 'vuex'
import axios from 'axios'

export default createStore({
  state: {
    users: [],
    channels: [],
    files: []
  },
  mutations: {
    setUsers (state, payload) {
      state.users = payload
    },
    setChannels (state, payload) {
      state.channels = payload
    },
    setFiles (state, payload) {
      state.files = payload
    }
  },
  actions: {
    async getChannels ({ commit }) {
      try {
        const url = process.env.VUE_APP_BACKEND_SERVER + '/channel'
        const data = await axios.get(url)
        commit('setChannels', data.data)
      } catch (error) {
        console.log('Error en getChannels', error)
      }
    },
    async getUsers ({ commit }) {
      try {
        const url = process.env.VUE_APP_BACKEND_SERVER + '/user'
        const data = await axios.get(url)
        commit('setUsers', data.data)
      } catch (error) {
        console.log('Error en getUsers', error)
      }
    },
    async getFiles ({ commit }) {
      try {
        const url = process.env.VUE_APP_BACKEND_SERVER + '/file'
        const data = await axios.get(url)
        commit('setFiles', data.data)
      } catch (error) {
        console.log('Error en getUsers', error)
      }
    }
  },
  getters: {
    getUserChannel: () => (user) => {
      return (user.channel != null) ? user.channel.name : ''
    },
    getFilesByChannel: (state) => {
      const files = {
        channels: [],
        quantities: []
      }
      for (const ch of state.channels) {
        let quantity = 0
        if (state.files != null) {
          for (const f of state.files) {
            if (f.channel.name === ch.name) {
              quantity = quantity + 1
            }
          }

          files.channels.push(ch.name)
          files.quantities.push(quantity)
        }
      }
      return files
    },
    getFilesByType: (state) => {
      const files = {
        types: [],
        quantities: []
      }
      if (state.files != null) {
        for (const f of state.files) {
          const fileExt = f.name.split('.').pop().toUpperCase()
          const item = files.types.find(t => t === fileExt)
          const index = files.types.indexOf(item)

          if (index >= 0) {
            files.quantities[index] = files.quantities[index] + 1
          } else {
            files.types.push(fileExt)
            files.quantities.push(1)
          }
        }
      }
      return files
    }
  },
  modules: {
  }
})
