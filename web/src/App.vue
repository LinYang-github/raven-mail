<template>
  <div class="app-container">
    <MailClient v-if="showMail" />
    <ChatWidget v-if="showIM" />
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import MailClient from './components/MailClient.vue'
import ChatWidget from './components/ChatWidget.vue'

const props = defineProps({
  modules: { type: [Array, String], default: () => ['mail', 'im'] }
})

const isEnabled = (moduleName) => {
  if (Array.isArray(props.modules)) {
    return props.modules.includes(moduleName)
  }
  return props.modules === moduleName
}

const showMail = computed(() => isEnabled('mail'))
const showIM = computed(() => isEnabled('im'))
</script>

<style>
/* Global reset */
html, body, #app {
  margin: 0;
  padding: 0;
  height: 100%;
  width: 100%;
  overflow: hidden;
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
}

.app-container {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  display: flex;
  background: #f5f7fa;
  overflow: hidden;
}
</style>
