<template>
  <div class="sky-container">
    <!-- 搜索栏 -->
    <div class="search-bar">
      <span>
        <el-select v-model="searchType" placeholder="搜索类型" style="width: 110px" @change="clearSearchKeyword">
          <el-option label="文件" value="file" />
          <el-option label="标签" value="label" />
        </el-select>
      </span>
      <span v-if="searchType === 'file'">
        <el-select
          clearable
          prefix-icon="el-icon-search"
          class="search-input"
          v-model="searchKeyword"
          filterable
          remote
          placeholder="搜索数据文件"
          :remote-method="remoteMethod"
          :loading="loading"
          @change="searchEmails"
        >
          <el-option
            v-for="item in fileOptions"
            :key="item.value"
            :label="item.label"
            :value="item.value"
          />
        </el-select>
      </span>
      <span v-if="searchType === 'label'">
        <el-select clearable class="search-input" v-model="searchKeyword" placeholder="搜索数据选择标签" @change="searchEmails">
          <el-option label="spam" value="spam" />
          <el-option label="ham" value="ham" />
        </el-select>
      </span>
    </div>

    <!-- 操作栏 -->
    <div class="action-bar">
      <div class="left-actions">
        <el-button type="primary" @click="exportSelectedEmails"><el-icon><Download /></el-icon>批量导出</el-button>
        
        <!-- <el-button icon="el-icon-refresh">刷新数据</el-button> -->
      </div>
      <div class="right-actions">
      </div>
    </div>

    <!-- 邮件表格 -->
    <div class="mail-table">
      <el-table
        :data="emails"
        stripe
        style="width: 100%"
        v-loading="isLoading"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="45" />
        <el-table-column prop="file" label="文件" width="150" sortable>
          <template #default="{row}">
            {{ row.file }}
          </template>
        </el-table-column>
        <el-table-column prop="content" label="内容" min-width="200" show-overflow-tooltip>
          <template #default="{row}">
            {{ row.content }}
          </template>
        </el-table-column>
        <el-table-column prop="label" label="标签" width="150">
          <template #default="{row}">
            {{ row.label }}
          </template>
        </el-table-column>
        <!-- <el-table-column label="操作" width="120">
          <template #default="{row}">
            <el-button link @click="showDetail(row)">详情</el-button>
          </template>
        </el-table-column> -->
      </el-table>
    </div>

    <!-- 分页 -->
    <div class="pagination-wrapper">
      <el-pagination
        v-model:currentPage="currentPage"
        v-model:pageSize="pageSize"
        :page-sizes="[10, 20, 50, 100, 200]"
        layout="prev, pager, next, total, sizes"
        :total="totalEmails"
        @update:currentPage="handlePageChange"
        @update:pageSize="handlePageSizeChange"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { saveAs } from 'file-saver';

const emails = ref([])
const totalEmails = ref(0)
const searchKeyword = ref('')
const fileOptions = ref([])
const searchType = ref('file')
const currentPage = ref(1)
const pageSize = ref(10)
const isLoading = ref(false)
const loading = ref(false)
const selectedEmails = ref([]); // 存储选中的邮件

// 切换搜索类型时清空搜索关键词
const clearSearchKeyword = () => {
  searchKeyword.value = ''
  fileOptions.value = [] // 清空文件选项
}

// 远程搜索方法
const remoteMethod = async (query) => {
  if (query !== '') {
    loading.value = true
    try {
      const response = await fetch(`/api/v1/data/search?type=file&keyword=${query}&limit=600&page=1`)
      const data = await response.json()
      fileOptions.value = data.results.map(item => ({
        value: item.file,
        label: item.file
      }))
    } catch (error) {
      console.error('获取文件选项失败:', error)
      ElMessage.error('获取文件选项失败')
    } finally {
      loading.value = false
    }
  } else {
    fileOptions.value = []
  }
}

// 搜索邮件
const searchEmails = async () => {
  if (searchKeyword.value) {
    isLoading.value = true
    try {
      const response = await fetch(`/api/v1/data/search?type=${searchType.value}&keyword=${searchKeyword.value}&limit=${pageSize.value}&page=${currentPage.value}`)
      const data = await response.json()
      emails.value = data.results || []
      totalEmails.value = data.total || 0
      currentPage.value = currentPage.value // 重置当前页码
    } catch (error) {
      console.error('搜索邮件失败:', error)
      ElMessage.error('搜索邮件失败')
    } finally {
      isLoading.value = false
    }
  }
}

const handleSelectionChange = (selection) => {
  selectedEmails.value = selection;
};

