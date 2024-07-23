<script setup lang="ts">
import { computed } from 'vue'
import type { FragmentType } from '@/gql'

import { ChatMessage_Fragment, ChatReaction_Fragment } from '@/api/chat.js'
import { useProfile } from '@/api/profile.js'
import { useFragment } from '@/gql'
import { ChatMessageReactionType } from '@/gql/graphql.js'
import { useReactions } from '~/composables/use-reactions.js'
import { useShowReactionsOnMessage } from '~/composables/use-show-reactions-on-message.js'
import { arrayUniqueBy } from '~/utils/array-unique.js'

type Props = {
  msg: FragmentType<typeof ChatMessage_Fragment>
}

const props = defineProps<Props>()
const unwrappedMessage = useFragment(ChatMessage_Fragment, props.msg)
const reactions = useFragment(ChatReaction_Fragment, unwrappedMessage.reactions)
const { showReactionsOnMessage } = useShowReactionsOnMessage()
const { addReaction } = useReactions()

const { data: profile } = useProfile().useData()

const mappedReactions = computed(() => {
  const uniqueReactions = arrayUniqueBy(
    reactions,
    (a, b) => a.reaction === b.reaction
  )

  return uniqueReactions
    .map((r) => {
      return {
        ...r,
        count: reactions.filter((reaction) => reaction.reaction === r.reaction)
          .length
      }
    })
    .sort((a, b) => b.count - a.count)
    .slice(0, 3)
})

const { openDialog } = useReactions()
</script>

<template>
  <div
    class="gap-2 items-center flex"
    v-if="showReactionsOnMessage"
  >
    <UiButton
      size="xs"
      class="hidden group-hover:block items-center"
      :disabled="!profile"
      @click="openDialog(props.msg)"
    >
      <Icon
        name="lucide:smile-plus"
        class="size-4"
      />
    </UiButton>
    <UiButton
      v-for="(reaction, index) of mappedReactions"
      :key="index"
      size="xs"
      class="relative rounded-full h-6 disabled:opacity-100"
      :class="{
        'p-0': reaction.type === ChatMessageReactionType.Emote
      }"
      variant="secondary"
      @click="addReaction(reaction.reaction, unwrappedMessage.id)"
      :disabled="
        reactions.some(
          (r) =>
            r.user.id === profile?.userProfile.id &&
            r.reaction === reaction.reaction
        )
      "
    >
      <span
        v-if="reaction.type === ChatMessageReactionType.Emoji"
        class="text-xs"
      >
        {{ reaction.reaction }}
      </span>
      <NuxtImg
        v-else-if="
          reaction.type === ChatMessageReactionType.Emote && 'emote' in reaction
        "
        :src="'https:' + reaction.emote.url"
        class="size-6 rounded-full"
      />

      <span
        class="absolute text-accent bottom-[-7px] right-[-5px] px-1 bg-accent-foreground rounded-full text-xs font-bold"
      >
        {{ reaction.count }}
      </span>
    </UiButton>
  </div>
</template>
