<template>
  <div class="sky-container">
    <!-- 搜索栏 -->
    <div class="search-bar">
      <el-input
          v-model="searchKeyword"
          placeholder="搜索邮件主题或内容"
          prefix-icon="el-icon-search"
          clearable
          @input="handleSearch"
          class="search-input"
      >
        <template #prepend>
          <el-select v-model="searchType" placeholder="搜索类型" style="width: 110px">
            <el-option label="全部" value="all" />
            <el-option label="主题" value="subject" />
            <el-option label="内容" value="body" />
          </el-select>
        </template>
      </el-input>
    </div>

    <!-- 操作栏 -->
    <div class="action-bar">
      <div class="left-actions">
        <el-button type="primary" @click="batchDownload" icon="el-icon-download">批量导出</el-button>
        <el-switch
            v-model="isDeduplication"
            active-text="智能去重"
            inactive-text="显示全部"
            active-color="#409EFF"
            @change="handleDeduplication"
        />
      </div>
      <div class="right-actions">
        <el-tag type="info">共 {{ filteredCount }} 封邮件</el-tag>
      </div>
    </div>

    <!-- 邮件表格 -->
    <div class="mail-table">
      <el-table
          :data="filteredEmails"
          @selection-change="handleSelectionChange"
          style="width: 100%"
          stripe
          v-loading="loading"
      >
        <el-table-column type="selection" width="55" align="center" />
        <el-table-column prop="date" label="日期" width="120" align="center">
          <template #default="{row}">
            <div class="date-cell">
              <span class="day">{{ formatDay(row.date) }}</span>
              <span class="time">{{ formatTime(row.date) }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="subject" label="邮件主题" min-width="200">
          <template #default="{row}">
            <div class="subject-cell">
              <el-icon :color="getPriorityColor(row.priority)" class="priority-icon">
                <component :is="getPriorityIcon(row.priority)" />
              </el-icon>
              <span class="subject-text">{{ row.subject }}</span>
              <el-tag v-if="row.attachments" size="mini" effect="dark" type="warning">
                <el-icon class="el-icon-paperclip"></el-icon>
                {{ row.attachments }}
              </el-tag>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="from" label="发件人" width="180" />
        <el-table-column label="操作" width="120" align="center">
          <template #default="{row}">
            <el-tooltip content="查看详情" placement="top">
              <el-button
                  size="mini"
                  type="text"
                  @click="viewDetail(row)"
                  icon="el-icon-view"
                  class="action-btn"
              >查看</el-button>
            </el-tooltip>
            <el-tooltip content="下载邮件" placement="top">
              <el-button
                  type="text"
                  @click="downloadEmail(row)"
                  icon="el-icon-download"
                  class="action-btn"
              >下载</el-button>
            </el-tooltip>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <!-- 分页 -->
    <div class="pagination-wrapper">
      <el-pagination
          :current-page="currentPage"
          :page-sizes="[5, 10, 20]"
          :page-size="pageSize"
          layout="prev, pager, next, total"
          :total="filteredEmails.length"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
      />
    </div>

    <!-- 邮件详情弹窗 -->
    <el-dialog
        v-model="dialogVisible"
        title="邮件详情"
        :width="dialogWidth"
        custom-class="mail-detail-dialog"
        top="5vh"
    >
      <div class="detail-content">
        <h3 class="detail-title">{{ selectedEmail.subject }}</h3>
        <div class="meta-info">
          <el-row :gutter="20">
            <el-col :span="12">
              <el-descriptions :column="1" border>
                <el-descriptions-item label="发件人">
                  <el-tag size="small">{{ selectedEmail.from }}</el-tag>
                </el-descriptions-item>
                <el-descriptions-item label="收件人">
                  <el-tag size="small" type="success">{{ selectedEmail.to }}</el-tag>
                </el-descriptions-item>
              </el-descriptions>
            </el-col>
            <el-col :span="12">
              <el-descriptions :column="1" border>
                <el-descriptions-item label="日期">
                  {{ formatFullDate(selectedEmail.date) }}
                </el-descriptions-item>
                <el-descriptions-item label="优先级">
                  <el-rate
                      v-model="selectedEmail.priority"
                      disabled
                      :colors="['#99A9BF', '#F7BA2A', '#FF9900']"
                      :max="3"
                  />
                </el-descriptions-item>
              </el-descriptions>
            </el-col>
          </el-row>
        </div>
        <el-divider content-position="left">邮件内容</el-divider>
        <div class="email-body">
          <el-scrollbar>
            <pre>{{ selectedEmail.body }}</pre>
          </el-scrollbar>
        </div>
      </div>
      <template #footer>
        <el-button @click="dialogVisible = false">关闭</el-button>
        <el-button type="primary" @click="handleReply">回复邮件</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { ArrowUp, Warning, Star } from '@element-plus/icons-vue';

// 模拟数据生成器
const generateMockData = () => {
  const mockSubjects = [
    '项目进度报告', '系统维护通知', '季度财务总结',
    '产品需求评审', '团队建设活动', '安全培训提醒'
  ];
  const mockFrom = ['项目经理', '系统运维', '财务部', '产品团队', '人力资源', 'IT支持'];

  return Array.from({length: 20}, (_, i) => ({
    id: i + 1,
    subject: `${mockSubjects[i % 6]} #${Math.floor(Math.random()*1000)}`,
    from: `${mockFrom[i % 6]} <${mockFrom[i % 6].toLowerCase()}@company.com>`,
    to: 'user@company.com',
    date: new Date(2025, 1, 23, 9 + i%8, i%60),
    body: `尊敬的同事：\n这是关于${mockSubjects[i % 6]}的详细内容...\n请及时查阅附件。\n\n此致\n敬礼`,
    priority: Math.floor(Math.random()*3) + 1,
    attachments: Math.random() > 0.7 ? `${Math.floor(Math.random()*3)+1}个附件` : null
  }));
};

// 响应式数据
const emails = ref(generateMockData());
const searchKeyword = ref('');
const searchType = ref('all');
const isDeduplication = ref(false);
const currentPage = ref(1);
const pageSize = ref(5);
const dialogVisible = ref(false);
const selectedEmail = ref({});
const dialogWidth = ref('60%');
const multipleSelection = ref([]);
const loading = ref(false);

// 计算属性
const filteredEmails = computed(() => {
  let result = emails.value.filter(email => {
    const keyword = searchKeyword.value.toLowerCase();
    if (!keyword) return true;

    const matchSubject = email.subject.toLowerCase().includes(keyword);
    const matchBody = email.body.toLowerCase().includes(keyword);

    return searchType.value === 'all'
        ? (matchSubject || matchBody)
        : (searchType.value === 'subject' ? matchSubject : matchBody);
  });

  if (isDeduplication.value) {
    const seen = new Set();
    result = result.filter(email => {
      const key = `${email.subject}_${email.from}`;
      return seen.has(key) ? false : (seen.add(key), true);
    });
  }

  return result.slice((currentPage.value - 1) * pageSize.value, currentPage.value * pageSize.value);
});

const filteredCount = computed(() => filteredEmails.value.length);

// 方法
const handleSearch = () => {
  currentPage.value = 1;
};

const handleDeduplication = () => {
  currentPage.value = 1;
};

const formatDay = (date) => {
  return date.toLocaleDateString('zh-CN', { month: '2-digit', day: '2-digit' });
};

const formatTime = (date) => {
  return date.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' });
};

const formatFullDate = (date) => {
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  });
};

