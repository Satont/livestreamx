<script setup lang="ts">
import 'vidstack/bundle'

import { computed } from 'vue'

import { useChat } from '@/api/chat.ts'
import { useStream } from '@/api/stream.ts'

const { channelData } = useChat()
const { data: streamData } = useStream().useStreamState()

const src = computed(() => {
  if (!channelData.value || !streamData.value?.streamInfo?.startedAt)
    return null
  return `${window.location.origin}/api/streams/${channelData.value!.fetchUserByName.id}/index.m3u8`
})
</script>

<template>
  <div
    v-if="!streamData?.streamInfo?.startedAt"
    class="flex items-center justify-center w-full h-full"
  >
    <div class="text-center">
      <div class="text-2xl font-bold text-accent-foreground">
        Stream is offline
      </div>
    </div>
  </div>

  <media-player
    ref="player"
    v-else-if="src"
    :src="src"
    class="overflow-hidden h-full w-full"
    storage="streamx-player-v3"
    :title="channelData?.fetchUserByName.name"
    playsInline
    autoPlay
    logLevel="debug"
    :controls="false"
    :live-edge-tolerance="4"
    stream-type="live"
  >
    <media-provider />
    <media-video-layout disable-time-slider />
  </media-player>
</template>
