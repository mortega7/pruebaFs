<template>
  <div class="card border-primary border-2 rounded-3">
    <div class="card-body">
      <div class="d-flex justify-content-center">
        <h3 class="text-primary"><i class="fa fa-user"></i> Clientes <span class="badge bg-primary rounded-pill mx-1">{{ (users != null) ? users.length : 0 }}</span></h3>
      </div>
      <div class="mt-3">
        <div class="row row-cols-1 row-cols-lg-2 row-cols-sm-1 g-2">
          <div class="col" v-for="user, index in users" :key="index">
            <div class="ms-3">&bull; {{ user.address }}</div>
            <span v-if="getUserChannel(user) != ''" class="ms-5">&bull; Canal: {{ getUserChannel(user) }}</span>
            <span v-else class="ms-5">&bull; Sin suscripción</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { useStore } from 'vuex'
import { computed, onMounted } from '@vue/runtime-core'

export default {
  setup () {
    const store = useStore()
    const users = computed(() => store.state.users)

    onMounted(async () => {
      await store.dispatch('getUsers').then(() => console.log('Users loaded...'))
    })

    setInterval(async () => {
      await store.dispatch('getUsers').then(() => console.log('Users updated...'))
    }, 30 * 1000)

    return {
      users
    }
  },
  methods: {
    getUserChannel (user) {
      return this.$store.getters.getUserChannel(user)
    }
  }
}
</script>
