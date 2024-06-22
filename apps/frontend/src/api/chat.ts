import { createGlobalState } from "@vueuse/core";
import { ref } from "vue";
import { useMutation, useQuery, useSubscription } from "@urql/vue";
import { graphql, useFragment } from "@/gql";
import { watch } from "vue";
import { MessageFragmentFragment , EmoteFragmentFragment } from "@/gql/graphql.js";

export const ChatMessage_Fragment = graphql(`
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
						...on MessageSegmentEmote {
								emote {
										id
										name
										url
										width
										height
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
		fragment EmoteFragment on Emote {
				id
				name
				url
		}
`)

export const useChat = createGlobalState(() => {
	const messages = ref<MessageFragmentFragment[]>([])

	const sub = useSubscription({
		query: graphql(`
			subscription NewChatMessages {
				chatMessages {
					...MessageFragment
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
					...MessageFragment
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

	const emotes = ref<EmoteFragmentFragment[]>([])

	const useQueryEmotes = useQuery({
		query: graphql(`
			query ChatEmotes {
				getEmotes {
					...EmoteFragment
				}
			}
		`),
		variables: {},
	})

	watch(useQueryEmotes.data, (data) => {
		console.log(emotes)
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