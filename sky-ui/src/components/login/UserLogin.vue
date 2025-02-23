<template>
  <div class="user-login">
    <el-button text @click="emit('toggle-mode')" class="mode-toggle">
      没有账号？注册
    </el-button>
    <h2 class="title">登录</h2>

    <div class="form-wrapper">
      <!-- 登录表单滑动容器 -->
      <transition name="slide-fade" mode="out-in">
        <div v-if="!showQRCode" class="login-form">
          <el-form @submit.prevent="emit('submit', formData)">
            <el-form-item>
              <el-input v-model="formData.username" placeholder="用户名" required />
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
            <el-row :gutter="20" class="options">
              <el-col :span="12">
                <el-checkbox style="color: #66b1ff;" v-model="rememberMe">记住密码</el-checkbox>
              </el-col>
              <el-col :span="12" style="text-align: right">
                <span href="#" @click.prevent="handleForgetPassword">忘记密码？</span>
              </el-col>
            </el-row>
            <el-form-item>
              <el-button @click="goTo('/home')" type="primary" native-type="submit" style="width: 100%">
                登录
              </el-button>
            </el-form-item>
          </el-form>
        </div>
      </transition>

      <!-- 二维码滑动容器 -->
      <transition name="qrcode-slide">
        <div v-if="showQRCode" class="qrcode-container">
          <div class="qrcode-mock">
            <div class="qr-pattern"></div>
          </div>
          <p>请使用手机扫描二维码登录</p>
          <el-button text @click="showQRCode = false" class="close-btn">
            <el-icon><Close /></el-icon>
          </el-button>
        </div>
      </transition>
    </div>

    <el-row :gutter="20" class="additional-options">
      <el-col :span="12" style="text-align: left">
        <el-button
            text
            @click="toggleQRCode"
            :type="showQRCode ? 'primary' : 'text'"
            class="scan-btn"
        >
          {{ showQRCode ? '返回密码登录' : '扫码登录' }}
        </el-button>
      </el-col>
      <el-col :span="12" style="text-align: right;">
        <el-button style="color: #66b1ff;" text @click="emit('admin-login')">管理员登录</el-button>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, defineEmits } from 'vue'
import { Close } from '@element-plus/icons-vue'
import { useRouter } from 'vue-router'

const route = useRouter()
const emit = defineEmits(['toggle-mode', 'submit', 'admin-login'])

const showQRCode = ref(false)
const formData = ref({ username: 'admin', password: '123456' })
const rememberMe = ref(false)

const toggleQRCode = () => {
  showQRCode.value = !showQRCode.value
}

const handleForgetPassword = () => {
  console.log('忘记密码处理')
}


const goTo = () => {
  route.push('/home')
}
</script>

<style scoped>
.user-login {
  position: relative;
}

.user-login .title {
  color: #F0F8FF;
  margin-bottom: 12px;
  font-size: 2em;
  font-weight: 700;
}


/* 按钮颜色调整 */
/* 按钮定位 */
.mode-toggle {
  position: absolute;
  top: 10px;
  right: 0px;
  transform: translateX(0);
  margin-left: auto;
  color: #66b1ff;
  z-index: 1; /* 确保按钮在最上层 */
}

/* 悬停效果保持 */
.mode-toggle:hover {
  color: #8cc5ff;
  right: 9px; /* 微调悬停位置保持视觉对齐 */
  transition: all 0.2s ease;
}

/* 表单容器布局 */
.form-wrapper {
  position: relative;
  width: 100%;
  min-height: 180px;
  overflow: hidden;
}

.login-form {
  width: 100%;
  transition: all 0.3s ease;
}

/* 二维码容器样式 */
.qrcode-container {
  position: absolute;
  right: 0;
  top: 0;
  width: 160px;
  background: rgba(255, 255, 255, 0.9);
  border-radius: 8px;
  padding: 16px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.qrcode-mock {
  width: 128px;
  height: 128px;
  background: #fff;
  border: 2px solid #ddd;
  border-radius: 8px;
  position: relative;
  overflow: hidden;
}

.qr-pattern {
  position: absolute;
  width: 100%;
  height: 100%;
  background:
      linear-gradient(90deg, #000 10%, transparent 10%, transparent 90%, #000 90%),
      linear-gradient(0deg, #000 10%, transparent 10%, transparent 90%, #000 90%);
  animation: qr-flicker 1.5s infinite alternate;
}

.qrcode-container p {
  font-size: 13px;
  color: #606266;
}
@keyframes qr-flicker {
  from { opacity: 0.9; }
  to { opacity: 1; }
}

/* 表单滑动动画 */
.slide-fade-enter-active {
  transition: all 0.3s ease-out;
}

.slide-fade-leave-active {
  transition: all 0.3s cubic-bezier(1, 0.5, 0.8, 1);
}

.slide-fade-enter-from,
.slide-fade-leave-to {
  transform: translateX(-20px);
  opacity: 0;
}

/* 二维码滑动动画 */
.qrcode-slide-enter-active {
  transition: all 0.3s ease-out;
}

.qrcode-slide-leave-active {
  transition: all 0.3s cubic-bezier(1, 0.5, 0.8, 1);
}

.qrcode-slide-enter-from,
.qrcode-slide-leave-to {
  transform: translateX(100%);
  opacity: 0;
}

.close-btn {
  position: absolute;
  top: 8px;
  right: 8px;
  padding: 4px;
  color: #666;
}

.close-btn:hover {
  color: #409eff;
}

.options .el-col-12 {
  display: flex;
  align-items: center;
}

.options .el-col-12:last-child {
  justify-content: flex-end;
  color: #66b1ff;
  font-size: 14px;
}

.scan-btn {
  transition: all 0.3s ease;
  color: #66b1ff;
}

</style>