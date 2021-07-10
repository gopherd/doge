package router

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/gopherd/doge/proto"
	"github.com/gopherd/doge/service/discovery"
)

type Target struct {
	Address  string
	Priority int
}

type Targets []Target

func (ts Targets) Len() int           { return len(ts) }
func (ts Targets) Less(i, j int) bool { return ts[i].Priority < ts[j].Priority }
func (ts Targets) Swap(i, j int)      { ts[i], ts[j] = ts[j], ts[i] }

type Router struct {
	Module  string     `json:"module"`
	MinType proto.Type `json:"min_type"`
	MaxType proto.Type `json:"max_type"`
	Targets Targets    `json:"targets"`
}

func (r Router) Key() string {
	if r.MinType == 0 && r.MaxType == 0 {
		return r.Module
	}
	return fmt.Sprintf("%s/%d/%d", r.Module, r.MinType, r.MaxType)
}

func (r Router) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *Router) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r Router) Less(r2 Router) bool {
	return r.less(r2.Module, r2.MinType)
}

func (r Router) less(mod string, typ proto.Type) bool {
	if r.Module == mod {
		return r.MinType < typ
	}
	return r.Module < mod
}

func (r Router) has(mod string, typ proto.Type) bool {
	if r.Module != mod {
		return false
	}
	if r.MinType == 0 && r.MaxType == 0 {
		return true
	}
	return typ >= r.MinType && typ <= r.MaxType
}

func (r Router) getTarget() string {
	return r.Targets[0].Address
}

const (
	prefix   = "message/router/"
	prefix0  = prefix + "0"
	cacheTTL = time.Second * 60
)

func Register(ctx context.Context, discovery discovery.Discovery, uid int64, router Router, ttl time.Duration) error {
	content, err := router.Marshal()
	if err != nil {
		return err
	}
	name := prefix + strconv.FormatInt(uid, 10)
	return discovery.Register(ctx, name, router.Key(), string(content), false, ttl)
}

func Unregister(ctx context.Context, discovery discovery.Discovery, uid int64, router Router) error {
	name := prefix + strconv.FormatInt(uid, 10)
	return discovery.Unregister(ctx, name, router.Key())
}

type routers struct {
	expires time.Time
	routers []Router
}

func (rs *routers) Len() int           { return len(rs.routers) }
func (rs *routers) Less(i, j int) bool { return rs.routers[i].Less(rs.routers[j]) }
func (rs *routers) Swap(i, j int)      { rs.routers[i], rs.routers[j] = rs.routers[j], rs.routers[i] }

func (rs *routers) autofix() {
	for i := len(rs.routers) - 1; i >= 0; i-- {
		if len(rs.routers[i].Targets) == 0 {
			rs.routers = append(rs.routers[:i], rs.routers[i+1:]...)
			continue
		}
		sort.Sort(rs.routers[i].Targets)
	}
	sort.Sort(rs)
}

func (rs *routers) lookup(mod string, typ proto.Type) string {
	i, j := 0, len(rs.routers)
	for i < j {
		h := int(uint(i+j) >> 1)
		if rs.routers[h].less(mod, typ) {
			i = h + 1
		} else {
			j = h
		}
	}
	if i == len(rs.routers) {
		if len(rs.routers) > 0 {
			if rs.routers[0].has(mod, typ) {
				return rs.routers[0].getTarget()
			}
		}
		return ""
	}
	if rs.routers[i].has(mod, typ) {
		return rs.routers[i].getTarget()
	}
	if rs.routers[0].has(mod, typ) {
		return rs.routers[0].getTarget()
	}
	return ""
}

type Cache struct {
	discovery discovery.Discovery

	mu      sync.RWMutex
	shared  *routers
	routers map[int64]*routers
}

func NewCache(discovery discovery.Discovery) *Cache {
	return &Cache{
		discovery: discovery,
		routers:   make(map[int64]*routers),
	}
}

func (cache *Cache) Init() error {
	if routers, err := cache.load(prefix0); err != nil {
		return err
	} else {
		cache.mu.Lock()
		defer cache.mu.Unlock()
		cache.shared = routers
	}
	return nil
}

func (cache *Cache) Add(uid int64, rs []Router) {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	if uid == 0 {
		cache.shared = &routers{
			expires: time.Now().Add(cacheTTL),
			routers: rs,
		}
	} else {
		cache.routers[uid] = &routers{
			expires: time.Now().Add(cacheTTL),
			routers: rs,
		}
	}
}

func (cache *Cache) Remove(uid int64) {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	delete(cache.routers, uid)
}

func (cache *Cache) Lookup(uid int64, mod string, typ proto.Type) (string, error) {
	if uid == 0 {
		shared, err := cache.getOrLoadShared()
		if err != nil {
			return "", err
		}
		return shared.lookup(mod, typ), nil
	}
	routers, err := cache.getOrLoad(uid)
	if err != nil {
		return "", err
	}
	return routers.lookup(mod, typ), nil
}

func (cache *Cache) load(name string) (*routers, error) {
	values, err := cache.discovery.ResolveAll(context.Background(), prefix0)
	if err != nil {
		return nil, err
	}
	var rs = new(routers)
	for _, v := range values {
		var r Router
		if err := r.Unmarshal([]byte(v)); err != nil {
			return nil, err
		}
		rs.routers = append(rs.routers, r)
	}
	rs.expires = time.Now().Add(cacheTTL)
	rs.autofix()
	return rs, nil
}

func (cache *Cache) getOrLoadShared() (*routers, error) {
	cache.mu.RLock()
	shared := cache.shared
	cache.mu.RUnlock()
	if shared == nil || shared.expires.Before(time.Now()) {
		return shared, nil
	}
	routers, err := cache.load(prefix0)
	if err != nil {
		return nil, err
	}
	cache.mu.Lock()
	defer cache.mu.Unlock()
	cache.shared = routers
	return routers, nil
}

func (cache *Cache) getOrLoad(uid int64) (*routers, error) {
	cache.mu.RLock()
	routers, ok := cache.routers[uid]
	cache.mu.RUnlock()
	if ok && routers.expires.Before(time.Now()) {
		return routers, nil
	}
	routers, err := cache.load(prefix + strconv.FormatInt(uid, 10))
	if err != nil {
		return nil, err
	}
	cache.mu.Lock()
	defer cache.mu.Unlock()
	cache.routers[uid] = routers
	return routers, nil
}
