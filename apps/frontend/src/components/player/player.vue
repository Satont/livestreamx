<script setup lang="ts">
import { useLocalStorage } from '@vueuse/core'
import { Fullscreen, Pause, Play } from 'lucide-vue-next'
import { ref, watchEffect } from 'vue'

import { useChat } from '@/api/chat.ts'
import { Button } from '@/components/ui/button'
import { Slider } from '@/components/ui/slider'
import { init } from './player.js'

const videoRef = ref<HTMLVideoElement | null>(null)
const messageRef = ref<HTMLDivElement | null>(null)

const { channelData } = useChat()

watchEffect(() => {
  if (
    !videoRef.value ||
    !messageRef.value ||
    !channelData.value?.fetchUserByName.id
  )
    return

  init({
    videoEl: videoRef.value,
    messageEl: messageRef.value,
    streamSrc: `${window.location.origin}/api/streams/${channelData.value.fetchUserByName.id}/`,
    controls: false,
    muted: false
  })
})

const paused = ref(false)
function switchPaused() {
  if (paused.value) {
    videoRef.value?.play()
  } else {
    videoRef.value?.pause()
  }

  paused.value = !paused.value
}

const volume = useLocalStorage('livestreamx-player-volume-v2', [30], {
  serializer: {
    read: (v: any) => (v ? JSON.parse(v) : null),
    write: (v: any) => JSON.stringify(v)
  }
})

watchEffect(() => {
  if (!videoRef.value) return

  videoRef.value.volume = volume.value / 100
})

const controlsVisible = ref(false)
</script>

<template>
  <div
    class="relative w-full h-full group"
    @touchstart.passive="controlsVisible = true"
    @touchend.passive="controlsVisible = false"
  >
    <video
      ref="videoRef"
      class="w-full h-full"
    />

    <div
      class="absolute group-hover:flex justify-between items-center bottom-0 w-full hidden bg-foreground/30 p-2 gap-4"
      :class="{ flex: controlsVisible }"
    >
      <Button
        size="xs"
        @click="switchPaused"
      >
        <component
          :is="paused ? Play : Pause"
          class="size-4"
        />
      </Button>

      <div class="flex gap-2 items-center flex-wrap gap-4">
        <Slider
          v-model="volume"
          class="w-64 md:w-80"
        />
        <Button
          size="xs"
          @click="videoRef?.requestFullscreen()"
          ><Fullscreen class="size-4"
        /></Button>
      </div>
    </div>

    <div
      ref="messageRef"
      class="absolute top-[50%] left-[50%] transform translate-x-[-50%] translate-y-[-50%] bg-red-900 text-white rounded-md"
    ></div>
  </div>
</template>
