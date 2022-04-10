<template>
  <div class="card border-primary border-2 rounded-3">
    <div class="card-body">
      <div class="d-flex justify-content-center">
        <h3 class="text-primary"><i class="fa fa-user"></i> Clients <span class="badge bg-primary rounded-pill mx-1">{{ getUsersLength }}</span></h3>
      </div>
      <div class="mt-3">
        <div v-if="getUsersLength > 0" class="row row-cols-1 row-cols-lg-2 row-cols-sm-1 g-2">
          <div class="col" v-for="user, index in users" :key="index">
            <div class="ms-3">&bull; {{ user.address }}</div>
            <span v-if="getUserChannel(user) != ''" class="ms-5">&bull; Channel: {{ getUserChannel(user) }}</span>
            <span v-else class="ms-5">&bull; No subscription</span>
          </div>
        </div>
        <div v-else>No clients connected!</div>
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
    const users = computed(() => store.state.moduleUsers.users)

    onMounted(async () => {
      await store.dispatch('moduleUsers/getUsers')
    })

    setInterval(async () => {
      await store.dispatch('moduleUsers/getUsers')
    }, process.env.VUE_APP_RELOAD_TIME * 1000)

    return {
      users
    }
  },
  computed: {
    ...mapGetters('moduleUsers', ['getUsersLength', 'getUserChannel'])
  }
}
</script>