const exportSelectedEmails = () => {
  if (selectedEmails.value.length === 0) {
    ElMessage.warning('请先选择要导出的邮件');
    return;
  }

  // 将选中的邮件数据转换为 JSON 字符串
  const jsonString = JSON.stringify(selectedEmails.value, null, 2);

  // 创建一个 Blob 对象
  const blob = new Blob([jsonString], { type: 'application/json;charset=utf-8;' });

  // 使用 file-saver 保存文件
  saveAs(blob, 'data.json');
};

// 获取邮件数据
const fetchEmails = async () => {
  isLoading.value = true
  try {
    const response = await fetch(`/api/v1/data?page=${currentPage.value}&limit=${pageSize.value}`)
    const data = await response.json()
    emails.value = data.results || []
    totalEmails.value = data.total || 0
  } catch (error) {
    console.error('获取邮件数据失败:', error)
    ElMessage.error('获取邮件数据失败')
  } finally {
    isLoading.value = false
  }
}

// 分页页码变化时触发
const handlePageChange = async () => {
  if (searchKeyword.value) {
    await searchEmails() // 如果有搜索关键词，则重新搜索
  } else {
    await fetchEmails() // 如果没有搜索关键词，则获取正常邮件数据
  }
}

// 分页每页大小变化时触发
const handlePageSizeChange = async () => {
  if (searchKeyword.value) {
    await searchEmails() // 如果有搜索关键词，则重新搜索
  } else {
    await fetchEmails() // 如果没有搜索关键词，则获取正常邮件数据
  }
}

const showDetail = async (email) => {
  console.log('selectedEmail:', email);

  try {
    // 构造 POST 请求
    const response = await fetch(`/api/v1/data/detail`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json', // 设置请求头，表明发送的是 JSON 数据
      },
      body: JSON.stringify({ content: email.content }), // 将 email.content 包装为 JSON 格式
    });

    // 检查响应状态
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    // 解析响应数据
    const data = await response.json();
    console.log('emailDetail:', data);
  } catch (error) {
    console.error('Failed to fetch email detail:', error);
  }
}

onMounted(() => {
  fetchEmails()
})
</script>
<style scoped lang="less">
.sky-container {
  padding: 20px;
  background-color: #f5f7fa;
  min-height: 80vh;
  
  .search-bar {
    margin-bottom: 20px;
    
    .search-input {
      max-width: 600px;
    }
  }
  
  .action-bar {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
    
    .right-actions {
      display: flex;
      gap: 10px;
    }
  }
  
  .mail-table {
    background-color: white;
    border-radius: 8px;
    padding: 10px;
    box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.05);
    
    .subject-cell {
      display: flex;
      align-items: center;
      
      .spam-icon {
        margin-right: 8px;
      }
      
      .spam-text {
        color: #F56C6C;
        font-weight: 500;
      }
    }
    
    .email-address {
      display: inline-block;
      max-width: 150px;
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }
    
    .danger-text {
      color: #F56C6C;
    }
    
    .warning-text {
      color: #E6A23C;
    }
    
    .success-text {
      color: #67C23A;
    }
  }
  
  .pagination-wrapper {
    margin-top: 20px;
    display: flex;
    justify-content: center;
  }
}

:deep(.mail-detail-dialog) {
  border-radius: 8px;
  
  .detail-content {
    .detail-title {
      margin-top: 0;
      color: #333;
    }
    
    .meta-info {
      margin-bottom: 20px;
    }
    
    .email-body {
      height: 50vh;
      background-color: #f9f9f9;
      border-radius: 4px;
      padding: 15px;
      
      pre {
        margin: 0;
        white-space: pre-wrap;
        font-family: 'Courier New', monospace;
        line-height: 1.6;
      }
      
      .spam-content {
        background-color: #fff6f6;
        padding: 10px;
        border-left: 3px solid #F56C6C;
      }
    }
    
    .attachments {
      margin-top: 20px;
      
      .attachment-list {
        display: flex;
        flex-wrap: wrap;
        gap: 10px;
      }
      
      .attachment-tag {
        cursor: pointer;
        user-select: none;
        
        &:hover {
          background-color: #f0f2f5;
        }
        
        .el-icon {
          margin-right: 5px;
        }
      }
    }
  }
}

@media screen and (max-width: 768px) {
  .sky-container {
    padding: 10px;
    
    .search-bar {
      .search-input {
        width: 100%;
      }
    }
    
    .action-bar {
      flex-direction: column;
      align-items: flex-start;
      gap: 10px;
      
      .right-actions {
        width: 100%;
        justify-content: space-between;
      }
    }
  }
  
  :deep(.mail-detail-dialog) {
    width: 95% !important;
    
    .el-row {
      flex-direction: column;
      
      .el-col {
        width: 100%;
        margin-bottom: 10px;
      }
    }
  }
}
</style>