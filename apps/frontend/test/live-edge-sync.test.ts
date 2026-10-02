import { describe, expect, test } from 'bun:test'

import {
  HLS_LIVE_CONFIG,
  LIVE_EDGE_RESUME_THRESHOLD_MS,
  LiveEdgeSyncController,
  MIN_LIVE_EDGE_SEEK_DISTANCE,
  NATIVE_LIVE_EDGE_OFFSET,
  resolveHlsLiveEdge,
  resolveNativeLiveEdge,
  type LiveEdgeAdapter
} from '../app/utils/live-edge-sync'

interface AdapterState {
  target: number | null
  currentTime: number | null
  paused: boolean
  seeking: boolean
  now: number
  seeks: number[]
}

function createAdapter(overrides: Partial<AdapterState> = {}) {
  const state: AdapterState = {
    target: 100,
    currentTime: 90,
    paused: false,
    seeking: false,
    now: 0,
    seeks: [],
    ...overrides
  }

  const adapter: LiveEdgeAdapter = {
    getTarget: () => state.target,
    getCurrentTime: () => state.currentTime,
    isPaused: () => state.paused,
    isSeeking: () => state.seeking,
    seek: (position) => state.seeks.push(position),
    now: () => state.now
  }

  return { adapter, state }
}

function createSeekable(ranges: Array<[number, number]>): TimeRanges {
  return {
    length: ranges.length,
    start: (index: number) => ranges[index][0],
    end: (index: number) => ranges[index][1]
  } as TimeRanges
}

describe('resolveHlsLiveEdge', () => {
  test('returns the finite live sync position once media is attached', () => {
    expect(
      resolveHlsLiveEdge({ liveSyncPosition: 123.5, media: {} as HTMLMediaElement })
    ).toBe(123.5)
  })

  test('returns null without an instance, media or finite position', () => {
    expect(resolveHlsLiveEdge(null)).toBeNull()
    expect(
      resolveHlsLiveEdge({ liveSyncPosition: 10, media: null })
    ).toBeNull()
    expect(
      resolveHlsLiveEdge({ liveSyncPosition: null, media: {} as HTMLMediaElement })
    ).toBeNull()
    expect(
      resolveHlsLiveEdge({
        liveSyncPosition: Number.POSITIVE_INFINITY,
        media: {} as HTMLMediaElement
      })
    ).toBeNull()
    expect(
      resolveHlsLiveEdge({
        liveSyncPosition: Number.NaN,
        media: {} as HTMLMediaElement
      })
    ).toBeNull()
  })
})

describe('resolveNativeLiveEdge', () => {
  test('stays the configured offset behind the last seekable position', () => {
    const media = { seekable: createSeekable([[0, 100]]) }
    expect(resolveNativeLiveEdge(media)).toBe(100 - NATIVE_LIVE_EDGE_OFFSET)
    expect(resolveNativeLiveEdge(media, 10)).toBe(90)
  })

  test('uses the last seekable range when the window starts later', () => {
    const media = { seekable: createSeekable([[0, 40], [60, 100]]) }
    expect(resolveNativeLiveEdge(media, 5)).toBe(95)
  })

  test('returns null without usable seekable data', () => {
    expect(resolveNativeLiveEdge(null)).toBeNull()
    expect(resolveNativeLiveEdge({ seekable: createSeekable([]) })).toBeNull()
    expect(
      resolveNativeLiveEdge({ seekable: createSeekable([[0, Number.POSITIVE_INFINITY]]) })
    ).toBeNull()
    expect(resolveNativeLiveEdge({ seekable: createSeekable([[98, 100]]) }, 3)).toBeNull()
  })
})

