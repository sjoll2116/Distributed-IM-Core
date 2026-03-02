<template>
  <div style="height: 100%; width: 100%">
    <el-table
      :data="userList"
      style="width: 100%"
      @selection-change="selectUsers"
      height="300"
    >
      <el-table-column type="selection" width="55" />
      <el-table-column prop="uuid" label="Uuid" width="200" />
      <el-table-column
        prop="nickname"
        label="昵称"
        width="120"
        show-overflow-tooltip
      />
      <el-table-column prop="telephone" label="电话" width="180" />
      <el-table-column prop="is_admin" label="管理员" width="120" >
        <template #default="scope">
          <el-button type="default" v-if="scope.row.is_admin == false"
            >普通用户</el-button
          >
          <el-button type="primary" v-if="scope.row.is_admin == true"
            >管理员</el-button
          >
        </template>
      </el-table-column>
      <el-table-column label="删除状态" width="120">
        <template #default="scope">
          <el-button type="default" v-if="scope.row.is_deleted == false"
            >未删除</el-button
          >
          <el-button type="primary" v-if="scope.row.is_deleted == true"
            >已删除</el-button
          >
        </template>
      </el-table-column>
    </el-table>
    <div class="footer">
      <el-button type="primary" @click="handleDeleteUsers">批量删除</el-button>
    </div>
  </div>
</template>

<script>
import { reactive, toRefs, onMounted } from "vue";
import http from "@/utils/axios";
export default {
  name: "DeleteUserModal",
  setup() {
    const data = reactive({
      userList: [],
      uuidList: [],
    });
    const getUserInfoList = async () => {
      try {
        const rsp = await http.post("/user/getUserInfoList");
        data.userList = rsp.data.data;
      } catch (error) {
        console.log(error);
      }
    };
    const selectUsers = (val) => {
      data.uuidList = val.map((item) => item.uuid);
    };
    const handleDeleteUsers = async () => {
      try {
        const req = {
           uuids: data.uuidList,
        };
        const rsp = await http.post("/user/deleteUsers", req);
        console.log(rsp.data);
        getUserInfoList();
      } catch (error) {
        console.log(error);
      }
    };
    onMounted(() => {
      getUserInfoList();
    });
    return {
      ...toRefs(data),
      selectUsers,
      handleDeleteUsers,
    };
  },
};
</script>

<style scoped>
.footer {
  display: flex;
  margin-top: 20px;
  flex-direction: row-reverse;
  align-items: center;
}
</style>
