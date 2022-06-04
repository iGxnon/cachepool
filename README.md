# cachepool

+ Require: golang 1.18 or above

### Usage:

1. cachepool

```go
func ExampleCachePool() {
	pool := NewDefault(nil)
	// Set
	pool.Cache.Set("foo", "bar", time.Minute*40)
	err := pool.Cache.Add("foo2", "bar2", cache.DefaultExpiration)
	if err != nil {
		// foo2 contains before
	}
	err = pool.Cache.Replace("foo", "barbar", cache.NoExpiration)
	if err != nil {
		// foo does not contain before
	}
	// Get
	_, _ = pool.Cache.Get("foo")
	_, _, _ = pool.Cache.GetWithExpiration("foo2")

	// increment and decrement
	pool.Cache.Set("foo3", 114514, cache.NoExpiration)

	_ = pool.Cache.Increment("foo3", 1919810) // then foo3 equals 2034324
	_ = pool.Cache.Decrement("foo3", 1919810)

	// use message queue, sync some cache
	conn, _ := amqp.Dial("amqp://guest:guest@localhost:5672/")
	ch, _ := conn.Channel()
	pool.UseMQ(ch, "cache1")

	// publish an importance message into cache (
	_ = cache.Publish(ch, "下北沢一番臭の伝説", struct {
		Age    int
		Prefix string
		Movie  string
	}{
		24,
		"野獣せんべい",
		"真夏の夜の银夢",
	}, time.Minute*5)

	time.Sleep(time.Second)

	// stop using message queue
	pool.StopMQ()
}

```

2. cache helper

```go
/*
------------------------
yee	| bar	 | foo
------------------------
Hello   |  1     | null
------------------------
Hi	|  2     | NOW()
------------------------
*/

type FooBar struct {
	Bar int64
	Yee string
	Foo sql.NullTime
}

func ExampleHelper() {
	dsn := "root:12345678@tcp(127.0.0.1:3306)/awesome?charset=utf8mb4&parseTime=True"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return
	}
	pool := NewDefault(db)
	got, err := helper.QueryRow[FooBar](pool, "SELECT * FROM t LIMIT 1 OFFSET 1")
	if err != nil {
		return
	}
	fmt.Printf("%#v\n", got)

	gots, err := helper.Query[map[string]any](pool, "SELECT * FROM t LIMIT 5")
	if err != nil {
		return
	}
	fmt.Println(gots)

	gotOnes, err := helper.Query[int32](pool, "SELECT bar FROM t LIMIT 5")
	if err != nil {
		return
	}
	fmt.Println(gotOnes)

	gotOnesNullable, err := helper.Query[sql.NullTime](pool, "SELECT foo FROM t LIMIT 5")
	if err != nil {
		return
	}
	fmt.Println(gotOnesNullable)
}

```