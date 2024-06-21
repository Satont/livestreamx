<script setup lang="ts">
import ThemeSwitcher from "@/components/theme-switcher.vue";
import ChatProfile from "@/components/chat-profile.vue";
import { useChat } from "@/api/chat.ts";
import ChatMessage from "@/components/chat-message.vue";
import { nextTick, ref, watch } from "vue";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import { useScroll } from "@vueuse/core";
import { useProfile } from "@/api/profile.ts";
import ChatSettings from "@/components/chat-settings.vue";
import { useStream } from "@/api/stream.ts";
import ChatViewers from "@/components/chat-viewers.vue";

const { data: profile } = useProfile();
const { messages, useSendMessage } = useChat()

const messageSender = useSendMessage()

const text = ref('')

const sendRetries = ref(0)

async function sendMessage() {
	if (!text.value) return;

	const msg = text.value.replace(/\s+/g, ' ').trim()
	if (!msg) return;
	if (msg.length > 700) {
		return;
	}

	while (sendRetries.value < 5) {
		const res = await messageSender.executeMutation({ opts: { text: msg }})
		if (res.error) {
			console.error(res.error)
			sendRetries.value++
		} else {
			text.value = ''
			sendRetries.value = 0
			break;
		}
	}
}

const messagesEl = ref<HTMLElement | null>(null)
const { y, arrivedState } = useScroll(messagesEl)

const scrollPaused = ref(false)

watch(arrivedState, (v) => {
	scrollPaused.value = !v.bottom
})

watch(messages, async () => {
	if (!messagesEl.value || scrollPaused.value) return;

	await nextTick()
	y.value = messagesEl.value?.scrollHeight
}, { immediate: true })
</script>

<template>
  <div class="flex h-full max-h-full flex-col">
    <div class="flex flex-row justify-between bg-secondary border-b-2 border-red-400 items-center px-4 min-w-48">
      <div class="flex items-center">
				<ChatViewers />
			</div>
			<div class="flex items-center">
				<ThemeSwitcher />
				<ChatProfile />
			</div>
    </div>
    <div ref="messagesEl" class="h-full max-w-96 relative flex flex-col overflow-y-auto pl-2">
      <ChatMessage
        v-for="message in messages"
        :key="message.id"
        :message="message"
      />
    </div>
    <div class="flex flex-col gap-2 bg-accent border-t-2 border-red-400 min-h-36 p-2">
      <Textarea v-model="text" @keydown.enter="sendMessage" :disabled="!profile" placeholder="Send message" />
			<div class="flex gap-2 place-self-end">
				<ChatSettings />
				<Button @click="sendMessage" size="sm" :disabled="!profile">Send message</Button>
			</div>
    </div>
  </div>
</template>

<style scoped>

</style>