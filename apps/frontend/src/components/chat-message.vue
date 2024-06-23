<script setup lang="ts">
import { Copy } from 'lucide-vue-next'

import { ChatMessage_Fragment } from '@/api/chat.ts'
import { useProfile } from '@/api/profile.ts'
import ChatMessageReactions from '@/components/chat-message-reactions.vue'
import { Button } from '@/components/ui/button'
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger
} from '@/components/ui/tooltip'
import { chatFontSize } from '@/composables/chat-font-size.js'
import { colorMode } from '@/composables/color-mode.ts'
import { showAvatars } from '@/composables/show-avatars.js'
import { showTimestamps } from '@/composables/show-timestamps.js'
import { FragmentType, useFragment } from '@/gql'
import {
  ChatEmote_FragmentFragment,
  MessageSegmentType
} from '@/gql/graphql.ts'
import { calculateColor } from '@/lib/color.js'

type Props = {
  msg: FragmentType<typeof ChatMessage_Fragment>
}
const props = defineProps<Props>()
const unwrappedMessage = useFragment(ChatMessage_Fragment, props.msg)

const { data: profile } = useProfile()

function correctColor(color: string) {
  return calculateColor(color, colorMode.value === 'dark')
}

function copyText() {
  navigator.clipboard.writeText(
    unwrappedMessage.segments.map((s) => s.content).join(' ')
  )
}
</script>

<template>
  <div
    :style="{ fontSize: `${chatFontSize}px` }"
    class="relative group hover:bg-accent hover:rounded p-0.5"
  >
    <p class="leading-7">
      <span
        v-if="showTimestamps"
        class="mr-1 opacity-50"
      >
        {{
          new Date(unwrappedMessage.createdAt).toLocaleTimeString('en', {
            hour12: false,
            hour: '2-digit',
            minute: '2-digit'
          })
        }}
      </span>
      <span>
        <span
          class="inline-flex align-sub"
          v-if="showAvatars"
        >
          <img
            :src="unwrappedMessage.sender.avatarUrl"
            class="size-4 rounded-full mr-1"
          />
        </span>
        <span
          class="font-bold"
          :style="{ color: correctColor(unwrappedMessage.sender.color) }"
        >
          {{ unwrappedMessage.sender.displayName }}
        </span>
      </span>
      <span>: </span>
      <span class="break-words">
        <template v-for="segment of unwrappedMessage.segments">
          <template v-if="segment.type === MessageSegmentType.Text">{{
            segment.content
          }}</template>
          <span
            v-else-if="
              segment.type === MessageSegmentType.Mention && 'user' in segment
            "
            :style="{ color: correctColor(segment.user.color) }"
            class="p-0.5 rounded"
            :class="{
              'bg-zinc-400': segment.user.id === profile?.userProfile.id
            }"
          >
            @{{ segment.user.displayName }}
          </span>
          <a
            v-else-if="segment.type === MessageSegmentType.Link"
            :href="segment.content"
            target="_blank"
            class="underline"
          >
            {{ segment.content }}
          </a>
          <template
            v-else-if="
              segment.type === MessageSegmentType.Emote && 'emote' in segment
            "
          >
            <Tooltip>
              <TooltipTrigger>
                <img
                  :src="(segment.emote as ChatEmote_FragmentFragment).url"
                  :style="{
                    width: `${(segment.emote as ChatEmote_FragmentFragment).width}px`,
                    height: `${(segment.emote as ChatEmote_FragmentFragment).height}px`
                  }"
                  class="scale-90 inline-block relative"
                />
              </TooltipTrigger>
              <TooltipContent>
                <div class="flex flex-col">
                  <img
                    :src="
                      (segment.emote as ChatEmote_FragmentFragment).url.replace(
                        '1x.webp',
                        '4x.webp'
                      )
                    "
                    :style="{
                      width: `${(segment.emote as ChatEmote_FragmentFragment).width * 2.5}px`,
                      height: `${(segment.emote as ChatEmote_FragmentFragment).height * 2.5}px`
                    }"
                  />
                  <h1 class="place-self-center text-lg font-bold">
                    {{ (segment.emote as ChatEmote_FragmentFragment).name }}
                  </h1>
                </div>
              </TooltipContent>
            </Tooltip>
          </template>
          {{ ' ' }}
        </template>
      </span>
    </p>

    <div class="absolute right-0 top-[-10px] group">
      <div class="flex gap-2">
        <ChatMessageReactions :msg="msg" />
        <div class="hidden group-hover:block">
          <Button
            @click="copyText"
            size="xs"
            variant="secondary"
          >
            <Copy class="size-4" />
          </Button>
        </div>
      </div>
    </div>
  </div>
</template>