describe('LiveEdgeSyncController', () => {
  test('seeks to the live edge when resuming after a long pause', () => {
    const { adapter, state } = createAdapter({ now: 5000 })
    const controller = new LiveEdgeSyncController(adapter)

    controller.pause(1000)

    expect(controller.play(1000 + LIVE_EDGE_RESUME_THRESHOLD_MS)).toBe(true)
    expect(state.seeks).toEqual([100])
  })

  test('does nothing when resuming after a short pause or without a pause', () => {
    const { adapter, state } = createAdapter({ now: 5000 })
    const controller = new LiveEdgeSyncController(adapter)

    controller.pause(4999)
    expect(controller.play(5000)).toBe(false)
    expect(controller.play(6000)).toBe(false)
    expect(state.seeks).toEqual([])
  })

  test('keeps the correction pending until the live edge is known', () => {
    const { adapter, state } = createAdapter({ target: null })
    const controller = new LiveEdgeSyncController(adapter)

    expect(controller.request()).toBe('pending')
    expect(controller.isPending).toBe(true)

    state.target = 100
    expect(controller.retry()).toBe('seeked')
    expect(controller.isPending).toBe(false)
    expect(state.seeks).toEqual([100])
  })

  test('ignores retries without a pending correction', () => {
    const { adapter, state } = createAdapter()
    const controller = new LiveEdgeSyncController(adapter)

    expect(controller.retry()).toBe('skipped')
    expect(state.seeks).toEqual([])
  })

  test('skips small drift by default and corrects it when resuming from a pause', () => {
    const { adapter, state } = createAdapter({ currentTime: 99 })
    const controller = new LiveEdgeSyncController(adapter)

    expect(controller.request()).toBe('skipped')
    expect(controller.request(0)).toBe('seeked')
    expect(state.seeks).toEqual([100])
    expect(MIN_LIVE_EDGE_SEEK_DISTANCE).toBeGreaterThan(0)
  })

  test('does not seek while the user keeps the player paused', () => {
    const { adapter, state } = createAdapter({ paused: true })
    const controller = new LiveEdgeSyncController(adapter)

    expect(controller.request()).toBe('pending')
    expect(controller.isPending).toBe(true)
    expect(state.seeks).toEqual([])

    state.paused = false
    expect(controller.retry()).toBe('seeked')
    expect(state.seeks).toEqual([100])
  })

  test('does not seek while another seek is in progress', () => {
    const { adapter, state } = createAdapter({ seeking: true })
    const controller = new LiveEdgeSyncController(adapter)

    expect(controller.request()).toBe('pending')

    state.seeking = false
    expect(controller.retry()).toBe('seeked')
    expect(state.seeks).toEqual([100])
  })

  test('does not seek when playback is already at or ahead of the target', () => {
    const { adapter, state } = createAdapter({ currentTime: 100 })
    const controller = new LiveEdgeSyncController(adapter)

    expect(controller.request()).toBe('skipped')
    expect(controller.isPending).toBe(false)
    expect(state.seeks).toEqual([])
  })

  test('cooldown prevents repeated seeks but keeps the pending correction', () => {
    const { adapter, state } = createAdapter({ now: 0 })
    const controller = new LiveEdgeSyncController(adapter)

    expect(controller.request()).toBe('seeked')

    state.now = 1000
    expect(controller.request()).toBe('pending')
    expect(state.seeks).toEqual([100])

    state.now = 3000
    expect(controller.retry()).toBe('seeked')
    expect(state.seeks).toEqual([100, 100])
  })

  test('reset clears pause and pending state', () => {
    const { adapter, state } = createAdapter({ target: null })
    const controller = new LiveEdgeSyncController(adapter)

    controller.pause(0)
    controller.request()
    controller.reset()

    expect(controller.isPending).toBe(false)

    state.target = 100
    expect(controller.retry()).toBe('skipped')
    expect(controller.play(5000)).toBe(false)
    expect(state.seeks).toEqual([])
  })

  test('dispose prevents any further seeks', () => {
    const { adapter, state } = createAdapter()
    const controller = new LiveEdgeSyncController(adapter)

    controller.dispose()

    expect(controller.request()).toBe('skipped')
    expect(controller.isPending).toBe(false)
    expect(state.seeks).toEqual([])
  })
})

describe('HLS_LIVE_CONFIG', () => {
  test('enables low latency with a bounded target and no stall creep', () => {
    expect(HLS_LIVE_CONFIG.lowLatencyMode).toBe(true)
    expect(HLS_LIVE_CONFIG.liveSyncDuration).toBeLessThan(
      HLS_LIVE_CONFIG.liveMaxLatencyDuration
    )
    expect(HLS_LIVE_CONFIG.liveSyncOnStallIncrease).toBe(0)
    expect(Number.isFinite(HLS_LIVE_CONFIG.liveMaxLatencyDuration)).toBe(true)
  })
})
