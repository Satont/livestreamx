<script setup lang="ts">
import { FragmentType, graphql, useFragment } from "@/gql";
import { showAvatars } from "@/composables/show-avatars.js";
import { showTimestamps } from "@/composables/show-timestamps.ts";
import { useProfile } from "@/api/profile.ts";
import { MessageSegmentType } from "@/gql/graphql.ts";

const ChatMessage_Fragment = graphql(`
    fragment MessageFragment on ChatMessage {
			id
			segments {
					type
					content
					...on MessageSegmentMention {
						user {
							id
							color
							displayName
						}
					}
			}
			sender {
					id
					avatarUrl
					color
					createdAt
					name
					displayName
			}
			createdAt
    }
`)

type Props = {
  message: FragmentType<typeof ChatMessage_Fragment>
}

const props = defineProps<Props>()

const data = useFragment(ChatMessage_Fragment, props.message)

const { data: profile } = useProfile()
</script>

<template>
	<div>
		<p>
			<span v-if="showTimestamps" class="mr-1 opacity-50">
				{{ new Date(data.createdAt)
					.toLocaleTimeString('en', { hour12: false, hour: '2-digit', minute:'2-digit' })
				}}
			</span>
			<span>
				<span class="inline-flex items-center" v-if="showAvatars">
					<img :src="data.sender.avatarUrl" class="size-4 rounded-full mr-1" />
				</span>
				<span class="font-bold" :style="{ color: data.sender.color }">
					{{ data.sender.displayName }}
				</span>
			</span>
			<span>: </span>
			<span class="break-words">
				<template v-for="segment of data.segments">
					<template v-if="segment.type === MessageSegmentType.Text">{{ segment.content }}</template>
					<span v-else-if="segment.type === MessageSegmentType.Mention && 'user' in segment">
						<span
							:style="{ color: segment.user.color }"
							class="bg-zinc-400 p-0.5 rounded"
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