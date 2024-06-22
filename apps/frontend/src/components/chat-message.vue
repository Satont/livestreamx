<script setup lang="ts">
import { showAvatars } from "@/composables/show-avatars.js";
import { showTimestamps } from "@/composables/show-timestamps.ts";
import { useProfile } from "@/api/profile.ts";
import { ChatMessage_FragmentFragment, ChatEmote_FragmentFragment, MessageSegmentType } from "@/gql/graphql.ts";
import { chatFontSize } from "@/composables/chat-font-size.ts";

type Props = {
  message: ChatMessage_FragmentFragment
}
defineProps<Props>()

const { data: profile } = useProfile()
</script>

<template>
	<div :style="{
		fontSize: `${chatFontSize}px`
	}">
		<p>
			<span v-if="showTimestamps" class="mr-1 opacity-50">
				{{ new Date(message.createdAt)
					.toLocaleTimeString('en', { hour12: false, hour: '2-digit', minute:'2-digit' })
				}}
			</span>
			<span>
				<span class="inline-flex align-sub" v-if="showAvatars">
					<img :src="message.sender.avatarUrl" class="size-4 rounded-full mr-1" />
				</span>
				<span class="font-bold" :style="{ color: message.sender.color }">
					{{ message.sender.displayName }}
				</span>
			</span>
			<span>: </span>
			<span class="break-words">
				<template v-for="segment of message.segments">
					<template v-if="segment.type === MessageSegmentType.Text">{{ segment.content }}</template>
					<span v-else-if="segment.type === MessageSegmentType.Mention && 'user' in segment">
						<span
							:style="{ color: segment.user.color }"
							class="p-0.5 rounded"
							:class="{ 'bg-zinc-400': segment.user.id === profile?.userProfile.id }"
						>
							@{{ segment.user.displayName }}
						</span>
					</span>
					<a
						v-else-if="segment.type === MessageSegmentType.Link"
						:href="segment.content"
						target="_blank"
						class="underline"
					>
						{{ segment.content }}
					</a>
					<img
						v-else-if="segment.type === MessageSegmentType.Emote && 'emote' in segment"
						:src="(segment.emote as ChatEmote_FragmentFragment).url"
						:style="{
							width: `${(segment.emote as ChatEmote_FragmentFragment).width}px`,
							height: `${(segment.emote as ChatEmote_FragmentFragment).height}px`
						}"
						class="inline-block relative"
					/>
					{{ ' ' }}
				</template>
			</span>
		</p>
	</div>
</template>