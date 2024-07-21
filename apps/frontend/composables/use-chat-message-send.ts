import { createGlobalState } from '@vueuse/core'
import { ref } from 'vue'

import { useChat } from '@/api/chat.js'

export const useChatMessageSend = createGlobalState(() => {
  const text = ref('')
  const sendRetries = ref(0)
  const replyTo = ref<string | null>(null)
  const textElement = ref<HTMLTextAreaElement | null>(null)
  const isSending = ref(false)
  const mentionNickname = ref<null | string>(null)

  const { useSendMessage, channelData } = useChat()
  const messageSender = useSendMessage()

  async function sendMessage() {
    if (!text.value || !channelData.value?.fetchUserByName || isSending.value) {
      return
    }
    isSending.value = true

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
            replyTo: replyTo.value,
            channelId: channelData.value.fetchUserByName.id
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

    isSending.value = false
    // hack for focus textarea after sending message
    setTimeout(() => {
      textElement.value?.focus()
    }, 100)
  }

  function setMentionNickName(nick: string | null) {
    mentionNickname.value = nick
  }

  return {
    text,
    replyTo,
    textElement,
    sendMessage,
    isSending,
    mentionNickname,
    setMentionNickName
  }
})
