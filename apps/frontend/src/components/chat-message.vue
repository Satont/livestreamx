<script setup lang="ts">
import { showAvatars } from "@/composables/show-avatars.js";
import { showTimestamps } from "@/composables/show-timestamps.ts";
import { useProfile } from "@/api/profile.ts";
import { MessageFragmentFragment, MessageSegmentType } from "@/gql/graphql.ts";

type Props = {
  message: MessageFragmentFragment
}

defineProps<Props>()

const { data: profile } = useProfile()
</script>

<template>
	<div>
		<p>
			<span v-if="showTimestamps" class="mr-1 opacity-50">
				{{ new Date(message.createdAt)
					.toLocaleTimeString('en', { hour12: false, hour: '2-digit', minute:'2-digit' })
				}}
			</span>
			<span>
				<span class="inline-flex items-center" v-if="showAvatars">
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
					{{ ' ' }}
				</template>
			</span>
		</p>
	</div>
</template>