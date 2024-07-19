import { createGlobalState, useLocalStorage } from '@vueuse/core'

export const useShowReactionsOnMessage = createGlobalState(() => {
  const showReactionsOnMessage = useLocalStorage(
    'livestreamx-show-reactions-on-message',
    true
  )

  return {
    showReactionsOnMessage
  }
})
