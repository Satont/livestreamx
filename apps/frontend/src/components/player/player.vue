<script setup lang="ts">
import { ref, watchEffect } from 'vue'

import { init } from './player.js'

const videoRef = ref<HTMLVideoElement | null>(null)
const messageRef = ref<HTMLDivElement | null>(null)

watchEffect(() => {
  if (!videoRef.value || !messageRef.value) return

  init({
    videoEl: videoRef.value,
    messageEl: messageRef.value,
    streamSrc: 'https://streamx.satont.dev/stream/',
    controls: true
  })
})
</script>

<template>
  <div v-bind="$attrs">
    <video
      ref="videoRef"
      class="w-full h-full"
    />
  </div>
  <div ref="messageRef"></div>
</template>
