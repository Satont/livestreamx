import { useLocalStorage } from '@vueuse/core'

export const reverseStreamChatDirection = useLocalStorage(
  'livestreamx-reverse-chat-stream-direction',
  false
)
