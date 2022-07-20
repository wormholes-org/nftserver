package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/garyburd/redigo/redis"
	"golang.org/x/crypto/sha3"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

const (
	queryRedisMaxactive   = 0
	queryRedisMaxidle     = 20
	queryCatchMax         = 50
	queryRedisIdleTimeout = time.Second * 100
	MainRedisRetry        = time.Second * 10
	ScanDirtyQuerysTime   = time.Second * 10
)

type Remote struct {
	conn    net.Conn
	chDatas chan []byte
	chping  chan []byte
}

type queryRedisCatch struct {
	pool     *redis.Pool
	Mux      sync.Mutex
	Dirty    map[string]bool
	CathFlag bool
	connMux  sync.Mutex
	listener net.Listener
	connPool map[string]Remote
	close    chan struct{}
	wFlag    bool
}

var qr *queryRedisCatch

func (c queryRedisCatch) Close() {
	c.pool.Close()
}

func (c queryRedisCatch) QueryHash(data string) string {
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write([]byte(data))
	hash := hasher.Sum(nil)
	return hexutil.Encode(hash)
}

func (c queryRedisCatch) QueryKey(exchanger, queryname, queryhash string) string {
	return c.QueryHash(exchanger + queryname + queryhash)
}

func (c *queryRedisCatch) SetRedisData(key string, datas []byte) error {
	conn := c.pool.Get()
	defer conn.Close()
	_, err := conn.Do("set", key, datas)
	return err
}

func (c *queryRedisCatch) HsetRedisData(hkey, key string, datas []byte) error {
	conn := c.pool.Get()
	defer conn.Close()
	_, err := conn.Do("hset", hkey, key, datas)
	return err
}

func (c *queryRedisCatch) SetExpireRedisData(key string, datas []byte, expire int) error {
	conn := c.pool.Get()
	defer conn.Close()
	_, err := conn.Do("set", ExchangeOwer+key, datas, "EX", expire)
	return err
}

func (c *queryRedisCatch) GetRedisData(key string) ([]byte, error) {
	conn := c.pool.Get()
	defer conn.Close()
	replay, err := conn.Do("get", key)
	if err != nil {
		log.Println("GetRedisData() get data err=", err, " key = ", key)
		return nil, err
	}
	if replay != nil {
		datas, ok := replay.([]byte)
		if ok {
			return datas, nil
		}
		log.Println("GetRedisData() convert data err=", err, " key = ", key)
	}
	return nil, nil
}

func (c *queryRedisCatch) HgetRedisData(hkey, key string) ([]byte, error) {
	conn := c.pool.Get()
	defer conn.Close()
	replay, err := conn.Do("hget", hkey, key)
	if err != nil {
		log.Println("GetRedisData() get data err=", err, " key = ", key)
		return nil, err
	}
	if replay != nil {
		datas, ok := replay.([]byte)
		if ok {
			return datas, nil
		}
		log.Println("GetRedisData() convert data err=", err, " key = ", key)
	}
	return nil, nil
}

func (c *queryRedisCatch) DelRedisData(key string) error {
	conn := c.pool.Get()
	defer conn.Close()
	_, err := conn.Do("del", key)
	if err != nil {
		log.Println("DelRedisData() err=", err, " key = ", key)
		return err
	}
	return err
}

func (c *queryRedisCatch) HdelRedisData(hkey, key string) error {
	conn := c.pool.Get()
	defer conn.Close()
	_, err := conn.Do("hdel", hkey, key)
	if err != nil {
		log.Println("HdelRedisData() err=", err, " hkey = ", hkey, " key = ", key)
		return err
	}
	return err
}

func (c *queryRedisCatch) CatchQueryData(queryname, querysql string, queryData interface{}) error {
	if !c.wFlag {
		return errors.New("read only redis.")
	}
	c.Mux.Lock()
	defer c.Mux.Unlock()
	fmt.Println("CatchQueryData() key = ", ExchangeOwer+"."+queryname+"."+querysql)
	ckey := c.QueryHash(ExchangeOwer + queryname + "count")
	hkey := c.QueryHash(ExchangeOwer + queryname)
	skey := c.QueryHash(querysql)
	catchCount := make(map[string]int)
	datas, err := c.GetRedisData(ckey)
	if err == nil && datas != nil {
		err = json.Unmarshal(datas, &catchCount)
		if err == nil {
			fmt.Println("CatchQueryData() Catch name=", ExchangeOwer+queryname, "Catch data count=", len(catchCount))
			for s, i := range catchCount {
				fmt.Println("CatchQueryData() key, count=", s, ":", i)
			}
			_, ok := catchCount[skey]
			if !ok {
				if len(catchCount) > queryCatchMax {
					k := 0
					var key string
					for s, i := range catchCount {
						if k == 0 {
							k = i
							key = s
						} else if i < k {
							k = i
							key = s
						}
					}
					delete(catchCount, key)
					c.HdelRedisData(hkey, key)
				}
			}
		}
	}
	queryDatas, _ := json.Marshal(queryData)
	err = c.HsetRedisData(hkey, skey, queryDatas)
	if err != nil {
		log.Println("CatchQueryData() CatchQueryData  err=", err)
	}
	catchCount[skey] = catchCount[skey] + 1
	datas, err = json.Marshal(catchCount)
	err = c.SetRedisData(ckey, datas)
	return err
}

