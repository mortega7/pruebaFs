import { createStore } from 'vuex'
import moduleUsers from '@/store/modules/users'
import moduleChannels from '@/store/modules/channels'
import moduleFiles from '@/store/modules/files'

export default createStore({
  modules: {
    moduleUsers,
    moduleChannels,
    moduleFiles
  }
})
