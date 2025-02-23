<template>
  <div class="enhanced-container">
    <!-- 双栏布局 -->
    <div class="enhanced-wrapper">
      <!-- 邮件详情列 -->
      <div class="enhanced-mail-list">
        <h3 class="mail-title">{{ currentDateLabel }} 邮件</h3>
        <el-scrollbar>
          <el-timeline>
            <el-timeline-item
                v-for="mail in filteredMails"
                :key="mail.id"
                :timestamp="formatTime(mail.time)"
                placement="top"
                :type="mail.type"
            >
              <el-card class="mail-card" @mouseenter="hoverMail = mail.id" @mouseleave="hoverMail = null">
                <div class="mail-header">
                  <el-tag :type="mailTypeMap[mail.type]?.color" effect="light">
                    {{ mailTypeMap[mail.type]?.label }}
                  </el-tag>
                  <h4 :class="{ 'highlight-title': hoverMail === mail.id }">{{ mail.subject }}</h4>
                </div>
                <div class="mail-body">
                  <p><el-icon><User /></el-icon> {{ mail.from }}</p>
                  <p class="content-preview">
                    <el-icon><Document /></el-icon> {{ mail.contentPreview }}
                  </p>
                </div>
              </el-card>
            </el-timeline-item>
          </el-timeline>
        </el-scrollbar>
      </div>

      <!-- 日历列 -->
      <div class="enhanced-calendar">
        <el-calendar v-model="currentDate" ref="calendarRef">
          <template #date-cell="{ data }">
            <div
                class="custom-cell"
                :class="{
                'has-mail': hasMail(data.day),
                'current-day': isToday(data.day)
              }"
            >
              <div class="day-header">
                <span class="day-number">{{ data.day.split('-')[2] }}</span>
                <span v-if="isWeekend(data.day)" class="weekend-tag">休</span>
              </div>
              <div v-if="hasMail(data.day)" class="event-indicator">
                <div class="dot"></div>
                <div class="glow"></div>
              </div>
            </div>
          </template>
        </el-calendar>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { User, Document } from '@element-plus/icons-vue'
import zhCn from 'element-plus/dist/locale/zh-cn.mjs'

// 日历实例引用
const calendarRef = ref(null)

const validateMailData = (mails) => {
  return mails.map(mail => ({
    ...mail,
    time: new Date(mail.time).toISOString() // 确保时间格式合法
  }))
}

// 模拟数据
const mailData = ref(validateMailData([
  {
    id: 1,
    time: '2025-02-23T09:30:00',
    subject: '项目进度更新',
    from: '项目经理',
    contentPreview: '本周项目进度正常推进，已完成前端界面开发...',
    type: 'work'
  },
  {
    id: 2,
    time: '2025-02-23T14:15:00',
    subject: '设计评审会议',
    from: '设计团队',
    contentPreview: '关于新版界面设计方案的讨论会议...',
    type: 'meeting'
  },
  {
    id: 3,
    time: '2025-02-24T16:00:00',
    subject: '系统维护通知',
    from: '运维部门',
    contentPreview: '计划于今晚进行系统维护，预计停机2小时...',
    type: 'notice'
  }
]))

// 响应式状态
const currentDate = ref(new Date())
const hoverMail = ref(null)

// 类型映射配置
const mailTypeMap = {
  work: { color: 'primary', label: '工作邮件' },
  meeting: { color: 'warning', label: '会议通知' },
  notice: { color: 'success', label: '系统公告' }
}

// 计算属性
const currentDateLabel = computed(() => {
  return currentDate.value.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
})

const filteredMails = computed(() => {
  return mailData.value.filter(mail =>
      new Date(mail.time).toDateString() === currentDate.value.toDateString()
  )
})

// 方法
const hasMail = (dateString) => {
  const targetDate = new Date(dateString)
  return mailData.value.some(mail => {
    const mailDate = mail.time ? new Date(mail.time) : null
    return mailDate && mailDate.toDateString() === targetDate.toDateString()
  })
}

const isToday = (dateString) => {
  return new Date(dateString).toDateString() === new Date().toDateString()
}

const isWeekend = (dateString) => {
  const day = new Date(dateString).getDay()
  return day === 0 || day === 6
}

const formatTime = (isoString) => {
  try {
    const date = new Date(isoString)
    return isNaN(date) ? '时间错误' : date.toLocaleTimeString('zh-CN', {
      hour: '2-digit',
      minute: '2-digit'
    })
  } catch {
    return '--:--'
  }
}

// 生命周期
onMounted(() => {
  calendarRef.value.locale = zhCn
})
</script>

