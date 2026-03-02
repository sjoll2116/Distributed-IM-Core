import axios from 'axios'
import store from '@/store/index.js'

// 创建 Axios 实例
const http = axios.create({
  baseURL: store.state.backendUrl,
  timeout: 15000,
})

// 请求拦截器：自动附加 JWT Token 到 Authorization Header
http.interceptors.request.use(
  (config) => {
    const token = sessionStorage.getItem('token')
    if (token) {
      config.headers['Authorization'] = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器：统一处理 401 未授权（Token 过期/无效）
http.interceptors.response.use(
  (response) => {
    return response
  },
  (error) => {
    if (error.response && error.response.status === 401) {
      // Token 失效，清除登录状态并跳转到登录页
      sessionStorage.removeItem('token')
      sessionStorage.removeItem('userInfo')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

export default http
