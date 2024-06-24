<script setup lang="ts">
import { formatDuration, intervalToDuration } from 'date-fns'
import { Users } from 'lucide-vue-next'
import { computed } from 'vue'

import { useStream } from '@/api/stream.ts'
import { Button } from '@/components/ui/button'
import {
  Popover,
  PopoverContent,
  PopoverTrigger
} from '@/components/ui/popover'
import { ScrollArea } from '@/components/ui/scroll-area'

const { data: streamState } = useStream().useStreamState()
</script>

<template>
  <Popover
    side="left"
    prioritizePosition
  >
    <PopoverTrigger as-child>
      <Button
        size="xs"
        variant="ghost"
        class="flex justify-start items-center gap-2"
      >
        <Users />
        <span>
          {{ streamState?.streamInfo?.viewers }}
        </span>
      </Button>
    </PopoverTrigger>
    <PopoverContent
      v-if="streamState?.streamInfo?.chatters"
      class="p-0.5"
    >
      <ScrollArea class="h-[200px] rounded-md flex flex-col">
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
      </ScrollArea>
    </PopoverContent>
  </Popover>
</template>

<style scoped></style>
