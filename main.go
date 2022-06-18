package main

import (
	"fmt"
	"time"

	"github.com/bluele/gcache"
)

func main() {

	/**
	* New : 指定緩存大小
	*
	* 緩存策略:
	* Simple: 根據先存入的先淘汰
	* LUR(Least Recentily Used): 替換原則將最近最少使用的內容先替換掉
	* LFU(Least Frequently Used): 替換原則先淘汰一定時間內被訪問次數最少的資料
	* ARC(Adaptive Replacement Cache): ARC介於LRU和LFU之間
	 */

	/** 回調function
	* 可參考function automaticallyLoadValue
	*
	* LoaderExpireFunc: 過時回調function
	* EvictedFunc: 淘汰回調function
	* PurgeVisitorFunc: 清除所有key回調function
	* AddedFunc: 新增key回調function
	* SerializeFunc: 對value序列化回調function
	* DeseriailzeFunc: 對value反序列化回調function
	 */

	/* 計數事件
	*
	* HitCount: 命中次數
	* MissCount: 沒命中次數
	* LookupCount: 查找次數
	* HitRate: 命中率
	 */

	sampleCode()
}

/* 簡單範本 */
func sampleCode() {
	gc := gcache.New(20).LRU().Build()

	gc.Set("key", "ok")
	value, err := gc.Get("key")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Get: ", value)
	fmt.Println("Hit", gc.HitCount())
}

/* 設定過時時間 */
func timeoutCache() {
	gc := gcache.New(20).LRU().Expiration(time.Second * 4).Build()

	gc.Set("timeout_key", "timeout ok")

	// 未過時取值
	outValue, err := gc.Get("timeout_key")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Time Out Get:", outValue)

	time.Sleep(time.Second * 5) // Wait for value to expire

	// 過時取值
	outValue, err = gc.Get("timeout_key")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Time Out Get:", outValue)

	/** 輸出結果
	*
	* Get:  timeout ok
	* Key not found.
	* Time Out Get: <nil>
	 */
}

/* 設定過時時間 第2種寫法 */
func timeoutCache2() {
	gc := gcache.New(20).LRU().Build()

	gc.SetWithExpire("timeout_key", "timeout ok", time.Second*4)

	// 未過時取值
	outValue, err := gc.Get("timeout_key")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Time Out Get:", outValue)

	time.Sleep(time.Second * 5) // Wait for value to expire

	// 過時取值
	outValue, err = gc.Get("timeout_key")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Time Out Get:", outValue)
}

/* load回調function */
func loadData() {
	gc := gcache.New(20).LRU().LoaderFunc(func(i interface{}) (interface{}, error) {
		return "ok", nil
	}).Build()

	value, err := gc.Get("key")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Get:", value)

	/** 輸出結果
	*
	* Get: ok
	 */
}

/* 回調function使用 */
func automaticallyLoadValue() {
	var evictCounter, loaderCounter, purgeCounter int

	gc := gcache.New(20).
		LRU().
		LoaderExpireFunc(func(i interface{}) (interface{}, *time.Duration, error) {
			loaderCounter++
			expire := 1 * time.Second
			return "ok", &expire, nil
		}).
		EvictedFunc(func(key, value interface{}) {
			evictCounter++
			fmt.Println("evicted key:", key)
		}).
		PurgeVisitorFunc(func(key, value interface{}) {
			purgeCounter++
			fmt.Println("purged key:", key)
		}).
		Build()

	value, err := gc.Get("key")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Get:", value)

	time.Sleep(1 * time.Second)

	value, err = gc.Get("key")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Get:", value)
	gc.Purge()
	if loaderCounter != evictCounter+purgeCounter {
		fmt.Println("bad")
	}
}
