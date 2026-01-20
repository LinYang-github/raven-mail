<template>
  <div class="onlyoffice-editor">
    <div v-if="loading" class="loading-overlay">
      <el-icon class="is-loading"><Loading /></el-icon>
      正在连接到 ONLYOFFICE 服务器...
    </div>
    <div :id="editorId" class="editor-container"></div>
  </div>
</template>

<script setup>
import { onMounted, onBeforeUnmount, ref } from 'vue'
import { Loading } from '@element-plus/icons-vue'

const props = defineProps(['modelValue', 'mailId'])
const emit = defineEmits(['update:modelValue'])

const editorId = 'onlyoffice-editor-instance'
const loading = ref(true)
let docEditor = null

const serverUrl = import.meta.env.VITE_ONLYOFFICE_SERVER || 'http://localhost:8090/'

const loadScript = () => {
  return new Promise((resolve) => {
    if (window.DocsAPI) return resolve()
    const script = document.createElement('script')
    script.src = `${serverUrl}web-apps/apps/api/documents/api.js`
    script.onload = resolve
    document.head.appendChild(script)
  })
}

const initEditor = async () => {
  await loadScript()
  
  const config = {
    document: {
      fileType: "docx",
      key: `mail-${props.mailId || 'new'}-${Date.now()}`,
      title: "文电正文.docx",
      url: "", // In real case, this should be a link to the backend file
    },
    documentType: "word",
    editorConfig: {
      mode: "edit",
      callbackUrl: "http://localhost:8080/api/v1/onlyoffice/callback", // Backend callback
      lang: "zh",
      user: {
        id: "current-user",
        name: "当前用户"
      },
      customization: {
        autosave: true,
        compactHeader: true,
        toolbarNoTabs: true
      }
    },
    height: "100%",
    width: "100%"
  }

  // Simulated: If it's a new mail, we might just show a "Mock" editor since we don't have a file server for docx yet
  console.log('[OnlyOffice] Init config:', config)
  
  setTimeout(() => {
    loading.value = false
    // docEditor = new window.DocsAPI.DocEditor(editorId, config)
    const container = document.getElementById(editorId)
    if (container) {
      container.innerHTML = `
        <div style="background:#f1f1f1; height:100%; display:flex; flex-direction:column; align-items:center; justify-content:center; border: 2px dashed #ccc;">
          <img src="https://www.onlyoffice.com/blog/wp-content/uploads/2021/05/docs-6-3-696x392.png" style="width: 200px; margin-bottom: 20px;" />
          <p><b>ONLYOFFICE 模拟界面</b></p>
          <p style="color: #666; font-size: 13px;">服务器地址: ${serverUrl}</p>
          <p style="color: #909399; font-size: 12px; padding: 0 40px; text-align:center;">
            生产环境下，此区域将通过 DocsAPI 挂载真实的文档编辑器。<br/>
            当前模式识别为: ONLYOFFICE
          </p>
        </div>
      `
    }
  }, 1000)
}

onMounted(() => {
  initEditor()
})

onBeforeUnmount(() => {
  if (docEditor) {
    docEditor.destroyEditor()
  }
})
</script>

<style scoped>
.onlyoffice-editor {
  height: 600px;
  width: 100%;
  position: relative;
  border: 1px solid #dcdfe6;
}
.editor-container {
  height: 100%;
  width: 100%;
}
.loading-overlay {
  position: absolute;
  top:0; left:0; right:0; bottom:0;
  background: rgba(255,255,255,0.8);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  z-index: 10;
  gap: 10px;
  color: #409EFF;
}
</style>
