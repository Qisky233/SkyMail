<template>
  <div class="day-calendar">
    <div class="time-column">
      <div
          v-for="hour in 24"
          :key="hour"
          class="hour-block"
      >{{ hour - 1 }}:00</div>
    </div>
    <div class="main-column">
      <div
          v-for="hour in 24"
          :key="hour"
          class="hour-row"
      >
        <div class="mail-events">
          <div
              v-for="mail in getMailsForHour(hour - 1)"
              :key="mail.id"
              class="mail-event"
          >
            <div class="event-time">{{ formatTime(mail.time) }}</div>
            <div class="event-title">{{ mail.subject }}</div>
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

const getMailsForHour = (hour) => {
  return props.mailData.filter(mail => {
    const d = new Date(mail.time)
    return d.getDate() === props.modelValue.getDate() && d.getHours() === hour
  })
}

const formatTime = (isoString) => {
  return new Date(isoString).toLocaleTimeString('zh-CN', {
    hour: '2-digit',
    minute: '2-digit'
  })
}
</script>

<style scoped>
.day-calendar {
  display: flex;
  height: 100%;
}

.time-column {
  width: 60px;
  border-right: 1px solid #ebeef5;
}

.hour-block {
  height: 60px;
  font-size: 12px;
  color: #909399;
}

.main-column {
  flex: 1;
}

.hour-row {
  height: 60px;
  border-bottom: 1px solid #f0f0f0;
  position: relative;
}

.mail-event {
  background: #e6f7ff;
  border-left: 3px solid #409EFF;
  padding: 8px;
  margin: 4px;
  border-radius: 4px;
}

.event-time {
  font-size: 12px;
  color: #606266;
}

.event-title {
  font-weight: 500;
}
</style>