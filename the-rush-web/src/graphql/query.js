import { gql } from '@apollo/client';

export const QUERY_PLAYERS = gql`
query PlayersWithArgs($args: PlayersArgs!) {
    players(args: $args) {
        total
        offset
        limit
        players {
            id
            createdAt
            name
            team
            position
            rushingAttempts
            rushingAttemptsPerGameAverage
            totalRushingYards
            rushingAverageYardsPerAttempt
            rushingYardsPerGame
            totalRushingTouchdowns
            longestRush {
                value
                isTouchdown
            }
            rushingFirstDowns
            rushingFirstDownsPercentage
            rushing20PlusYardsEach
            rushing40PlusYardsEach
            rushingFumbles
        }
    }
}`;