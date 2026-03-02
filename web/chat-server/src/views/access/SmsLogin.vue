<template>
  <div class="login-wrap">
    <div class="login-card">
      <div class="login-brand">
        <div class="brand-icon">📲</div>
        <h1 class="brand-title">手机验证</h1>
        <p class="brand-subtitle">验证码快速登录</p>
      </div>
      <el-form :model="smsData" class="login-form">
        <el-form-item>
          <el-input v-model="smsData.telephone" placeholder="手机号码" prefix-icon="User" size="large" />
        </el-form-item>
        <el-form-item>
          <div class="sms-input-row">
            <el-input v-model="smsData.code" placeholder="验证码" prefix-icon="Message" size="large" />
            <el-button type="primary" class="code-btn" @click="handleSendCode" :disabled="timer > 0">
               {{ timer > 0 ? timer + 's' : '获取' }}
            </el-button>
          </div>
        </el-form-item>
        <el-button class="login-btn" size="large" @click="handleSmsLogin" :loading="loading">登 录</el-button>
      </el-form>
      <div class="login-footer">
        <button class="link-btn" @click="handleToLogin">返回账号登录</button>
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
  name: "SmsLogin",
  setup() {
    const data = reactive({
      smsData: { telephone: "", code: "" },
      timer: 0,
    });
    const loading = ref(false);
    const router = useRouter();
    const store = useStore();

    const handleSendCode = async () => {
       if (!data.smsData.telephone) { ElMessage.error("请先填写手机号"); return; }
       try {
          await http.post("/auth/sendCode", { telephone: data.smsData.telephone });
          ElMessage.success("验证码已发送");
          data.timer = 60;
          const interval = setInterval(() => {
             data.timer--;
             if (data.timer <= 0) clearInterval(interval);
          }, 1000);
       } catch (error) { ElMessage.error("发送验证码失败"); }
    };

    const handleSmsLogin = async () => {
      try {
        if (!data.smsData.telephone || !data.smsData.code) { ElMessage.error("请完善信息"); return; }
        loading.value = true;
        const response = await http.post("/auth/smsLogin", data.smsData);
        loading.value = false;
        if (response.data.code === 0 || response.data.code === 200) {
           store.commit("setToken", response.data.data.token);
           store.commit("setUserInfo", response.data.data);
           router.push("/chat/sessionlist");
        } else { ElMessage.error(response.data.message); }
      } catch (error) { loading.value = false; ElMessage.error("登录失败"); }
    };

    const handleToLogin = () => { router.push("/login"); };

    return { ...toRefs(data), loading, handleSendCode, handleSmsLogin, handleToLogin };
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
.sms-input-row { display: flex; gap: 10px; width: 100%; }
.code-btn { width: 100px; }
.login-btn { width: 100%; margin-top: 8px; border: none; border-radius: 12px; font-weight: 600; font-size: 16px; letter-spacing: 6px; color: #fff; background: linear-gradient(135deg, #7c3aed, #6366f1); transition: all 0.3s ease; }
.login-btn:hover { background: linear-gradient(135deg, #6d28d9, #4f46e5); transform: translateY(-1px); box-shadow: 0 8px 24px rgba(99, 102, 241, 0.35); }
.login-footer { margin-top: 24px; text-align: center; }
.link-btn { background: none; border: none; cursor: pointer; color: rgba(255,255,255,0.55); font-size: 13px; transition: color 0.2s; }
.link-btn:hover { color: #a78bfa; }
</style>
