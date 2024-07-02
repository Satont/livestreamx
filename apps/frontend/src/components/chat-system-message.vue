<script setup lang="ts">
import { computed, FunctionalComponent } from 'vue'

import { SystemMessage_Fragment, useChat } from '@/api/chat.ts'
import SevenTv from '@/assets/images/seventv.svg?component'
import { chatFontSize } from '@/composables/chat-font-size.ts'
import { FragmentType, useFragment } from '@/gql'
import { ChatEmote_FragmentFragment, SystemMessageType } from '@/gql/graphql.ts'

type Props = {
  msg: FragmentType<typeof SystemMessage_Fragment>
}

const props = defineProps<Props>()
const { emotes } = useChat()

const unwrappedMessage = useFragment(SystemMessage_Fragment, props.msg)

const data = computed<{
  title: string
  icon?: FunctionalComponent
}>(() => {
  switch (unwrappedMessage.type) {
    case SystemMessageType.EmoteAdded:
      return {
        title: 'Emote added',
        icon: SevenTv
      }
    case SystemMessageType.EmoteRemoved:
      return {
        title: 'Emote removed',
        icon: SevenTv
      }
    case SystemMessageType.KickMessage:
      return {
        title: 'Kick message',
        icon: SevenTv
      }
    default:
      return {
        title: 'System message'
      }
  }
})

const removedEmote = computed(() => {
  if (
    unwrappedMessage.type !== SystemMessageType.EmoteRemoved ||
    !('emoteId' in unwrappedMessage)
  ) {
    return null
  }

  return emotes.value.find((emote) => emote.id === unwrappedMessage.emoteId)
})
</script>

<template>
  <div
    :style="{ fontSize: `${chatFontSize}px` }"
    class="relative group p-0.5 flex flex-col rounded mb-1"
    :class="{
      'bg-cyan-700/90':
        unwrappedMessage.type === SystemMessageType.EmoteAdded ||
        unwrappedMessage.type === SystemMessageType.EmoteRemoved
    }"
  >
    <div class="text-sm px-1 flex items-center gap-1">
      <component
        v-if="data.icon"
        :is="data.icon"
        class="size-4"
      />
      {{ data.title }}
      <template v-if="'actor' in unwrappedMessage">
        by {{ unwrappedMessage.actor.displayName }}
      </template>
    </div>
    <div class="p-1 bg-accent rounded-b">
      <template
        v-if="
          unwrappedMessage.type === SystemMessageType.EmoteAdded &&
          'emote' in unwrappedMessage
        "
      >
        <img
          :src="(unwrappedMessage.emote as ChatEmote_FragmentFragment).url"
          :style="{
            width: `${(unwrappedMessage.emote as ChatEmote_FragmentFragment).width}px`,
            height: `${(unwrappedMessage.emote as ChatEmote_FragmentFragment).height}px`
          }"
          class="scale-90 inline-block relative"
        />
        {{ (unwrappedMessage.emote as ChatEmote_FragmentFragment).name }}
      </template>
      <template
        v-if="
          unwrappedMessage.type === SystemMessageType.EmoteRemoved &&
          removedEmote
        "
      >
        <img
          :src="removedEmote.url"
          :style="{
            width: `${removedEmote.width}px`,
            height: `${removedEmote.height}px`
          }"
          class="scale-90 inline-block relative"
        />
        {{ removedEmote.name }}
      </template>
    </div>
  </div>
</template>
