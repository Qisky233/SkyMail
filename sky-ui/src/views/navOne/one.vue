<template>
  <div class="email-form">
    <el-form :model="emailForm" label-width="100px" style="color: #333;">
      <el-form-item label="收件人邮箱" required>
        <el-input 
          placeholder="请输入收件人邮箱" 
          v-model="emailForm.to"
          clearable
        ></el-input>
      </el-form-item>
      
      <el-form-item label="邮件主题" required>
        <el-input 
          placeholder="请输入邮件主题" 
          v-model="emailForm.subject"
          clearable
        ></el-input>
      </el-form-item>
      
      <el-divider content-position="center">邮件正文内容</el-divider>
      
      <span>
        <el-button 
          type="primary" 
          @click="sendEmail"
          :loading="isSending"
        >
          {{ isSending ? '发送中...' : '发送邮件' }}
        </el-button>
      </span>

      <Editor 
        v-model="emailForm.body" 
        :min-height="480" 
        placeholder="请输入邮件正文..."
      />
    </el-form>
  </div>
</template>

<script setup>
import Editor from './MyEditor/index.vue';
import { ref } from 'vue';
import { ElMessage } from 'element-plus';

// 表单数据
const emailForm = ref({
  to: '2192767718@qq.com',
  subject: '',
  body: ''
});

// 发送状态
const isSending = ref(false);

// 验证邮箱格式
const validateEmail = (email) => {
  const re = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  return re.test(email);
};

// 发送邮件
const sendEmail = async () => {
  try {
    // 验证表单
    if (!emailForm.value.to.trim()) {
      ElMessage.warning('收件人邮箱不能为空');
      return;
    }
    
    if (!validateEmail(emailForm.value.to)) {
      ElMessage.warning('请输入有效的邮箱地址');
      return;
    }
    
    if (!emailForm.value.subject.trim()) {
      ElMessage.warning('邮件主题不能为空');
      return;
    }
    
    if (!emailForm.value.body.trim()) {
      ElMessage.warning('邮件正文不能为空');
      return;
    }

    isSending.value = true;
    
    // 发送请求
    const response = await fetch('/api/v1/mail', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        to: emailForm.value.to,
        subject: emailForm.value.subject,
        body: emailForm.value.body
      })
    });

    const data = await response.json();
    
    if (response.ok) {
      ElMessage.success(data.message || '邮件发送成功');
      // 清空表单
      emailForm.value = {
        to: '',
        subject: '',
        body: ''
      };
    } else {
      throw new Error(data.message || '邮件发送失败');
    }
  } catch (error) {
    console.error('邮件发送失败:', error);
    ElMessage.error(error.message || '邮件发送失败，请重试');
  } finally {
    isSending.value = false;
  }
};
</script>

<style scoped lang="less">
.email-form {
  background: rgba(255, 255, 255, 0.88);
  backdrop-filter: blur(8px);
  border-radius: 12px;
  box-shadow:
      0 8px 32px rgba(144, 202, 249, 0.1),
      inset 0 1px 1px rgba(255, 255, 255, 0.15);

  :deep(.el-form) {
    .el-form-item__label {
      color: #2c3e50;
      font-weight: 500;
      letter-spacing: 0.5px;
    }

    .el-input {
      &__wrapper {
        background: rgba(231, 239, 254, 0.3);
        border-radius: 8px;
        box-shadow: 0 2px 6px rgba(144, 202, 249, 0.05);
        transition: all 0.3s ease;

        &:hover {
          box-shadow: 0 4px 12px rgba(144, 202, 249, 0.15);
        }

        &.is-focus {
          box-shadow: 0 0 0 2px rgba(92, 142, 229, 0.2);
        }
      }

      &__inner::placeholder {
        color: #7f8c8d;
        opacity: 0.8;
      }
    }
  }

  .el-divider {
    margin: 24px 0;
    border-color: rgba(144, 202, 249, 0.3);
    &__text {
      color: #5C6AC4;
      font-size: 14px;
      letter-spacing: 1px;
    }
  }

  .el-button {
    background: linear-gradient(135deg, #8AAEFF 0%, #B6C8FD 100%);
    color: white;
    border: none;
    border-radius: 6px;
    padding: 12px 32px;
    box-shadow: 0 4px 12px rgba(138, 174, 255, 0.2);
    transition: all 0.3s ease;
    margin-bottom: 6px;

    &:hover {
      transform: translateY(-2px);
      box-shadow: 0 6px 16px rgba(138, 174, 255, 0.3);
    }

    &:active {
      transform: translateY(0);
    }

    &.is-loading {
      opacity: 0.8;
    }
  }
}

@media screen and (max-width: 768px) {
  .email-form {
    border-radius: 8px;

    :deep(.el-form) {
      .el-form-item__label {
        font-size: 14px;
      }

      .el-button {
        width: 100%;
      }
    }
  }
}
</style>