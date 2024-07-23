<script setup lang="ts">
import { breakpointsTailwind, useBreakpoints } from '@vueuse/core'

import { ChatReaction_Fragment, useChat } from '~/api/chat'
import { useReactions } from '~/composables/use-reactions.js'
import { useFragment } from '~/gql'
import { ChatMessageReactionType } from '~/gql/graphql'

const { dialogOpened, currentMessage, addReaction } = useReactions()

const reactions = computed(() => {
  if (!currentMessage.value) return null

  return useFragment(ChatReaction_Fragment, currentMessage.value.reactions)
})
const { emotes } = useChat()

const emotesSearchTerm = ref('')
const filteredEmotes = computed(() => {
  return emotes.value.filter((emote) =>
    emote.name.toLowerCase().includes(emotesSearchTerm.value.toLowerCase())
  )
})

const breakpoints = useBreakpoints(breakpointsTailwind)
const smallerThanLg = breakpoints.smaller('lg')
</script>

<template>
  <UiDialog v-model:open="dialogOpened">
    <UiDialogContent
      class="p-0 w-[600px]"
      hide-close
    >
      <UiDialogClose
        class="absolute right-[-15px] top-[-10px]"
        as-child
      >
        <UiButton
          size="xs"
          variant="secondary"
          class="rounded-full"
        >
          <Icon
            name="lucide:x"
            class="size-4"
          />
        </UiButton>
      </UiDialogClose>
      <UiTabs
        default-value="emotes"
        orientation="vertical"
        class="flex gap-2 p-2 h-[350px]"
      >
        <UiTabsList class="flex flex-col h-auto justify-start">
          <UiTabsTrigger
            value="emotes"
            class="w-full"
          >
            Emotes
          </UiTabsTrigger>
          <UiTabsTrigger
            value="reactions"
            class="w-full"
          >
            Reactions
          </UiTabsTrigger>
        </UiTabsList>

        <UiSeparator orientation="vertical" />

        <UiTabsContent
          value="emotes"
          class="w-full"
        >
          <div class="flex flex-col gap-2 h-[320px]">
            <UiInput
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
                  @click="addReaction(item.name)"
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
        </UiTabsContent>
        <UiTabsContent
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
                <NuxtImg
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
        </UiTabsContent>
      </UiTabs>
    </UiDialogContent>
  </UiDialog>
</template>

<style scoped></style>
