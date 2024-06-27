<script setup lang="ts">
import { ref, watchEffect } from 'vue'

import { init } from './player.js'

const videoRef = ref<HTMLVideoElement | null>(null)
const messageRef = ref<HTMLDivElement | null>(null)

const streamUrl = import.meta.env.VITE_STREAM_URL

watchEffect(() => {
  if (!videoRef.value || !messageRef.value) return

  init({
    videoEl: videoRef.value,
    messageEl: messageRef.value,
    streamSrc: streamUrl ?? 'http://localhost:8889/mystream/',
    controls: true
  })
})
</script>

<template>
  <div class="relative w-full h-full">
    <video
      ref="videoRef"
      class="w-full h-full"
    />
    <div
      ref="messageRef"
      class="absolute top-[50%] left-[50%] transform translate-x-[-50%] translate-y-[-50%] bg-red-900 text-white rounded-md"
    ></div>
  </div>
</template>
