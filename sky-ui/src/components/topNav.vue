<template>
  <div class="topNav">
    <!-- 左侧品牌区 -->
    <div class="brand-area">
      <img class="logo" src="@/assets//logo.png" alt="SkyMail">
      <div class="user-info">
        <div class="username">root<span class="email">&lt;root@ganxy03.cn&gt;</span></div>
        <div class="quick-links">
          <a v-for="link in leftLinks" :key="link.text" :href="link.url">{{ link.text }}</a>
        </div>
      </div>
    </div>

    <!-- 中央搜索区 -->
    <div class="search-area">
      <el-input
          placeholder="在云端寻找邮件..."
          class="search-input"
          :class="{ 'mobile-search': isMobile }"
      >
        <template #prefix>
          <el-icon class="search-icon"><Search /></el-icon>
        </template>
      </el-input>
    </div>

    <!-- 右侧操作区 -->
    <div class="action-area">
<!--      <el-switch-->
<!--          v-model="isLight"-->
<!--          :active-icon="Sun"-->
<!--          :inactive-icon="Moon"-->
<!--          class="theme-switch"-->
<!--          @change="toggleTheme"-->
<!--      />-->
      <div class="action-links">
        <a v-for="link in rightLinks" :key="link.text" :href="link.url">{{ link.text }}</a>
      </div>
      <el-icon class="menu-toggle" @click="toggleDrawer"><Fold /></el-icon>
    </div>

    <!-- 移动端抽屉 -->
    <el-drawer
        v-model="drawerVisible"
        :with-header="false"
        direction="rtl"
        size="75vw"
        class="mobile-drawer"
    >
      <NavList />
    </el-drawer>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue';
import NavList from '@/components/navList.vue';

const isLight = ref(false);
const drawerVisible = ref(false);
const isMobile = computed(() => window.innerWidth < 768);

const leftLinks = [
  { text: '邮箱首页', url: '#' },
  { text: '设置', url: '#' },
  { text: '换肤', url: '#' }
];

const rightLinks = [
  { text: '帮助中心', url: '#' },
  { text: '退出', url: '#' }
];

const toggleDrawer = () => {
  drawerVisible.value = !drawerVisible.value;
};

const toggleTheme = () => {
  document.documentElement.style.setProperty('--theme-hue', isLight.value ? '240deg' : '210deg');
};
</script>

<style scoped lang="less">
.topNav {
  height: 64px;
  padding: 0 24px;
  display: flex;
  align-items: center;
  width: calc(100% - 48px);
  justify-content: space-between;
  //background: rgba(255, 255, 255, 0.88);
  //backdrop-filter: blur(12px) saturate(180%);
  border-bottom: 1px solid rgba(144, 202, 249, 0.15);
  box-shadow: 0 4px 20px rgba(31, 38, 135, 0.05);
  position: relative;
  z-index: 1000;

  .brand-area {
    display: flex;
    align-items: center;
    gap: 16px;

    .logo {
      height: 40px;
      width: auto;
      transition: transform 0.3s ease;
      &:hover {
        transform: rotate(-5deg) scale(1.05);
      }
      scale: 1.3;
    }

    .user-info {
      .username {
        font-weight: 600;
        color: #1e5494;
        font-size: 15px;
        .email {
          font-size: 12px;
          opacity: 0.8;
          margin-left: 8px;
        }
      }

      .quick-links {
        margin-top: 4px;
        a {
          color: #5C6AC4;
          font-size: 12px;
          margin-right: 12px;
          position: relative;
          &::after {
            content: '';
            position: absolute;
            bottom: -2px;
            left: 0;
            width: 0;
            height: 1px;
            background: currentColor;
            transition: width 0.3s ease;
          }
          &:hover::after {
            width: 100%;
          }
        }
      }
    }
  }

  .search-area {
    .search-input {
      :deep(.el-input__wrapper) {
        background: rgba(255, 255, 255, 0.6) !important; /* 半透明白色背景 */
        border: 1px solid rgba(144, 202, 249, 0.2); /* 添加淡蓝色边框 */
        border-radius: 20px;
        box-shadow:
            0 2px 8px rgba(144, 202, 249, 0.1),
            inset 0 1px 2px rgba(255, 255, 255, 0.2); /* 添加内发光 */
        transition: all 0.3s ease;

        &:hover {
          background: rgba(255, 255, 255, 0.8) !important;
          box-shadow:
              0 4px 12px rgba(144, 202, 249, 0.2),
              inset 0 2px 3px rgba(255, 255, 255, 0.3);
        }

        /* 输入文字颜色调整 */
        .el-input__inner {
          color: #2c3e50;
          &::placeholder {
            color: #7f8c8d;
          }
        }
      }

      .search-icon {
        color: #5C6AC4 !important; /* 调整图标颜色增强对比 */
        filter: drop-shadow(0 1px 1px rgba(144, 202, 249, 0.3));
      }
    }
  }


  .action-area {
    display: flex;
    align-items: center;
    gap: 24px;
    .theme-switch {
      --el-switch-border-color: #B6C8FD;
      --el-switch-on-color: rgba(138, 174, 255, 0.2);
    }
    .action-links {
      a {
        color: #1753b5;
        font-size: 13px;
        margin-left: 16px;
        position: relative;
        &::before {
          content: '';
          position: absolute;
          top: 50%;
          left: -8px;
          width: 4px;
          height: 4px;
          background: currentColor;
          border-radius: 50%;
          opacity: 0;
          transition: opacity 0.3s ease;
        }
        &:hover::before {
          opacity: 1;
        }
      }
    }
    .menu-toggle {
      display: none;
      font-size: 24px;
      color: #1e5494;
      cursor: pointer;
      transition: transform 0.3s ease;
      &:hover {
        transform: rotate(180deg);
      }
    }
  }

  @media screen and (max-width: 768px) {
    .topNav {
      height: 56px;

      .brand-area {
        .logo {
          height: 32px;
        }
        .user-info .quick-links {
          display: none;
        }
      }

      .search-area {
        display: none;
        &.mobile-search {
          display: block;
          position: absolute;
          top: 100%;
          left: 0;
          right: 0;
          padding: 8px;
          background: inherit;
          backdrop-filter: inherit;
        }
      }

      .action-area {
        gap: 12px;
        .action-links {
          display: none;
        }
        .menu-toggle {
          display: block;
        }
      }
    }
  }

}
</style>