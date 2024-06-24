<script setup lang="ts">
import { formatDuration, intervalToDuration } from 'date-fns'
import { computed } from 'vue'

import { useStream } from '@/api/stream.ts'

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
    <!--    <svg-->
    <!--      width="24"-->
    <!--      height="24"-->
    <!--      viewBox="0 0 24 24"-->
    <!--      fill="none"-->
    <!--      xmlns="http://www.w3.org/2000/svg"-->
    <!--      class="absolute top-[-12px] left-[-15px]"-->
    <!--      :class="{-->
    <!--        'text-green-400': !!uptime,-->
    <!--        'text-red-400': !uptime-->
    <!--      }"-->
    <!--      aria-hidden-->
    <!--    >-->
    <!--      <circle-->
    <!--        cx="12"-->
    <!--        cy="12"-->
    <!--        r="5"-->
    <!--        fill="currentColor"-->
    <!--      />-->
    <!--      <circle-->
    <!--        cx="12"-->
    <!--        cy="12"-->
    <!--        r="6"-->
    <!--        fill="currentColor"-->
    <!--        style="-->
    <!--          animation: ping 1s cubic-bezier(0, 0, 0.2, 1) infinite;-->
    <!--          transform-origin: center center;-->
    <!--        "-->
    <!--        v-if="!!uptime"-->
    <!--      />-->
    <!--    </svg>-->
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
