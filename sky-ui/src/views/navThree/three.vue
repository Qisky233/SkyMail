<template>
  <div class="spam-checker-container">
    <el-card shadow="hover" class="main-card">
      <!-- 标题区 -->
      <div class="header">
        <h2 class="title">
          <el-icon class="icon"><Message /></el-icon>
          垃圾邮件检测
        </h2>
        <el-tag type="info" effect="dark" round>
          准确率: 0.9749 | 召回率: 0.9902
          <el-tooltip placement="top">
            <template #content>
              <div style="max-width: 300px;">
                <div>最佳阈值结果:</div>
                <div>阈值: 5.00</div>
                <div>准确率: 0.9749</div>
                <div>精确率: 0.9727</div>
                <div>召回率: 0.9902</div>
                <div>F1分数: 0.9813</div>
                <img 
                  src="@/assets/Figure_1.png" 
                  style="width: 100%; margin-top: 8px;"
                  alt="预测结果"
                />
              </div>
            </template>
            <el-icon class="info-icon" style="margin-left: 8px; cursor: pointer;">
              <InfoFilled />
            </el-icon>
          </el-tooltip>
        </el-tag>
      </div>

      <el-divider border-style="dashed" />

      <el-alert title="训练数据基于TREC06邮件数据集，" type="warning" />
      <el-alert title="主要针对常见垃圾邮件场景，对普通邮件中的无意义内容未做专门优化。" type="warning" />
      <br>

      <!-- 检测表单区 -->
      <el-form :model="form" label-width="100px">
        <el-form-item label="邮件内容" prop="content">
          <el-input
            v-model="form.content"
            type="textarea"
            :rows="8"
            placeholder="请输入待检测的邮件内容或直接粘贴解析后的.eml文件内容"
            resize="none"
            class="content-input"
          />
        </el-form-item>

        <!-- <el-form-item label="附件分析">
          <el-upload
            action=""
            :auto-upload="false"
            :on-change="handleFileChange"
            :show-file-list="false"
          >
            <el-button type="primary" plain>
              <el-icon><Upload /></el-icon>上传邮件文件(.eml)
            </el-button>
          </el-upload>
          <div v-if="fileInfo.name" class="file-info">
            <el-icon><Document /></el-icon>
            {{ fileInfo.name }} ({{ fileInfo.size }}KB)
          </div>
        </el-form-item> -->

        <el-form-item>
          <el-button 
            type="primary" 
            @click="submitCheck" 
            :loading="isChecking"
            class="check-button"
          >
            <el-icon><Search /></el-icon>
            {{ isChecking ? '检测中...' : '开始检测' }}
          </el-button>
        </el-form-item>
      </el-form>

      <!-- 检测结果区 -->
      <!-- <div v-if="result" class="result-section">
        <el-divider content-position="left">
          <span class="result-title">检测结果</span>
        </el-divider>

        <div class="result-content">
          <el-alert 
            :title="resultTitle" 
            :type="resultType" 
            :closable="false"
            class="result-alert"
          >
            <template #default>
              <div class="result-details">
                <div class="detail-item">
                  <span class="label">垃圾邮件概率:</span>
                  <el-progress 
                    :percentage="result.probability * 100" 
                    :status="result.probability > 0.7 ? 'exception' : 'success'"
                    :stroke-width="16"
                    :text-inside="true"
                  />
                </div>
                <div class="detail-item">
                  <span class="label">检测模型:</span>
                  <el-tag>朴素贝叶斯算法</el-tag>
                </div>
                <div class="detail-item">
                  <span class="label">风险特征:</span>
                  <div class="features">
                    <el-tag 
                      v-for="(feature, index) in result.risk_features" 
                      :key="index"
                      type="danger"
                      size="small"
                    >
                      {{ feature }}
                    </el-tag>
                  </div>
                </div>
              </div>
            </template>
          </el-alert>

          <div v-if="result.suggestions" class="suggestions">
            <h4><el-icon><Warning /></el-icon> 安全建议</h4>
            <ul>
              <li v-for="(item, index) in result.suggestions" :key="index">
                {{ item }}
              </li>
            </ul>
          </div>
        </div>
      </div> -->
    </el-card>

    <!-- 历史记录区 -->
    <el-card shadow="never" class="history-card">
      <template #header>
        <div class="history-header">
          <span>检测历史</span>
          <el-button 
            link
            @click="fetchHistory"
            :loading="isLoadingHistory"
          >
            <el-icon><Refresh /></el-icon>刷新
          </el-button>
        </div>
      </template>

      <el-table 
        :data="historyData" 
        stripe 
        style="width: 100%"
        v-loading="isLoadingHistory"
      >
        <el-table-column prop="time" label="检测时间" width="180" />
        <el-table-column prop="content" label="内容摘要">
          <template #default="{row}">
            <el-tooltip :content="row.content" placement="top">
              <span class="content-preview">{{ row.content }}</span>
            </el-tooltip>
          </template>
        </el-table-column>
        <el-table-column prop="result" label="检测结果" width="120">
          <template #default="{row}">
            <el-tag :type="row.result == 'spam'  ? 'danger' : 'success'" size="small">
              {{ row.result ? '垃圾邮件' : '正常邮件' }}
            </el-tag>
          </template>
        </el-table-column>
        <!-- <el-table-column label="操作" width="80">
          <template #default="{row}">
            <el-button 
              link
              size="small" 
              @click="showDetail(row)"
            >
              详情
            </el-button>
          </template>
        </el-table-column> -->
      </el-table>
      <!-- 分页 -->
      <div class="pagination-wrapper">
        <el-pagination
          v-model:currentPage="currentPage"
          v-model:pageSize="pageSize"
          :page-sizes="[10, 20, 50, 100, 200]"
          layout="prev, pager, next, total, sizes"
          :total="totalHistory"
          @update:currentPage="handlePageChange"
          @update:pageSize="handlePageSizeChange"
        />
      </div>
    </el-card>

  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { 
  Message, 
  Search, 
  Upload, 
  Document, 
  Warning, 
  Refresh 
} from '@element-plus/icons-vue'

