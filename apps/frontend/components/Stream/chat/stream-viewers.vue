<script setup lang="ts">
import { useStream } from '@/api/stream.js'

const { data: streamState } = useStream().useStreamState()
</script>

<template>
  <UiPopover>
    <UiPopoverTrigger as-child>
      <UiButton
        size="sm"
        variant="ghost"
        class="flex justify-start items-center gap-2"
      >
        <Icon
          name="lucide:users"
          class="size-5"
        />
        <span>
          {{ streamState?.streamInfo?.viewers }}
        </span>
      </UiButton>
    </UiPopoverTrigger>
    <UiPopoverContent
      v-if="streamState?.streamInfo?.chatters"
      class="p-0.5"
    >
      <UiScrollArea class="h-[200px] rounded-md flex flex-col">
        <a
          v-for="chatter of streamState.streamInfo.chatters"
          :key="chatter.user.id"
          class="flex items-center gap-2 hover:bg-accent p-2 rounded"
          :href="`https://twitch.tv/${chatter.user.name}`"
          target="_blank"
        >
          <img
            :src="chatter.user.avatarUrl"
            class="size-6 rounded-full"
          />
          <span class="font-bold">{{ chatter.user.displayName }}</span>
        </a>
      </UiScrollArea>
    </UiPopoverContent>
  </UiPopover>
</template>

<style scoped></style>
