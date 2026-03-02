<template>
  <div style="height: 100%; width: 100%">
    <el-table
      :data="groupList"
      style="width: 100%"
      @selection-change="selectGroups"
      height="300"
    >
      <el-table-column type="selection" width="55" />
      <el-table-column prop="uuid" label="Uuid" width="200" />
      <el-table-column
        prop="name"
        label="群名称"
        width="120"
        show-overflow-tooltip
      />
      <el-table-column prop="owner_id" label="群主id" width="200" />
      <el-table-column prop="status" label="禁用状态" width="120" >
        <template #default="scope">
          <el-button type="default" v-if="scope.row.status == 0"
            >未禁用</el-button
          >
          <el-button type="primary" v-if="scope.row.status == 1"
            >已禁用</el-button
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
      <el-button type="primary" @click="handleDisableGroups">批量禁用</el-button>
      <el-button type="primary" @click="handleEnableGroups">批量解封</el-button>
    </div>
  </div>
</template>

<script>
import { reactive, toRefs, onMounted } from "vue";
import http from "@/utils/axios";
export default {
  name: "DisableGroupModal",
  setup() {
    const data = reactive({
      groupList: [],
      uuidList: [],
    });
    const getGroupInfoList = async () => {
      try {
        const rsp = await http.post("/group/getGroupInfoList");
        data.groupList = rsp.data.data;
      } catch (error) {
        console.log(error);
      }
    };
    const selectGroups = (val) => {
      data.uuidList = val.map((item) => item.uuid);
    };
    const handleDisableGroups = async () => {
      try {
        const req = { uuids: data.uuidList };
        await http.post("/group/disableGroups", req);
        getGroupInfoList();
      } catch (error) {
        console.log(error);
      }
    };
    const handleEnableGroups = async () => {
      try {
        const req = { uuids: data.uuidList };
        await http.post("/group/enableGroups", req);
        getGroupInfoList();
      } catch (error) {
        console.log(error);
      }
    };
    onMounted(() => {
      getGroupInfoList();
    });
    return {
      ...toRefs(data),
      selectGroups,
      handleDisableGroups,
      handleEnableGroups,
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
