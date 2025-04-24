<template>
  <div class="common-layout">
    <el-container>
      <el-header class="topNav">
        <TopNav />
      </el-header>
      <el-container>
        <el-aside class="sideNav" width="200px">
          <NavList />
        </el-aside>
        <el-main class="content">
          <Three />
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script setup>
import Three from './three.vue'
import TopNav from '@/components/topNav.vue'
import NavList from '@/components/navList.vue'
</script>

<style scoped>
.common-layout {
  height: 100vh;
  position: relative;
  background: linear-gradient(135deg, #e3f2fd 0%, #bbdefb 100%); /* 基础蓝色渐变层 */
}

.common-layout:before {
  content: "";
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background:
      linear-gradient(45deg, rgba(255,255,255,0.15) 10%, transparent 30%),
      linear-gradient(-45deg, rgba(255,255,255,0.1) 20%, transparent 40%); /* 双层斜向渐变纹理 */
  pointer-events: none;
}

.common-layout:after {
  content: "";
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  height: 40%;
  background: linear-gradient(0deg, rgba(144,202,249,0.2) 0%, transparent 60%); /* 底部浅色提升层次 */
  pointer-events: none;
}
.topNav {
    /* background-color: #409EFF; */
    height: 60px;
    /*width: 100%;*/
}

.sideNav {
    width: 200px;
    height: calc(100vh - 60px);
    /* background-color: #57a1ec; */
}

.content {
  width: calc(100vw - 200px);
  height: calc(100vh - 60px - 40px);
  background: #fff;
  margin: 20px;

  /* 新增样式 */
  border-radius: 12px;
  box-shadow:
      0 8px 32px rgba(145, 186, 238, 0.12),
      0 2px 6px rgba(200, 220, 250, 0.08);
  border: 1px solid rgba(144, 202, 249, 0.15);
  position: relative;
  overflow-y: auto;
  overflow-x: hidden;
}

/* 添加云纹装饰线 */
.content::before {
  content: "";
  position: absolute;
  top: 0;
  left: -10%;
  width: 120%;
  height: 40px;
  background: repeating-linear-gradient(
      -45deg,
      transparent,
      transparent 15px,
      rgba(144, 202, 249, 0.05) 15px,
      rgba(144, 202, 249, 0.05) 30px
  );
  animation: cloudLine 40s linear infinite;
  /*box-shadow:*/
  /*    0 12px 40px rgba(92, 142, 229, 0.15),*/
  /*    0 4px 12px rgba(200, 220, 250, 0.1);*/
}

@keyframes cloudLine {
  0% { transform: translateX(-15px) }
  100% { transform: translateX(0) }
}

@media screen and (max-width: 500px) {
    .common-layout {
      background: linear-gradient(180deg, #e3f2fd 0%, #bbdefb 100%); /* 移动端改为垂直渐变 */
    }

    .common-layout:before {
      background:
          linear-gradient(0deg, rgba(255,255,255,0.15) 10%, transparent 30%),
          linear-gradient(180deg, rgba(255,255,255,0.1) 20%, transparent 40%); /* 移动端改为垂直纹理 */
    }

    .common-layout:after {
      height: 30%;
      background: linear-gradient(0deg, rgba(144,202,249,0.15) 0%, transparent 80%); /* 减弱底部渐变 */
    }

    .sideNav {
        display: none;
    }

    .content {
      margin: 10px;
      border-radius: 8px;
      /* 移动端保持白色主背景 */
      background: #fff;
    }

    .content::before {
      height: 30px; /* 缩小装饰线高度 */
      background-size: 20px 20px; /* 加密云纹间距 */
    }

    .common-layout:before {
      background: linear-gradient(to bottom right, rgba(255, 255, 255, 0.1), rgba(255, 255, 255, 0.4));
    }
}

@media screen and (max-width: 500px) {

}
</style>