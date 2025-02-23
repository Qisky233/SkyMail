<template>
  <div class="month-calendar">
    <div class="weekdays">
      <div v-for="day in weekdays" :key="day" class="weekday">{{ day }}</div>
    </div>
    <div class="days-grid">
      <div
          v-for="(date, index) in visibleDates"
          :key="index"
          :class="['day-cell', getDayClass(date)]"
          @click="handleDateClick(date)"
      >
        <div class="day-number">{{ date.getDate() }}</div>
        <div class="mail-dot" v-if="hasMail(date)"></div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  modelValue: Date,
  mailData: Array
})

const formatDate = (date) => {
  if (!(date instanceof Date)) {
    console.error('Invalid date format:', date)
    return new Date() // 返回默认日期
  }
  return date
}

const weekdays = ['日', '一', '二', '三', '四', '五', '六']

const visibleDates = computed(() => {
  const dates = []
  const current = new Date(props.modelValue)
  current.setDate(1)

  // 填充上月
  const firstDay = current.getDay()
  for(let i = firstDay; i > 0; i--) {
    const date = new Date(current)
    date.setDate(date.getDate() - i)
    dates.push(date)
  }

  // 本月
  while(current.getMonth() === props.modelValue.getMonth()) {
    dates.push(new Date(current))
    current.setDate(current.getDate() + 1)
  }

  // 填充下月
  while(dates.length % 7 !== 0) {
    dates.push(new Date(current))
    current.setDate(current.getDate() + 1)
  }

  return dates
})

const getDayClass = (date) => {
  return {
    'current-month': date.getMonth() === props.modelValue.getMonth(),
    'today': isToday(date)
  }
}

const hasMail = (date) => {
  return props.mailData.some(mail => {
    const d = new Date(mail.time)
    return d.toDateString() === date.toDateString()
  })
}

const handleDateClick = (date) => {
  emit('update:modelValue', date)
}
</script>

<style scoped lang="less">
.month-calendar {
  height: 100%;
}

.weekdays {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  margin-bottom: 8px;
  font-weight: 500;
  color: #909399;
}

.days-grid {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: 4px;
  height: calc(100% - 30px);
}

.day-cell {
  aspect-ratio: 1;
  background: #fff;
  border: 1px solid #ebeef5;
  border-radius: 4px;
  padding: 8px;
  cursor: pointer;

&:hover {
   background: #f5f7fa;
 }

&.current-month {
   background: #fff;
 }

&:not(.current-month) {
   background: #fafafa;
   color: #c0c4cc;
 }

&.today {
   border-color: #409EFF;
 }
}

.mail-dot {
  width: 6px;
  height: 6px;
  background: #409EFF;
  border-radius: 50%;
  margin-top: 4px;
}
</style>