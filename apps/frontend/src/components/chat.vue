<script setup lang="ts">
import { useScroll } from '@vueuse/core'
import { Pause } from 'lucide-vue-next'
import { nextTick, ref, watch } from 'vue'

import { ChatMessage_Fragment, useChat } from '@/api/chat.ts'
import ChatMessageForm from '@/components/chat-message-form.vue'
import ChatMessage from '@/components/chat-message.vue'
import ChatProfile from '@/components/chat-profile.vue'
import ChatViewers from '@/components/chat-viewers.vue'
import ThemeSwitcher from '@/components/theme-switcher.vue'
import { TooltipProvider } from '@/components/ui/tooltip'
import { useFragment } from '@/gql'

const { messages } = useChat()

const messagesEl = ref<HTMLElement | null>(null)
const { y, arrivedState } = useScroll(messagesEl)

const scrollPaused = ref(false)

watch(arrivedState, (v) => {
  scrollPaused.value = !v.bottom
})

watch(
  messages,
  async () => {
    if (!messagesEl.value || scrollPaused.value) return

    await nextTick()
    y.value = messagesEl.value?.scrollHeight
  },
  { immediate: true }
)
</script>

<template>
  <div class="flex h-full max-h-full flex-col">
    <div
      class="flex flex-row justify-between bg-secondary border-b-2 border-red-400 items-center px-4 min-w-48"
    >
      <div class="flex items-center">
        <ChatViewers />
      </div>
      <div class="flex items-center">
        <ThemeSwitcher />
        <ChatProfile />
      </div>
    </div>
    <TooltipProvider
      :delay-duration="150"
      :skip-delay-duration="100"
    >
      <div
        ref="messagesEl"
        class="h-full relative flex flex-col overflow-y-auto px-2"
      >
        <ChatMessage
          v-for="message in messages"
          :key="useFragment(ChatMessage_Fragment, message).id"
          :msg="message"
        />
      </div>
      <div
        v-if="scrollPaused"
        class="sticky w-full bottom-0 bg-zinc-700 place-self-center flex items-center justify-center"
      >
        <Pause />
        <span class="text-xl">Scroll paused</span>
      </div>
    </TooltipProvider>
    <ChatMessageForm />
  </div>
</template>

<style scoped></style>
