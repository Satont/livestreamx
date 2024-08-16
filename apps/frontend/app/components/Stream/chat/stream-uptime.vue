<script setup lang="ts">
import { formatDuration, intervalToDuration } from 'date-fns'
import { computed } from 'vue'

import { useStream } from '@/api/stream.js'

const { data: streamState } = useStream().useStreamState()

const uptime = computed(() => {
  if (!streamState.value?.streamInfo?.startedAt) {
    return ''
  }

  const startedAt = new Date(streamState.value.streamInfo.startedAt)
  const now = new Date()

  const duration = intervalToDuration({
    start: startedAt,
    end: now
  })

  return formatDuration(duration)
})
</script>

<template>
  <span class="relative">
    {{ uptime || 'Offline' }}
  </span>
</template>

<style>
@keyframes ping {
  75%,
  to {
    transform: scale(2);
    opacity: 0;
  }
}
.animate-ping {
  animation: ping 1s cubic-bezier(0, 0, 0.2, 1) infinite;
}
@keyframes pulse {
  50% {
    opacity: 0.5;
  }
}
.animate-pulse {
  animation: pulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite;
}
</style>
