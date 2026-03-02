<template>
  <router-view />
</template>

<script>
import { onMounted } from "vue";
import { useStore } from "vuex";
import http from "@/utils/axios";
export default {
  name: "App",
  setup() {
    const store = useStore();
    const getUserInfo = async () => {
      try {
        const req = {
          uuid: store.state.userInfo.uuid,
        };
        const rsp = await http.post("/user/getUserInfo", req);
        if (rsp.data.code === 0 || rsp.data.code === 200) {
          if (!rsp.data.data.avatar.startsWith("http")) {
            rsp.data.data.avatar = store.state.backendUrl + rsp.data.data.avatar;
          }
          store.commit("setUserInfo", rsp.data.data);
        } else {
          console.error(rsp.data.message);
        }
      } catch (error) {
        console.log(error);
      }
    };
    onMounted(() => {
      if (store.state.userInfo.uuid) {
        getUserInfo();
        const wsUrl =
          store.state.wsUrl + "/wss?client_id=" + store.state.userInfo.uuid;
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
      }
    });
  },
};
</script>

<style>
@import url('https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700&display=swap');

* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}
</style>