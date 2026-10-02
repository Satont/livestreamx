/**
 * Live edge synchronization policy for the stream player.
 *
 * The controller is intentionally DOM-free: media access is provided through
 * `LiveEdgeAdapter`, which keeps the pending/seek logic unit-testable.
 */

/** Minimum pause duration (ms) after which resuming playback returns to the live edge. */
export const LIVE_EDGE_RESUME_THRESHOLD_MS = 1000

/**
 * Minimum distance (seconds) behind the live edge that justifies a recovery seek.
 * Resuming from a pause requests a correction with `0`, so any positive drift is fixed.
 */
export const MIN_LIVE_EDGE_SEEK_DISTANCE = 2

/** Native live playback (e.g. Safari HLS) resumes this far behind the last seekable position. */
export const NATIVE_LIVE_EDGE_OFFSET = 3

/** Minimum time (ms) between two player-triggered live edge seeks. */
export const LIVE_EDGE_SEEK_COOLDOWN_MS = 3000

/**
 * hls.js latency policy for the LL-HLS ladder served by OvenMediaEngine.
 *
 * - `liveSyncDuration` is the target distance from the live edge. OME publishes a
 *   `PART-HOLD-BACK` of 1.5s (ChunkDuration 0.5s); 3s keeps a safety margin.
 * - `liveMaxLatencyDuration` makes hls.js hard seek back once playback falls
 *   further than 8s behind, which bounds latency after stalls or tab throttling.
 * - `maxLiveSyncPlaybackRate` lets the latency controller catch up smoothly
 *   instead of only seeking.
 * - `liveSyncOnStallIncrease: 0` stops hls.js from raising the target latency
 *   after every buffer stall, which used to leave the player permanently behind.
 */
export const HLS_LIVE_CONFIG = {
  lowLatencyMode: true,
  liveSyncDuration: 3,
  liveMaxLatencyDuration: 8,
  maxLiveSyncPlaybackRate: 1.15,
  liveSyncOnStallIncrease: 0
} as const

/** Delays (ms) between attempts to restart hls.js loading after fatal network errors. */
export const HLS_NETWORK_RESTART_DELAYS_MS = [2000, 5000, 10000]

/** Minimal shape of an hls.js instance needed to resolve the live edge. */
export interface HlsLike {
  liveSyncPosition: number | null
  media: HTMLMediaElement | null
}

/** Minimal shape of a media element needed to resolve the native live edge. */
export interface MediaLike {
  seekable: TimeRanges
}

/**
 * Resolves the hls.js live edge position. Returns null while the playlist is not
 * ready or when hls.js cannot compute a finite position.
 */
export function resolveHlsLiveEdge(hls: HlsLike | null | undefined): number | null {
  if (!hls?.media) {
    return null
  }
  const position = hls.liveSyncPosition
  return position != null && Number.isFinite(position) ? position : null
}

/**
 * Resolves the live edge for native playback (e.g. Safari HLS) from the last
 * seekable range, staying `offset` seconds behind it to avoid seek failures.
 */
export function resolveNativeLiveEdge(
  media: MediaLike | null | undefined,
  offset: number = NATIVE_LIVE_EDGE_OFFSET
): number | null {
  const seekable = media?.seekable
  if (!seekable || seekable.length === 0) {
    return null
  }
  const range = seekable.length - 1
  const end = seekable.end(range)
  const start = seekable.start(range)
  if (!Number.isFinite(end)) {
    return null
  }
  const target = end - offset
  return target > start ? target : null
}

/** Media access required by `LiveEdgeSyncController`. */
export interface LiveEdgeAdapter {
  /** Live edge target in seconds, or null while it cannot be determined. */
  getTarget(): number | null
  /** Current playback position in seconds, or null while media is not ready. */
  getCurrentTime(): number | null
  isPaused(): boolean
  isSeeking(): boolean
  seek(position: number): void
  now(): number
}

export interface LiveEdgeSyncOptions {
  /** Pause duration (ms) that triggers a correction on resume. */
  resumeThresholdMs?: number
  /** Default drift (seconds) that justifies a recovery seek. */
  minSeekDistance?: number
  /** Cooldown (ms) between two player-triggered seeks. */
  cooldownMs?: number
}

export type LiveEdgeAttemptResult = 'seeked' | 'pending' | 'skipped'

/**
 * Decides when playback should seek back to the live edge.
 *
 * A correction can be requested by a pause resume, by returning to the tab, or
 * by a network recovery. When the live edge position or media is not ready yet,
 * the intent is kept as pending and retried by `retry()` on fresh playlist or
 * media events.
 */
export class LiveEdgeSyncController {
  private pausedAt: number | null = null
  private pending = false
  private pendingMinDistance: number
  private lastSeekAt = Number.NEGATIVE_INFINITY
  private disposed = false

  constructor(
    private readonly adapter: LiveEdgeAdapter,
    private readonly options: LiveEdgeSyncOptions = {}
  ) {
    this.pendingMinDistance = options.minSeekDistance ?? MIN_LIVE_EDGE_SEEK_DISTANCE
  }

  get isPending(): boolean {
    return this.pending && !this.disposed
  }

  /** Records the moment the user paused playback. */
  pause(at: number = this.adapter.now()): void {
    this.pausedAt = at
  }

  /**
   * Handles the `play` event. When the pause was long enough, asks for a
   * correction with no minimum drift, so any pause delay is removed.
   */
  play(at: number = this.adapter.now()): boolean {
    const pausedFor = this.pausedAt === null ? 0 : at - this.pausedAt
    this.pausedAt = null

    if (pausedFor < (this.options.resumeThresholdMs ?? LIVE_EDGE_RESUME_THRESHOLD_MS)) {
      return false
    }
    return this.request(0) === 'seeked'
  }

  /**
   * Requests a live edge correction. `minDistance` is the drift (seconds) that
   * justifies the seek; pass `0` to correct any positive drift.
   */
  request(
    minDistance: number = this.options.minSeekDistance ?? MIN_LIVE_EDGE_SEEK_DISTANCE
  ): LiveEdgeAttemptResult {
    if (this.disposed) {
      return 'skipped'
    }
    this.pending = true
    this.pendingMinDistance = minDistance
    return this.attempt(minDistance)
  }

  /** Retries a pending correction after fresh playlist or media events. */
  retry(): LiveEdgeAttemptResult {
    if (this.disposed || !this.pending) {
      return 'skipped'
    }
    return this.attempt(this.pendingMinDistance)
  }

  /** Clears all pending state, e.g. after switching the stream. */
  reset(): void {
    this.pausedAt = null
    this.pending = false
    this.pendingMinDistance = this.options.minSeekDistance ?? MIN_LIVE_EDGE_SEEK_DISTANCE
    this.lastSeekAt = Number.NEGATIVE_INFINITY
  }

  dispose(): void {
    this.disposed = true
    this.reset()
  }

  private attempt(minDistance: number): LiveEdgeAttemptResult {
    if (this.adapter.isPaused() || this.adapter.isSeeking()) {
      return 'pending'
    }

    const target = this.adapter.getTarget()
    const currentTime = this.adapter.getCurrentTime()
    if (target === null || currentTime === null) {
      return 'pending'
    }

    const distance = target - currentTime
    if (distance <= 0 || distance < minDistance) {
      this.pending = false
      return 'skipped'
    }

    const now = this.adapter.now()
    if (now - this.lastSeekAt < (this.options.cooldownMs ?? LIVE_EDGE_SEEK_COOLDOWN_MS)) {
      return 'pending'
    }

    this.adapter.seek(target)
    this.lastSeekAt = now
    this.pending = false
    return 'seeked'
  }
}
