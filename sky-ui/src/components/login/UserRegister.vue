<template>
  <div class="user-register">
    <el-button text @click="emit('toggle-mode')" class="mode-toggle">
      已有账号？登录
    </el-button>
    <h2>注册</h2>
    <el-form @submit.prevent="emit('submit', formData)">
      <el-form-item>
        <el-input
            v-model="formData.username"
            placeholder="账号"
            suffix="@ganxy03.cn"
            required
        />
      </el-form-item>
      <el-form-item>
        <el-input
            v-model="formData.password"
            type="password"
            placeholder="密码"
            show-password
            required
        />
      </el-form-item>
      <el-form-item>
        <el-input
            v-model="formData.confirmPassword"
            type="password"
            placeholder="重复密码"
            show-password
            required
        />
      </el-form-item>
      <el-form-item>
        <el-input v-model="formData.phone" placeholder="手机号" required />
      </el-form-item>
      <el-form-item>
        <el-row style="width: 100%">
          <el-col :span="16" style="text-align: left">
            <el-input
                v-model="formData.verificationCode"
                placeholder="验证码"
                required
            />
          </el-col>
          <el-col :span="8" style="text-align: right;">
            <el-button
                type="primary"
                @click="sendVerificationCode"
                :disabled="sending"
                style="width: 100%"
            >
              {{ sendBtnText }}
            </el-button>
          </el-col>
        </el-row>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" native-type="submit" style="width: 100%">
          注册
        </el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'

const emit = defineEmits(['toggle-mode', 'submit'])

const formData = ref({
  username: '',
  password: '',
  confirmPassword: '',
  phone: '',
  verificationCode: ''
})

const sending = ref(false)
const countdown = ref(120)

const sendBtnText = computed(() =>
    sending.value ? `${countdown.value}s后重发` : '发送验证码'
)

const sendVerificationCode = () => {
  if (!sending.value) {
    sending.value = true
    const timer = setInterval(() => {
      countdown.value--
      if (countdown.value <= 0) {
        clearInterval(timer)
        sending.value = false
        countdown.value = 120
      }
    }, 1000)
  }
}
</script>

<style scoped>
.user-register {
  position: relative;
}
.mode-toggle {
  position: absolute;
  top: 10px;
  right: 0px;
  transform: translateX(0);
  margin-left: auto;
  color: #66b1ff;
  z-index: 1; /* 确保按钮在最上层 */
}

.mode-toggle:hover {
  color: #8cc5ff;
  right: 9px; /* 微调悬停位置保持视觉对齐 */
  transition: all 0.2s ease;
}

h2 {
  font-size: 2em;
  font-weight: 700;
  margin-bottom: 12px;
  color: aliceblue;
}

:deep(.el-form-item) {
  width: 100%;
}
</style>