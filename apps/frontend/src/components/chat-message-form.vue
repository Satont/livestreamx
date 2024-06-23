<script setup lang="ts">
import { UseVirtualList } from '@vueuse/components'
import { Smile } from 'lucide-vue-next'
import { computed, onMounted, ref } from 'vue'

import { ChatMessage_Fragment, useChat } from '@/api/chat.ts'
import { useProfile } from '@/api/profile.js'
import ChatSettings from '@/components/chat-settings.vue'
import { Button } from '@/components/ui/button'
import {
  Command,
  CommandEmpty,
  CommandInput,
  CommandItem,
  CommandList
} from '@/components/ui/command'
import {
  Popover,
  PopoverContent,
  PopoverTrigger
} from '@/components/ui/popover'
import { Textarea } from '@/components/ui/textarea'
import { useChatMessageSend } from '@/composables/use-chat-message-send.ts'
import { useFragment } from '@/gql'
import Mention from './mention.vue'

const { emotes, messages } = useChat()
const unwrappedMessages = computed(() =>
  useFragment(ChatMessage_Fragment, messages.value)
)

const { data: profile } = useProfile()

const textElement = ref<HTMLTextAreaElement | null>(null)
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
</script>

<template>
  <div
    class="flex flex-col gap-2 bg-accent border-t-2 border-red-400 min-h-36 p-2 relative"
  >
    <Mention
      :keys="mentionKeys"
      :items="mentionItems"
      offset="6"
      insert-space
      @open="(k: typeof mentionKey) => (mentionKey = k)"
      :omit-key="mentionKey !== '@'"
      :item-height="mentionKey === '@' ? 22 : 48"
      :map-insert="mapInsert"
    >
      <Textarea
        id="chat-messages-form-textarea"
        v-model="text"
        placeholder="Send message..."
        @keydown.enter="sendMessage"
        @paste="console.log"
        class="pr-12"
        @keyup="updateCarretPosition"
        @click="updateCarretPosition"
        @focus="updateCarretPosition"
        :disabled="!profile"
        :maxlength="700"
      />

      <Popover v-model:open="emoteMenuOpened">
        <PopoverTrigger as-child>
          <Button
            class="absolute right-2 top-2"
            variant="ghost"
            size="sm"
          >
            <Smile />
          </Button>
        </PopoverTrigger>
        <PopoverContent class="h-[340px] w-96 p-0 mb-2">
          <Command v-model:search-term="emoteMenuSearchTerm">
            <CommandInput
              class="h-9"
              placeholder="Search emote..."
            />
            <CommandEmpty>No emote found.</CommandEmpty>
            <CommandList class="overflow-hidden">
              <UseVirtualList
                :list="emotesMenuEmotes"
                :options="{
                  itemHeight: 44
                }"
                height="340px"
              >
                <template #default="props">
                  <CommandItem
                    class="flex items-center justify-between pr-4"
                    :value="props.data.name"
                    @select="(e) => insertEmoteInText(e.detail.value)"
                  >
                    <span>{{ props.data.name }}</span>
                    <img
                      :src="props.data.url"
                      class="h-8 max-w-12"
                    />
                  </CommandItem>
                </template>
              </UseVirtualList>
            </CommandList>
          </Command>
        </PopoverContent>
      </Popover>

      <template #no-result>
        <div class="dim">No result</div>
      </template>

      <template #item-@="{ item }">
        <Button
          class="w-full"
          variant="ghost"
          size="sm"
        >
          {{ item.data.label }}
        </Button>
      </template>

      <template #item-:="{ item }">
        <div class="flex items-center gap-2 cursor-pointer">
          <img
            :src="item.data.url"
            class="size-10"
          />
          <span>{{ item.data.label }}</span>
        </div>
      </template>
    </Mention>

    <div class="flex gap-2 place-self-end">
      <ChatSettings />
      <Button
        @click="sendMessage"
        size="sm"
        :disabled="!profile"
        >Send message</Button
      >
    </div>
  </div>
</template>
