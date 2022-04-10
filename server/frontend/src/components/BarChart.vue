<template>
  <Bar :chart-options="chartOptions" :chart-data="chartData" :chart-id="chartId" :dataset-id-key="datasetIdKey" :height="height" />
</template>

<script>
import { Bar } from 'vue-chartjs'
import { Chart as ChartJS, Title, Tooltip, Legend, BarElement, CategoryScale, LinearScale } from 'chart.js'

ChartJS.register(Title, Tooltip, Legend, BarElement, CategoryScale, LinearScale)

export default {
  name: 'BarChart',
  components: {
    Bar
  },
  props: ['channelFiles'],
  setup () {
    return {
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
          label: 'File shares per channel',
          backgroundColor: '#0D6EFD',
          data: this.channelFiles.quantities
        }]
      }
    }
  }
}
</script>
