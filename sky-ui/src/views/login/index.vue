<template>
  <div class="welcome">
    <h1 v-if="step === 1" :class="{'h1-fadeInLeft': step === 1}">{{ msg }}</h1>
    <h1 v-if="step === 2" :class="{'h1-fadeInRight': step === 2}">欢迎来到 EasyMail</h1>
    <LoginRegister class="loginRegister" v-loading="loading" v-if="step === 3" @submit="handleLoginRegister" />
  </div>
</template>

<script>
import { ElLoading } from 'element-plus'
import { ref, onMounted } from 'vue';
import LoginRegister from '../../components/login/loginRegister.vue';

export default {
  components: {
    LoginRegister
  },
  setup() {
    const step = ref(1);
    const msg = ref('Hello World!');

    const loading = ref(false);

    const startLoading = () => {
      loading.value = true;
      const loadingInstance = ElLoading.service({
        lock: true,
        text: '加载中...',
        background: 'rgba(0, 0, 0, 0.6)',
      });
      setTimeout(() => {
        loadingInstance.close();
        loading.value = false;
      }, 2000);
    };
    const handleNextStep = () => {
      if (step.value < 3) {
        step.value++;
      }
    };

    const handleLoginRegister = (credentials) => {
      // 处理登录/注册逻辑
      console.log('Credentials:', credentials);
    };

    onMounted(() => {
      setTimeout(handleNextStep, 3500); // 3.5秒后显示 "欢迎来到 EasyMail"
      setTimeout(startLoading, 7000); // 再过3.5秒显示登录/注册窗口
      setTimeout(handleNextStep, 9000); // 再过3秒切换到登录/注册窗口
    });

    return {
      step,
      msg,
      handleLoginRegister,


      loading,
      startLoading,
    };
  }
};
</script>

<style scoped>
.welcome {
  height: 100vh;
  width: 100%;
  background: url('../../assets/images/background/pexels-eberhardgross-730981.jpg') 50% no-repeat;
  background-size: cover;
  display: flex;
  justify-content: center;
  align-items: center;
  flex-direction: column;
  text-align: center;
}

.welcome h1 {
  font-size: 2.5em;
  color: #fff;

}

.loginRegister {

}

/* 应用fadeInLeft动画 */
.h1-fadeInLeft {
  animation: fadeInLeft 3.5s ease-in-out;
}

/* 应用fadeInRight动画 */
.h1-fadeInRight {
  animation: fadeInRight 3.5s ease-in-out;
}

/* 定义fadeInLeft动画 */
@keyframes fadeInLeft {
  from {
    opacity: 0;
    /* transform: translate3d(-50%, 0, 2); */
    transform: translateX(-20px);
  }

  to {
    opacity: 1;
    /* transform: none; */
    transform: translateX(0);
  }
}

/* 定义fadeInRight动画 */
@keyframes fadeInRight {
  from {
    opacity: 0;
    /* transform: translate3d(50%, 0, 0); */

    transform: translateX(20px);
  }

  to {
    opacity: 1;
    /* transform: none; */

    transform: translateX(0);
  }
}



@media screen and (max-width: 768px) {
  .welcome {
    background: url('../../assets/images/background/pexels-eberhardgross-1366919.jpg') 50% no-repeat;
    background-size: cover;
  }

  .loginRegister {
    width: 96%;
  }
}
</style>
