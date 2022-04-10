<template>
  <div>
    <div class="container mt-5">
      <h2>Dashboard</h2>
      <hr class="text-primary" />
      <div class="dashboard-list-top">
        <UsersComponent />
        <ChannelsComponent />
      </div>
      <div class="dashboard-list-bottom">
        <FilesComponent />
      </div>
    </div>
  </div>
</template>

<script>
import UsersComponent from '@/components/Users.vue'
import ChannelsComponent from '@/components/Channels.vue'
import FilesComponent from '@/components/Files.vue'

export default {
  components: {
    UsersComponent,
    ChannelsComponent,
    FilesComponent
  }
}
</script>

<style>
.dashboard-list-top {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 3rem;
  margin: 3rem 0;
}
@media (max-width: 999px) {
  .dashboard-list-top {
    grid-template-columns: repeat(1, 1fr);
    gap: 1rem;
    margin: 1rem 0;
  }
  .dashboard-list-top .card:not(:first-child) {
    margin-top: 2rem;
  }
}
@media (min-width: 1000px) {
  .dashboard-list-top {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 3rem;
    margin: 3rem 0;
  }
  .dashboard-list-top .card, .dashboard-list-top .card:not(:first-child) {
    margin-top: 0px;
  }
}
.dashboard-list-bottom {
  margin: 3rem 0;
}
</style>
