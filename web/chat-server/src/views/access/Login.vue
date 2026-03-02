<template>
  <div class="login-wrap">
    <div class="login-card">
      <div class="login-brand">
        <div class="brand-icon">💬</div>
        <h1 class="brand-title">KamaChat</h1>
        <p class="brand-subtitle">安全 · 即时 · 分布式</p>
      </div>
      <el-form
        :model="loginData"
        class="login-form"
      >
        <el-form-item>
          <el-input
            v-model="loginData.telephone"
            placeholder="手机号码"
            prefix-icon="User"
            size="large"
            class="login-input"
          />
        </el-form-item>
        <el-form-item>
          <el-input
            type="password"
            v-model="loginData.password"
            placeholder="密码"
            prefix-icon="Lock"
            size="large"
            show-password
            class="login-input"
          />
        </el-form-item>
        <el-button
          class="login-btn"
          size="large"
          @click="handleLogin"
          :loading="loading"
        >登 录</el-button>
      </el-form>
      <div class="login-footer">
        <button class="link-btn" @click="handleRegister">注册账号</button>
        <span class="divider">|</span>
        <button class="link-btn" @click="handleSmsLogin">验证码登录</button>
      </div>
    </div>
  </div>
</template>

<script>
import { reactive, toRefs, ref } from "vue";
import http from "@/utils/axios";
import { useRouter } from "vue-router";
import { ElMessage } from "element-plus";
import { useStore } from "vuex";
export default {
  name: "Login",
  setup() {
    const data = reactive({
      loginData: {
        telephone: "",
        password: "",
      },
    });
    const loading = ref(false);
    const router = useRouter();
    const store = useStore();
    const handleLogin = async () => {
      try {
        if (!data.loginData.telephone || !data.loginData.password) {
          ElMessage.error("请填写完整登录信息。");
          return;
        }
        if (!checkTelephoneValid()) {
          ElMessage.error("请输入有效的手机号码。");
          return;
        }
        loading.value = true;
        const response = await http.post("/login", data.loginData);
        loading.value = false;
        if (response.data.code === 0 || response.data.code === 200) {
          if (response.data.data.status === 1) {
            ElMessage.error("该账号已被封禁，请联系管理员。");
            return;
          }
          ElMessage.success(response.data.message || "登录成功");

          // 存储 JWT Token
          if (response.data.data.token) {
            store.commit("setToken", response.data.data.token);
          }

          if (!response.data.data.avatar.startsWith("http")) {
            response.data.data.avatar =
              store.state.backendUrl + response.data.data.avatar;
          }
          store.commit("setUserInfo", response.data.data);

          // 创建 WebSocket 连接
          const wsUrl =
            store.state.wsUrl + "/wss?client_id=" + response.data.data.uuid;
          store.state.socket = new WebSocket(wsUrl);
          store.state.socket.onopen = () => {
            console.log("WebSocket连接已打开");
          };
          store.state.socket.onmessage = (message) => {
            console.log("收到消息：", message.data);
          };
          store.state.socket.onclose = () => {
            console.log("WebSocket连接已关闭");
          };
          store.state.socket.onerror = () => {
            console.log("WebSocket连接发生错误");
          };
          router.push("/chat/sessionlist");
        } else {
          ElMessage.error(response.data.message);
        }
      } catch (error) {
        loading.value = false;
        ElMessage.error("登录失败，请检查网络连接");
      }
    };
    const checkTelephoneValid = () => {
      const regex = /^1[3456789]\d{9}$/;
      return regex.test(data.loginData.telephone);
    };
    const handleRegister = () => {
      router.push("/register");
    };
    const handleSmsLogin = () => {
      router.push("/smsLogin");
    };

    return {
      ...toRefs(data),
      loading,
      handleLogin,
      handleRegister,
      handleSmsLogin,
    };
  },
};
</script>

<style scoped>
.login-wrap {
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #0f0c29, #302b63, #24243e);
  font-family: 'Inter', -apple-system, BlinkMacSystemFont, sans-serif;
}

.login-card {
  width: 420px;
  padding: 48px 40px;
  border-radius: 24px;
  background: rgba(255, 255, 255, 0.08);
  backdrop-filter: blur(24px);
  -webkit-backdrop-filter: blur(24px);
  border: 1px solid rgba(255, 255, 255, 0.12);
  box-shadow: 0 32px 64px rgba(0, 0, 0, 0.3);
}

.login-brand {
  text-align: center;
  margin-bottom: 36px;
}

.brand-icon {
  font-size: 48px;
  margin-bottom: 12px;
}

.brand-title {
  font-size: 28px;
  font-weight: 700;
  color: #ffffff;
  margin: 0;
  letter-spacing: 2px;
}

.brand-subtitle {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.5);
  margin-top: 6px;
  letter-spacing: 4px;
}

.login-form {
  margin-top: 20px;
}

.login-form :deep(.el-input__wrapper) {
  background: rgba(255, 255, 255, 0.06);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  box-shadow: none;
  transition: all 0.3s ease;
}

.login-form :deep(.el-input__wrapper:hover),
.login-form :deep(.el-input__wrapper.is-focus) {
  border-color: rgba(139, 92, 246, 0.6);
  background: rgba(255, 255, 255, 0.1);
  box-shadow: 0 0 0 3px rgba(139, 92, 246, 0.15);
}

.login-form :deep(.el-input__inner) {
  color: #ffffff;
}

.login-form :deep(.el-input__inner::placeholder) {
  color: rgba(255, 255, 255, 0.35);
}

.login-form :deep(.el-input__prefix .el-icon) {
  color: rgba(255, 255, 255, 0.4);
}

.login-btn {
  width: 100%;
  margin-top: 8px;
  border: none;
  border-radius: 12px;
  font-weight: 600;
  font-size: 16px;
  letter-spacing: 6px;
  color: #ffffff;
  background: linear-gradient(135deg, #7c3aed, #6366f1);
  transition: all 0.3s ease;
}

.login-btn:hover {
  background: linear-gradient(135deg, #6d28d9, #4f46e5);
  transform: translateY(-1px);
  box-shadow: 0 8px 24px rgba(99, 102, 241, 0.35);
}

.login-footer {
  display: flex;
  justify-content: center;
  align-items: center;
  margin-top: 24px;
  gap: 12px;
}

.link-btn {
  background: none;
  border: none;
  cursor: pointer;
  color: rgba(255, 255, 255, 0.55);
  font-size: 13px;
  transition: color 0.2s;
}

.link-btn:hover {
  color: #a78bfa;
}

.divider {
  color: rgba(255, 255, 255, 0.2);
  font-size: 12px;
}
</style>
