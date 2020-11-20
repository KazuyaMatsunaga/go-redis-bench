package rbbench_test

import (
	"testing"

	"github.com/boltdb/bolt"
	redigo "github.com/garyburd/redigo/redis"
	redis "github.com/go-redis/redis"
)

var redisOpts = &redis.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
}

func BenchmarkRedisSet(b *testing.B) {
	client := redis.NewClient(redisOpts)
	defer client.Close()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		client.Set("key"+string(i), "value", 0).Err()
	}
}

func BenchmarkRedisGet(b *testing.B) {
	client := redis.NewClient(redisOpts)
	defer client.Close()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		client.Get("key" + string(i)).Val()
	}
}

func BenchmarkRedigoSet(b *testing.B) {
	conn, _ := redigo.Dial("tcp", "localhost:6379")
	defer conn.Close()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		conn.Do("SET", "key"+string(i), "value")
	}
}

func BenchmarkRedigoGet(b *testing.B) {
	conn, _ := redigo.Dial("tcp", "localhost:6379")
	defer conn.Close()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		redigo.String(conn.Do("GET", "key"+string(i)))
	}
}

func BenchmarkBoltSet(b *testing.B) {
	db, _ := bolt.Open("bolt.db", 0600, nil)
	defer db.Close()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		db.Update(func(tx *bolt.Tx) error {
			b, _ := tx.CreateBucketIfNotExists([]byte("bucket"))
			b.Put([]byte("key"+string(i)), []byte("value"))
			return nil
		})
	}
}

func BenchmarkBoltGet(b *testing.B) {
	db, _ := bolt.Open("bolt.db", 0600, nil)
	defer db.Close()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		db.View(func(tx *bolt.Tx) error {
			_ = string(tx.Bucket([]byte("bucket")).Get([]byte("key" + string(i))))
			return nil
		})
	}
}