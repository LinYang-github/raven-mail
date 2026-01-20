<template>
  <el-input
    v-model="content"
    type="textarea"
    :rows="15"
    placeholder="请输入正文内容..."
    @input="handleInput"
  />
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
:deep(.el-textarea__inner) {
  height: 100%;
  border: none;
  border-radius: 0;
}
:deep(.el-textarea) {
  height: 100%;
}
</style>
