<script setup lang="ts">
import { SmilePlus, X } from 'lucide-vue-next'
import { computed, ref } from 'vue'
// @ts-ignore
import { RecycleScroller } from 'vue-virtual-scroller'

import {
  ChatMessage_Fragment,
  ChatReaction_Fragment,
  useChat
} from '@/api/chat.js'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogTrigger
} from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Separator } from '@/components/ui/separator'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { ChatMessageReactionType } from '@/gql/graphql.ts'
import { arrayUniqueBy } from '@/lib/array-unique.js'

import 'vue-virtual-scroller/dist/vue-virtual-scroller.css'

import { breakpointsTailwind, useBreakpoints } from '@vueuse/core'

import { useProfile } from '@/api/profile.js'
import { useShowReactionsOnMessage } from '@/composables/use-show-reactions-on-message.ts'
import { FragmentType, useFragment } from '@/gql'

type Props = {
  msg: FragmentType<typeof ChatMessage_Fragment>
}

const props = defineProps<Props>()
const unwrappedMessage = useFragment(ChatMessage_Fragment, props.msg)
const reactions = useFragment(ChatReaction_Fragment, unwrappedMessage.reactions)
const { showReactionsOnMessage } = useShowReactionsOnMessage()

const { data: profile } = useProfile().useData()
const dialogOpen = ref(false)
const { emotes, useReactionAddMutation } = useChat()

const emotesSearchTerm = ref('')
const filteredEmotes = computed(() => {
  return emotes.value.filter((emote) =>
    emote.name.toLowerCase().includes(emotesSearchTerm.value.toLowerCase())
  )
})

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

const breakpoints = useBreakpoints(breakpointsTailwind)
const smallerThanLg = breakpoints.smaller('lg')

const reactionAddMutation = useReactionAddMutation()
async function handleAddReaction(name: string) {
  try {
    const { error } = await reactionAddMutation.executeMutation({
      messageId: unwrappedMessage.id,
      content: name
    })

    if (error) {
      throw new Error(error.toString())
    }

    dialogOpen.value = false
  } catch (e) {
    console.log(e)
  }
}
</script>

<template>
  <div
    class="gap-2 items-center flex"
    v-if="showReactionsOnMessage"
  >
    <Button
      v-for="(reaction, index) of mappedReactions"
      :key="index"
      size="xs"
      class="relative rounded-full h-6 disabled:opacity-100"
      :class="{
        'p-0': reaction.type === ChatMessageReactionType.Emote
      }"
      variant="secondary"
      @click="handleAddReaction(reaction.reaction)"
      :disabled="
        reactions.some(
          (r) =>
            r.user.id === profile?.userProfile.user.id &&
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
      <img
        v-else-if="
          reaction.type === ChatMessageReactionType.Emote && 'emote' in reaction
        "
        :src="reaction.emote.url"
        class="size-6 rounded-full"
      />

      <span
        class="absolute text-accent bottom-[-7px] right-[-5px] px-1 bg-accent-foreground rounded-full text-xs font-bold"
      >
        {{ reaction.count }}
      </span>
    </Button>
  </div>

  <Dialog v-model:open="dialogOpen">
    <DialogTrigger :disabled="!profile">
      <Button
        size="xs"
        class="hidden group-hover:block"
        :disabled="!profile"
      >
        <SmilePlus class="size-4" />
      </Button>
    </DialogTrigger>
    <DialogContent
      class="p-0 max-w-full w-[600px]"
      disable-default-close
    >
      <DialogClose
        class="absolute right-[-15px] top-[-10px]"
        as-child
      >
        <Button
          size="xs"
          variant="secondary"
          class="rounded-full"
        >
          <X class="size-4" />
        </Button>
      </DialogClose>
      <Tabs
        default-value="emotes"
        orientation="vertical"
        class="flex gap-2 p-2 h-[350px]"
      >
        <TabsList class="flex flex-col h-auto justify-start">
          <TabsTrigger
            value="emotes"
            class="w-full"
          >
            Emotes
          </TabsTrigger>
          <TabsTrigger
            value="reactions"
            class="w-full"
          >
            Reactions
          </TabsTrigger>
        </TabsList>

        <Separator orientation="vertical" />

        <TabsContent
          value="emotes"
          class="w-full"
        >
          <div class="flex flex-col gap-2 h-[320px]">
            <Input
              v-model="emotesSearchTerm"
              placeholder="Search emote..."
            />
            <RecycleScroller
              class="h-full"
              :items="filteredEmotes"
              :item-size="60"
              :item-secondary-size="79"
              :gridItems="smallerThanLg ? 4 : 6"
              key-field="name"
            >
              <template #default="{ item }">
                <div
                  class="flex flex-col items-center text-ellipsis overflow-hidden cursor-pointer"
                  style="height: 60px; width: 79px"
                  @click="handleAddReaction(item.name)"
                >
                  <img
                    :src="item.url"
                    class="size-8"
                  />
                  <span>{{ item.name }}</span>
                </div>
              </template>
            </RecycleScroller>
          </div>
        </TabsContent>
        <TabsContent
          value="reactions"
          class="h-full w-full overflow-y-auto"
        >
          <div class="grid grid-cols-2 gap-4 grid-flow-row-dense w-full">
            <div
              v-for="(reaction, index) of reactions"
              :key="index"
              class="flex justify-between bg-accent p-2 rounded items-center"
            >
              <div class="flex items-center gap-2">
                <img
                  :src="reaction.user.avatarUrl"
                  class="size-8 rounded-full"
                />
                <span>{{ reaction.user.displayName }}</span>
              </div>

              <img
                v-if="
                  reaction.type === ChatMessageReactionType.Emote &&
                  'emote' in reaction
                "
                :src="reaction.emote.url"
                class="size-8 rounded-full"
              />
              <span v-else>
                {{ reaction.reaction }}
              </span>
            </div>
          </div>
        </TabsContent>
      </Tabs>
    </DialogContent>
  </Dialog>
</template>
