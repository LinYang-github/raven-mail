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
import { userStore } from '../../store/user'

const props = defineProps(['modelValue', 'mailId'])
const emit = defineEmits(['update:modelValue'])

const editorId = 'onlyoffice-editor-instance'
const loading = ref(true)
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

// Key 中包含 SessionID，彻底物理隔离缓存
const docKey = ref(`s-${userStore.sessionId}-m-${props.mailId || 'new'}-${Math.random().toString(36).substring(7)}`)

const initEditor = async () => {
  await loadScript()
  
  const effectiveBackendUrl = backendUrl
  
  const config = {
    document: {
      fileType: "docx",
      key: docKey.value,
      title: "文电正文.docx",
      // URL 携带 session_id 参数
      url: `${effectiveBackendUrl}/api/v1/onlyoffice/template?key=${docKey.value}&session_id=${userStore.sessionId}`, 
    },
    documentType: "word",
    editorConfig: {
      mode: "edit",
      // Callback 携带 session_id 参数，方便后端保存到正确目录
      callbackUrl: `${effectiveBackendUrl}/api/v1/onlyoffice/callback?session_id=${userStore.sessionId}`,
      lang: "zh",
      user: {
        id: userStore.id,
        name: userStore.name
      },
      customization: {
        autosave: true,
        compactHeader: true,
        toolbarNoTabs: false
      }
    },
    height: "100%",
    width: "100%"
  }

  console.log('[OnlyOffice] Initializing with server:', serverUrl)
  
  try {
    loading.value = false
    if (window.DocsAPI) {
      docEditor = new window.DocsAPI.DocEditor(editorId, config)
      // 将 docKey 同步给父组件，作为“内容”标识，解决校验问题
      emit('update:modelValue', docKey.value)
    } else {
      throw new Error('DocsAPI not found')
    }
  } catch (err) {
    console.error('[OnlyOffice] Init failed:', err)
    const container = document.getElementById(editorId)
    if (container) {
      container.innerHTML = `
        <div style="background:#fff2f0; height:100%; display:flex; flex-direction:column; align-items:center; justify-content:center; border: 2px dashed #ffccc7; color: #ff4d4f; padding: 20px; text-align:center;">
          <p><b>ONLYOFFICE 服务连接失败</b></p>
          <p style="font-size: 13px;">请检查浏览器是否能访问: <br/> ${serverUrl}web-apps/apps/api/documents/api.js</p>
          <p style="font-size: 12px; color: #999;">错误信息: ${err.message}</p>
        </div>
      `
    }
  }
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
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  position: relative;
  overflow: hidden;
}
.editor-container {
  flex: 1;
  width: 100%;
  height: 100%; /* DocsAPI 强依赖此属性完成 Iframe 初始化 */
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
