package ratelimit

import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

const (
	NotExpireTTL = 24 * time.Hour * 365 *5 //不过期的时间
)

type cacheData struct {
	key string
	value interface{}
	expireAt time.Time
}
func newCacheData(key string, value interface{}, expireTime time.Time) *cacheData {
	return &cacheData{key: key, value: value, expireAt:expireTime}
}


type Cache struct {
	maxsize int  //限制最大key的数量,0 代表不限制
	elist *list.List
	edata map[string]*list.Element
	lock sync.Mutex
}

type CacheOption func(g *Cache)
type CacheOptions []CacheOption

func(opts CacheOptions) apply(c *Cache){
	for _,fn:=range opts{
		fn(c)
	}
}

//设置了最大长度
func WithMaxSize(size int) CacheOption{
	return func(c *Cache) {
		if size > 0 {
			c.maxsize=size
		}
	}
}


func NewCache(opts ...CacheOption) *Cache {
	cache:= &Cache{elist:list.New(),edata: make(map[string]*list.Element),
		maxsize: 0}
	CacheOptions(opts).apply(cache)
	cache.clear()
	return cache

}
//获取缓存
func(this *Cache) Get(key string) interface{}{
	this.lock.Lock()
	defer this.lock.Unlock()
	if v,ok:=this.edata[key];ok{
		if time.Now().After(v.Value.(*cacheData).expireAt) {
			this.removeItem(v)
			return nil
		}
		this.elist.MoveToFront(v)
		return v.Value.(*cacheData).value

	}
	return nil
}
func(this *Cache) Set(key string ,newv interface{}, expireTTL time.Duration){
	this.lock.Lock()
	defer this.lock.Unlock()
	var setExpire time.Time
	if expireTTL == 0 {
		setExpire = time.Now().Add(NotExpireTTL)  // 设置不过期时间
	} else {
		setExpire = time.Now().Add(expireTTL)
	}
	newCache:= newCacheData(key,newv, setExpire)
	if v,ok:=this.edata[key];ok{
		 v.Value=newCache
		 this.elist.MoveToFront(v)
	}else{
		this.edata[key]=this.elist.PushFront(newCache)
		// 判断长度是否溢出 ,如果是：末尾淘汰一个缓存
		if this.maxsize > 0 && len(this.edata) > this.maxsize{
			this.removeOldest()
		}
	}
}
func(this *Cache) Print(){
	ele:=this.elist.Front()
	if ele==nil{
		return
	}
	for{
		fmt.Println(this.Get(ele.Value.(*cacheData).key))
		ele=ele.Next()
		if ele==nil{
			break
		}
	}
}
// 删除最后一个元素
func(this *Cache) removeOldest(){

	back:=this.elist.Back()
	if back==nil{
		return
	}
	this.removeItem(back)

}
func(this *Cache) removeItem(ele *list.Element){
	key:=ele.Value.(*cacheData).key
	delete(this.edata,key) //删除map里面的key
	this.elist.Remove(ele)
}

func (this *Cache) Len() int  {
	return len(this.edata)
}

func(this *Cache) removeExpire()  {
	this.lock.Lock()
	defer this.lock.Unlock()

	for _, v := range this.edata {
		if time.Now().After(v.Value.(*cacheData).expireAt) {
			this.removeItem(v)
		}
	}
}

func(this *Cache)clear() {
	go func() {

		for {
			this.removeExpire()
			time.Sleep(time.Second * 1)
		}

	}()
}


