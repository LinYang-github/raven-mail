<template>
  <el-dialog
    v-model="visible"
    title="新建文电"
    width="600px"
    :before-close="handleClose"
    class="compose-dialog"
    align-center
    top="5vh"
  >
    <el-form :model="form" ref="formRef" label-position="top">
      <el-form-item label="收件人">
        <el-input v-model="form.to" placeholder="用户ID, 逗号分隔" />
      </el-form-item>
      
      <el-form-item label="抄送">
        <el-input v-model="form.cc" placeholder="可选" />
      </el-form-item>
      
      <el-form-item label="主题">
        <el-input v-model="form.subject" placeholder="邮件主题" />
      </el-form-item>
      
      <el-form-item label="正文">
        <el-input
          v-model="form.content"
          type="textarea"
          :rows="10"
          placeholder="在此输入正文..."
        />
      </el-form-item>
      
      <el-form-item>
        <div class="attachment-area">
          <el-upload
            v-model:file-list="fileList"
            action="#"
            :auto-upload="false"
            multiple
            class="upload-demo"
            drag
          >
            <el-icon class="el-icon--upload"><Paperclip /></el-icon>
            <div class="el-upload__text">
              拖拽文件到此处 或 <em>点击上传</em>
            </div>
          </el-upload>
        </div>
      </el-form-item>
    </el-form>
    
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="handleClose">取消</el-button>
        <el-button type="primary" :loading="loading" @click="handleSubmit">
          发送 <el-icon class="el-icon--right"><Promotion /></el-icon>
        </el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script setup>
import { reactive, ref, computed } from 'vue'
import { sendMail } from '../services/api'
import { ElMessage } from 'element-plus'
import { Paperclip, Promotion } from '@element-plus/icons-vue'

const props = defineProps(['modelValue'])
const emit = defineEmits(['update:modelValue', 'success'])

const visible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const loading = ref(false)
const fileList = ref([])
const form = reactive({
  to: '',
  cc: '',
  subject: '',
  content: ''
})

const handleClose = () => {
  visible.value = false
}

const handleSubmit = async () => {
  if (!form.to || !form.subject || !form.content) {
    ElMessage.warning('请填写必要字段')
    return
  }
  
  loading.value = true
  try {
    const formData = new FormData()
    formData.append('to', form.to)
    formData.append('cc', form.cc)
    formData.append('subject', form.subject)
    formData.append('content', form.content)
    
    fileList.value.forEach(file => {
      formData.append('attachments', file.raw)
    })

    await sendMail(formData)
    ElMessage.success('发送成功')
    emit('success')
    handleClose()
    
    // Reset
    form.to = ''
    form.cc = ''
    form.subject = ''
    form.content = ''
    fileList.value = []
  } catch (err) {
    ElMessage.error('发送失败')
    console.error(err)
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.compose-dialog :deep(.el-dialog__body) {
  padding: 20px 30px 10px;
  max-height: 65vh;
  overflow-y: auto;
}

.compose-dialog :deep(.el-form-item__label) {
  font-weight: 500;
  color: #606266;
  margin-bottom: 8px;
}

.compose-dialog :deep(.el-input__wrapper),
.compose-dialog :deep(.el-textarea__inner) {
  box-shadow: 0 0 0 1px #dcdfe6;
  border-radius: 6px;
  padding: 8px 12px;
  transition: all 0.2s;
}

.compose-dialog :deep(.el-input__wrapper:hover),
.compose-dialog :deep(.el-textarea__inner:hover) {
  box-shadow: 0 0 0 1px #409EFF;
}

.compose-dialog :deep(.el-input__wrapper.is-focus),
.compose-dialog :deep(.el-textarea__inner:focus) {
  box-shadow: 0 0 0 2px rgba(64, 158, 255, 0.2) !important;
  border-color: #409EFF;
}

.attachment-area {
  width: 100%;
  margin-top: 10px;
}

.upload-demo :deep(.el-upload-dragger) {
  padding: 20px;
  border: 1px dashed #dcdfe6;
  border-radius: 6px;
  height: auto;
  transition: all 0.3s;
}

.upload-demo :deep(.el-upload-dragger:hover) {
  border-color: #409EFF;
  background-color: #f5f7fa;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding-top: 10px;
}
</style>
