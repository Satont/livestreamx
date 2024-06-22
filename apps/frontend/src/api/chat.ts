import { createGlobalState } from "@vueuse/core";
import { ref } from "vue";
import { useMutation, useQuery, useSubscription } from "@urql/vue";
import { graphql, useFragment } from "@/gql";
import { watch } from "vue";
import { ChatMessage_FragmentFragment , ChatEmote_FragmentFragment } from "@/gql/graphql.js";

export const ChatMessage_Fragment = graphql(`
    fragment ChatMessage_Fragment on ChatMessage {
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
						...on MessageSegmentEmote {
								emote {
										...ChatEmote_Fragment
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

export const ChatEmote_Fragment = graphql(`
		fragment ChatEmote_Fragment on Emote {
				id
				name
				url
				width
				height
		}
`)

export const useChat = createGlobalState(() => {
	const messages = ref<ChatMessage_FragmentFragment[]>([])

	const sub = useSubscription({
		query: graphql(`
			subscription NewChatMessages {
				chatMessages {
					...ChatMessage_Fragment
				}
			}
		`),
		variables: {},
	})

	watch(sub.data, (data) => {
		if (!data) return

		const fragment = useFragment(ChatMessage_Fragment, data.chatMessages)
		messages.value = [...messages.value, fragment]
	})

	const initialMessages = useQuery({
		query: graphql(`
			query ChatMessages {
				chatMessagesLatest {
					...ChatMessage_Fragment
				}
			}
		`),
		variables: {},
	})

	watch(initialMessages.data, (data) => {
		if (!data) return

		const fragments = useFragment(ChatMessage_Fragment, data.chatMessagesLatest)
		messages.value = fragments
	})

	const useSendMessage = () => useMutation(graphql(`
		mutation SendMessage($opts: SendMessageInput!) {
			sendMessage(input: $opts)
		}
	`))

	const emotes = ref<ChatEmote_FragmentFragment[]>([])

	const useQueryEmotes = useQuery({
		query: graphql(`
			query ChatEmotes {
				getEmotes {
					...ChatEmote_Fragment
				}
			}
		`),
		variables: {},
	})

	watch(useQueryEmotes.data, (data) => {
		if (!data) return

		const fragments = useFragment(ChatEmote_Fragment, data.getEmotes)
		emotes.value = fragments
	})

	return {
		messages,
		useSendMessage,
		emotes,
	}
})