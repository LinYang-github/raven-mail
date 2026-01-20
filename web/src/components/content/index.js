import { defineAsyncComponent } from 'vue'

const MODE = import.meta.env.VITE_MAIL_CONTENT_MODE || 'text'

export const EditorDriver = defineAsyncComponent(() => {
  switch (MODE) {
    case 'rich':
      return import('./RichTextEditor.vue')
    case 'onlyoffice':
      return import('./OnlyOfficeEditor.vue')
    default:
      return import('./PlainTextEditor.vue')
  }
})

export const getPreviewDriver = (contentType) => {
  const type = contentType || 'text'
  return defineAsyncComponent(() => {
    switch (type) {
      case 'rich':
        return import('./RichTextPreview.vue')
      case 'onlyoffice':
        return import('./OnlyOfficePreview.vue')
      default:
        return import('./PlainTextPreview.vue')
    }
  })
}
