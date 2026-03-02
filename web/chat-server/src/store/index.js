import { createStore } from 'vuex'

export default createStore({
  state: {
    // 本地开发地址（根据实际部署环境修改）
    backendUrl: 'https://localhost:8000',
    wsUrl: 'wss://localhost:8000',
    userInfo: (sessionStorage.getItem('userInfo') && JSON.parse(sessionStorage.getItem('userInfo'))) || {},
    socket: null,
  },
  getters: {
  },
  mutations: {
    setUserInfo(state, userInfo) {
      state.userInfo = userInfo;
      sessionStorage.setItem('userInfo', JSON.stringify(userInfo));
    },
    setToken(state, token) {
      sessionStorage.setItem('token', token);
    },
    cleanUserInfo(state) {
      state.userInfo = {};
      sessionStorage.removeItem('userInfo');
      sessionStorage.removeItem('token');
    }
  },
  actions: {
  },
  modules: {
  }
})
