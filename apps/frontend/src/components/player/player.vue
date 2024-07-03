<script setup lang="ts">
import 'vidstack/bundle'

import { computed } from 'vue'

import { useChat } from '@/api/chat.ts'

const { channelData } = useChat()

const videoSource = computed(() => {
  if (!channelData.value) return null

  return {
    src: `${window.location.origin}/api/streams/read/${channelData.value.fetchUserByName.id}/index.m3u8`,
    type: 'application/x-mpegURL'
  }
})
</script>

<template>
  <div v-if="!videoSource">
    <div class="flex items-center justify-center w-full h-full">
      <div class="text-center">
        <div class="text-2xl font-bold text-gray-500">No stream available</div>
      </div>
    </div>
  </div>

  <media-player
    class="overflow-hidden"
    v-else-if="videoSource"
    storage="streamx-player-v3"
    :title="channelData?.fetchUserByName.name"
    :src="videoSource"
    playsInline
    autoPlay
    stream-type="live"
    disableTimeSlider
  >
    <media-provider />

    <media-video-layout />
  </media-player>
</template>
