import axios from 'axios'

export default {
  namespaced: true,
  state: {
    files: []
  },
  mutations: {
    setFiles (state, payload) {
      state.files = payload
    }
  },
  actions: {
    async getFiles ({ commit }) {
      try {
        const url = process.env.VUE_APP_BACKEND_SERVER + '/file'
        const data = await axios.get(url)
        commit('setFiles', data.data)
      } catch (error) {
        commit('setFiles', [])
        console.log('Error in getUsers', error)
      }
    }
  },
  getters: {
    getFilesLength: (state) => {
      return (state.files != null) ? state.files.length : 0
    },
    getFilesByChannel: (state, getters, rootState) => {
      const files = {
        channels: [],
        quantities: []
      }
      if (getters.getFilesLength > 0) {
        for (const ch of rootState.moduleChannels.channels) {
          const quantity = state.files.filter(fl => fl.channel.name === ch.name).length
          files.channels.push(ch.name)
          files.quantities.push(quantity)
        }
      }
      return files
    },
    getFilesByType: (state, getters) => {
      const files = {
        types: [],
        quantities: []
      }
      if (getters.getFilesLength > 0) {
        for (const fl of state.files) {
          const fileExt = fl.name.split('.').pop().toUpperCase()
          const item = files.types.find(tp => tp === fileExt)
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
  }
}
