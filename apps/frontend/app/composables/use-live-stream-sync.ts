import type Hls from 'hls.js'
import type { ErrorData } from 'hls.js'
import { isHLSProvider, isVideoProvider } from 'vidstack'
import type { HLSProvider, MediaProviderAdapter } from 'vidstack'
import type { MediaPlayerElement } from 'vidstack/elements'
import { onBeforeUnmount, onMounted, watch, type Ref } from 'vue'

import {
  HLS_NETWORK_RESTART_DELAYS_MS,
  LiveEdgeSyncController,
  MIN_LIVE_EDGE_SEEK_DISTANCE,
  NATIVE_LIVE_EDGE_OFFSET,
  resolveHlsLiveEdge,
  resolveNativeLiveEdge
} from '~/utils/live-edge-sync'

type Timer = ReturnType<typeof setTimeout>

/**
 * Keeps live playback close to the live edge.
 *
 * hls.js bounds latency on its own through the config applied in `Player.vue`,
 * this composable adds the event driven corrections:
 *
 * - resuming from a pause seeks back to the live edge;
 * - returning to the tab, bfcache restore and network recovery re-check drift;
 * - pending corrections are retried on fresh playlists/media events;
 * - fatal network errors restart hls.js loading with a capped backoff.
 *
 * A correction is only applied to playing media, so a user paused player is
 * never resumed implicitly.
 */
export function useLiveStreamSync(player: Ref<MediaPlayerElement | null>) {
  let media: HTMLMediaElement | null = null
  let hls: Hls | null = null
  let disposeInstanceCallback: (() => void) | null = null
  let disposeHlsListeners: (() => void) | null = null
  let networkRestartTimer: Timer | null = null
  let networkRestartAttempts = 0
  let loadingStopped = false

  function getMediaElement(): HTMLMediaElement | null {
    const element = player.value
    if (!element) {
      return null
    }
    const provider = element.provider
    if (provider && isVideoProvider(provider)) {
      return provider.video
    }
    return element.querySelector('video')
  }

  function resolveLiveEdgeTarget(): number | null {
    const element = player.value
    if (!element) {
      return null
    }
    const provider = element.provider
    if (provider && isHLSProvider(provider)) {
      return resolveHlsLiveEdge(provider.instance)
    }
    return resolveNativeLiveEdge(getMediaElement(), NATIVE_LIVE_EDGE_OFFSET)
  }

  const controller = new LiveEdgeSyncController({
    getTarget: resolveLiveEdgeTarget,
    getCurrentTime: () => getMediaElement()?.currentTime ?? null,
    isPaused: () => getMediaElement()?.paused ?? true,
    isSeeking: () => getMediaElement()?.seeking ?? false,
    seek: (position) => {
      const element = getMediaElement()
      if (element) {
        element.currentTime = position
      }
    },
    now: () => Date.now()
  })

  function onMediaReady() {
    controller.retry()
  }

  function requestCorrection(minDistance: number = MIN_LIVE_EDGE_SEEK_DISTANCE) {
    const element = getMediaElement()
    if (!element || element.paused) {
      return
    }
    controller.request(minDistance)
  }

  function clearNetworkRestartTimer() {
    if (networkRestartTimer !== null) {
      clearTimeout(networkRestartTimer)
      networkRestartTimer = null
    }
  }

  function restartHlsLoading() {
    if (!hls || !loadingStopped) {
      return
    }
    hls.startLoad()
  }

  function scheduleNetworkRestart() {
    if (
      networkRestartTimer !== null ||
      networkRestartAttempts >= HLS_NETWORK_RESTART_DELAYS_MS.length
    ) {
      return
    }
    const delay = HLS_NETWORK_RESTART_DELAYS_MS[networkRestartAttempts]
    networkRestartAttempts += 1
    networkRestartTimer = setTimeout(() => {
      networkRestartTimer = null
      restartHlsLoading()
    }, delay)
  }

  function attachHlsInstance(provider: HLSProvider, instance: Hls) {
    const ctor = provider.ctor
    if (!ctor) {
      return
    }
    disposeHlsListeners?.()
    disposeHlsListeners = null
    hls = instance

    const onError = (_event: string, data: ErrorData) => {
      if (!data.fatal || data.type !== ctor.ErrorTypes.NETWORK_ERROR) {
        return
      }
      loadingStopped = true
      scheduleNetworkRestart()
    }
    const onLevelLoaded = () => {
      loadingStopped = false
      networkRestartAttempts = 0
      controller.retry()
    }
    const onLevelUpdated = () => {
      controller.retry()
    }

    instance.on(ctor.Events.ERROR, onError)
    instance.on(ctor.Events.LEVEL_LOADED, onLevelLoaded)
    instance.on(ctor.Events.LEVEL_UPDATED, onLevelUpdated)

    disposeHlsListeners = () => {
      instance.off(ctor.Events.ERROR, onError)
      instance.off(ctor.Events.LEVEL_LOADED, onLevelLoaded)
      instance.off(ctor.Events.LEVEL_UPDATED, onLevelUpdated)
    }
  }

  function detachProvider() {
    clearNetworkRestartTimer()
    disposeHlsListeners?.()
    disposeHlsListeners = null
    disposeInstanceCallback?.()
    disposeInstanceCallback = null

    if (media) {
      media.removeEventListener('canplay', onMediaReady)
      media.removeEventListener('playing', onMediaReady)
      media.removeEventListener('seeked', onMediaReady)
    }

    media = null
    hls = null
    loadingStopped = false
    networkRestartAttempts = 0
  }

  function attachProvider(provider: MediaProviderAdapter | null) {
    detachProvider()

    if (!provider) {
      return
    }

    media = isVideoProvider(provider)
      ? provider.video
      : (player.value?.querySelector('video') ?? null)

    if (media) {
      media.addEventListener('canplay', onMediaReady)
      media.addEventListener('playing', onMediaReady)
      media.addEventListener('seeked', onMediaReady)
    }

    if (isHLSProvider(provider)) {
      disposeInstanceCallback = provider.onInstance((instance) => {
        attachHlsInstance(provider, instance)
      })
    }
  }

  /** Clears pending corrections, e.g. after the stream source changed. */
  function reset() {
    controller.reset()
    loadingStopped = false
    networkRestartAttempts = 0
    clearNetworkRestartTimer()
  }

  function onVisibilityChange() {
    if (document.visibilityState !== 'visible') {
      return
    }
    requestCorrection()
  }

  function onPageShow() {
    requestCorrection()
  }

  function onOnline() {
    if (loadingStopped) {
      clearNetworkRestartTimer()
      networkRestartAttempts = 0
      restartHlsLoading()
    }
    requestCorrection()
  }

  watch(player, (element, previousElement) => {
    if (element === previousElement) {
      return
    }
    detachProvider()
    controller.reset()
  })

  onMounted(() => {
    document.addEventListener('visibilitychange', onVisibilityChange)
    window.addEventListener('pageshow', onPageShow)
    window.addEventListener('online', onOnline)
  })

  onBeforeUnmount(() => {
    document.removeEventListener('visibilitychange', onVisibilityChange)
    window.removeEventListener('pageshow', onPageShow)
    window.removeEventListener('online', onOnline)
    detachProvider()
    controller.dispose()
  })

  return {
    attachProvider,
    reset,
    onPause: () => controller.pause(),
    onPlay: () => controller.play()
  }
}
