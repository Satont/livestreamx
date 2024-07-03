<script setup lang="ts">
import { useScroll } from '@vueuse/core'
import { Pause } from 'lucide-vue-next'
import { computed, nextTick, ref, watch } from 'vue'

import {
  ChatMessage_Fragment,
  SystemMessage_Fragment,
  useChat
} from '@/api/chat.ts'
import ChatMessageForm from '@/components/chat-message-form.vue'
import ChatMessage from '@/components/chat-message.vue'
import ChatSystemMessage from '@/components/chat-system-message.vue'
import StreamUptime from '@/components/stream-uptime.vue'
import StreamViewers from '@/components/stream-viewers.vue'
import { TooltipProvider } from '@/components/ui/tooltip'
import { useChatMessageSend } from '@/composables/use-chat-message-send.ts'
import { FragmentType, useFragment } from '@/gql'

const { messages, systemMessages } = useChat()
const unwrappedMessages = computed(() =>
  useFragment(ChatMessage_Fragment, messages.value)
)

const allMessages = computed(() => {
  return [
    ...unwrappedMessages.value,
    ...systemMessages.value
  ].sort((a, b) => {
    if (!('createdAt' in a) || !('createdAt' in b)) return 0

    return new Date(a.createdAt).getTime() - new Date(b.createdAt).getTime()
  })
})

const { replyTo } = useChatMessageSend()

const messagesEl = ref<HTMLElement | null>(null)
const { y, arrivedState } = useScroll(messagesEl)

const scrollPaused = ref(false)

watch(arrivedState, (v) => {
  scrollPaused.value = !v.bottom
})

watch(
  allMessages,
  async () => {
    if (!messagesEl.value || scrollPaused.value) return

    await scrollToBottom()
  },
  { immediate: true }
)

async function scrollToBottom() {
  await nextTick()
  if (!messagesEl.value) return

  y.value = messagesEl.value.scrollHeight
}

const replyingTo = computed(() => {
  if (!replyTo.value) return null

  return unwrappedMessages.value.find((m) => m.id === replyTo.value) as
    | FragmentType<typeof ChatMessage_Fragment>
    | undefined
})
</script>

<template>
  <div
    class="relative flex h-full max-h-full flex-col lg:border-l-2 border-t-2 lg:border-t-0 border-border"
  >
    <div
      class="flex flex-row justify-between bg-[#111111] border-b-2 border-border items-center px-4 py-2 min-w-48"
    >
      <StreamUptime class="text-md font-semibold" />
      <StreamViewers />
    </div>
    <TooltipProvider
      :delay-duration="150"
      :skip-delay-duration="100"
    >
      <div
        ref="messagesEl"
        class="h-full relative flex flex-col overflow-y-auto px-2 bg-[#111111]"
      >
        <template v-for="message in allMessages">
          <ChatMessage
            v-if="!('type' in message)"
            :msg="message as FragmentType<typeof ChatMessage_Fragment>"
            @reply="scrollToBottom"
          />
          <ChatSystemMessage
            v-else
            :msg="message as FragmentType<typeof SystemMessage_Fragment>"
          />
        </template>
      </div>
      <div
        v-if="scrollPaused || replyingTo"
        class="sticky w-full bottom-0"
        @click="scrollToBottom"
      >
        <div
          v-if="scrollPaused"
          class="bg-accent place-self-center flex items-center justify-center cursor-pointer"
        >
          <Pause />
          <span class="text-xl">Scroll paused</span>
        </div>
        <div
          v-if="replyingTo"
          class="bg-accent"
        >
          <ChatMessage
            :msg="replyingTo"
            is-reply
          />
        </div>
      </div>
    </TooltipProvider>
    <ChatMessageForm />
  </div>
</template>

<style scoped></style>
