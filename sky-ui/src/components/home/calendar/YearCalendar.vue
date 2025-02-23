<template>
  <div class="year-calendar">
    <div
        v-for="month in 12"
        :key="month"
        class="month-block"
        @click="handleMonthClick(month)"
    >
      <div class="month-header">
        {{ month }}月
        <span class="mail-count">{{ getMonthMailCount(month) }}</span>
      </div>
      <div class="day-grid">
        <div
            v-for="d in daysInMonth"
            :key="d"
            :class="['day-cell', { 'has-mail': hasMail(month, d) }]"
        ></div>
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

const emit = defineEmits(['update:modelValue'])

const currentYear = computed(() => props.modelValue.getFullYear())

const daysInMonth = computed(() => {
  return new Date(currentYear.value, 1, 0).getDate()
})

const getMonthMailCount = (month) => {
  return props.mailData.filter(mail => {
    const d = new Date(mail.time)
    return d.getFullYear() === currentYear.value && d.getMonth() + 1 === month
  }).length
}

const handleMonthClick = (month) => {
  const newDate = new Date(currentYear.value, month - 1)
  emit('update:modelValue', newDate)
}
</script>

<style scoped lang="less">
.year-calendar {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
}

.month-block {
  background: #f8f9fc;
  border-radius: 8px;
  padding: 12px;
  cursor: pointer;
  transition: all 0.2s;

&:hover {
   transform: translateY(-2px);
   box-shadow: 0 4px 12px rgba(0,0,0,0.1);
 }
}

.month-header {
  font-weight: 500;
  color: #606266;
  margin-bottom: 8px;
}

.day-grid {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: 2px;
}

.day-cell {
  aspect-ratio: 1;
  background: #fff;
  border-radius: 2px;

&.has-mail {
   background: #e6f7ff;
 }
}
</style>