func (c *queryRedisCatch) GetCatchData(queryname, querysql string, queryData interface{}) error {
	c.Mux.Lock()
	defer c.Mux.Unlock()
	fmt.Println("GetCatchData() key = ", ExchangeOwer+"."+queryname+"."+querysql)
	hkey := c.QueryHash(ExchangeOwer + queryname)
	skey := c.QueryHash(querysql)
	datas, err := c.HgetRedisData(hkey, skey)
	if err != nil {
		log.Println("GetCatchData() get cathed data err=", err, " exchanger + queryname + queryhash = ", querysql+queryname+querysql)
		return err
	}
	if datas != nil {
		err = json.Unmarshal(datas, queryData)
		if err != nil {
			log.Println("GetCatchData() miss catch  err=", err)
			return err
		}
		ckey := c.QueryHash(ExchangeOwer + queryname + "count")
		log.Println("GetCatchData() ckey=", ExchangeOwer+":"+queryname+"count")
		catchCount := make(map[string]int)
		cdatas, err := c.GetRedisData(ckey)
		if err == nil && datas != nil {
			err = json.Unmarshal(cdatas, &catchCount)
			if err == nil {
				fmt.Println("GetCatchData() Catch name=", ExchangeOwer+":"+queryname, "   Catch data count=", len(catchCount))
				for s, i := range catchCount {
					fmt.Println("GetCatchData() key, count=", s, ":", i)
				}
				catchCount[skey] = catchCount[skey] + 1
				fmt.Println("GetCatchData() queryname=", queryname, " count=", catchCount[skey])
				datas, err = json.Marshal(catchCount)
				err = c.SetRedisData(ckey, datas)
			}
		}
		return nil
	}
	return errors.New("no catch data.")
}
func (c *queryRedisCatch) ClearSingleQueryCatch(queryname string) {
	ckey := c.QueryHash(ExchangeOwer + queryname + "count")
	hkey := c.QueryHash(ExchangeOwer + queryname)
	catchCount := make(map[string]int)
	datas, err := c.GetRedisData(ckey)
	if err == nil && datas != nil {
		err = json.Unmarshal(datas, &catchCount)
		if err == nil {
			fmt.Println("ClearSingleQueryCatch() Catch name=", ExchangeOwer+queryname, "Catch data count=", len(catchCount))
			for key, _ := range catchCount {
				c.HdelRedisData(hkey, key)
			}
			c.DelRedisData(ckey)
		}
	}
}

func (c *queryRedisCatch) SaveDirtyQuerys(querynames []string) error {
	dirtyKey := c.QueryHash(ExchangeOwer + "dirty")
	dirtyValue := make(map[string]uint64)
	datas, err := c.GetRedisData(dirtyKey)
	if err == nil && datas != nil {
		err = json.Unmarshal(datas, &dirtyValue)
		if err != nil {
			log.Println("SaveDirtyQuerys() Unmarshal err=", err)
		}
	}
	for _, queryname := range querynames {
		dirtyValue[queryname] = dirtyValue[queryname] + 1
	}
	datas, _ = json.Marshal(dirtyValue)
	err = c.SetRedisData(dirtyKey, datas)
	if err != nil {
		log.Println("SaveDirtyQuerys() SetRedisData err=", err)
	}
	return err
}

func (c *queryRedisCatch) GetDirtyQuerys() (map[string]uint64, error) {
	dirtyKey := c.QueryHash(ExchangeOwer + "dirty")
	dirtyValue := make(map[string]uint64)
	datas, err := c.GetRedisData(dirtyKey)
	if err == nil && datas != nil {
		err = json.Unmarshal(datas, &dirtyValue)
		return dirtyValue, nil
	}
	return nil, errors.New("no catch dirty data.")
}

