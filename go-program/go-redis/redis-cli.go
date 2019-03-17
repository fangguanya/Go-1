package go_redis

import (
	"angenalZZZ/go-program/api-config"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

/**
数据库连接 redis : Client
*/
var Cli redis.Conn
var Pool *redis.Pool
var cliAddr string

// 初始化Cli
func InitCli() {
	if Pool != nil {
		return
	}

	// config
	api_config.Check("REDIS_ADDR")
	api_config.Check("REDIS_PWD")
	api_config.Check("REDIS_DB")
	cliAddr = os.Getenv("REDIS_ADDR")
	i, e := strconv.Atoi(os.Getenv("REDIS_DB"))
	if e != nil {
		i = 0
	}
	// client
	opt := redis.DialClientName("redis-cli")
	opt = redis.DialUseTLS(false)
	// password
	password := os.Getenv("REDIS_PWD")
	if len(password) > 0 {
		opt = redis.DialPassword(password)
	}
	// db number
	if i > 0 && i < 16 {
		opt = redis.DialDatabase(i)
	}
	// managed Pool
	Pool = &redis.Pool{
		MaxActive: 5,
		MaxIdle:   5,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", cliAddr, opt)
		},
	}

	// new client
	Cli = Pool.Get()

	//// check
	if e := Cli.Err(); e != nil {
		log.Fatal(e) // 中断程序时输出
	}
}

// 数据库 Redis Cli close
func ShutdownCli() {
	log.Println("缓存数据库 Redis Cli closing..")
	if Cli != nil {
		if e := Cli.Close(); e != nil {
			log.Fatal(e) // 中断程序时输出
		}
	}
	if Pool != nil {
		if e := Pool.Close(); e != nil {
			log.Fatal(e) // 中断程序时输出
		}
	}
}

// 测试
func TestCli() {
	InitCli()
	log.Printf("缓存数据库 Redis Cli testing.. Addr: %s\n\n", cliAddr)

	// redis : new Cli
	c := Pool.Get()
	defer func() { _ = c.Close() }()
	rand.Seed(time.Now().UnixNano())

	// 写入数据 Set
	key, val := fmt.Sprintf("timestamp%d%d", time.Now().Unix(), rand.Intn(1000)), "hello"
	if _, e := c.Do("SET", key, val, "EX", 60, "NX"); e != nil {
		log.Printf(" redis Set: Err\n [%s] %v\n", key, e)
	} else {
		log.Printf(" redis Set: Ok\n [%s] %s\n", key, val)
	}

	// 读取数据 Get
	valSaved, e := c.Do("GET", key)
	if valSaved == nil {
		log.Printf(" redis Get: Nil\n [%s] does not exist\n", key)
	} else if e != nil {
		log.Printf(" redis Get: Err\n [%s] %v\n", key, e)
	} else {
		log.Printf(" redis Get: Ok\n [%s] %s\n", key, valSaved)
	}

	// 删除数据 Del
	_, e = c.Do("DEL", key)
	if e != nil {
		log.Printf(" redis Del: Err\n [%s] %v\n", key, e)
	} else {
		log.Printf(" redis Del: Ok\n [%s]\n", key)
	}

	// 写入数据?当key不存在时+过期时间 SET key value EX 10 NX
	key, val = fmt.Sprintf("timestamp%d%d", time.Now().Unix(), rand.Intn(1000)), "values"
	_, e = c.Do("SET", key, val, "EX", 10, "NX")
	if e != nil {
		log.Printf(" redis SetNX: Err\n [%s] %v\n", key, e)
	} else {
		log.Printf(" redis SetNX: Ok\n [%s] %s\n", key, val)
	}

	// 读取数组+分页+排序 SORT list0 LIMIT 0 6 ASC, sort list0 desc alpha
	key = fmt.Sprintf("list%d%d", time.Now().Unix(), rand.Intn(1000))
	for i := range [6]int{1} {
		v := rand.Intn(1000) + i
		// LPUSH list0 1
		if _, e = c.Do("LPUSH", key, v); e != nil {
			log.Printf(" redis LPush: Err\n [%s] %v\n", key, e)
		} else {
			log.Printf(" redis LPush: Ok\n [%s] %d\n", key, v)
		}
	}
	arr, err := c.Do("SORT", key, "LIMIT", 0, 6, "ASC")
	if err != nil {
		log.Printf(" redis Sort: Err\n [%s] %v\n", key, e)
	} else {
		log.Printf(" redis Sort: Ok\n [%s] %v\n", key, arr)
	}

	// 读取有序集合中指定分数区间的成员列表 ZRANGEBYSCORE zset0 -inf +inf WITHSCORES LIMIT 0 6 [WITHSCORES:输出分数]
	key = fmt.Sprintf("zset%d%d", time.Now().Unix(), rand.Intn(1000))
	for i := range [6]int{1} {
		score, member := rand.Float64(), fmt.Sprintf("member%d", rand.Intn(100)+i)
		// ZADD zset0 1 member1
		if _, e = c.Do("ZADD", key, score, member); e != nil {
			log.Printf(" redis ZAdd: Err\n [%s] %v\n", key, e)
		} else {
			log.Printf(" redis ZAdd: Ok\n [%s] %s=%f\n", key, member, score)
		}
	}
	set1, er1 := c.Do("ZRANGEBYSCORE", key, "-inf", "+inf", "WITHSCORES", "LIMIT", 0, 6)
	if er1 != nil {
		log.Printf(" redis ZRANGEBYSCORE: Err\n [%s] %v\n", key, er1)
	} else {
		log.Printf(" redis ZRANGEBYSCORE: Ok\n [%s] %v\n", key, set1)
	}

	// 计算: 给定有序集的交集,并将该交集(结果集)储存起来 http://www.runoob.com/redis/sorted-sets-zinterstore.html
	// ZINTERSTORE out 2 zset01 zset02 WEIGHTS 2 3 AGGREGATE SUM
	e = c.Send("ZADD", "zset01", rand.Intn(100), "A", "NX")
	e = c.Send("ZADD", "zset01", rand.Intn(100), "B", "NX")
	e = c.Send("ZADD", "zset01", rand.Intn(100), "C", "NX")
	e = c.Send("ZADD", "zset02", rand.Intn(100), "A", "NX")
	e = c.Send("ZADD", "zset02", rand.Intn(100), "B", "NX")
	e = c.Send("ZADD", "zset02", rand.Intn(100), "C", "NX")
	e = c.Flush()
	_, e = c.Receive()
	_, e = c.Receive()
	_, e = c.Receive()
	_, e = c.Receive()
	_, e = c.Receive()
	_, e = c.Receive()
	// 交集的目标key
	key = "zset0102"
	set2, er2 := c.Do("ZINTERSTORE", key, 2, "zset01", "zset02", "WEIGHTS", 0, 100, "AGGREGATE", "SUM")
	if er1 != nil {
		log.Printf(" redis ZINTERSTORE: Err\n [%s] %v\n", key, er2)
	} else {
		log.Printf(" redis ZINTERSTORE: Ok\n [%s] %v\n", key, set2)
	}

	// 计算: EVAL "return {KEYS[1],ARGV[1]}" 1 "key" "hello"
	//_, er3 := c.Do("EVAL", "return {KEYS[1],ARGV[1]}", 1, "key", "hello")
}