package go_redis

import (
	"log"
	"math/rand"
	"time"

	api_config "github.com/angenalZZZ/Go/go-program/api-config"
	go_type "github.com/angenalZZZ/Go/go-program/go-type"
	"github.com/gomodule/redigo/redis"
)

/**
数据库连接 redis : Client
*/
var Cli redis.Conn
var CliPoll *redis.Pool
var cliOpt redis.DialOption
var cliAddr string

// 初始化配置
func init() {
	// config
	cliAddr = api_config.Config.RedisCli.Addr
	// client
	cliOpt = redis.DialClientName("redis-cli")
	// db number
	cliOpt = redis.DialDatabase(api_config.Config.RedisCli.Db)
	cliOpt = redis.DialUseTLS(false)
	// password
	password := api_config.Config.RedisCli.Pwd
	if len(password) > 0 {
		cliOpt = redis.DialPassword(password)
	}
}
func initCli() {
	if CliPoll != nil {
		return
	}

	// managed Pool
	CliPoll = &redis.Pool{
		MaxActive: 5,
		MaxIdle:   5,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", cliAddr, cliOpt)
		},
	}

	// new client
	Cli = CliPoll.Get()

	// check
	if e := Cli.Err(); e != nil {
		log.Fatal(e) // 中断程序时输出
	}
	if e := Cli.Send("ping"); e != nil {
		log.Fatal(e) // 中断程序时输出
	}
}

// 数据库 Redis Cli close
func ShutdownCli() {
	if CliPoll != nil {
		//log.Println("缓存数据库 Redis Cli closing..")
		if Cli != nil {
			if e := Cli.Close(); e != nil {
				log.Fatal(e) // 中断程序时输出
			}
		}
		if e := CliPoll.Close(); e != nil {
			log.Fatal(e) // 中断程序时输出
		}
	}
}

// 测试
func TestCli() {
	//log.Println("缓存数据库 Redis Cli initCli..")
	initCli()
	log.Printf("缓存数据库 Redis Cli testing.. Addr: %s\n\n", cliAddr)

	// redis : new Cli
	c := CliPoll.Get()
	defer func() { _ = c.Close() }()
	// 时间戳：以秒计
	//timestamp := time.Now().Unix()
	// 时间戳：以毫秒计
	timestampNano := time.Now().UnixNano()
	rand.Seed(timestampNano)
	// 查找所有符合给定模式( pattern)的 key
	if keys, e := c.Do("keys", "timestamp*"); e != nil {
		log.Printf(" redis keys: Err\n %v\n", e)
	} else {
		ks := go_type.BytesToStrings(keys)
		log.Printf(" redis keys:   %v\n", ks)
		for _, k := range ks {
			t, _ := c.Do("type", k)
			log.Printf(" redis key:type %s => %v\n", k, t)
		}
	}

	/************ String（字符串）*************/
	TestCli_string(c)

	/************ Hash（哈希）*************/
	TestCli_hash(c)

	/************ List（列表）*************/
	TestCli_list(c)

	/************ List（列表）*************/
	TestCli_list(c)

	/************ Set（集合）*************/
	TestCli_set(c)

	/************ ZSet（有序集合）*************/
	TestCli_zset(c)

	// 计算: EVAL "return {KEYS[1],ARGV[1]}" 1 "key" "hello"
	//_, er3 := c.Do("EVAL", "return {KEYS[1],ARGV[1]}", 1, "key", "hello")
}
