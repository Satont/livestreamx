import { isHLSProvider, isVideoProvider } from 'vidstack'
import type { MediaPlayerElement } from 'vidstack/elements'
import { onBeforeUnmount, shallowRef, watch, type Ref } from 'vue'

import {
  EMPTY_PLAYER_STATS,
  PLAYER_STATS_INTERVAL_MS,
  readPlayerStats,
  type PlayerStats
} from '~/utils/player-stats'

/** Polls playback statistics of the stream player while it is mounted. */
export function usePlayerStats(player: Ref<MediaPlayerElement | null>) {
  const stats = shallowRef<PlayerStats>(EMPTY_PLAYER_STATS)
  let timer: ReturnType<typeof setInterval> | null = null

  function refresh() {
    const element = player.value
    if (!element) {
      stats.value = EMPTY_PLAYER_STATS
      return
    }
    const provider = element.provider
    const media =
      provider && isVideoProvider(provider)
        ? provider.video
        : element.querySelector('video')
    const hls = provider && isHLSProvider(provider) ? provider.instance : null

    stats.value = readPlayerStats({
      media,
      hls,
      viewport: { width: element.clientWidth, height: element.clientHeight }
    })
  }

  function stop() {
    if (timer !== null) {
      clearInterval(timer)
      timer = null
    }
  }

  watch(
    player,
    (element) => {
      stop()
      refresh()
      if (element) {
        timer = setInterval(refresh, PLAYER_STATS_INTERVAL_MS)
      }
    },
    { immediate: true }
  )

  onBeforeUnmount(stop)

  return stats
}
