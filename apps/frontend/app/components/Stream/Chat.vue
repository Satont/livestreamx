<script setup lang="ts">
import { useScroll } from '@vueuse/core'
import { computed, nextTick, ref, watch } from 'vue'
import type { FragmentType } from '@/gql'

import {
  ChatMessage_Fragment,
  SystemMessage_Fragment,
  useChat
} from '@/api/chat.js'
import { useFragment } from '@/gql'
import ChatMessageForm from '~/components/Stream/chat/chat-message-form.vue'
import ChatMessage from '~/components/Stream/chat/chat-message.vue'
import ChatSystemMessage from '~/components/Stream/chat/chat-system-message.vue'
import ReactionsModal from '~/components/Stream/chat/reactions-modal.vue'
import StreamUptime from '~/components/Stream/chat/stream-uptime.vue'
import StreamViewers from '~/components/Stream/chat/stream-viewers.vue'
import { useChatMessageSend } from '~/composables/use-chat-message-send.js'

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
    class="relative flex h-full max-h-full flex-col lg:border-l-2 border-t-2 lg:border-t-0 border-border text-accent-foreground dark:bg-[#111111]"
  >
    <div
      class="flex flex-row justify-between border-b-2 border-border items-center px-4 py-2 min-w-48"
    >
      <StreamUptime class="text-md font-semibold" />
      <StreamViewers />
    </div>
    <UiTooltipProvider
      :delay-duration="150"
      :skip-delay-duration="100"
    >
      <div
        ref="messagesEl"
        class="h-full relative flex flex-col overflow-y-auto px-2 dark:bg-[#111111]"
        :style="{ fontSize: `${chatFontSize}px` }"
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
          <Icon name="lucide:pause" />
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
    </UiTooltipProvider>
    <ChatMessageForm />
  </div>

  <ReactionsModal />
</template>
