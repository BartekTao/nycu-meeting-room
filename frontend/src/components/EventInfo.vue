<template>
    <div class="card text-bg-info mb-3 white-text" :style="showDivStyle">
      <div class="card-header">詳細資訊</div>
      <div class="card-body">
        <h5 class="card-title">預約人：{{ reservator }}</h5>
        <h5 class="card-title">會議名稱：{{ title }}</h5>
        <h5 class="card-title">開始時間：{{ start_time }}</h5>
        <h5 class="card-title">結束時間：{{ end_time }}</h5>
      </div>
    </div>
  </template>
  
  <script>
  export default {
    data() {
      return {
        showDivStyle: {
          display: 'none',
          position: 'absolute',
          maxWidth: '18rem',
          zIndex: 1000,
          left: '0px',
          top: '0px'
        },
        time_period: [],
        reservator: '',
        start_time: '',
        end_time: '',
        title: '',
      };
    },
    methods: {
      showDiv(data) {
        this.showDivStyle.left = data.event.pageX + 'px';
        this.showDivStyle.top = data.event.pageY + 'px';
        this.reservator = data.unit.name;
        this.title = data.unit.title
        if (data.unit.name) {
          this.showDivStyle.display = 'block';
          this.start_time = `${data.unit.startHours}:${data.unit.startMinutes.toString().padStart(2, '0')}`;
          this.end_time = `${data.unit.endHours}:${data.unit.endMinutes.toString().padStart(2, '0')}`;
        }

        // this.start_time = `${data.unit.startHours}:${data.unit.startMinutes.toString().padStart(2, '0')}`;
        // this.end_time = `${data.unit.endHours}:${data.unit.endMinutes}`;
      },
      hideDiv() {
        this.showDivStyle.display = 'none';
      }
    },
    mounted() {
      this.time_period = this.$names;
    }
  }
  </script>
  
  <style scoped>
  .text-bg-info {
    color: white;
  }
  </style>
  