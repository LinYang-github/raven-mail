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

const route = useRoute()

// Determine capabilities based on modules prop OR route
const props = defineProps({
  modules: { type: [Array, String], default: () => ['mail', 'im'] }
})

const showMail = computed(() => {
  return route.path.startsWith('/mail')
})

const showIM = computed(() => {
  // IM is always enabled if in 'im' module OR if route is /mail (floating widget)
  // Or simply: enable if route is /im OR /mail
  return route.path.startsWith('/im') || route.path.startsWith('/mail') 
})
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
