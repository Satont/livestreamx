<script setup lang="ts">
import { VideoPlayer } from '@videojs-player/vue'
import videojs from 'video.js'

import 'videojs-persist'

import { computed, shallowRef } from 'vue'

import 'video.js/dist/video-js.css'

import { useChat } from '@/api/chat.ts'

type Player = ReturnType<typeof videojs>

const player = shallowRef<Player>()

const { channelData } = useChat()

const videoSource = computed(() => {
  if (!channelData.value) return null

  return {
    src: `${window.location.origin}/api/streams/read/${channelData.value.fetchUserByName.id}/index.m3u8`,
    type: 'application/x-mpegURL'
  }
})

const handleMounted = (payload: { player: Player }) => {
  player.value = payload.player
  // @ts-ignore
  player.value.persist()
}

const handleReady = () => {
  const { vhs } = player.value?.tech() as any
  vhs.xhr.beforeRequest = (options: any) => {
    return options
  }
}
</script>

<template>
  <div v-if="!videoSource">
    <div class="flex items-center justify-center w-full h-full">
      <div class="text-center">
        <div class="text-2xl font-bold text-gray-500">No stream available</div>
      </div>
    </div>
  </div>

  <video-player
    v-else
    class="w-full h-full video-player vjs-theme-forest"
    poster="/images/example/4.jpg"
    crossorigin="anonymous"
    playsinline
    autoplay="any"
    controls
    liveui
    :muted="false"
    :sources="[videoSource]"
    :control-bar="{
      progressControl: false,
      currentTimeDisplay: false,
      remainingTimeDisplay: false
    }"
    :html5="{
      vhs: {
        maxPlaylistRetries: Infinity
      },
      nativeAudioTracks: false,
      nativeVideoTracks: false
    }"
    @mounted="handleMounted"
    @ready="handleReady"
  />
</template>