func (c *queryRedisCatch) SetDirtyFlag(querynames []string) {
	/*c.Mux.Lock()
	defer c.Mux.Unlock()
	fmt.Println("SetDirtyFlag() querys = ", querynames)
	for _, queryname := range querynames {
		key := c.QueryHash(ExchangeOwer + queryname)
		err := c.DelRedisData(key)
		if err != nil {
			log.Println("SetDirtyFlag() DelRedisData key=", key, " err=", err)
			continue
		}
	}*/
	/*if !LimitWritesDatabase {
		c.ClearClientCatch(querynames)
	}*/
	c.Mux.Lock()
	defer c.Mux.Unlock()
	fmt.Println("SetDirtyFlag() querys = ", querynames)
	for _, queryname := range querynames {
		ckey := c.QueryHash(ExchangeOwer + queryname + "count")
		hkey := c.QueryHash(ExchangeOwer + queryname)
		catchCount := make(map[string]int)
		datas, err := c.GetRedisData(ckey)
		if err == nil && datas != nil {
			err = json.Unmarshal(datas, &catchCount)
			if err == nil {
				fmt.Println("SetDirtyFlag() Catch name=", ExchangeOwer+queryname, "Catch data count=", len(catchCount))
				for key, _ := range catchCount {
					c.HdelRedisData(hkey, key)
				}
				c.DelRedisData(ckey)
			}
		}
	}
	if !LimitWritesDatabase {
		err := c.SaveDirtyQuerys(querynames)
		if err != nil {
			log.Println("SetDirtyFlag() SaveDirtyQuerys err=", err)
		}
	}
}

func (c *queryRedisCatch) scanDirtyQuery(mc *queryRedisCatch) {
	if !c.wFlag {
		return
	}
	mdirtyQuerys, err := mc.GetDirtyQuerys()
	if err == nil {
		dirtyQuerys, err := c.GetDirtyQuerys()
		if err == nil {
			flush := false
			for s, u := range mdirtyQuerys {
				if dirtyQuerys[s] != u {
					dirtyQuerys[s] = u
					c.ClearSingleQueryCatch(s)
					flush = true
					fmt.Println("scanDirtyQuery() clear query=", s, " :count=", u)
				}
			}
			if flush {
				dirtyKey := c.QueryHash(ExchangeOwer + "dirty")
				datas, _ := json.Marshal(dirtyQuerys)
				err = c.SetRedisData(dirtyKey, datas)
			}
		} else {
			dirtyQuerys = make(map[string]uint64)
			for s, u := range mdirtyQuerys {
				dirtyQuerys[s] = u
				c.ClearSingleQueryCatch(s)
			}
			dirtyKey := c.QueryHash(ExchangeOwer + "dirty")
			datas, _ := json.Marshal(dirtyQuerys)
			err = c.SetRedisData(dirtyKey, datas)
		}
	}
}

func (c *queryRedisCatch) ScanDirtyQuerys(mc *queryRedisCatch) {
	if !c.wFlag {
		log.Println("ScanDirtyQuerys() read only redis.")
		return
	}
	ticker := time.NewTicker(ScanDirtyQuerysTime)
	for {
		select {
		case <-ticker.C:
			c.scanDirtyQuery(mc)
		}
	}
}

func (c *queryRedisCatch) SetFlag(querynames string) {
	c.Mux.Lock()
	defer c.Mux.Unlock()
	fmt.Println("SetFlag() querys = ", querynames)
	key := c.QueryHash(ExchangeOwer + querynames)
	err := c.DelRedisData(key)
	if err != nil {
		log.Println("SetDirtyFlag() DelRedisData key=", key, " err=", err)
	}
}

func (c *queryRedisCatch) ClearClientCatch(clear []string) {
	clearData, _ := json.Marshal(clear)
	for _, remote := range c.connPool {
		go func() { remote.chDatas <- clearData }()
	}
}

func (c *queryRedisCatch) StartClientCatch() {
	go func() {
		for {
			conn, err := net.Dial("tcp", MainRedisCatchSvr)
			if err != nil {
				log.Println("ClientCatchServe() Dial err=", err)
				time.Sleep(time.Second * 2)
				continue
			}
			for {
				fmt.Println("ClientCatchServe() remote=", conn.RemoteAddr().String())
				fmt.Println("ClientCatchServe() locate=", conn.LocalAddr().String())
				conn.SetReadDeadline(time.Now().Add(5 * time.Second))
				buf := make([]byte, 2000)
				n, err := conn.Read(buf)
				if err != nil {
					if err.Error() == "EOF" {
						log.Println("ClientCatchServe() link redail  err=", err)
						break
					}
					if strings.Contains(err.Error(), "i/o timeout") {
						conn.Write([]byte("ping"))
						continue
					}
				}
				if n == 0 {
					log.Println("ClientCatchServe() link redail  err=", err)
					break
				}
				var pong string
				pong = string(buf[:len("pong")])
				if pong == "pong" {
					continue
				}
				clearData := string(buf[:n])
				clearData = clearData[strings.Index(clearData, "[") : strings.Index(clearData, "]")+1]
				var dirtyData []string
				err = json.Unmarshal([]byte(clearData), &dirtyData)
				if err != nil {
					for _, queryname := range dirtyData {
						key := c.QueryHash(ExchangeOwer + queryname)
						err := c.DelRedisData(key)
						if err != nil {
							log.Println("SetDirtyFlag() DelRedisData key=", key, " err=", err)
							continue
						}
					}
				}
				fmt.Println("ClientCatchServe() received clear catch string=", clearData)
			}
			conn.Close()
		}
	}()
}

