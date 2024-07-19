<script setup lang="ts">
import { computed } from 'vue'

import { ChatMessage_Fragment, useChat } from '@/api/chat.js'
import { useProfile } from '@/api/profile.js'
import { chatFontSize } from '@/composables/chat-font-size.js'
import { reverseStreamChatDirection } from '@/composables/reverse-stream-chat-direction.js'
import { showAvatars } from '@/composables/show-avatars.js'
import { showTimestamps } from '@/composables/show-timestamps.js'
import { useShowReactionsOnMessage } from '@/composables/use-show-reactions-on-message.js'
import { useFragment } from '@/gql'

const fontSize = computed({
  get() {
    return [chatFontSize.value]
  },
  set(values: number[]) {
    chatFontSize.value = values.at(0)
  }
})

const min = 10
const max = 50

const { useData, useUpdateMutation } = useProfile()
const { data: profile, executeQuery: refetchProfile } = useData()
const updateUser = useUpdateMutation()
const { showReactionsOnMessage } = useShowReactionsOnMessage()

const { messages } = useChat()

async function handleColorChange(e: Event) {
  const newValue = (e.target as HTMLInputElement).value

  await updateUser.executeMutation({
    input: {
      color: newValue
    }
  })
  await refetchProfile({ requestPolicy: 'network-only' })

  // update colors in chat for that user
  for (const message of messages.value) {
    const msgFragment = useFragment(ChatMessage_Fragment, message)
    msgFragment.sender.color = newValue
  }

  useToast().toast({
    title: 'Color updated',
    description: 'Your nickname color has been updated',
    variant: 'success'
  })
}

function focusColorPicker() {
  const colorPicker = document.getElementById('user-profile-color-picker')
  colorPicker?.click()
}
</script>

<template>
  <UiPopover>
    <UiPopoverTrigger asChild>
      <UiButton
        size="sm"
        variant="secondary"
        class="dark:text-stone-300/80 p-1.5 h-8"
      >
        <Icon
          name="lucide:cog"
          class="size-5"
        />
      </UiButton>
    </UiPopoverTrigger>
    <UiPopoverContent class="p-2 w-80">
      <div class="w-full flex flex-col">
        <UiButton
          @click="showAvatars = !showAvatars"
          size="sm"
          class="flex gap-2 justify-between"
          variant="ghost"
        >
          <span> Show avatars </span>
          <UiSwitch
            :checked="showAvatars"
            class="data-[state=unchecked]:bg-zinc-600"
          />
        </UiButton>
        <UiButton
          @click="showTimestamps = !showTimestamps"
          size="sm"
          class="flex gap-2 justify-between"
          variant="ghost"
        >
          <span> Show time </span>
          <UiSwitch
            :checked="showTimestamps"
            class="data-[state=unchecked]:bg-zinc-600"
          />
        </UiButton>

        <UiButton
          @click="showReactionsOnMessage = !showReactionsOnMessage"
          size="sm"
          class="flex gap-2 justify-between"
          variant="ghost"
        >
          <span>Show reactions on message</span>
          <UiSwitch
            :checked="showReactionsOnMessage"
            class="data-[state=unchecked]:bg-zinc-600"
          />
        </UiButton>

        <UiButton
          @click="reverseStreamChatDirection = !reverseStreamChatDirection"
          size="sm"
          class="flex gap-2 justify-between"
          variant="ghost"
        >
          <span> Swap chat position left/right </span>
          <UiSwitch
            :checked="reverseStreamChatDirection"
            class="data-[state=unchecked]:bg-zinc-600"
          />
        </UiButton>

        <UiButton
          size="sm"
          variant="ghost"
          class="w-full flex justify-between"
          @click="focusColorPicker"
          :disabled="!profile"
        >
          <span>Change nickname color</span>
          <input
            id="user-profile-color-picker"
            type="color"
            :value="profile?.userProfile.color"
            class="size-6"
            :disabled="!profile"
            @change="handleColorChange"
          />
        </UiButton>
        <UiSeparator class="my-4" />
        <div class="flex flex-col gap-2 px-3">
          <div class="flex justify-between items-center">
            <h1 class="text-lg text-slate-900 font-medium dark:text-slate-200">
              Font size
            </h1>
            <UiNumberField
              :min
              :max
              v-model="chatFontSize"
              class="w-36"
            />
          </div>
          <UiSlider
            v-model="fontSize"
            :min
            :max
            :step="1"
            class="w-full"
          />
          <div class="flex justify-between">
            <span>{{ min }}px</span>
            <span>{{ max }}px</span>
          </div>
        </div>
      </div>
    </UiPopoverContent>
  </UiPopover>
</template>
