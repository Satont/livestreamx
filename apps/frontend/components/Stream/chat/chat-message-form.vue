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

const { textElement, isSending, mentionNickname } = useChatMessageSend()
const { data: profile } = useProfile().useData()

onMounted(() => {
  textElement.value = document.getElementById(
    'chat-messages-form-textarea'
  ) as HTMLTextAreaElement
})

const { text, sendMessage } = useChatMessageSend()

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

watch(mentionNickname, (v) => {
  if (!v) return

  text.value += `@${v}`
})

const showSuggestions = ref(false)
const suggestionInput = ref('')
const filteredSuggestions = computed(() => {
  if (!text.value) return emotesMenuEmotes.value

  const mentionText = text.value.split(' ').pop()!.substring(1)
  return emotesMenuEmotes.value.filter((e) =>
    e.name.toLowerCase().includes(mentionText.toLowerCase())
  )
})

const selectedEmoteIndex = ref(0)
watch(filteredSuggestions, () => {
  selectedEmoteIndex.value = 0
})
const emoteMentionRegexp = /:\w*$/

async function handleInput(e: InputEvent) {
  const value = (e.target as HTMLTextAreaElement)!.value
  const mentionMatch = value.match(emoteMentionRegexp)

  if (mentionMatch) {
    showSuggestions.value = true
    suggestionInput.value = mentionMatch[0]
  } else {
    showSuggestions.value = false
    await nextTick()
    textElement.value?.focus()
  }
}

const virtualListRef = ref<any>(null)

function handleEmoteClick(index: number) {
  selectedEmoteIndex.value = index
  insertEmote()
}

async function insertEmote() {
  const emoteMentionPosition = text.value.lastIndexOf(suggestionInput.value)
  text.value =
    text.value.slice(0, emoteMentionPosition) +
    filteredSuggestions.value[selectedEmoteIndex.value].name

  showSuggestions.value = false
  await nextTick()
  textElement.value?.focus()
}

async function handleKeyDown(e: KeyboardEvent) {
  if (e.key === 'ArrowDown') {
    e.preventDefault()
    if (selectedEmoteIndex.value < filteredSuggestions.value.length - 1) {
      selectedEmoteIndex.value++
      virtualListRef.value.scrollTo(selectedEmoteIndex.value)
    }
  } else if (e.key === 'ArrowUp') {
    e.preventDefault()

    if (showSuggestions.value) {
      if (selectedEmoteIndex.value > 0) {
        selectedEmoteIndex.value--
        virtualListRef.value.scrollTo(selectedEmoteIndex.value)
      }
    } else if (!text.value.length) {
      insertLatestMessage()
    }
  } else if (e.key === 'Enter') {
    if (showSuggestions.value) {
      e.preventDefault()

      await insertEmote()
    } else {
      await sendMessage()
    }
  }
}
</script>

<template>
  <div class="flex flex-col gap-2.5 dark:bg-[#111111] p-2.5 relative">
    <div>
      <UiPopover :open="showSuggestions">
        <UiPopoverContent class="w-96 p-0">
          <UiCommand>
            <UiCommandEmpty>No emote found.</UiCommandEmpty>
            <UiCommandList class="overflow-hidden">
              <UseVirtualList
                ref="virtualListRef"
                :list="filteredSuggestions"
                :options="{
                  itemHeight: 44
                }"
                height="300px"
              >
                <template #default="props">
                  <UiCommandItem
                    class="flex items-center justify-between pr-4"
                    :value="props.data.name"
                    :class="{
                      ['bg-accent text-accent-foreground']:
                        selectedEmoteIndex === props.index
                    }"
                    @click="handleEmoteClick(props.index)"
                    disable-highlight
                  >
                    <span>{{ props.data.name }}</span>
                    <img
                      :src="props.data.url"
                      class="h-8 max-w-12"
                    />
                  </UiCommandItem>
                </template>
              </UseVirtualList>
            </UiCommandList>
          </UiCommand>
        </UiPopoverContent>
        <UiPopoverAnchor>
          <UiTextarea
            id="chat-messages-form-textarea"
            v-model="text"
            placeholder="Send a message"
            @paste="console.log"
            class="pr-12 min-h-8 max-h-20 resize-none bg-[#181818] focus-visible:ring-offset-0 focus-visible:ring-[#4D4D4D] transition-[box-shadow,border-color,background-color] hover:border-white/20 border-white/15 px-3 text-white dark:placeholder:text-white/50 rounded-md focus-visible:border-white/15 focus-visible:bg-[#111111]"
            :rows="isSmall ? 1 : 3"
            @keyup="updateCarretPosition"
            @input="handleInput"
            @click="updateCarretPosition"
            @focus="updateCarretPosition"
            @keydown="handleKeyDown"
            :disabled="chatLocked"
            :maxlength="700"
          />
        </UiPopoverAnchor>
      </UiPopover>
    </div>

    <UiButton
      class="absolute right-4 top-4"
      size="xs"
      variant="ghost"
      @click="
        () => {
          showSuggestions = !showSuggestions
        }
      "
    >
      <svg
        xmlns="http://www.w3.org/2000/svg"
        width="32"
        height="32"
        viewBox="0 0 24 24"
        class="size-5 text-stone-300/80"
      >
        <path
          fill="currentColor"
          d="M17 3.34a10 10 0 1 1-14.995 8.984L2 12l.005-.324A10 10 0 0 1 17 3.34m-1.8 10.946a1 1 0 0 0-1.414.014a2.5 2.5 0 0 1-3.572 0a1 1 0 0 0-1.428 1.4a4.5 4.5 0 0 0 6.428 0a1 1 0 0 0-.014-1.414M9.01 9l-.127.007A1 1 0 0 0 9 11l.127-.007A1 1 0 0 0 9.01 9m6 0l-.127.007A1 1 0 0 0 15 11l.127-.007A1 1 0 0 0 15.01 9"
        />
      </svg>
    </UiButton>

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