func (c *queryRedisCatch) ServerCatch(remote string) {
	close := make(chan struct{})
	rm := c.connPool[remote]
	go func() {
		for {
			rm.conn.SetReadDeadline(time.Now().Add(20 * time.Second))
			buf := make([]byte, 100)
			n, err := rm.conn.Read(buf)
			if err != nil {
				if err.Error() == "EOF" {
					log.Println("ServerCatch() close  err=", err)
					c.connMux.Lock()
					delete(c.connPool, remote)
					c.connMux.Unlock()
					rm.conn.Close()
					close <- struct{}{}
					return
				}
			}
			if n == 0 {
				log.Println("ServerCatch() close err=", err)
				c.connMux.Lock()
				delete(c.connPool, remote)
				c.connMux.Unlock()
				rm.conn.Close()
				close <- struct{}{}
				return
			}
			var ping string
			ping = string(buf[:len("ping")])
			if ping == "ping" {
				rm.chping <- []byte("pong")
				continue
			}
		}
	}()
close:
	for {
		select {
		case data := <-rm.chDatas:
			rm.conn.Write(data)
		case p := <-rm.chping:
			fmt.Println("ServerCatch() ", rm.conn.RemoteAddr().String(), " receive ping.")
			rm.conn.Write(p)
		case <-c.close:
			fmt.Println("ServerCatch() link error close.")
			break close
		case <-close:
			fmt.Println("ServerCatch() ", rm.conn.RemoteAddr().String(), " close.")
			break close
		}
	}
}

func (c *queryRedisCatch) ServeListener() {
	for {
		for {
			conn, err := c.listener.Accept()
			if err != nil {
				log.Println("ServeListener() server err=", err)
				c.close <- struct{}{}
				time.Sleep(time.Second * 2)
				c.connPool = make(map[string]Remote)
				break
			}
			fmt.Println("ServeListener() connect remote addr=", conn.RemoteAddr().String())
			c.connMux.Lock()
			c.connPool[conn.RemoteAddr().String()] = Remote{conn: conn, chDatas: make(chan []byte, 100), chping: make(chan []byte)}
			c.connMux.Unlock()
			go c.ServerCatch(conn.RemoteAddr().String())
		}
		for {
			listener, err := net.Listen("tcp", MainRedisCatchSvr)
			if err != nil {
				log.Println("InitCatchSvr() listen err=", err)
				time.Sleep(time.Second * 10)
				continue
			}
			c.listener = listener
		}
	}
}

func (c *queryRedisCatch) StartCatchServe() error {
	listener, err := net.Listen("tcp", MainRedisCatchSvr)
	if err != nil {
		log.Println("InitCatchSvr() listen err=", err)
		return err
	}
	c.listener = listener
	go c.ServeListener()
	return nil
}

func NewQueryCatch(redisIP, passwd string) error {
	qr = new(queryRedisCatch)
	qr.pool = &redis.Pool{
		MaxActive:   queryRedisMaxactive,
		MaxIdle:     queryRedisMaxidle,
		IdleTimeout: queryRedisIdleTimeout,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", redisIP, redis.DialPassword(passwd))
		},
	}
	qr.CathFlag = true
	qr.connPool = make(map[string]Remote)
	qr.close = make(chan struct{})
	err := qr.SetRedisData("testtest", []byte("testtest"))
	if err != nil {
		qr.wFlag = false
	} else {
		qr.wFlag = true
		qr.DelRedisData("testtest")
	}

	/*if !LimitWritesDatabase {
		err := qr.StartCatchServe()
		return err
	} else {
		qr.StartClientCatch()
	}*/
	return nil
}

func GetRedisCatch() *queryRedisCatch {
	if qr == nil {
		log.Fatalln("qr is nil.")
	}
	return qr
}
