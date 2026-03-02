<template>
  <div class="manager-wrap">
    <el-container class="manager-container">
      <el-aside width="200px">
        <el-menu
          default-active="1"
          class="el-menu-vertical"
          @select="handleSelect"
        >
          <el-menu-item index="1">
            <el-icon><User /></el-icon>
            <span>用户管理</span>
          </el-menu-item>
          <el-menu-item index="2">
            <el-icon><ChatDotRound /></el-icon>
            <span>群聊管理</span>
          </el-menu-item>
          <el-menu-item index="3">
            <el-icon><Setting /></el-icon>
            <span>系统设置</span>
          </el-menu-item>
        </el-menu>
      </el-aside>
      <el-main>
        <div class="main-content">
          <div v-if="activeIndex == '1'" class="content-section">
             <div class="toolbar">
                <el-button type="danger" @click="showDeleteUserModal = true">删除用户</el-button>
                <el-button type="warning" @click="showDisableUserModal = true">封禁用户</el-button>
                <el-button type="primary" @click="showSetAdminModal = true">权限管理</el-button>
             </div>
             <el-table :data="userList" stripe style="width: 100%">
                <el-table-column prop="uuid" label="ID" width="200" />
                <el-table-column prop="nickname" label="昵称" />
                <el-table-column prop="telephone" label="电话" />
                <el-table-column prop="is_admin" label="管理员">
                   <template #default="scope">
                      {{ scope.row.is_admin ? '是' : '否' }}
                   </template>
                </el-table-column>
                <el-table-column prop="status" label="状态">
                   <template #default="scope">
                      <el-tag :type="scope.row.status == 0 ? 'success' : 'danger'">
                         {{ scope.row.status == 0 ? '正常' : '封禁' }}
                      </el-tag>
                   </template>
                </el-table-column>
             </el-table>
          </div>
          <div v-if="activeIndex == '2'" class="content-section">
             <div class="toolbar">
                <el-button type="danger" @click="showDeleteGroupModal = true">解散群聊</el-button>
                <el-button type="warning" @click="showDisableGroupModal = true">封禁群聊</el-button>
             </div>
             <el-table :data="groupList" stripe style="width: 100%">
                <el-table-column prop="uuid" label="群ID" width="200" />
                <el-table-column prop="name" label="群名称" />
                <el-table-column prop="owner_id" label="创建者" />
                <el-table-column prop="status" label="状态">
                   <template #default="scope">
                      <el-tag :type="scope.row.status == 0 ? 'success' : 'danger'">
                         {{ scope.row.status == 0 ? '正常' : '封禁' }}
                      </el-tag>
                   </template>
                </el-table-column>
             </el-table>
          </div>
        </div>
      </el-main>
    </el-container>

    <!-- Modals -->
    <el-dialog v-model="showDeleteUserModal" title="批量删除用户" width="70%">
       <DeleteUserModal />
    </el-dialog>
    <el-dialog v-model="showDisableUserModal" title="批量封禁用户" width="70%">
       <DisableUserModal />
    </el-dialog>
    <el-dialog v-model="showSetAdminModal" title="管理员权限管理" width="70%">
       <SetAdminModal />
    </el-dialog>
    <el-dialog v-model="showDeleteGroupModal" title="批量解散群聊" width="70%">
       <DeleteGroupModal />
    </el-dialog>
    <el-dialog v-model="showDisableGroupModal" title="批量封禁群聊" width="70%">
       <DisableGroupModal />
    </el-dialog>

  </div>
</template>

<script>
import { reactive, toRefs, onMounted } from "vue";
import http from "@/utils/axios";
import DeleteUserModal from "@/components/DeleteUserModal.vue";
import DisableUserModal from "@/components/DisableUserModal.vue";
import SetAdminModal from "@/components/SetAdminModal.vue";
import DeleteGroupModal from "@/components/DeleteGroupModal.vue";
import DisableGroupModal from "@/components/DisableGroupModal.vue";

export default {
  name: "Manager",
  components: {
    DeleteUserModal,
    DisableUserModal,
    SetAdminModal,
    DeleteGroupModal,
    DisableGroupModal
  },
  setup() {
    const data = reactive({
      activeIndex: "1",
      userList: [],
      groupList: [],
      showDeleteUserModal: false,
      showDisableUserModal: false,
      showSetAdminModal: false,
      showDeleteGroupModal: false,
      showDisableGroupModal: false
    });

    const handleSelect = (index) => {
      data.activeIndex = index;
      if (index == "1") fetchUsers();
      if (index == "2") fetchGroups();
    };

    const fetchUsers = async () => {
      try {
        const rsp = await http.post("/user/getUserInfoList");
        data.userList = rsp.data.data;
      } catch (error) { console.error(error); }
    };

    const fetchGroups = async () => {
      try {
        const rsp = await http.post("/group/getGroupInfoList");
        data.groupList = rsp.data.data;
      } catch (error) { console.error(error); }
    };

    onMounted(() => {
      fetchUsers();
    });

    return {
      ...toRefs(data),
      handleSelect
    };
  }
};
</script>

<style scoped>
.manager-wrap { height: 100vh; background: #f5f7fa; }
.manager-container { height: 100%; }
.el-aside { background: #fff; border-right: 1px solid #e6e6e6; }
.toolbar { margin-bottom: 20px; display: flex; gap: 10px; }
.main-content { padding: 20px; }
</style>
