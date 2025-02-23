<template>
  <div class="mail-list-container">
    <el-table :data="paginatedData">

      <el-table-column prop="modTime" label="日期" width="180"></el-table-column>
      <!-- 文件名列应用了自定义样式 -->
      <el-table-column prop="name" label="文件名" :show-overflow-tooltip="true">
        <template #default="scope">
          <span class="ellipsis">{{ scope.row.name }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="size" label="大小" width="100"></el-table-column>
      <el-table-column label="操作" width="150">
        <template #default="scope">
          <el-button type="text" @click="fetchMailDetails(scope.row.name)">详情</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
        :current-page="currentPage"
        :page-sizes="[5, 10, 20, 30, 40]"
        :page-size="pageSize"
        layout="prev, pager, next, total"
        :total="allData.length"
        class="pagination-container"
    />

    <el-divider content-position="center">邮件详情</el-divider>
    <pre class="mail-details" v-if="selectedMailDetails">{{ selectedMailDetails }}</pre>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue';

// 模拟12条收件箱数据
const mockInboxData = Array.from({ length: 22 }, (_, i) => ({
  modTime: `2024-05-${String(i+1).padStart(2, '0')} 09:${String(i*5).padStart(2, '0')}`,
  name: `邮件_${String(i+1).padStart(2, '0')}.eml`,
  size: `${(Math.random()*512 + 128).toFixed(0)}KB`
}));

// 模拟邮件详情数据
const mockMailDetails = {
  '邮件_01.eml': {
    from: 'service@skymail.com',
    to: 'user@example.com',
    subject: '欢迎使用SkyMail',
    content: '感谢您注册我们的服务...'
  },
  '邮件_02.eml': {
    from: 'notify@cloud.com',
    to: 'user@example.com',
    subject: '系统升级通知',
    content: '我们将于今晚进行系统维护...'
  },
};

const allData = ref(mockInboxData);
const currentPage = ref(1);
const pageSize = ref(5);
const selectedMailDetails = ref(null);

// 排序和分页计算属性保持不变
const sortedData = computed(() => [...allData.value].sort((a, b) =>
    new Date(b.modTime) - new Date(a.modTime)
));

const paginatedData = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value;
  return sortedData.value.slice(start, start + pageSize.value);
});

// 修改获取详情的方法
const fetchMailDetails = (fileName) => {
  selectedMailDetails.value = mockMailDetails[fileName] || '详情信息未找到';
};

// 分页处理方法保持不变
const handleSizeChange = (newSize) => {
  pageSize.value = newSize;
  currentPage.value = 1;
};

const handleCurrentChange = (newPage) => {
  currentPage.value = newPage;
};
</script>

<style scoped lang="less">
.mail-list-container {
  background: rgba(255, 255, 255, 0.88);
  backdrop-filter: blur(8px);
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 8px 32px rgba(144, 202, 249, 0.1);

  :deep(.el-table) {
    height: calc(35vh + 48px);

    --el-table-header-bg-color: rgba(231, 239, 254, 0.3);
    --el-table-border-color: rgba(144, 202, 249, 0.15);

    .el-table__body-wrapper {
      overflow-y: auto;
    }
    th.el-table__cell {
      color: #5C6AC4;
      font-weight: 500;
    }

    tr:hover td {
      background: rgba(144, 202, 249, 0.05) !important;
    }

    .el-button--text {
      color: #8AAEFF;
      &:hover {
        color: #5C6AC4;
      }
    }
  }

  .pagination-container {
    margin-top: 24px;

    :deep(.el-pagination) {
      --el-pagination-bg-color: transparent;
      --el-pagination-button-disabled-bg-color: transparent;

      .btn-prev,
      .btn-next,
      .el-pager li {
        background: rgba(231, 239, 254, 0.3);
        border-radius: 6px;
        margin: 0 4px;
        min-width: 32px;

        &.active {
          background: linear-gradient(135deg, #8AAEFF 0%, #B6C8FD 100%);
          color: white;
        }
      }
    }
  }

  .mail-details {
    margin-top: 32px;
    padding: 20px;
    background: rgba(247, 250, 255, 0.9);
    border-radius: 8px;
    border: 1px solid rgba(144, 202, 249, 0.2);
    font-family: Monaco, Consolas, monospace;
    font-size: 14px;
    line-height: 1.6;
    white-space: pre-wrap;
    max-height: 400px;
    overflow-y: auto;
  }

  .el-divider {
    margin: 32px 0;
    border-color: rgba(144, 202, 249, 0.3);
    &__text {
      color: #5C6AC4;
      font-size: 14px;
      letter-spacing: 1px;
    }
  }

  @media screen and (max-width: 768px) {
    padding: 16px;
    border-radius: 8px;

    :deep(.el-table) {
      font-size: 14px;

      th.el-table__cell {
        font-size: 13px;
      }
    }

    .pagination-container {
      :deep(.el-pagination) {
        flex-wrap: wrap;

        .el-pagination__total,
        .el-pagination__sizes {
          flex: 1 1 100%;
          margin-bottom: 12px;
        }

        .btn-prev,
        .btn-next,
        .el-pager li {
          min-width: 28px;
          height: 28px;
        }
      }
    }

    .mail-details {
      padding: 12px;
      font-size: 13px;
    }
  }
}
</style>