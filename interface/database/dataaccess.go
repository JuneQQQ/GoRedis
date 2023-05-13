package database

type DataAccess interface {
	GetEntity(key string) (*DataEntity, bool)
	SetEntity(key string, entity *DataEntity) int
	SetIfExists(key string, entity *DataEntity) int
	SetIfAbsent(key string, entity *DataEntity) int
	Remove(key string)
	RemoveBulk(keys ...string) (deleted int)
}
