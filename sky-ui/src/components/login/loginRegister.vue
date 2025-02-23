<template>
  <div class="loginRegister">
    <transition
        enter-active-class="animate__animated animate__fadeInUp"
        mode="out-in"
    >
      <el-card class="form-container">
        <component
            :is="currentComponent"
            :key="componentKey"
            @toggle-mode="toggleMode"
            @admin-login="showAdminLogin"
            @back-to-user="showUserLogin"
            @submit="handleSubmit"
        />
      </el-card>
    </transition>
  </div>
</template>

<script setup>
import { ref, computed, defineEmits } from 'vue'
import UserLogin from './UserLogin.vue'
import UserRegister from './UserRegister.vue'
import AdminLogin from './AdminLogin.vue'

const emit = defineEmits(['back-to-user', 'submit'])

const isAdminMode = ref(false)
const isRegisterMode = ref(false)

const componentKey = computed(() =>
    `${isAdminMode.value ? 'admin' : 'user'}-${isRegisterMode.value ? 'register' : 'login'}`
)

const currentComponent = computed(() => {
  if (isAdminMode.value) return AdminLogin
  return isRegisterMode.value ? UserRegister : UserLogin
})

const toggleMode = () => {
  isRegisterMode.value = !isRegisterMode.value
  isAdminMode.value = false
}

const showAdminLogin = () => {
  isAdminMode.value = true
  isRegisterMode.value = false
}

const showUserLogin = () => {
  isAdminMode.value = false
}

const handleSubmit = (data) => {
  emit('submit', data)
}

</script>

<style scoped>
.loginRegister {
  position: relative;
  height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
}

.form-container {
  background: rgba(255, 255, 255, 0.2);
  backdrop-filter: blur(2px);
  -webkit-backdrop-filter: blur(2px);
  padding: 20px 48px;
  border-radius: 10px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
  width: 480px;
  min-height: 320px;
  text-align: center;
  overflow: auto;
  border: none;
  position: relative;
  z-index: 0;
}

.form-container::before {
  background: url('../../assets/images/background/pexels-ifreestock-619419.jpg') 50% no-repeat;
  background-size: cover;
  content: '';
  position: absolute;
  inset: 0;
  z-index: -1;
  opacity: 0.5;
}

.animate__animated {
  animation-duration: 0.5s;
  animation-fill-mode: both;
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@media screen and (max-width: 768px) {
  .form-container {
    padding: 20px 8px;
  }
}

</style>