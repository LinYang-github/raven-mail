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

import { userStore } from './store/user'

const props = defineProps({
  modules: { type: [Array, String], default: () => ['mail', 'im'] }
})

// Use userStore.modules which is reactive and updated via GlobalState
const isEnabled = (moduleName) => {
  const currentModules = userStore.modules || props.modules
  if (Array.isArray(currentModules)) {
    return currentModules.includes(moduleName)
  }
  return currentModules === moduleName
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
