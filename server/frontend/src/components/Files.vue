<template>
  <div class="card border-primary border-2 rounded-3" v-if="getFilesLength > 0">
    <div class="card-body">
      <div class="d-flex justify-content-center">
        <h3 class="text-primary"><i class="fa fa-file"></i> File Management/Statistics <span class="badge bg-primary rounded-pill mx-1">{{ getFilesLength }}</span></h3>
      </div>
      <div class="row">
        <div class="col-lg-5 col-sm-12 mt-3">
          <div class="card rounded-8">
            <div class="card-header">Shared</div>
            <div class="card-body">
              <p class="card-text">
                <ul class="list-group list-group-flush">
                  <li v-for="file, index in files" :key="index" class="d-flex justify-content-between align-items-start mt-2">
                    <div class="me-auto">
                      <span class="text-file-name">{{ file.name }}</span>
                    </div>
                    <div class="ms-auto mx-3">
                      <span class="badge bg-light text-dark border border-dark">{{ file.channel.name }}</span>
                    </div>
                    <button class="btn btn-sm btn-outline-primary btn-download" @click="downloadFile(file)">Download</button>
                  </li>
                </ul>
              </p>
            </div>
          </div>
        </div>
        <div class="col-lg-7 col-sm-12 mt-3">
          <div class="card rounded-8">
            <div class="card-header">Shared by Channel</div>
            <div class="card-body">
              <BarChart />
            </div>
          </div>
          <div class="card rounded-8 mt-3">
            <div class="card-header">Shared by Extension</div>
            <div class="card-body">
              <PieChart />
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { mapGetters, mapActions, mapState } from 'vuex'
import BarChart from '@/components/BarChart.vue'
import PieChart from '@/components/PieChart.vue'

export default {
  components: {
    BarChart,
    PieChart
  },
  mounted () {
    this.getFiles()
    this.interval = setInterval(this.getFiles, process.env.VUE_APP_RELOAD_TIME)
  },
  destroy () {
    clearInterval(this.interval)
  },
  data () {
    return {
      interval: null
    }
  },
  computed: {
    ...mapState('moduleFiles', ['files']),
    ...mapGetters('moduleFiles', ['getFilesLength'])
  },
  methods: {
    ...mapActions('moduleFiles', ['getFiles']),
    async downloadFile (file) {
      const a = document.createElement('a')
      a.href = 'data:' + file.type + ';base64,' + file.data
      a.download = file.name
      a.click()
    }
  }
}
</script>

<style>
span.text-file-name {
  font-size: 0.85em;
}
button.btn-download {
  font-size: .75em;
  line-height: 1.1em;
  margin-top: 0.2em;
  vertical-align: baseline;
}
</style>
