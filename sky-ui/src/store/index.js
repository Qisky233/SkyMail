// store/index.js
import { createStore } from 'vuex';

export default createStore({
  state: {
    sideNavIndex: 0 // 初始值为 0
  },
  mutations: {
    // 定义一个 mutation 来更新 sideNavIndex
    setSideNavIndex(state, index) {
      state.sideNavIndex = index;
    }
  },
  actions: {
    // 定义一个 action 来触发 mutation
    updateSideNavIndex({ commit }, index) {
      commit('setSideNavIndex', index);
    }
  },
  getters: {
    // 提供一个 getter 来获取当前的 sideNavIndex
    getCurrentSideNavIndex: (state) => state.sideNavIndex
  }
});