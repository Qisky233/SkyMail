<template>
  <div class="mail-list-container">
    <!-- 邮件列表表格 -->
    <transition name="slide-left">
      <div v-if="!isDetailsVisible" class="mail-table">
        <el-table :data="paginatedData" v-loading="isLoading">
          <!-- 表格列保持不变 -->
          <el-table-column prop="modTime" label="日期" width="180"></el-table-column>
          <el-table-column prop="name" label="文件名" :show-overflow-tooltip="true">
            <template #default="scope">
              <span class="ellipsis">{{ scope.row.name }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="size" label="大小" width="100"></el-table-column>
          <el-table-column label="操作" width="150">
            <template #default="scope">
              <el-button link @click="showMailDetails(scope.row.name)">详情</el-button>
            </template>
          </el-table-column>
        </el-table>

        <!-- 分页组件保持不变 -->
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
      </div>
    </transition>

    <!-- 邮件详情区域 -->
    <transition name="slide-right">
      <div v-if="isDetailsVisible" class="mail-details-container">
        <div class="mail-details-header">
          <el-button link @click="hideMailDetails">返回</el-button>
        </div>
        <div class="mail-details">
          <el-tag :type="selectedMailDetails.meta.isSpam ? 'danger' : 'success'" size="large">
            {{ selectedMailDetails.meta.isSpam ? '垃圾邮件' : '正常邮件' }}
          </el-tag>
          <el-divider content-position="center">邮件详情</el-divider>
          
          <!-- 邮件头部信息 -->
          <div class="mail-headers">
            <h4>邮件头部信息:</h4>
            <div v-for="(value, key) in selectedMailDetails.data.headers" :key="key">
              <strong>{{ key }}:</strong> {{ value }}
            </div>
          </div>
          
          <!-- 邮件正文 -->
          <div class="mail-body">
            <h4>邮件内容:</h4>
            <div v-for="(part, index) in selectedMailDetails.data.parts" :key="index">
              <div v-if="part.content_type === 'text/html'" v-html="part.body"></div>
              <pre v-else-if="part.content_type === 'text/plain'">{{ part.body }}</pre>
            </div>
          </div>
          
          <!-- 元数据 -->
          <div class="mail-meta">
            <h4>邮件信息:</h4>
            <pre>{{ JSON.stringify(selectedMailDetails.meta, null, 2) }}</pre>
          </div>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';

// 响应式数据
const allData = ref([]);
const currentPage = ref(1);
const pageSize = ref(5);
const selectedMailDetails = ref(null);
const isLoading = ref(false);
const isDetailsVisible = ref(false); // 控制详情区域是否显示

// 获取收件箱列表
const fetchInboxData = async () => {
  try {
    isLoading.value = true;
    const response = await fetch('/api/v1/inbox');
    if (!response.ok) throw new Error('获取邮件列表失败');
    const data = await response.json();
    allData.value = data.files || [];
  } catch (error) {
    console.error('获取邮件列表出错:', error);
    // 可以添加用户提示
  } finally {
    isLoading.value = false;
  }
};

// 获取邮件详情
const fetchMailDetails = async (fileName) => {
  try {
    isLoading.value = true;
    const response = await fetch(`/api/v1/detail?fileName=${encodeURIComponent(fileName)}`);
    if (!response.ok) throw new Error('获取邮件详情失败');
    const data = await response.json();
    selectedMailDetails.value = data;
  } catch (error) {
    console.error('获取邮件详情出错:', error);
    selectedMailDetails.value = null;
    // 可以添加用户提示
  } finally {
    isLoading.value = false;
  }
};

// 显示邮件详情
const showMailDetails = async (fileName) => {
  await fetchMailDetails(fileName);
  isDetailsVisible.value = true;
};

// 隐藏邮件详情
const hideMailDetails = () => {
  isDetailsVisible.value = false;
};

// 排序和分页计算属性保持不变
const sortedData = computed(() => [...allData.value].sort((a, b) => 
  new Date(b.modTime) - new Date(a.modTime)
));

const paginatedData = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value;
  return sortedData.value.slice(start, start + pageSize.value);
});

// 分页处理方法保持不变
const handleSizeChange = (newSize) => {
  pageSize.value = newSize;
  currentPage.value = 1;
};

const handleCurrentChange = (newPage) => {
  currentPage.value = newPage;
};

// 组件挂载时获取数据
onMounted(() => {
  fetchInboxData();
});
</script>

<style scoped lang="less">
/* 邮件列表容器样式 */
.mail-list-container {
  background: rgba(255, 255, 255, 0.88);
  backdrop-filter: blur(8px);
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 8px 32px rgba(144, 202, 249, 0.1);
  position: relative;
  overflow: visible;

  /* 表格深度样式 */
  :deep(.el-table) {
    height: calc(35vh + 48px);
    --el-table-header-bg-color: rgba(231, 239, 254, 0.3);
    --el-table-border-color: rgba(144, 202, 249, 0.15);

    .el-table__body-wrapper {
      overflow-y: auto; /* 表格内容区域垂直滚动 */
    }

    th.el-table__cell {
      color: #5C6AC4;
      font-weight: 500;
    }

    /* 行悬停效果 */
    tr:hover td {
      background: rgba(144, 202, 249, 0.05) !important;
    }

    /* 按钮样式 */
    .el-button--text {
      color: #8AAEFF;
      &:hover {
        color: #5C6AC4;
      }
    }
  }

  /* 分页容器样式 */
  .pagination-container {
    margin-top: 24px;

    :deep(.el-pagination) {
      --el-pagination-bg-color: transparent;
      --el-pagination-button-disabled-bg-color: transparent;

      /* 分页按钮样式 */
      .btn-prev,
      .btn-next,
      .el-pager li {
        background: rgba(231, 239, 254, 0.3);
        border-radius: 6px;
        margin: 0 4px;
        min-width: 32px;

        /* 当前页按钮样式 */
        &.active {
          background: linear-gradient(135deg, #8AAEFF 0%, #B6C8FD 100%);
          color: white;
        }
      }
    }
  }

  /* 邮件详情样式 */
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
    overflow-y: auto;
  }

  /* 分隔线样式 */
  .el-divider {
    margin: 32px 0;
    border-color: rgba(144, 202, 249, 0.3);
    &__text {
      color: #5C6AC4;
      font-size: 14px;
      letter-spacing: 1px;
    }
  }

  /* 移动端响应式样式 */
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

/* 表格深度样式 */
.mail-table {
  transition: transform 0.3s ease;
}

/* 邮件详情容器样式 */
.mail-details-container {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(247, 250, 255, 0.9);
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 8px 32px rgba(144, 202, 249, 0.1);
  transition: transform 0.3s ease;
}

/* 邮件详情头部样式 */
.mail-details-header {
  margin-bottom: 24px;
}

/* 滑动动画 */
.slide-left-enter-active, .slide-left-leave-active {
  transition: transform 0.5s cubic-bezier(0.25, 0.8, 0.25, 1);
}
.slide-left-enter-from {
  transform: translateX(-100%);
}
.slide-left-enter-to {
  transform: translateX(0);
}
.slide-left-leave-from {
  transform: translateX(0);
}
.slide-left-leave-to {
  transform: translateX(100%);
}

.slide-right-enter-active, .slide-right-leave-active {
  transition: transform 0.5s cubic-bezier(0.25, 0.8, 0.25, 1);
}
.slide-right-enter-from {
  transform: translateX(100%);
}
.slide-right-enter-to {
  transform: translateX(0);
}
.slide-right-leave-from {
  transform: translateX(0);
}
.slide-right-leave-to {
  transform: translateX(100%);
}
</style>