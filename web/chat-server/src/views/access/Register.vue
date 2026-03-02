<template>
  <div class="login-wrap">
    <div class="login-card">
      <div class="login-brand">
        <div class="brand-icon">📝</div>
        <h1 class="brand-title">注册账号</h1>
        <p class="brand-subtitle">加入 KamaChat 社交网络</p>
      </div>
      <el-form :model="registerData" class="login-form">
        <el-form-item>
          <el-input v-model="registerData.telephone" placeholder="手机号码" prefix-icon="User" size="large" />
        </el-form-item>
        <el-form-item>
          <el-input v-model="registerData.nickname" placeholder="昵称" prefix-icon="Postcard" size="large" />
        </el-form-item>
        <el-form-item>
          <el-input type="password" v-model="registerData.password" placeholder="设置密码" prefix-icon="Lock" size="large" show-password />
        </el-form-item>
        <el-button class="login-btn" size="large" @click="handleRegister" :loading="loading">注 册</el-button>
      </el-form>
      <div class="login-footer">
        <button class="link-btn" @click="handleToLogin">已有账号？立即登录</button>
      </div>
    </div>
  </div>
</template>

<script>
import { reactive, toRefs, ref } from "vue";
import http from "@/utils/axios";
import { useRouter } from "vue-router";
import { ElMessage } from "element-plus";

export default {
  name: "Register",
  setup() {
    const data = reactive({
      registerData: {
        telephone: "",
        nickname: "",
        password: "",
      },
    });
    const loading = ref(false);
    const router = useRouter();

    const handleRegister = async () => {
      try {
        if (!data.registerData.telephone || !data.registerData.nickname || !data.registerData.password) {
          ElMessage.error("请填写完整信息。");
          return;
        }
        loading.value = true;
        const response = await http.post("/register", data.registerData);
        loading.value = false;
        if (response.data.code === 0 || response.data.code === 200) {
          ElMessage.success("注册成功！即将跳转登录页");
          setTimeout(() => router.push("/login"), 1500);
        } else {
          ElMessage.error(response.data.message);
        }
      } catch (error) {
        loading.value = false;
        ElMessage.error("注册失败，网络异常");
      }
    };

    const handleToLogin = () => {
      router.push("/login");
    };

    return {
      ...toRefs(data),
      loading,
      handleRegister,
      handleToLogin,
    };
  },
};
</script>

<style scoped>
/* 共享 Login.vue 的样式系统 */
.login-wrap { height: 100vh; display: flex; align-items: center; justify-content: center; background: linear-gradient(135deg, #0f0c29, #302b63, #24243e); font-family: 'Inter', sans-serif; }
.login-card { width: 420px; padding: 48px 40px; border-radius: 24px; background: rgba(255, 255, 255, 0.08); backdrop-filter: blur(24px); border: 1px solid rgba(255, 255, 255, 0.12); box-shadow: 0 32px 64px rgba(0,0,0,0.3); }
.login-brand { text-align: center; margin-bottom: 36px; }
.brand-icon { font-size: 48px; margin-bottom: 12px; }
.brand-title { font-size: 28px; font-weight: 700; color: #fff; margin: 0; letter-spacing: 2px; }
.brand-subtitle { font-size: 13px; color: rgba(255,255,255,0.5); margin-top: 6px; letter-spacing: 4px; }
.login-btn { width: 100%; margin-top: 8px; border: none; border-radius: 12px; font-weight: 600; font-size: 16px; letter-spacing: 6px; color: #fff; background: linear-gradient(135deg, #7c3aed, #6366f1); transition: all 0.3s ease; }
.login-btn:hover { background: linear-gradient(135deg, #6d28d9, #4f46e5); transform: translateY(-1px); box-shadow: 0 8px 24px rgba(99, 102, 241, 0.35); }
.login-footer { margin-top: 24px; text-align: center; }
.link-btn { background: none; border: none; cursor: pointer; color: rgba(255,255,255,0.55); font-size: 13px; transition: color 0.2s; }
.link-btn:hover { color: #a78bfa; }
</style>
