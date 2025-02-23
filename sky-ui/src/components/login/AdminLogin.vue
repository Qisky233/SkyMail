<template>
  <div class="admin-login">
    <el-button text @click="emit('back-to-user')" class="mode-toggle">
      &lt;&lt; 返回
    </el-button>
    <h2>管理员登录</h2>
    <el-form @submit.prevent="emit('submit', formData)">
      <el-form-item>
        <el-input
            v-model="formData.username"
            placeholder="用户名"
        />
      </el-form-item>
      <el-form-item>
        <el-input
            v-model="formData.password"
            type="password"
            placeholder="密码"
            show-password
        />
      </el-form-item>
      <el-form-item>
        <el-button @click="goTo('/home')" type="primary" nstyle="width: 100%">
          登录
        </el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup>
import { ref, onMounted, defineEmits } from 'vue'
import { useRouter } from 'vue-router'

const route = useRouter()
const emit = defineEmits(['back-to-user', 'submit'])

const formData = ref({
  username: 'admin111',
  password: '123456'
})

const goTo = () => {
  route.push('/home')
}

</script>

<style scoped>
.mode-toggle {
  position: absolute;
  top: 10px;
  left: 10px;
  color: #409eff;
}

h2 {
  font-size: 2em;
  font-weight: 700;
  margin-bottom: 12px;
  color: aliceblue;
}
</style>