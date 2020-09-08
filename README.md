# The Rush

A simple GraphQL server for NFL players data.

## Installation & Running

- After starting the backend, go to `http://localhost:8080/playground` for the GraphQL playground.
- After starting the web UI, go to `http://localhost:3000` for the application page.

### Backend - Run with Docker

The backend can be started using Docker with the following commands `(port = 8080)`:

```
zsh> cd {$path_to_project}
zsh> docker-compose up
```

The backend should be **stopped properly** using:

```
zsh> docker-compose down
```

Note that the project does not use Docker volumes, **data does not persist on disk after** `docker-compose down`.

### Backend - Run Manually

If Docker is not available, the backend can also be started manually with the following steps.

- Requirements:
    - [Go v1.14+](https://golang.org/)
    - [Redis v6.0+](https://redis.io/)
    - [GNU Make v3.81+](https://www.gnu.org/software/make/)
- Since Go v1.14+ is required, `go mod` is used by default. Otherwise, `$GOPATH` should be appropriately configured.
- Make sure a local Redis server is running, default to `localhost:6379` with no password.
- Use environment variables `REDIS_URL` and `REDIS_PASSWORD` if the Redis configuration is different.
- Ingest data into the Redis server, and start the GraphQL server `(port = 8080)`:

```
zsh> cd {$path_to_project}
zsh> make ingest
zsh> make run
```

### Web UI - Run Manually

The web UI can be started manually with the following steps.

- Requirements:
    - [Yarn v1.22.4+](https://yarnpkg.com/)
- Change to `the-rush-web` sub-directory.
- Use `yarn` to start the web UI server `(port = 3000)`:

```
zsh> cd {$path_to_project}/the-rush-web
zsh> yarn install
zsh> yarn start
```

- Using `npm` should be similar, but it is not tested.

## GraphQL API

- The GraphQL API should be self-explanatory in the `graph/schema.graphqls` schema file.
- An example query with variables:

```graphql
query PlayersWithArgs($args: PlayersArgs!) {
  players(args: $args) {
    total
    offset
    limit
    players {
      id
      name
      team
      position
      totalRushingYards
      totalRushingTouchdowns
      longestRush {
        value
        isTouchdown
      }
    }
  }
}
```

```json
{
  "args": {
    "page": {
      "offset": 0,
      "limit": 10
    },
    "order": {
      "orderBy": "TOTAL_RUSHING_YARDS",
      "order": "DESC"
    }
  }
}
```

## Design

- The project has a relatively easy data model.
- A general SQL database (MySQL or PostgreSQL) can easily handle the required sorting/paging functionalities.
- We can even use [Dgraph](https://dgraph.io/) since the server exposes a GraphQL API.
- Simply passing the GraphQL query down to Dgraph can do the work.
- In this project, Redis is used. Here are some tradeoffs using Redis:
  - Redis is blazingly fast since it is an in-memory rich-data-type key-value store.
  - Since Redis operations does not require disk access, queries can return in a few milliseconds.
  - Redis is schemaless. We need to manage the data model and indexes ourselves.
  - Redis does not support full SQL, rich SQL operations cannot be achieved.
- Due to those tradeoffs above, in this project:
  - Each record has 1 primary index `id`, stored in Redis as a `Hash` with the format `id -> JSON(record)`.
  - Each record has 4 secondary indexes `createdAt`, `totalRushingTouchdowns`, `totalRushingYards`, and `longestRush`.
  - Each secondary index in Redis is a `SortedSet`.
  - Each record has a set of _poor-man's_ full-text search indexes at the whole word level.
  - For example, player with name `Joe Doe` has two entries: `by:name:joe -> id` and `by:name:doe -> id`.
  - Thus, user needs to input a whole word (i.e. `doe` instead of `do`) in order to search.
  - The design does not support partial search at the letter level.
  - Full-text search indexes are also `SortedSet`s so that they are easier to operated with other secondary indexes.
