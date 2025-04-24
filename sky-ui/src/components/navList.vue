<template>
  <div>
    <ul class="navList">
      <li
          v-for="(item, index) in navItems"
          :key="index"
          class="item"
          :class="{ active: isCurrentRoute(index) }"
          @click="goToRoute(item)"
      >
        {{ item }}
      </li>
    </ul>
  </div>
</template>

<script setup>
import { computed, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';

// 定义导航项
const navItems = ['写信', '收信', '邮件检测', '数据集'];

// 定义路由路径
const navRoutes = ['/send', '/inbox', '/spam', '/data'];

// 获取当前路由
const route = useRoute();
const router = useRouter();

// 判断当前路由是否匹配
const isCurrentRoute = (index) => {
  return route.path === navRoutes[index];
};

// 跳转到对应路由
const goToRoute = (item) => {
  const index = navItems.indexOf(item);
  router.push(navRoutes[index]);
};
</script>

<style scoped>
.navList {
  padding: 12px 8px;
  list-style: none;
  background: rgba(255, 255, 255, 0.08);
  backdrop-filter: blur(8px);
  border-radius: 8px;
  margin: 8px;
  box-shadow:
      inset 0 1px 1px rgba(255, 255, 255, 0.1),
      0 2px 6px rgba(31, 38, 135, 0.05);
}

.item {
  margin: 6px 0;
  padding: 14px 20px;
  cursor: pointer;
  border-radius: 6px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  background:
      linear-gradient(90deg,
      transparent 0%,
      rgba(144, 202, 249, 0.05) 50%,
      transparent 100%
      );
  position: relative;
}

.item::before {
  content: '';
  position: absolute;
  top: 0;
  left: -4px;
  width: 3px;
  height: 100%;
  background: #8AAEFF;
  opacity: 0;
  transition: opacity 0.3s;
}

.item:hover {
  background-color: rgba(101, 167, 224, 0.15);
  transform: translateX(4px);
  box-shadow: 2px 2px 8px rgba(92, 142, 229, 0.1);
}

.item.active {
  background: linear-gradient(45deg,
  rgba(138, 174, 255, 0.15) 0%,
  rgba(182, 200, 253, 0.2) 100%
  );
  color: #5C6AC4;
  font-weight: 500;
}

.item.active::before {
  opacity: 1;
}

.item.active::after {
  content: '';
  position: absolute;
  top: 50%;
  right: 12px;
  transform: translateY(-50%);
  width: 8px;
  height: 8px;
  background: #8AAEFF;
  border-radius: 50%;
  box-shadow: 0 0 8px rgba(138, 174, 255, 0.4);
}

@media screen and (max-width: 500px) {
  .navList {
    margin: 4px;
    padding: 8px 4px;
    backdrop-filter: blur(4px);
  }

  .item {
    padding: 12px 16px;
    margin: 4px 0;
  }

  .item.active::after {
    right: 8px;
    width: 6px;
    height: 6px;
  }
}
</style>