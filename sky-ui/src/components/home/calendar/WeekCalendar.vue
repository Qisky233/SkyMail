<template>
  <div class="week-calendar">
    <div class="time-axis">
      <div
          v-for="hour in 24"
          :key="hour"
          class="hour-block"
      >{{ hour - 1 }}:00</div>
    </div>
    <div class="days-container">
      <div
          v-for="(day, index) in weekDays"
          :key="index"
          class="day-column"
      >
        <div class="day-header">{{ formatDay(day) }}</div>
        <div class="time-slots">
          <div
              v-for="hour in 24"
              :key="hour"
              class="time-slot"
          >
            <div
                v-for="mail in getMailsForHour(day, hour - 1)"
                :key="mail.id"
                class="mail-event"
            >
              {{ mail.subject }}
            </div>
          </div>
        </div>
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

const weekDays = computed(() => {
  const start = new Date(props.modelValue)
  start.setDate(start.getDate() - start.getDay())
  return Array.from({ length: 7 }, (_, i) => {
    const d = new Date(start)
    d.setDate(d.getDate() + i)
    return d
  })
})

const getMailsForHour = (day, hour) => {
  return props.mailData.filter(mail => {
    const d = new Date(mail.time)
    return d.getDate() === day.getDate() && d.getHours() === hour
  })
}

const formatDay = (date) => {
  return date.toLocaleDateString('zh-CN', { weekday: 'short' })
}
</script>

<style scoped>
.week-calendar {
  display: flex;
  height: 100%;
}

.time-axis {
  width: 60px;
  border-right: 1px solid #ebeef5;
}

.hour-block {
  height: 60px;
  font-size: 12px;
  color: #909399;
}

.days-container {
  flex: 1;
  display: flex;
  overflow-x: auto;
}

.day-column {
  flex: 1;
  min-width: 150px;
  border-right: 1px solid #ebeef5;
}

.time-slot {
  height: 60px;
  border-bottom: 1px solid #f0f0f0;
  position: relative;
}

.mail-event {
  background: #e6f7ff;
  border-left: 3px solid #409EFF;
  padding: 4px;
  margin: 2px;
  border-radius: 4px;
  font-size: 12px;
}
</style>