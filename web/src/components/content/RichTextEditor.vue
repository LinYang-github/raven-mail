<template>
  <div class="rich-editor-placeholder">
    <div class="toolbar">
      <el-button-group>
        <el-button size="small">B</el-button>
        <el-button size="small">I</el-button>
        <el-button size="small">U</el-button>
      </el-button-group>
      <span class="mode-tag">富文本模式 (Simulated)</span>
    </div>
    <el-input
      v-model="content"
      type="textarea"
      :rows="14"
      placeholder="此处应集成富文本编辑器 (如 Tiptap/WangEditor)"
      @input="handleInput"
    />
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'

const props = defineProps(['modelValue'])
const emit = defineEmits(['update:modelValue'])

const content = ref(props.modelValue)

watch(() => props.modelValue, (val) => {
  content.value = val
})

const handleInput = () => {
  emit('update:modelValue', content.value)
}
</script>

<style scoped>
.rich-editor-placeholder {
  border: 1px solid #dcdfe6;
  border-radius: 4px;
}
.toolbar {
  padding: 8px;
  background: #f5f7fa;
  border-bottom: 1px solid #dcdfe6;
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.mode-tag {
  font-size: 12px;
  color: #909399;
}
</style>
