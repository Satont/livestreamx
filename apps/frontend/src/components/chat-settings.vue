<script setup lang="ts">
import { Cog } from 'lucide-vue-next'
import { computed } from 'vue'
import { toast } from 'vue-sonner'

import { ChatMessage_Fragment, useChat } from '@/api/chat.ts'
import { useProfile } from '@/api/profile.ts'
import { Button } from '@/components/ui/button'
import {
  NumberField,
  NumberFieldContent,
  NumberFieldDecrement,
  NumberFieldIncrement,
  NumberFieldInput
} from '@/components/ui/number-field'
import {
  Popover,
  PopoverContent,
  PopoverTrigger
} from '@/components/ui/popover'
import { Separator } from '@/components/ui/separator'
import { Slider } from '@/components/ui/slider'
import { Switch } from '@/components/ui/switch'
import { chatFontSize } from '@/composables/chat-font-size.ts'
import { reverseStreamChatDirection } from '@/composables/reverse-stream-chat-direction.ts'
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

  toast.success('Color updated', {
    description: 'Your nickname color has been updated',
    dismissible: true
  })
}

function focusColorPicker() {
  const colorPicker = document.getElementById('user-profile-color-picker')
  colorPicker?.click()
}
</script>

<template>
  <Popover>
    <PopoverTrigger asChild>
      <Button
        size="sm"
        variant="secondary"
        class="dark:text-stone-300/80 p-1.5 h-8"
      >
        <Cog class="size-5" />
      </Button>
    </PopoverTrigger>
    <PopoverContent class="p-2 w-80">
      <div class="w-full flex flex-col">
        <Button
          @click="showAvatars = !showAvatars"
          size="sm"
          class="flex gap-2 justify-between"
          variant="ghost"
        >
          <span> Show avatars </span>
          <Switch
            :checked="showAvatars"
            class="data-[state=unchecked]:bg-zinc-600"
          />
        </Button>
        <Button
          @click="showTimestamps = !showTimestamps"
          size="sm"
          class="flex gap-2 justify-between"
          variant="ghost"
        >
          <span> Show time </span>
          <Switch
            :checked="showTimestamps"
            class="data-[state=unchecked]:bg-zinc-600"
          />
        </Button>

        <Button
          @click="showReactionsOnMessage = !showReactionsOnMessage"
          size="sm"
          class="flex gap-2 justify-between"
          variant="ghost"
        >
          <span>Show reactions on message</span>
          <Switch
            :checked="showReactionsOnMessage"
            class="data-[state=unchecked]:bg-zinc-600"
          />
        </Button>

        <Button
          @click="reverseStreamChatDirection = !reverseStreamChatDirection"
          size="sm"
          class="flex gap-2 justify-between"
          variant="ghost"
        >
          <span> Swap chat position left/right </span>
          <Switch
            :checked="reverseStreamChatDirection"
            class="data-[state=unchecked]:bg-zinc-600"
          />
        </Button>

        <Button
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
        </Button>
        <Separator class="my-4" />
        <div class="flex flex-col gap-2 px-3">
          <div class="flex justify-between items-center">
            <h1 class="text-lg text-slate-900 font-medium dark:text-slate-200">
              Font size
            </h1>
            <NumberField
              :min
              :max
              v-model="chatFontSize"
              class="w-36"
            >
              <NumberFieldContent>
                <NumberFieldDecrement />
                <NumberFieldInput />
                <NumberFieldIncrement />
              </NumberFieldContent>
            </NumberField>
          </div>
          <Slider
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
    </PopoverContent>
  </Popover>

  <!--	<DropdownMenu>-->
  <!--		<DropdownMenuTrigger asChild>-->
  <!--			<Button size="sm" variant="ghost"><Settings /></Button>-->
  <!--		</DropdownMenuTrigger>-->
  <!--		<DropdownMenuContent>-->
  <!--			<DropdownMenuCheckboxItem-->
  <!--				v-model:checked="showAvatars"-->
  <!--			>-->
  <!--				Show Avatars-->
  <!--			</DropdownMenuCheckboxItem>-->
  <!--			<DropdownMenuCheckboxItem-->
  <!--				v-model:checked="showTimestamps"-->
  <!--			>-->
  <!--				Show Time-->
  <!--			</DropdownMenuCheckboxItem>-->
  <!--			<DropdownMenuItem>-->
  <!--				<Slider-->
  <!--					v-model="fontSize"-->
  <!--					:max="50"-->
  <!--					:step="1"-->
  <!--					:class="cn('w-3/5', $attrs.class ?? '')"-->
  <!--				/>-->
  <!--			</DropdownMenuItem>-->
  <!--		</DropdownMenuContent>-->
  <!--	</DropdownMenu>-->
</template>