// 表单数据
const form = ref({
  content: ''
})

// 文件信息
const fileInfo = ref({
  name: '',
  size: 0
})

// 检测状态
const isChecking = ref(false)

// 检测结果
const spamResult = ref(null)

// 历史记录相关
const totalHistory = ref(0)
const historyData = ref([])
const isLoadingHistory = ref(false)

// 分页相关
const currentPage = ref(1)
const pageSize = ref(10)

// 计算属性
// const resultTitle = computed(() => {
//   if (!spamResult.value.result) return ''
//   return result.value.result  == 'spam'
//     ? '检测到垃圾邮件'
//     : '未检测到垃圾邮件特征'
// })

const resultType = computed(() => {
  if (result.value.result == 'ham') return 'info'
  return result.valu.result == 'spam' ? 'error' : 'success'
})

// 文件上传处理
const handleFileChange = (file) => {
  fileInfo.value = {
    name: file.name,
    size: (file.size / 1024).toFixed(1)
  }
  
  const reader = new FileReader()
  reader.onload = (e) => {
    form.value.content = e.target.result
  }
  reader.readAsText(file.raw)
}

// 提交检测
const submitCheck = async () => {
  if (!form.value.content.trim()) {
    ElMessage.warning('请输入邮件内容')
    return
  }

  isChecking.value = true
  
  try {
    // 调用检测接口
    const response = await fetch('/api/v1/spam', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        content: form.value.content
      })
    })
    
    const data = await response.json()
    
    if (data) {
      // console.log('检测结果:', data)
      spamResult.value = data
      fetchHistory() // 刷新历史记录
      isChecking.value = false
    } else {
      throw new Error(data.message || '检测失败')
    }
  } catch (error) {
    console.error('检测出错:', error)
    ElMessage.error(error.message || '检测服务异常，请稍后重试')
  } finally {
    isChecking.value = false
  }
}

