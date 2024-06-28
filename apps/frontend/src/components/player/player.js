import Hls from 'hls.js'

const retryPause = 2000

let video
let message
let streamUrl

let defaultControls = false

const setMessage = (str) => {
  if (str !== '') {
    video.controls = false
  } else {
    video.controls = defaultControls
  }
  message.innerText = str
}

const loadStream = () => {
  // always prefer hls.js over native HLS.
  // this is because some Android versions support native HLS
  // but don't support fMP4s.
  if (Hls.isSupported()) {
    const hls = new Hls({
      maxLiveSyncPlaybackRate: 1.5
    })

    hls.on(Hls.Events.ERROR, (evt, data) => {
      if (data.fatal) {
        hls.destroy()

        if (data.details === 'manifestIncompatibleCodecsError') {
          setMessage(
            'stream makes use of codecs which are incompatible with this browser or operative system'
          )
        } else if (data.response && data.response.code === 404) {
          setMessage('stream not found, retrying in some seconds')
        } else {
          setMessage(data.error + ', retrying in some seconds')
        }

        setTimeout(() => loadStream(video), retryPause)
      }
    })

    hls.on(Hls.Events.MEDIA_ATTACHED, () => {
      const u = new URL('index.m3u8', streamUrl)
      hls.loadSource(u.toString())
    })

    hls.on(Hls.Events.MANIFEST_PARSED, () => {
      setMessage('')
      video.play()
    })

    hls.attachMedia(video)
  } else if (video.canPlayType('application/vnd.apple.mpegurl')) {
    // since it's not possible to detect timeout errors in iOS,
    // wait for the playlist to be available before starting the stream

    const u = new URL('index.m3u8', streamUrl)
    fetch(u.toString()).then(() => {
      video.src = 'index.m3u8'
      video.play()
    })
  }
}

export const init = ({
  videoEl,
  messageEl,
  streamSrc,
  controls = true,
  muted = true,
  autoplay = true,
  playsinline = true
}) => {
  video = videoEl
  message = messageEl
  streamUrl = streamSrc

  video.controls = controls
  video.muted = muted
  video.autoplay = autoplay
  video.playsInline = playsinline
  defaultControls = controls

  loadStream()
}
