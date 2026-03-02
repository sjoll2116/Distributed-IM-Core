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
      <el-table-column prop="is_admin" label="管理员" width="80" >
        <template #default="scope">
          <el-button type="default" v-if="scope.row.is_admin == false"
            >普通用户</el-button
          >
          <el-button type="primary" v-if="scope.row.is_admin == true"
            >管理员</el-button
          >
        </template>
      </el-table-column>
      <el-table-column label="删除状态" width="90">
        <template #default="scope">
          <el-button type="default" v-if="scope.row.is_deleted == false"
            >未删除</el-button
          >
          <el-button type="primary" v-if="scope.row.is_deleted == true"
            >已删除</el-button
          >
        </template>
      </el-table-column>
      <el-table-column label="禁用状态" width="90">
        <template #default="scope">
          <el-button type="default" v-if="scope.row.status == 0"
            >未禁用</el-button
          >
          <el-button type="primary" v-if="scope.row.status == 1"
            >已禁用</el-button
          >
        </template>
      </el-table-column>
    </el-table>
    <div class="footer">
      <el-button type="primary" @click="handleDisableUsers">批量禁用</el-button>
      <el-button type="primary" @click="handleEnableUsers">批量解封</el-button>
    </div>
  </div>
</template>

<script>
import { reactive, toRefs, onMounted } from "vue";
import http from "@/utils/axios";
export default {
  name: "DisableUserModal",
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
    const handleDisableUsers = async () => {
      try {
        const req = { uuids: data.uuidList };
        await http.post("/user/disableUsers", req);
        getUserInfoList();
      } catch (error) {
        console.log(error);
      }
    };
    const handleEnableUsers = async () => {
      try {
        const req = { uuids: data.uuidList };
        await http.post("/user/enableUsers", req);
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
      handleDisableUsers,
      handleEnableUsers,
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
  gap: 10px;
}
</style>
