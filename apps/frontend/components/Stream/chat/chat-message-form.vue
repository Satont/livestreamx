<script setup lang="ts">
import { UseVirtualList } from '@vueuse/components'
import { breakpointsTailwind, useBreakpoints } from '@vueuse/core'
import { computed, onMounted, ref } from 'vue'

import { ChatMessage_Fragment, useChat } from '@/api/chat.js'
import { useProfile } from '@/api/profile.js'
import { useChatMessageSend } from '@/composables/use-chat-message-send.js'
import { useFragment } from '@/gql'
import ChatSettings from '~/components/Stream/chat/chat-settings.vue'

const { emotes, messages } = useChat()
const unwrappedMessages = computed(() =>
  useFragment(ChatMessage_Fragment, messages.value)
)

const { textElement, isSending } = useChatMessageSend()
const { data: profile } = useProfile().useData()

onMounted(() => {
  textElement.value = document.getElementById(
    'chat-messages-form-textarea'
  ) as HTMLTextAreaElement
})

const { text, sendMessage } = useChatMessageSend()

const emotesForMention = computed(() => {
  return emotes.value.map((e) => ({
    label: e.name,
    url: e.url,
    value: e.name
  }))
})

const usersForMention = computed(() => {
  const mappedUsersFromMessages = unwrappedMessages.value
    .map((m) => ({
      label: m.sender.displayName,
      color: m.sender.color,
      value: m.sender.displayName
    }))
    .filter((v, i, a) => a.findIndex((t) => t.label === v.label) === i)

  return [
    ...mappedUsersFromMessages
  ]
})

const mentionKey = ref<'@' | ':'>('@')
function mapInsert(item: { label: string }) {
  return item.label + ' '
}
const mentionKeys = ['@', ':']
const mentionItems = computed(() => {
  return mentionKey.value === '@'
    ? usersForMention.value
    : emotesForMention.value
})

const emoteMenuOpened = ref(false)
const emoteMenuSearchTerm = ref('')
const emotesMenuEmotes = computed(() => {
  if (!emoteMenuSearchTerm.value) return emotes.value

  return emotes.value.filter((e) =>
    e.name.toLowerCase().includes(emoteMenuSearchTerm.value.toLowerCase())
  )
})
const currentCarretPosition = ref(0)
function updateCarretPosition(e: KeyboardEvent | MouseEvent) {
  if (!e.target) return
  if (e.type === 'focus') return

  const target = e.target as HTMLTextAreaElement
  currentCarretPosition.value = target.selectionStart
}
function insertEmoteInText(value: unknown) {
  if (typeof value !== 'string') return

  const pos = currentCarretPosition.value

  const newText = text.value.slice(0, pos) + value + ' ' + text.value.slice(pos)
  text.value = newText
  emoteMenuOpened.value = false
  textElement.value?.focus()
  currentCarretPosition.value = pos + value.length + 1
}

const breakPoints = useBreakpoints(breakpointsTailwind)
const isSmall = breakPoints.smallerOrEqual('lg')

const chatLocked = computed(() => {
  return !profile.value || isSending.value
})

function insertLatestMessage() {
  const latestMessageFromCurrentUser = unwrappedMessages.value
    .slice()
    .reverse()
    .find((m) => m.sender.id === profile.value?.userProfile.id)

  if (!latestMessageFromCurrentUser) return

  text.value = latestMessageFromCurrentUser.segments
    .map((s) => s.content)
    .join(' ')
}
</script>

<template>
  <div class="flex flex-col gap-2.5 dark:bg-[#111111] p-2.5 relative">
    <UiTextarea
      id="chat-messages-form-textarea"
      v-model="text"
      placeholder="Send a message"
      @keydown.enter="sendMessage"
      @paste="console.log"
      class="pr-12 min-h-8 max-h-20 resize-none bg-[#181818] focus-visible:ring-offset-0 focus-visible:ring-[#4D4D4D] transition-[box-shadow,border-color,background-color] hover:border-white/20 border-white/15 px-3 text-white dark:placeholder:text-white/50 rounded-md focus-visible:border-white/15 focus-visible:bg-[#111111]"
      :rows="isSmall ? 1 : 3"
      @keyup="updateCarretPosition"
      @click="updateCarretPosition"
      @focus="updateCarretPosition"
      @keydown.up="insertLatestMessage"
      :disabled="chatLocked"
      :maxlength="700"
    />

    <div class="flex gap-2 place-self-end">
      <ChatSettings />
      <UiButton
        @click="sendMessage"
        size="sm"
        :disabled="chatLocked"
        class="bg-blue-500 hover:bg-blue-600 text-white h-8 px-3"
      >
        Send
      </UiButton>
    </div>
  </div>
</template>
