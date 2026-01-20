<template>
  <div class="onlyoffice-preview-wrap" :class="{ 'is-fullscreen': isFullScreen }">
    <!-- 全屏模式下的顶部栏 -->
    <div v-if="isFullScreen" class="fullscreen-toolbar">
      <span class="doc-title">文电查阅: {{ docTitle || '正文内容' }}</span>
      <el-button type="primary" plain size="small" @click="toggleFullScreen">
        <el-icon><Aim /></el-icon> 退出全屏
      </el-button>
    </div>

    <div class="preview-scroll-area">
      <div class="preview-container">
        <!-- 悬浮全屏按钮 (非全屏模式下显示) -->
        <div v-if="!isFullScreen" class="floating-actions">
          <el-tooltip content="全屏阅读" placement="left">
            <el-button circle icon="FullScreen" @click="toggleFullScreen"></el-button>
          </el-tooltip>
        </div>

        <div v-if="loading" class="loading-overlay">
          <el-icon class="is-loading"><Loading /></el-icon>
          正在读取正文...
        </div>
        
        <div :key="renderKey" class="viewer-wrapper">
          <div :id="viewerId" class="viewer-instance"></div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount, watch } from 'vue'
import { Loading, FullScreen, Aim } from '@element-plus/icons-vue'

const props = defineProps(['content', 'title'])
const docTitle = ref(props.title)
const loading = ref(true)
const isFullScreen = ref(false)
const renderKey = ref(0)
const viewerId = `onlyoffice-viewer-${Math.random().toString(36).substring(7)}`
let docEditor = null

const serverUrl = import.meta.env.VITE_ONLYOFFICE_SERVER || 'http://localhost:8090/'
const backendUrl = import.meta.env.VITE_BACKEND_URL || 'http://localhost:8080'

const toggleFullScreen = () => {
  isFullScreen.value = !isFullScreen.value
}

// 监听 Esc 键退出全屏
const handleEsc = (e) => {
  if (e.key === 'Escape' && isFullScreen.value) {
    isFullScreen.value = false
  }
}

watch(isFullScreen, (val) => {
  if (val) {
    document.body.style.overflow = 'hidden'
    window.addEventListener('keydown', handleEsc)
  } else {
    document.body.style.overflow = ''
    window.removeEventListener('keydown', handleEsc)
  }
  // 通知编辑器容器大小变化
  setTimeout(() => {
    window.dispatchEvent(new Event('resize'))
  }, 100)
})

const loadScript = () => {
  return new Promise((resolve) => {
    if (window.DocsAPI) return resolve()
    const script = document.createElement('script')
    script.src = `${serverUrl}web-apps/apps/api/documents/api.js`
    script.onload = resolve
    document.head.appendChild(script)
  })
}

const initViewer = async () => {
  if (!props.content) return
  
  loading.value = true
  renderKey.value++ 
  
  if (docEditor) {
    docEditor.destroyEditor()
    docEditor = null
  }

  await loadScript()
  
  const config = {
    type: "desktop", 
    document: {
      fileType: "docx",
      key: props.content + "_v" + Date.now(), 
      title: "文电正文.docx",
      url: `${backendUrl}/api/v1/onlyoffice/template?key=${props.content}`,
    },
    documentType: "word",
    editorConfig: {
      mode: "view",
      lang: "zh",
      customization: {
        toolbar: false,
        header: false,
        statusBar: false,
        leftMenu: false,
        rightMenu: false,
        comments: false,
        help: false,
        about: false,
        compactHeader: true,
        autosave: false,
        chat: false,
        zoom: 100
      }
    },
    height: "100%",
    width: "100%"
  }

  try {
    if (window.DocsAPI) {
      setTimeout(() => {
        docEditor = new window.DocsAPI.DocEditor(viewerId, config)
        loading.value = false
      }, 50)
    }
  } catch (err) {
    console.error('[OnlyOffice] Load failed:', err)
    loading.value = false
  }
}

watch(() => props.content, () => {
  initViewer()
})

onMounted(() => {
  initViewer()
})

onBeforeUnmount(() => {
  if (docEditor) docEditor.destroyEditor()
  window.removeEventListener('keydown', handleEsc)
  document.body.style.overflow = ''
})
</script>

<style scoped>
.onlyoffice-preview-wrap {
  width: 100%;
  height: 100%;
  background: #f0f2f5;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  position: relative;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

/* 全屏样式 */
.onlyoffice-preview-wrap.is-fullscreen {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  z-index: 3000;
  padding: 0;
}

.fullscreen-toolbar {
  height: 50px;
  background: #fff;
  border-bottom: 1px solid #dcdfe6;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.05);
}

.doc-title {
  font-weight: 600;
  color: #303133;
}

.preview-scroll-area {
  flex: 1;
  overflow-y: auto;
  padding: 30px;
  display: flex;
  justify-content: center;
  scroll-behavior: smooth;
}

.is-fullscreen .preview-scroll-area {
  padding: 20px 0;
}

.preview-container {
  width: 900px; 
  min-height: 1160px; /* 模拟 A4 标准高度 */
  background: #fff;
  border: 1px solid #dcdfe6;
  box-shadow: 0 10px 30px rgba(0,0,0,0.08);
  position: relative;
  overflow: hidden;
  border-radius: 2px;
}

.floating-actions {
  position: absolute;
  top: 20px;
  right: -60px; /* 悬浮在纸张右侧 */
  z-index: 100;
  transition: right 0.3s;
}

.preview-container:hover .floating-actions {
  right: 20px;
}

.viewer-wrapper {
  width: 100%;
  height: 100%;
}

.viewer-instance {
  width: 100%;
  height: 100%;
}

:deep(iframe) {
  border: none !important;
}

.loading-overlay {
  position: absolute;
  top: 0; left: 0; right: 0; bottom: 0;
  background: #fff;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  z-index: 20;
  gap: 12px;
  color: #409eff;
}
</style>
