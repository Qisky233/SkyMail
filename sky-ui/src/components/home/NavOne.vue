<template>
  <div class="email-form">
    <el-form :model="emailForm" label-width="100px" style="color: #333;">
      <el-form-item label="收件人邮箱">
        <el-input placeholder="请输入收件人邮箱" v-model="emailForm.to"></el-input>
      </el-form-item>
      <el-form-item label="主题">
        <el-input placeholder="请输入邮件主题" v-model="emailForm.subject"></el-input>
      </el-form-item>
      <el-divider content-position="center">邮件正文内容</el-divider>
      <el-button type="text" @click="sendEmail">发送邮件</el-button>

      <Editor v-model="emailForm.content" :min-height="480" />
    </el-form>
  </div>
</template>

<script setup>
import Editor from '../MyEditor/index.vue';
import axios from 'axios';
import { ref } from 'vue';
import { ElMessage } from 'element-plus'; // 确保引入 ElMessage

// 定义表单数据
const emailForm = ref({
  to: '2192767718@qq.com',
  subject: '',
  content: ''
});

// 发送邮件的方法
const sendEmail = async () => {
  // 检查邮件正文是否为空
  if(!emailForm.value.subject.trim()) {
    ElMessage({
      message: '邮件标题不能为空',
      type: 'warning',
      duration: 3000
    });
    return; // 如果为空，终止发送邮件
  } else if(!emailForm.value.to.trim()) {
    ElMessage({
      message: '邮件收件人不能为空',
      type: 'warning',
      duration: 3000
    });
    return; // 如果为空，终止发送邮件
  } else if (!emailForm.value.content.trim()) {
    ElMessage({
      message: '邮件正文不能为空',
      type: 'warning',
      duration: 3000
    });
    return; // 如果为空，终止发送邮件
  }

  try {
    console.log('Sending email to:', emailForm.value.to);
    console.log('Subject:', emailForm.value.subject);
    console.log('Content:', emailForm.value.content);

    // 使用axios发送POST请求
    const response = await axios.post('/api/mail', {
      to: emailForm.value.to,
      subject: emailForm.value.subject,
      body: emailForm.value.content
    });

    if (response.data.message === 'Mail sent successfully') {
      ElMessage({
        message: '邮件发送成功',
        type: 'success',
        duration: 3000
      });

      // 清空表单
      emailForm.value = {
        to: '',
        subject: '',
        content: ''
      };
    } else {
      ElMessage({
        message: '邮件发送失败',
        type: 'error',
        duration: 3000
      });
    }
  } catch (error) {
    console.error('Error occurred during sendEmail:', error);
    ElMessage({
      message: '邮件发送失败，请重试或联系管理员。',
      type: 'error',
      duration: 3000
    });
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
  margin: 16px;

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
    transition: transform 0.3s ease;
    margin-bottom: 6px;

    &:hover {
      transform: translateY(-2px);
      box-shadow: 0 6px 16px rgba(138, 174, 255, 0.3);
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
        padding: 12px;
      }
    }
  }
}
</style>