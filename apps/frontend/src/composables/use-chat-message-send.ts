import { createGlobalState } from '@vueuse/core'
import { ref } from 'vue'

import { useChat } from '@/api/chat.ts'

export const useChatMessageSend = createGlobalState(() => {
  const text = ref('')
  const sendRetries = ref(0)
  const replyTo = ref<string | null>(null)
  const textElement = ref<HTMLTextAreaElement | null>(null)

  const { useSendMessage } = useChat()
  const messageSender = useSendMessage()

  async function sendMessage() {
    if (!text.value) return

    const msg = text.value.replace(/\s+/g, ' ').trim()
    if (!msg) return
    if (msg.length > 700) {
      return
    }

    while (sendRetries.value < 5) {
      try {
        const res = await messageSender.executeMutation({
          opts: {
            text: msg,
            replyTo: replyTo.value
          }
        })
        if (res.error) {
          console.error(res.error)
          sendRetries.value++
          await new Promise((r) => setTimeout(r, 200))
        } else {
          text.value = ''
          replyTo.value = null
          sendRetries.value = 0
          break
        }
      } catch {
        sendRetries.value++
        await new Promise((r) => setTimeout(r, 200))
      }
    }
  }

  return {
    text,
    replyTo,
    textElement,
    sendMessage
  }
})