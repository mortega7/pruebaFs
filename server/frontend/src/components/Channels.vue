<template>
  <div class="card border-primary border-2 rounded-3">
    <div class="card-body">
      <div class="d-flex justify-content-center">
        <h3 class="text-primary"><i class="fa fa-users"></i> Channels <span class="badge bg-primary rounded-pill mx-1">{{ getChannelsLength }}</span></h3>
      </div>
      <div class="mt-3">
        <div v-if="getChannelsLength > 0"  class="row row-cols-1 row-cols-lg-2 row-cols-sm-1 g-2">
          <div class="col" v-for="channel, index in channels" :key="index">
            <span class="ms-3">&bull; {{ channel.name }}</span>
          </div>
        </div>
        <div v-else>No channels created!</div>
      </div>
    </div>
  </div>
</template>

<script>
import { useStore, mapGetters } from 'vuex'
import { computed, onMounted } from '@vue/runtime-core'

export default {
  setup () {
    const store = useStore()
    const channels = computed(() => store.state.moduleChannels.channels)

    onMounted(async () => {
      await store.dispatch('moduleChannels/getChannels')
    })

    setInterval(async () => {
      await store.dispatch('moduleChannels/getChannels')
    }, process.env.VUE_APP_RELOAD_TIME * 1000)

    return {
      channels
    }
  },
  computed: {
    ...mapGetters('moduleChannels', ['getChannelsLength'])
  }
}
</script>