// 获取历史记录
const fetchHistory = async () => {
  isLoadingHistory.value = true
  try {
    const response = await fetch(`/api/v1/spam/history?page=${currentPage.value}&limit=${pageSize.value}`)
    const data = await response.json()
    if (data) {
      historyData.value = data.data.items || []
      totalHistory.value = data.data.total || 0
      console.log('历史记录:', historyData.value)
    }
  } catch (error) {
    console.error('获取历史记录失败:', error)
  } finally {
    isLoadingHistory.value = false
  }
}

const handlePageChange = async () => {
  await fetchHistory()
}

// 分页每页大小变化时触发
const handlePageSizeChange = async () => {
  await fetchHistory()
}

// 查看详情
const showDetail = (row) => {
  form.value.content = row.content
  result.value = {
    is_spam: row.is_spam,
    probability: row.probability,
    risk_features: row.risk_features || [],
    suggestions: row.suggestions || []
  }
  
  // 滚动到结果区域
  setTimeout(() => {
    const el = document.querySelector('.result-section')
    if (el) el.scrollIntoView({ behavior: 'smooth' })
  }, 100)
}

// 初始化获取历史记录
fetchHistory()
</script>

<style scoped lang="less">
.spam-checker-container {
  max-width: 1000px;
  margin: 0 auto;
  padding: 20px;

  .header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 10px;

    .title {
      display: flex;
      align-items: center;
      margin: 0;
      color: #333;

      .icon {
        margin-right: 10px;
        font-size: 24px;
        color: var(--el-color-primary);
      }
    }
  }

  .main-card {
    margin-bottom: 20px;
    border-radius: 8px;
    border: 1px solid var(--el-border-color-light);

    :deep(.el-card__body) {
      padding: 20px;
    }

    .content-input {
      :deep(.el-textarea__inner) {
        font-family: 'Courier New', monospace;
        font-size: 14px;
      }
    }

    .file-info {
      margin-top: 8px;
      padding: 6px 10px;
      background-color: var(--el-color-info-light-9);
      border-radius: 4px;
      font-size: 13px;
      color: var(--el-color-info);

      .el-icon {
        margin-right: 5px;
      }
    }

    .check-button {
      width: 180px;
      height: 40px;
      font-size: 16px;

      .el-icon {
        margin-right: 5px;
      }
    }

    .result-section {
      margin-top: 30px;

      .result-title {
        font-size: 16px;
        font-weight: bold;
        color: #333;
      }

      .result-alert {
        margin-bottom: 20px;

        :deep(.el-alert__title) {
          font-size: 16px;
        }
      }

      .result-details {
        .detail-item {
          display: flex;
          align-items: center;
          margin-bottom: 12px;

          .label {
            width: 100px;
            font-size: 14px;
            color: #666;
          }

          .features {
            .el-tag {
              margin-right: 6px;
              margin-bottom: 6px;
            }
          }
        }
      }

      .suggestions {
        margin-top: 20px;
        padding: 15px;
        background-color: var(--el-color-warning-light-9);
        border-radius: 4px;

        h4 {
          margin-top: 0;
          margin-bottom: 10px;
          display: flex;
          align-items: center;
          color: var(--el-color-warning-dark);

          .el-icon {
            margin-right: 5px;
          }
        }

        ul {
          margin: 0;
          padding-left: 20px;
          color: var(--el-color-warning-dark);

          li {
            margin-bottom: 5px;
            font-size: 14px;
          }
        }
      }
    }
  }

  .history-card {
    border-radius: 8px;

    .pagination-wrapper {
      margin-top: 20px;
      display: flex;
      justify-content: center;
    }

    .history-header {
      display: flex;
      align-items: center;
      justify-content: space-between;
    }

    .content-preview {
      display: inline-block;
      width: 200px;
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
      font-size: 13px;
      color: #666;
    }
  }
}

@media screen and (max-width: 768px) {
  .spam-checker-container {
    padding: 10px;

    .main-card {
      :deep(.el-card__body) {
        padding: 15px;
      }
    }

    .header {
      flex-direction: column;
      align-items: flex-start;

      .title {
        margin-bottom: 10px;
      }
    }

    .result-details {
      .detail-item {
        flex-direction: column;
        align-items: flex-start !important;

        .label {
          margin-bottom: 8px;
          width: auto !important;
        }
      }
    }
  }
}
</style>