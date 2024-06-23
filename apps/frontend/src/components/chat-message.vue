<script setup lang="ts">
import { showAvatars } from "@/composables/show-avatars.js";
import { showTimestamps } from "@/composables/show-timestamps.js";
import { useProfile } from "@/api/profile.ts";
import { ChatMessage_FragmentFragment, ChatEmote_FragmentFragment, MessageSegmentType } from "@/gql/graphql.ts";
import { chatFontSize } from "@/composables/chat-font-size.js";
import {
	Tooltip,
	TooltipContent,
	TooltipTrigger
} from '@/components/ui/tooltip'
import { calculateColor } from "@/lib/color.js";
import { colorMode } from "@/composables/color-mode.ts";
import { Button } from "@/components/ui/button";
import { Copy } from 'lucide-vue-next'

type Props = {
  message: ChatMessage_FragmentFragment
}
const props = defineProps<Props>()

const { data: profile } = useProfile()

function correctColor(color: string) {
	return calculateColor(color, colorMode.value === 'dark')
}

function copyText() {
	navigator.clipboard.writeText(props.message.segments.map(s => s.content).join(' '))
}
</script>

<template>
	<div
		:style="{ fontSize: `${chatFontSize}px` }"
		class="relative group hover:bg-accent hover:rounded p-0.5"
	>
		<p class="leading-7">
			<span v-if="showTimestamps" class="mr-1 opacity-50">
				{{ new Date(message.createdAt)
					.toLocaleTimeString('en', { hour12: false, hour: '2-digit', minute:'2-digit' })
				}}
			</span>
			<span>
				<span class="inline-flex align-sub" v-if="showAvatars">
					<img :src="message.sender.avatarUrl" class="size-4 rounded-full mr-1" />
				</span>
				<span class="font-bold" :style="{ color: correctColor(message.sender.color) }">
					{{ message.sender.displayName }}
				</span>
			</span>
			<span>: </span>
			<span class="break-words">
				<template v-for="segment of message.segments">
					<template v-if="segment.type === MessageSegmentType.Text">{{ segment.content }}</template>
					<span
						v-else-if="segment.type === MessageSegmentType.Mention && 'user' in segment"
						:style="{ color: correctColor(segment.user.color) }"
						class="p-0.5 rounded"
						:class="{ 'bg-zinc-400': segment.user.id === profile?.userProfile.id }"
					>
						@{{ segment.user.displayName }}
					</span>
					<a
						v-else-if="segment.type === MessageSegmentType.Link"
						:href="segment.content"
						target="_blank"
						class="underline"
					>
						{{ segment.content }}
					</a>
					<template v-else-if="segment.type === MessageSegmentType.Emote && 'emote' in segment">
						<Tooltip>
							<TooltipTrigger>
								<img
									:src="(segment.emote as ChatEmote_FragmentFragment).url"
									:style="{
										width: `${(segment.emote as ChatEmote_FragmentFragment).width}px`,
										height: `${(segment.emote as ChatEmote_FragmentFragment).height}px`
									}"
									class="scale-90 inline-block relative"
								/>
							</TooltipTrigger>
							<TooltipContent>
								<div class="flex flex-col">
									<img
										:src="(segment.emote as ChatEmote_FragmentFragment).url.replace('1x.webp', '4x.webp')"
										:style="{
											width: `${(segment.emote as ChatEmote_FragmentFragment).width * 2.5}px`,
											height: `${(segment.emote as ChatEmote_FragmentFragment).height * 2.5}px`
										}"
									/>
									<h1 class="place-self-center text-lg font-bold">{{ (segment.emote as ChatEmote_FragmentFragment).name }}</h1>
								</div>
							</TooltipContent>
						</Tooltip>
					</template>
					{{ ' ' }}
				</template>
			</span>
		</p>

		<Button class="hidden group-hover:block absolute right-0 top-0" @click="copyText" size="xs" variant="secondary">
			<Copy class="size-4" />
		</Button>
	</div>
</template>