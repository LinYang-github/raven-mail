<template>
  <div class="onlyoffice-preview">
    <div class="placeholder">
      <el-icon size="40"><Document /></el-icon>
      <div class="text">文电正文 (ONLYOFFICE 文档)</div>
      <el-button type="primary" plain @click="dialogVisible = true">在线阅读</el-button>
    </div>

    <!-- 预览弹窗 -->
    <el-dialog
      v-model="dialogVisible"
      title="在线预览"
      width="90%"
      top="5vh"
      destroy-on-close
      @opened="handleOpened"
    >
      <div class="viewer-container">
        <div v-if="loading" class="loading-overlay">
          <el-icon class="is-loading"><Loading /></el-icon>
          正在从 ONLYOFFICE 服务器加载文档...
        </div>
        <div :id="viewerId" class="viewer-instance"></div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onBeforeUnmount } from 'vue'
import { Document, Loading } from '@element-plus/icons-vue'

const props = defineProps(['content']) // 这里 content 存放的是 docKey

const dialogVisible = ref(false)
const loading = ref(true)
const viewerId = 'onlyoffice-viewer-instance'
let docEditor = null

const serverUrl = import.meta.env.VITE_ONLYOFFICE_SERVER || 'http://localhost:8090/'
const backendUrl = import.meta.env.VITE_BACKEND_URL || 'http://localhost:8080'

const loadScript = () => {
  return new Promise((resolve) => {
    if (window.DocsAPI) return resolve()
    const script = document.createElement('script')
    script.src = `${serverUrl}web-apps/apps/api/documents/api.js`
    script.onload = resolve
    document.head.appendChild(script)
  })
}

const handleOpened = async () => {
  loading.value = true
  await loadScript()
  
  const config = {
    document: {
      fileType: "docx",
      key: props.content, // 使用保存的 docKey
      title: "文电正文.docx",
      url: `${backendUrl}/api/v1/onlyoffice/template?key=${props.content}`, // 实际应指向该邮件对应的物理文件
    },
    documentType: "word",
    editorConfig: {
      mode: "view", // 关键：预览模式
      lang: "zh",
      user: { id: "viewer", name: "阅读者" },
      customization: {
        compactHeader: true,
        toolbarNoTabs: true
      }
    },
    height: "100%",
    width: "100%"
  }

  try {
    if (window.DocsAPI) {
      docEditor = new window.DocsAPI.DocEditor(viewerId, config)
      loading.value = false
    }
  } catch (err) {
    console.error('[OnlyOffice] Preview failed:', err)
  }
}

onBeforeUnmount(() => {
  if (docEditor) docEditor.destroyEditor()
})
</script>

<style scoped>
.onlyoffice-preview {
  padding: 40px;
  background: #f8f9fb;
  border-radius: 8px;
  border: 1px solid #ebedf0;
  text-align: center;
}
.placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 15px;
  color: #606266;
}
.text {
  font-weight: 600;
  margin-bottom: 5px;
}
.viewer-container {
  height: 75vh;
  position: relative;
  border: 1px solid #dcdfe6;
}
.viewer-instance {
  height: 100%;
  width: 100%;
}
.loading-overlay {
  position: absolute;
  top:0; left:0; right:0; bottom:0;
  background: white;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  z-index: 10;
  gap: 10px;
}
</style>
