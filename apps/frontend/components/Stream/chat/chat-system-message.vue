<script setup lang="ts">
import { computed } from 'vue'
import type { FragmentType } from '@/gql'
import type { ChatEmote_FragmentFragment } from '@/gql/graphql.js'
import type { FunctionalComponent } from 'vue'

import { SystemMessage_Fragment } from '@/api/chat.js'
import { chatFontSize } from '@/composables/chat-font-size.js'
import { useFragment } from '@/gql'
import { SystemMessageType } from '@/gql/graphql.js'

type Props = {
  msg: FragmentType<typeof SystemMessage_Fragment>
}

const props = defineProps<Props>()

const unwrappedMessage = useFragment(SystemMessage_Fragment, props.msg)

const data = computed<{
  title: string
  icon?: FunctionalComponent
}>(() => {
  switch (unwrappedMessage.type) {
    case SystemMessageType.EmoteAdded:
      return {
        title: 'Emote added'
        // icon: SevenTv
      }
    case SystemMessageType.EmoteRemoved:
      return {
        title: 'Emote removed'
        // icon: SevenTv
      }
    case SystemMessageType.KickMessage:
      return {
        title: 'Kick message'
        // icon: SevenTv
      }
    case SystemMessageType.UserJoined:
      return {
        title: 'User joined'
        // icon: MessageSquareShare
      }
    default:
      return {
        title: 'System message'
      }
  }
})
</script>

<template>
  <div
    :style="{ fontSize: `${chatFontSize}px` }"
    class="relative group p-0.5 flex flex-col rounded mb-1"
    :class="{
      'bg-cyan-700/90':
        unwrappedMessage.type === SystemMessageType.EmoteAdded ||
        unwrappedMessage.type === SystemMessageType.EmoteRemoved,
      'bg-green-700/90': unwrappedMessage.type === SystemMessageType.UserJoined
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
      <template v-if="'emote' in unwrappedMessage">
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
      <div
        v-if="
          'user' in unwrappedMessage &&
          unwrappedMessage.type === SystemMessageType.UserJoined
        "
        class="flex gap-2 p-2 items-center"
      >
        <NuxtImg
          class="size-7 rounded-full"
          :src="unwrappedMessage.user?.avatarUrl"
        />
        <span>{{ unwrappedMessage.user?.displayName }}</span>
      </div>
    </div>
  </div>
</template>