<style scoped lang="less">
.enhanced-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
  padding: 24px;

  display: flex;
  flex-direction: column;
  overflow: hidden;

  .enhanced-wrapper {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 24px;
    height: 100%;
    max-width: 1400px;
    margin: 0 auto;

    flex: 1;
    min-height: 0;

    .enhanced-mail-list,
    .enhanced-calendar {
      height: 100%;
      overflow: hidden;
      display: flex;
      flex-direction: column;

      > *:not(.el-scrollbar) {
        flex-shrink: 0;
      }

      .el-scrollbar {
        flex: 1;
        min-height: 0;
      }
    }
  }

  .enhanced-mail-list {
    background: rgba(255,255,255,0.9);
    border-radius: 12px;
    padding: 24px;
    box-shadow: 0 8px 32px rgba(0,0,0,0.1);
    backdrop-filter: blur(10px);

    .mail-title {
      font-size: 1.5rem;
      color: #2c3e50;
      margin-bottom: 24px;
      padding-bottom: 12px;
      border-bottom: 2px solid #409EFF;
    }

    .mail-card {
      margin: 12px 0;
      transition: transform 0.3s ease;

      &:hover {
        transform: translateX(8px);
      }

      .mail-header {
        display: flex;
        align-items: center;
        gap: 12px;
        margin-bottom: 12px;

        h4 {
          margin: 0;
          transition: color 0.3s ease;
        }

        .highlight-title {
          color: #409EFF;
        }
      }

      .mail-body {
        p {
          display: flex;
          align-items: center;
          gap: 8px;
          color: #666;
          margin: 8px 0;

          .el-icon {
            font-size: 1.1em;
          }
        }
      }
    }
  }

  .enhanced-calendar {
    background: rgba(255,255,255,0.9);
    border-radius: 12px;
    padding: 24px;
    box-shadow: 0 8px 32px rgba(0,0,0,0.1);
    backdrop-filter: blur(10px);

    :deep(.el-calendar-table) {
      height: calc(100% - 60px); // 为头部预留空间
      .el-calendar__body {
        padding: 8px;
      }

      // 单元格高度控制
      td {
        padding: 2px !important;
        .custom-cell {
          height: 0; // 使用padding-bottom实现等比例高度
          padding-bottom: 75%;
          position: relative;
          overflow: visible;

          > div {
            position: absolute;
            top: 4px;
            left: 4px;
            right: 4px;
            bottom: 4px;
          }
        }
      }
    }

    :deep(.el-calendar__header) {
      background: #409EFF;
      color: white;
      padding: 16px;
      border-radius: 8px;

      &::after {
        content: "第" attr(data-week) "周";
        font-size: 0.9em;
        opacity: 0.8;
      }
    }

    .custom-cell {
      height: 80px;
      padding: 8px;
      border-radius: 6px;
      background: rgba(255,255,255,0.9);
      transition: all 0.3s ease;
      border: 1px solid #eee;

      &:hover {
        box-shadow: 0 4px 12px rgba(64,158,255,0.2);
        transform: translateY(-2px);
      }

      .day-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
      }

      .day-number {
        font-size: 1.2em;
        font-weight: 500;
      }

      .weekend-tag {
        font-size: 0.8em;
        color: #ff4d4f;
        background: rgba(255,77,79,0.1);
        padding: 2px 6px;
        border-radius: 4px;
      }

      &.current-day {
        background: #e6f7ff;
        border-color: #409EFF;
      }

      .event-indicator {
        position: absolute;
        bottom: 4px;
        right: 4px;

        &::after {
          content: "有邮件";
          position: absolute;
          bottom: 100%;
          right: 0;
          background: #ff4d4f;
          color: white;
          padding: 4px 8px;
          border-radius: 4px;
          font-size: 12px;
          white-space: nowrap;
          opacity: 0;
          transition: opacity 0.3s;
          pointer-events: none;
        }

        &:hover {
          .dot {
            transform: scale(1.5);
            box-shadow: 0 0 8px rgba(255,77,79,0.5);
          }

          &::after {
            opacity: 1;
          }
        }

        .dot {
          width: 8px;
          height: 8px;
          background: #ff4d4f;
          border-radius: 50%;
          position: relative;
          z-index: 2;
          transition: all 0.3s cubic-bezier(0.68, -0.55, 0.27, 1.55);
        }

        .glow {
          position: absolute;
          width: 16px;
          height: 16px;
          background: rgba(255,77,79,0.2);
          border-radius: 50%;
          animation: pulse 1.5s infinite;
          top: -4px;
          left: -4px;
        }
      }
    }
  }
}

@keyframes pulse {
  0% { transform: scale(0.8); opacity: 0.5; }
  50% { transform: scale(1.2); opacity: 0.2; }
  100% { transform: scale(0.8); opacity: 0.5; }
}

@media (max-width: 768px) {
  .enhanced-container {
    padding: 12px;

    .enhanced-wrapper {
      grid-template-columns: 1fr;
      gap: 16px;

      flex-direction: column;

      > * {
        width: 100% !important;
        margin-bottom: 12px;

        &:last-child {
          margin-bottom: 0;
        }
      }
    }

    .enhanced-mail-list {
      order: 2;
      padding: 16px;
    }

    .enhanced-calendar {
      order: 1;
      padding: 16px;

      padding: 8px !important;

      :deep(.el-calendar__header) {
        padding: 8px !important;
        font-size: 14px;
      }

      .custom-cell {
        .day-number {
          font-size: 1em !important;
        }

        .weekend-tag {
          display: none; // 移动端隐藏周末标记
        }
      }

      // 缩小单元格间距
      td {
        padding: 1px !important;
      }
    }
  }
}

// 优化日历头部显示
:deep(.el-calendar__header) {
  display: flex;
  justify-content: space-between;
  align-items: center;

  // 添加月份快速切换
  .el-icon {
    cursor: pointer;
    padding: 4px;
    border-radius: 4px;

    &:hover {
      background: rgba(255,255,255,0.2);
    }
  }
}
</style>