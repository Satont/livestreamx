import { useQuery, useSubscription } from '@urql/vue'
import { createGlobalState } from '@vueuse/core'

import { graphql } from '@/gql'
import { useChat } from './chat.js'

export const useStream = createGlobalState(() => {
  const { channelData } = useChat()

  const useStreamState = () =>
    useSubscription({
      query: graphql(`
        subscription StreamState($channelID: UUID!) {
          streamInfo(channelID: $channelID) {
            chatters {
              user {
                id
                displayName
                color
                createdAt
                name
                avatarUrl
              }
            }
            viewers
            startedAt
          }
        }
      `),
      get variables() {
        return {
          channelID: channelData.value?.fetchUserByName.id
        }
      },
      pause() {
        return !channelData.value?.fetchUserByName.id
      }
    })

  return {
    useStreamState
  }
})

export const useStreamsList = () => {
  return useQuery({
    query: graphql(`
      query StreamsList {
        streams {
          viewers
          viewers
          channel {
            name
            id
            displayName
            avatarUrl
          }
          startedAt
          thumbnailUrl
        }
      }
    `),
    variables: {}
  })
}
