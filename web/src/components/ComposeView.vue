<template>
  <div class="compose-view">
    <div class="compose-header">
      <h2>新建文电</h2>
      <div class="actions">
        <el-button @click="$emit('cancel')">取消</el-button>
        <el-button type="primary" :loading="loading" @click="handleSubmit">
          发送 <el-icon class="el-icon--right"><Promotion /></el-icon>
        </el-button>
      </div>
    </div>

    <div class="compose-form">
      <el-form :model="form" ref="formRef" label-position="top" class="main-form">
        <div class="form-row">
          <el-form-item label="收件人" class="flex-item">
            <el-select
              v-model="form.toList"
              multiple
              filterable
              remote
              reserve-keyword
              placeholder="搜索并选择收件人"
              :remote-method="searchToUsers"
              :loading="searchLoading"
              style="width: 100%"
              @focus="() => searchToUsers('')"
            >
              <el-option
                v-for="item in toUserOptions"
                :key="item.id"
                :label="item.name"
                :value="item.id"
              >
                <span style="float: left">{{ item.name }}</span>
                <span style="float: right; color: #8492a6; font-size: 13px">{{ item.dept }}</span>
              </el-option>
            </el-select>
          </el-form-item>
          
          <el-form-item label="抄送" class="flex-item">
            <el-select
              v-model="form.ccList"
              multiple
              filterable
              remote
              reserve-keyword
              placeholder="可选"
              :remote-method="searchCcUsers"
              :loading="searchLoading"
              style="width: 100%"
              @focus="() => searchCcUsers('')"
            >
              <el-option
                v-for="item in ccUserOptions"
                :key="item.id"
                :label="item.name"
                :value="item.id"
              >
                <span style="float: left">{{ item.name }}</span>
                <span style="float: right; color: #8492a6; font-size: 13px">{{ item.dept }}</span>
              </el-option>
            </el-select>
          </el-form-item>
        </div>
        
        <el-form-item label="主题">
          <el-input v-model="form.subject" placeholder="请输入主题" />
        </el-form-item>

        <el-form-item label="正文">
          <EditorDriver v-model="form.content" :mailId="null" />
        </el-form-item>
        
        <el-form-item label="附件">
          <el-upload
            v-model:file-list="fileList"
            action="#"
            :auto-upload="false"
            multiple
            class="upload-inline"
            :show-file-list="true"
          >
            <el-button link>
              <el-icon><Paperclip /></el-icon> 添加附件
            </el-button>
          </el-upload>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref, onMounted } from 'vue'
import { sendMail } from '../services/api'
import { userStore } from '../store/user'
import { EditorDriver } from './content'
import { ElMessage } from 'element-plus'
import { Paperclip, Promotion } from '@element-plus/icons-vue'

const emit = defineEmits(['cancel', 'success'])

const loading = ref(false)
const searchLoading = ref(false)
const toUserOptions = ref([])
const ccUserOptions = ref([])
const fileList = ref([])

const form = reactive({
  toList: [],
  ccList: [],
  subject: '',
  content: ''
})

const fetchDefaultOptions = async () => {
  if (userStore.fetchUsers) {
    const results = await userStore.fetchUsers('')
    toUserOptions.value = results
    ccUserOptions.value = results
  }
}

onMounted(() => {
  fetchDefaultOptions()
})

const searchToUsers = async (query) => {
  searchLoading.value = true
  try {
    const results = userStore.fetchUsers ? await userStore.fetchUsers(query) : []
    toUserOptions.value = results
  } finally {
    searchLoading.value = false
  }
}

const searchCcUsers = async (query) => {
  searchLoading.value = true
  try {
    const results = userStore.fetchUsers ? await userStore.fetchUsers(query) : []
    ccUserOptions.value = results
  } finally {
    searchLoading.value = false
  }
}

const handleSubmit = async () => {
  const isOnlyOffice = (import.meta.env.VITE_MAIL_CONTENT_MODE === 'onlyoffice')
  
  if (form.toList.length === 0 || !form.subject || (!isOnlyOffice && !form.content)) {
    ElMessage.warning('请填写必要字段')
    return
  }

  loading.value = true

  // ONLYOFFICE 模式下，发送前强行触发一次服务端保存
  if (isOnlyOffice && form.content) {
    try {
      await fetch(`${import.meta.env.VITE_BACKEND_URL}/api/v1/onlyoffice/forcesave?key=${form.content}`, { method: 'POST' })
      // 给一点缓冲区让 ONLYOFFICE 完成回调（通常 1-2 秒足够，比 10 秒好得多）
      await new Promise(resolve => setTimeout(resolve, 1500))
    } catch (err) {
      console.warn('[OnlyOffice] Force save trigger failed', err)
    }
  }
  
  try {
    const formData = new FormData()
    formData.append('to', form.toList.join(','))
    formData.append('cc', form.ccList.join(','))
    formData.append('subject', form.subject)
    formData.append('content', form.content)
    formData.append('content_type', import.meta.env.VITE_MAIL_CONTENT_MODE || 'text')
    
    fileList.value.forEach(file => {
      formData.append('attachments', file.raw)
    })

    await sendMail(formData)
    ElMessage.success('发送成功')
    emit('success')
  } catch (err) {
    ElMessage.error('发送失败')
    console.error(err)
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.compose-view {
  flex: 1;
  display: flex;
  flex-direction: column;
  height: 100%;
  background: white;
}

.compose-header {
  padding: 16px 24px;
  border-bottom: 1px solid #ebedf0;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.compose-header h2 {
  margin: 0;
  font-size: 18px;
  color: #1a1a1a;
}

.compose-form {
  flex: 1;
  padding: 30px 50px;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  align-items: center; /* Center the form */
  background: #fff;
}

.main-form {
  width: 100%;
  max-width: 900px; /* Limit max width */
  display: flex;
  flex-direction: column;
  height: 100%;
}

.form-row {
  display: flex;
  gap: 20px;
}

.flex-item {
  flex: 1;
}

.editor-container {
  flex: 1;
  margin: 10px 0;
  display: flex;
  flex-direction: column;
  min-height: 200px;
}

.content-editor {
  flex: 1;
  display: flex; /* Helps textarea fill height */
}

:deep(.el-textarea__inner) {
  height: 100% !important;
  min-height: 200px;
  padding: 16px;
  font-size: 14px;
  line-height: 1.6;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  resize: none;
}

:deep(.el-textarea__inner:focus) {
  border-color: #409EFF;
}

.attachment-bar {
  padding-top: 10px;
  border-top: 1px solid #f0f0f0;
}

:deep(.el-form-item) {
  margin-bottom: 18px;
}
</style>
