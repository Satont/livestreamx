<script setup lang="ts">
import { ref, watchEffect } from 'vue'

import { init } from './player.js'

const videoRef = ref<HTMLVideoElement | null>(null)
const messageRef = ref<HTMLDivElement | null>(null)

const isDev = import.meta.env.DEV

watchEffect(() => {
  if (!videoRef.value || !messageRef.value) return

  init({
    videoEl: videoRef.value,
    messageEl: messageRef.value,
    streamSrc: isDev
      ? 'http://localhost:8889/mystream/'
      : 'https://streamx.satont.dev/stream/',
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
