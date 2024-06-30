<script setup lang="ts">
import { useStreamsList } from '@/api/stream.ts'

const { data: streams } = useStreamsList()
</script>

<template>
  <div class="w-full h-full justify-center mt-4 container">
    <h1 class="text-2xl">Browse channels</h1>

    <h1
      class="text-xl text-accent-foreground text-center"
      v-if="!streams?.streams.length"
    >
      No one online
    </h1>

    <div
      v-else
      class="flex pt-4 gap-4"
    >
      <router-link
        v-for="stream of streams?.streams"
        :key="stream.channel.id"
        class="flex flex-col border border-border rounded bg-accent relative w-80"
        :to="{
          name: 'Channel',
          params: { channelName: stream.channel.name }
        }"
      >
        <svg
          width="24"
          height="24"
          viewBox="0 0 24 24"
          fill="none"
          xmlns="http://www.w3.org/2000/svg"
          class="absolute top-[-12px] left-[-15px] text-green-400"
          aria-hidden
        >
          <circle
            cx="12"
            cy="12"
            r="5"
            fill="currentColor"
          />
          <circle
            cx="12"
            cy="12"
            r="6"
            fill="currentColor"
            style="
              animation: ping 1s cubic-bezier(0, 0, 0.2, 1) infinite;
              transform-origin: center center;
            "
          />
        </svg>
        <div class="p-4 h-auto">
          {{ stream.channel.displayName }}
        </div>
      </router-link>
    </div>
  </div>
</template>

<style scoped></style>
