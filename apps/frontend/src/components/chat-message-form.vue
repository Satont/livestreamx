<script setup lang="ts">
import { Button } from "@/components/ui/button";
import ChatSettings from "@/components/chat-settings.vue";
import { computed, ref } from "vue";
import { useProfile } from "@/api/profile.js";
import { useChat } from "@/api/chat.ts";
import { Textarea } from "@/components/ui/textarea";

// @ts-ignore
import Mention from "./mention.vue";

const { useSendMessage, emotes, messages } = useChat()
const { data: profile } = useProfile();

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

const emotesForMention = computed(() => {
	return emotes.value.map(e => ({
		label: e.name,
		url: e.url,
		value: e.name
	}))
})

const usersForMention = computed(() => {
	const mappedUsersFromMessages = messages.value
		.map(m => ({ label: m.sender.displayName, color: m.sender.color, value: m.sender.displayName }))
		.filter((v, i, a) => a.findIndex(t => (t.label === v.label)) === i)

	return [
		...mappedUsersFromMessages,
	]
})

const mentionKey = ref<'@' | ':'>('@')
function mapInsert(item: { label: string }) {
	return item.label + ' '
}
const mentionKeys = ['@', ':']
const mentionItems = computed(() => {
	return mentionKey.value === '@' ? usersForMention.value : emotesForMention.value
})
</script>

<template>
	<div class="flex flex-col gap-2 bg-accent border-t-2 border-red-400 min-h-36 p-2">
		<Mention
			:keys="mentionKeys"
			:items="mentionItems"
			offset="6"
			insert-space
			@open="(k: typeof mentionKey) => mentionKey = k"
			:omit-key="mentionKey !== '@'"
			:item-height="mentionKey === '@' ? 22 : 48"
			:map-insert="mapInsert"
		>
			<Textarea
				v-model="text"
				placeholder="Send message..."
			/>

			<template #no-result>
				<div class="dim">
					No result
				</div>
			</template>

			<template #item-@="{ item }">
				<Button class="w-full" variant="ghost" size="sm">
					{{ item.data.label }}
				</Button>
			</template>

			<template #item-:="{ item }">
				<div class="flex items-center gap-2 cursor-pointer">
					<img :src="item.data.url" class="size-10" />
					<span>{{ item.data.label }}</span>
				</div>
			</template>
		</Mention>

		<div class="flex gap-2 place-self-end">
			<ChatSettings />
			<Button @click="sendMessage" size="sm" :disabled="!profile">Send message</Button>
		</div>
	</div>
</template>
