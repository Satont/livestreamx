import { useSubscription } from '@urql/vue'
import { createGlobalState } from '@vueuse/core'

import { graphql } from '@/gql'

export const useStream = createGlobalState(() => {
  const useStreamState = () =>
    useSubscription({
      query: graphql(`
        subscription StreamState {
          streamInfo {
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
          }
        }
      `),
      variables: {}
    })

  return {
    useStreamState
  }
})
