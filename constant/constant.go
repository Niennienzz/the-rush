package constant

type RedisKey string

const (
	RedisKeyPlayerRecords RedisKey = "player:records"
)

func (x RedisKey) String() string {
	return string(x)
}

type RedisIndexKey string

const (
	RedisIndexKeyPlayerInterimPrefix            RedisIndexKey = "player:index:interim:%s"
	RedisIndexKeyPlayerByNamePrefix             RedisIndexKey = "player:index:by:name:%s"
	RedisIndexKeyPlayerByCreatedAt              RedisIndexKey = "player:index:by:created:at"
	RedisIndexKeyPlayerByLongestRush            RedisIndexKey = "player:index:by:longest:rush"
	RedisIndexKeyPlayerByTotalRushingTouchdowns RedisIndexKey = "player:index:by:total:rushing:touchdowns"
	RedisIndexKeyPlayerByTotalRushingYards      RedisIndexKey = "player:index:by:total:rushing:yards"
)

func (x RedisIndexKey) String() string {
	return string(x)
}

const MaxPageLimit int = 100