const getPriorityColor = (priority) => {
  return ['#67C23A', '#E6A23C', '#F56C6C'][priority - 1];
};

const getPriorityIcon = (priority) => {
  return [ArrowUp, Warning, Star][priority - 1];
};

const viewDetail = (row) => {
  selectedEmail.value = row;
  dialogVisible.value = true;
};

const downloadEmail = (row) => {
  const dataStr = JSON.stringify(row, null, 2);
  const blob = new Blob([dataStr], { type: 'application/json' });
  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = `${row.subject}.json`;
  a.click();
  URL.revokeObjectURL(url);
};

const batchDownload = () => {
  if (multipleSelection.value.length === 0) {
    ElMessage.warning('请选择要导出的邮件');
    return;
  }
  const dataStr = JSON.stringify(multipleSelection.value, null, 2);
  const blob = new Blob([dataStr], { type: 'application/json' });
  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = `批量邮件_${new Date().toISOString()}.json`;
  a.click();
  URL.revokeObjectURL(url);
};

const handleSizeChange = (size) => {
  pageSize.value = size;
};

const handleCurrentChange = (page) => {
  currentPage.value = page;
};

// 响应式对话框宽度
const adjustDialogWidth = () => {
  dialogWidth.value = window.innerWidth < 768 ? '95%' : '60%';
};

onMounted(() => {
  window.addEventListener('resize', adjustDialogWidth);
  adjustDialogWidth();
});
</script>

<style scoped lang="less">
.sky-container {
  background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
  padding: 24px;
  border-radius: 12px;
  box-shadow: 0 8px 32px rgba(0,0,0,0.1);
}

.search-bar {
  margin-bottom: 20px;
  .search-input {
    :deep(.el-input-group__prepend) {
      background-color: #fff;
      border-right: none;
    }
  }
}

.action-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  .left-actions {
    display: flex;
    align-items: center;
    gap: 16px;
  }
}

.mail-table {
  background: rgba(255,255,255,0.9);
  border-radius: 8px;
  padding: 16px;
  box-shadow: 0 4px 16px rgba(0,0,0,0.08);

  .date-cell {
    display: flex;
    flex-direction: column;
    align-items: center;
    .day {
      font-size: 14px;
      color: #606266;
    }
    .time {
      font-size: 12px;
      color: #909399;
    }
  }

  .subject-cell {
    display: flex;
    align-items: center;
    gap: 8px;
    .priority-icon {
      font-size: 16px;
    }
    .subject-text {
      flex: 1;
      overflow: hidden;
      text-overflow: ellipsis;
    }
  }

  .action-btn {
    padding: 8px;
    font-size: 16px;
    &:hover {
      color: #409EFF;
    }
  }
}

.pagination-wrapper {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}

.mail-detail-dialog {
  .detail-content {
    .detail-title {
      margin: 0 0 20px;
      color: #303133;
    }

    .meta-info {
      margin-bottom: 20px;
    }

    .email-body {
      height: 60vh;
      pre {
        white-space: pre-wrap;
        font-family: inherit;
        line-height: 1.6;
      }
    }
  }
}

@media (max-width: 768px) {
  .sky-container {
    padding: 12px;
  }

  .action-bar {
    flex-direction: column;
    gap: 12px;
    .left-actions {
      width: 100%;
      justify-content: space-between;
    }
  }

  .mail-table {
    .el-table__body-wrapper {
      overflow-x: auto;
    }
  }
}
</style>