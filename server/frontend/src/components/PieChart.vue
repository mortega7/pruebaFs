<template>
  <Pie :chart-options="chartOptions" :chart-data="chartData" :chart-id="chartId" :dataset-id-key="datasetIdKey" :height="height" />
</template>

<script>
import { Pie } from 'vue-chartjs'
import { Chart as ChartJS, Title, Tooltip, Legend, ArcElement, CategoryScale } from 'chart.js'

ChartJS.register(Title, Tooltip, Legend, ArcElement, CategoryScale)

export default {
  name: 'PieChart',
  components: {
    Pie
  },
  props: ['typeFiles'],
  setup () {
    return {
      chartOptions: {
        responsive: true,
        maintainAspectRatio: false
      },
      chartId: 'pie-chart',
      datasetIdKey: 'label',
      height: 300
    }
  },
  computed: {
    chartData () {
      return {
        labels: this.typeFiles.types,
        datasets: [{
          backgroundColor: ['#0275D8', '#F0AD4E', '#5CB85C', '#D9534F', '#ADB5BD', '#6610F2', '#FFC107', '#41B883', '#DD1B16', '#00D8FF', '#E46651', '#0D6EFD', '#212529', '#FD7E14'],
          data: this.typeFiles.quantities
        }]
      }
    }
  }
}
</script>
