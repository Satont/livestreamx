import { createGlobalState } from '@vueuse/core'
import { ref } from 'vue'

export const useProfileModalState = createGlobalState(() => {
  const opened = ref(false)

  return {
    opened
  }
})
