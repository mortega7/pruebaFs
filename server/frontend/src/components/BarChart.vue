<template>
  <Bar :chart-options="chartOptions" :chart-data="chartData" :chart-id="chartId" :dataset-id-key="datasetIdKey" :height="height" />
</template>

<script>
import { useStore } from 'vuex'
import { onMounted, ref } from '@vue/runtime-core'
import { Bar } from 'vue-chartjs'
import { Chart as ChartJS, Title, Tooltip, Legend, BarElement, CategoryScale, LinearScale } from 'chart.js'

ChartJS.register(Title, Tooltip, Legend, BarElement, CategoryScale, LinearScale)

export default {
  name: 'BarChart',
  components: {
    Bar
  },
  setup () {
    const store = useStore()
    const channelFiles = ref([])

    onMounted(() => {
      setTimeout(() => {
        channelFiles.value = store.getters.getFilesByChannel
      }, 1 * 1000)
    })

    setInterval(() => {
      channelFiles.value = store.getters.getFilesByChannel
    }, process.env.VUE_APP_RELOAD_TIME * 1000)

    return {
      channelFiles,
      chartOptions: {
        responsive: true,
        maintainAspectRatio: false
      },
      chartId: 'bar-chart',
      datasetIdKey: 'label',
      height: 300
    }
  },
  computed: {
    chartData () {
      return {
        labels: this.channelFiles.channels,
        datasets: [{
          barThickness: 30,
          minBarThickness: 15,
          label: 'Archivos compartidos por canal',
          backgroundColor: '#0D6EFD',
          data: this.channelFiles.quantities
        }]
      }
    }
  }
}
</script>
