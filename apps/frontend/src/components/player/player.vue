<script setup lang="ts">
import 'vidstack/bundle'

import { isHLSProvider } from 'vidstack'
import { computed } from 'vue'
import type { MediaProviderChangeEvent } from 'vidstack'

import { useChat } from '@/api/chat.ts'
import { useStream } from '@/api/stream.ts'

const { channelData } = useChat()
const { data: streamData } = useStream().useStreamState()

const streamingServiceAddr = import.meta.env.DEV
  ? 'http://localhost:8888'
  : 'http://78.46.90.174:8888'

const src = computed(() => {
  if (!channelData.value || !streamData.value?.streamInfo?.startedAt) {
    return null
  }
  return `${streamingServiceAddr}/${channelData.value!.fetchUserByName.name}/index.m3u8`
})

function onProviderChange(event: MediaProviderChangeEvent) {
  const provider = event.detail
  if (isHLSProvider(provider)) {
    provider.library = () => import('hls.js')
    provider.config = {
      maxLiveSyncPlaybackRate: 1.5
    }
  }
}
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
    streamType="live"
    viewType="video"
    :loop="false"
    @provider-change="onProviderChange"
  >
    <media-provider />
    <media-video-layout />
  </media-player>
</template>
