<script setup lang="ts">
import 'vidstack/bundle'

import { isHLSProvider } from 'vidstack'
import { computed, onBeforeUnmount, ref, watch } from 'vue'
import type { MediaPlayerElement, MediaProviderChangeEvent, VideoQuality } from 'vidstack'

import { useChat } from '~/api/chat.js'
import { useStream } from '~/api/stream.js'

const { channelData } = useChat()
const { data: streamData } = useStream().useStreamState()

const player = ref<MediaPlayerElement | null>(null)

// dev runs the player against the local ome instance directly,
// production goes through caddy: /mtx/* -> ome:3333
const streamingServiceAddr = import.meta.env.DEV
  ? 'http://127.0.0.1:3334'
  : `${window.location.origin}/mtx`

const src = computed(() => {
  if (!channelData.value || !streamData.value?.streamInfo?.startedAt) {
    return null
  }
  return `${streamingServiceAddr}/app/${channelData.value.fetchUserByName.name}/master.m3u8`
})

function onProviderChange(event: MediaProviderChangeEvent) {
  const provider = event.detail
  if (isHLSProvider(provider)) {
    provider.library = () => import('hls.js')
    provider.config = {
      maxLiveSyncPlaybackRate: 1.5,
      // hls.js starts auto quality from a bandwidth estimate capped at 5 Mbps
      // (abrEwmaDefaultEstimateMax), so it picks 720p even when the viewer can
      // handle the source rendition. Assume the top of the ladder is playable
      // and let ABR drop the level if the measured bandwidth can't keep up.
      abrEwmaDefaultEstimate: Number.POSITIVE_INFINITY
    }
  }
}

// after a pause we always want to be live again, otherwise the user keeps
// watching with a delay equal to the pause duration
const LIVE_EDGE_RESUME_THRESHOLD_MS = 1000

let pausedAt: number | null = null

function onPlayerPause() {
  pausedAt = Date.now()
}

function onPlayerPlay() {
  const pausedFor = pausedAt === null ? 0 : Date.now() - pausedAt
  pausedAt = null

  if (pausedFor < LIVE_EDGE_RESUME_THRESHOLD_MS) {
    return
  }

  seekToLiveEdge()
}

// `player.seekToLiveEdge()` is a no-op for live streams without a dvr window,
// which is how we play ome ll-hls here, so ask hls.js where live is right now.
// it keeps `liveSyncPosition` fresh while paused because it continues to
// refresh the playlist in the background.
function seekToLiveEdge() {
  const el = player.value
  const provider = el?.provider

  if (provider && isHLSProvider(provider)) {
    const hls = provider.instance
    const position = hls?.liveSyncPosition
    if (hls?.media && position != null && Number.isFinite(position)) {
      hls.media.currentTime = position
      return
    }
  }

  el?.seekToLiveEdge()
}

const qualities = ref<VideoQuality[]>([])
const selectedQuality = ref<VideoQuality | null>(null)
const isAutoQuality = ref(true)
const isQualityMenuOpen = ref(false)

const sortedQualities = computed(() =>
  [...qualities.value].sort(
    (a, b) => b.height - a.height || (b.bitrate ?? 0) - (a.bitrate ?? 0)
  )
)

const currentQualityLabel = computed(() => {
  if (isAutoQuality.value || !selectedQuality.value) {
    return 'Auto'
  }
  return qualityLabel(selectedQuality.value).split(' · ')[0]
})

function qualityLabel(quality: VideoQuality) {
  const height = quality.height ? `${quality.height}p` : 'Source'
  const bitrate = quality.bitrate
    ? ` · ${(quality.bitrate / 1e6).toFixed(1)} Mbps`
    : ''
  return `${height}${bitrate}`
}

function syncQualityState() {
  const list = player.value?.qualities
  if (!list) {
    return
  }
  qualities.value = list.toArray()
  selectedQuality.value = list.selected
  isAutoQuality.value = list.auto
}

function selectQuality(quality: VideoQuality) {
  quality.selected = true
  isQualityMenuOpen.value = false
  syncQualityState()
}

function selectAutoQuality() {
  player.value?.qualities.autoSelect()
  isQualityMenuOpen.value = false
  syncQualityState()
}

function attachQualityListeners(el: MediaPlayerElement) {
  el.qualities.addEventListener('change', syncQualityState)
  el.qualities.addEventListener('auto-change', syncQualityState)
  el.qualities.addEventListener('add', syncQualityState)
  el.qualities.addEventListener('remove', syncQualityState)
}

function detachQualityListeners(el: MediaPlayerElement) {
  el.qualities.removeEventListener('change', syncQualityState)
  el.qualities.removeEventListener('auto-change', syncQualityState)
  el.qualities.removeEventListener('add', syncQualityState)
  el.qualities.removeEventListener('remove', syncQualityState)
}

watch(player, (el, prevEl) => {
  pausedAt = null
  if (prevEl) {
    detachQualityListeners(prevEl)
  }
  if (el) {
    attachQualityListeners(el)
    syncQualityState()
  }
})

onBeforeUnmount(() => {
  if (player.value) {
    detachQualityListeners(player.value)
  }
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
    streamType="live"
    viewType="video"
    :loop="false"
    @pause="onPlayerPause"
    @play="onPlayerPlay"
    @provider-change="onProviderChange"
  >
    <media-provider />
    <media-video-layout />

    <div
      v-if="qualities.length > 1"
      class="quality-menu absolute right-2 top-2 z-20"
      :data-open="isQualityMenuOpen || undefined"
      @click.stop
      @pointerdown.stop
    >
      <button
        type="button"
        class="flex items-center gap-1 rounded-md bg-black/70 px-2 py-1 text-xs font-medium text-white backdrop-blur transition hover:bg-black/85"
        @click="isQualityMenuOpen = !isQualityMenuOpen"
      >
        <Icon
          name="lucide:settings-2"
          class="size-3.5"
        />
        {{ currentQualityLabel }}
      </button>

      <ul
        v-if="isQualityMenuOpen"
        class="absolute right-0 top-8 min-w-36 overflow-hidden rounded-md border border-white/10 bg-black/90 py-1 text-xs text-white shadow-lg backdrop-blur"
      >
        <li>
          <button
            type="button"
            class="flex w-full items-center justify-between gap-3 px-3 py-1.5 text-left transition hover:bg-white/10"
            :class="{ 'text-primary': isAutoQuality }"
            @click="selectAutoQuality"
          >
            Auto
            <Icon
              v-if="isAutoQuality"
              name="lucide:check"
              class="size-3.5"
            />
          </button>
        </li>

        <li
          v-for="quality of sortedQualities"
          :key="quality.id"
        >
          <button
            type="button"
            class="flex w-full items-center justify-between gap-3 px-3 py-1.5 text-left transition hover:bg-white/10"
            :class="{ 'text-primary': !isAutoQuality && selectedQuality?.id === quality.id }"
            @click="selectQuality(quality)"
          >
            {{ qualityLabel(quality) }}
            <Icon
              v-if="!isAutoQuality && selectedQuality?.id === quality.id"
              name="lucide:check"
              class="size-3.5"
            />
          </button>
        </li>
      </ul>
    </div>
  </media-player>
</template>

<style scoped>
/* the menu is a part of the player element, so it stays visible in fullscreen.
   mimic the default layout behaviour: hide it together with the controls. */
media-player:not([data-controls]) .quality-menu:not([data-open]) {
  opacity: 0;
  pointer-events: none;
}
</style>
