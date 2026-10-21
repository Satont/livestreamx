/**
 * Playback statistics for the stream player, shown under the player and in
 * the "Video stats" overlay.
 *
 * Reading is DOM-free: the media element, hls.js instance and viewport are
 * passed in through minimal shapes, which keeps it unit-testable.
 */

/** How often (ms) the player refreshes its statistics. */
export const PLAYER_STATS_INTERVAL_MS = 1000

/** Minimal shape of an hls.js level needed for statistics. */
export interface StatsLevelLike {
  width: number
  height: number
  bitrate: number
  frameRate: number
  videoCodec?: string
  audioCodec?: string
  details?: { partTarget: number }
}

/** Minimal shape of an hls.js instance needed for statistics. */
export interface StatsHlsLike {
  levels: StatsLevelLike[]
  currentLevel: number
  latency: number
  targetLatency: number | null
  bandwidthEstimate: number
}

/** Minimal shape of a media element needed for statistics. */
export interface StatsMediaLike {
  currentTime: number
  playbackRate: number
  videoWidth: number
  videoHeight: number
  buffered: TimeRanges
  seekable: TimeRanges
  getVideoPlaybackQuality?: () => {
    droppedVideoFrames: number
    totalVideoFrames: number
  }
}

export interface PlayerStats {
  /** Distance (seconds) from the live edge. */
  latency: number | null
  /** Latency (seconds) hls.js is steering towards. */
  targetLatency: number | null
  /** Forward buffer (seconds) ahead of the playback position. */
  buffer: number | null
  playbackRate: number | null
  /** Height of the rendition being played, e.g. `1080`. */
  levelHeight: number | null
  levelFrameRate: number | null
  /** Declared bitrate (bits/s) of the rendition being played. */
  levelBitrate: number | null
  /** hls.js bandwidth estimate (bits/s). */
  bandwidthEstimate: number | null
  /** Decoded video resolution. */
  videoSize: [number, number] | null
  /** Player element size in CSS pixels. */
  viewportSize: [number, number] | null
  droppedFrames: number | null
  totalFrames: number | null
  codecs: string | null
  protocol: string | null
}

export const EMPTY_PLAYER_STATS: PlayerStats = {
  latency: null,
  targetLatency: null,
  buffer: null,
  playbackRate: null,
  levelHeight: null,
  levelFrameRate: null,
  levelBitrate: null,
  bandwidthEstimate: null,
  videoSize: null,
  viewportSize: null,
  droppedFrames: null,
  totalFrames: null,
  codecs: null,
  protocol: null
}

function finiteOrNull(
  value: number | null | undefined,
  min = 0
): number | null {
  return value != null && Number.isFinite(value) && value >= min ? value : null
}

function positiveOrNull(value: number | null | undefined): number | null {
  return value != null && Number.isFinite(value) && value > 0 ? value : null
}

function sizeOrNull(width: number, height: number): [number, number] | null {
  return width > 0 && height > 0 ? [width, height] : null
}

/** Seconds buffered ahead of `currentTime`, or null when it is not buffered. */
export function forwardBuffer(media: StatsMediaLike): number | null {
  const { buffered, currentTime } = media
  for (let i = 0; i < buffered.length; i++) {
    if (currentTime >= buffered.start(i) && currentTime <= buffered.end(i)) {
      return buffered.end(i) - currentTime
    }
  }
  return null
}

/** Native playback latency, measured from the last seekable position. */
export function nativeLatency(media: StatsMediaLike): number | null {
  const { seekable } = media
  if (seekable.length === 0) {
    return null
  }
  return finiteOrNull(seekable.end(seekable.length - 1) - media.currentTime)
}

export function readPlayerStats(input: {
  media: StatsMediaLike | null
  hls: StatsHlsLike | null
  viewport: { width: number; height: number } | null
}): PlayerStats {
  const { media, hls, viewport } = input
  const stats: PlayerStats = { ...EMPTY_PLAYER_STATS }

  if (viewport) {
    stats.viewportSize = sizeOrNull(
      Math.round(viewport.width),
      Math.round(viewport.height)
    )
  }

  if (media) {
    stats.buffer = forwardBuffer(media)
    stats.playbackRate = finiteOrNull(media.playbackRate)
    stats.videoSize = sizeOrNull(media.videoWidth, media.videoHeight)
    const quality = media.getVideoPlaybackQuality?.()
    if (quality) {
      stats.droppedFrames = quality.droppedVideoFrames
      stats.totalFrames = quality.totalVideoFrames
    }
  }

  if (!hls) {
    if (media) {
      stats.latency = nativeLatency(media)
      stats.protocol = 'HLS (native)'
    }
    return stats
  }

  stats.latency = finiteOrNull(hls.latency)
  stats.targetLatency = finiteOrNull(hls.targetLatency)
  stats.bandwidthEstimate = positiveOrNull(hls.bandwidthEstimate)

  const level = hls.levels[hls.currentLevel]
  if (level) {
    stats.levelHeight = positiveOrNull(level.height)
    stats.levelFrameRate = positiveOrNull(level.frameRate)
    stats.levelBitrate = positiveOrNull(level.bitrate)
    const codecs = [level.videoCodec, level.audioCodec].filter(Boolean)
    stats.codecs = codecs.length ? codecs.join(', ') : null
    stats.protocol = level.details?.partTarget ? 'LL-HLS' : 'HLS'
  } else {
    stats.protocol = 'HLS'
  }

  return stats
}

/** `1080p60` style rendition label; the frame rate is only shown above 30 fps. */
export function formatRendition(
  height: number | null,
  frameRate: number | null
): string | null {
  if (!height) {
    return null
  }
  const fps = frameRate && frameRate > 30 ? Math.round(frameRate) : ''
  return `${height}p${fps}`
}

export function formatSeconds(value: number | null, digits = 2): string {
  return value === null ? '—' : `${value.toFixed(digits)} s`
}

export function formatBitrate(bitsPerSecond: number | null): string {
  if (bitsPerSecond === null) {
    return '—'
  }
  if (bitsPerSecond >= 1e6) {
    return `${(bitsPerSecond / 1e6).toFixed(1)} Mbps`
  }
  return `${Math.round(bitsPerSecond / 1e3)} Kbps`
}

export function formatSize(size: [number, number] | null): string {
  return size === null ? '—' : `${size[0]}×${size[1]}`
}
