<script setup lang="ts">
// import { WebRTCPlayer } from '@eyevinn/webrtc-player'
// import { useLocalStorage } from '@vueuse/core'
// import { onMounted, ref, watch } from 'vue'
//
// const volume = useLocalStorage('livestreamx-player-volume', 0.2)
//
// const player = ref<WebRTCPlayer | null>(null)
// const playerRef = ref<HTMLVideoElement | null>(null)
//
// async function initPlayer() {
//   while (true) {
//     if (!playerRef.value) {
//       continue
//     }
//     try {
//       player.value?.unload()
//       player.value = null
//       playerRef.value.src = ''
//
//       player.value = new WebRTCPlayer({
//         video: playerRef.value,
//         type: 'whep',
//         statsTypeFilter: '^candidate-*|^inbound-rtp',
//         timeoutThreshold: 2000,
//         debug: true
//       })
//       await player.value.load(new URL('https://streamx.satont.dev/stream/whep'))
//       // await player.value.load(new URL('http://localhost:8889/mystream/whep'))
//       player.value.unmute()
//
//       player.value.on('no-media', () => {
//         initPlayer()
//       })
//       player.value.on('media-recovered', () => {
//         console.log('media recovered')
//       })
//
//       break
//     } catch {
//       await new Promise((resolve) => setTimeout(resolve, 1000))
//     }
//   }
// }
//
// onMounted(async () => {
//   const video = document.getElementById('video') as HTMLVideoElement
//   playerRef.value = video
//
//   initPlayer()
// })
//
// watch(player, () => {
//   if (!player.value) return
//
//   if (playerRef.value) {
//     playerRef.value.volume = volume.value
//   }
// })
//
// watch(
//   playerRef,
//   (v) => {
//     if (!v) return
//
//     v.onvolumechange = (ev) => {
//       const target = ev.target as HTMLVideoElement
//       if (target.muted) {
//         return
//       }
//
//       volume.value = v.volume
//     }
//   },
//   { immediate: true }
// )
//
const autoplay = !import.meta.env.DEV
</script>

<template>
  <iframe
    :src="`https://streamx.satont.dev/stream/?autoplay=${autoplay.toString()}&muted=false&controls=true&playsinline=false`"
    scrolling="no"
    class="w-full h-full"
  />
  <!--  <video-->
  <!--    ref="playerRef"-->
  <!--    id="video"-->
  <!--    class="h-full w-full"-->
  <!--    controls-->
  <!--    :autoplay="autoplay"-->
  <!--    playsinline="true"-->
  <!--  />-->
</template>
