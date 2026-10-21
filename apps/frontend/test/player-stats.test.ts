import { describe, expect, test } from 'bun:test'

import {
  EMPTY_PLAYER_STATS,
  formatBitrate,
  formatRendition,
  formatSeconds,
  formatSize,
  forwardBuffer,
  readPlayerStats,
  type StatsHlsLike,
  type StatsMediaLike
} from '../app/utils/player-stats'

function createRanges(ranges: Array<[number, number]>): TimeRanges {
  return {
    length: ranges.length,
    start: (index: number) => ranges[index][0],
    end: (index: number) => ranges[index][1]
  } as TimeRanges
}

function createMedia(overrides: Partial<StatsMediaLike> = {}): StatsMediaLike {
  return {
    currentTime: 50,
    playbackRate: 1,
    videoWidth: 1920,
    videoHeight: 1080,
    buffered: createRanges([[40, 52.5]]),
    seekable: createRanges([[0, 53]]),
    getVideoPlaybackQuality: () => ({ droppedVideoFrames: 3, totalVideoFrames: 900 }),
    ...overrides
  }
}

function createHls(overrides: Partial<StatsHlsLike> = {}): StatsHlsLike {
  return {
    levels: [
      {
        width: 1280,
        height: 720,
        bitrate: 3_000_000,
        frameRate: 30,
        videoCodec: 'avc1.4d401f',
        audioCodec: 'mp4a.40.2'
      },
      {
        width: 1920,
        height: 1080,
        bitrate: 8_000_000,
        frameRate: 60,
        videoCodec: 'avc1.64002a',
        audioCodec: 'mp4a.40.2',
        details: { partTarget: 0.5 }
      }
    ],
    currentLevel: 1,
    latency: 2.4,
    targetLatency: 3,
    bandwidthEstimate: 25_000_000,
    ...overrides
  }
}

describe('forwardBuffer', () => {
  test('measures the range that contains the playback position', () => {
    const media = createMedia({
      currentTime: 50,
      buffered: createRanges([
        [10, 20],
        [45, 51.5]
      ])
    })
    expect(forwardBuffer(media)).toBe(1.5)
  })

  test('returns null when the playback position is not buffered', () => {
    expect(forwardBuffer(createMedia({ buffered: createRanges([[10, 20]]) }))).toBeNull()
  })
})

describe('readPlayerStats', () => {
  test('reads the playing hls.js level and media state', () => {
    const stats = readPlayerStats({
      media: createMedia(),
      hls: createHls(),
      viewport: { width: 1372.4, height: 771.8 }
    })

    expect(stats).toEqual({
      latency: 2.4,
      targetLatency: 3,
      buffer: 2.5,
      playbackRate: 1,
      levelHeight: 1080,
      levelFrameRate: 60,
      levelBitrate: 8_000_000,
      bandwidthEstimate: 25_000_000,
      videoSize: [1920, 1080],
      viewportSize: [1372, 772],
      droppedFrames: 3,
      totalFrames: 900,
      codecs: 'avc1.64002a, mp4a.40.2',
      protocol: 'LL-HLS'
    })
  })

  test('drops values hls.js cannot measure yet', () => {
    const stats = readPlayerStats({
      media: null,
      hls: createHls({
        currentLevel: -1,
        latency: Number.NaN,
        targetLatency: null,
        bandwidthEstimate: Number.POSITIVE_INFINITY
      }),
      viewport: null
    })

    expect(stats).toEqual({ ...EMPTY_PLAYER_STATS, protocol: 'HLS' })
  })

  test('measures native playback latency from the seekable range', () => {
    const stats = readPlayerStats({
      media: createMedia({ getVideoPlaybackQuality: undefined }),
      hls: null,
      viewport: null
    })

    expect(stats.latency).toBe(3)
    expect(stats.protocol).toBe('HLS (native)')
    expect(stats.droppedFrames).toBeNull()
    expect(stats.levelHeight).toBeNull()
  })
})

describe('formatting', () => {
  test('formats renditions with a frame rate above 30 fps', () => {
    expect(formatRendition(1080, 60)).toBe('1080p60')
    expect(formatRendition(720, 30)).toBe('720p')
    expect(formatRendition(720, null)).toBe('720p')
    expect(formatRendition(null, 60)).toBeNull()
  })

  test('formats seconds, bitrates and sizes', () => {
    expect(formatSeconds(2.456)).toBe('2.46 s')
    expect(formatSeconds(2.456, 1)).toBe('2.5 s')
    expect(formatSeconds(null)).toBe('—')
    expect(formatBitrate(7_987_000)).toBe('8.0 Mbps')
    expect(formatBitrate(640_000)).toBe('640 Kbps')
    expect(formatBitrate(null)).toBe('—')
    expect(formatSize([1920, 1080])).toBe('1920×1080')
    expect(formatSize(null)).toBe('—')
  })
})
