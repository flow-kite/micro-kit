package mconsul

import (
	"github.com/hashicorp/consul/api"
	"github.com/pkg/errors"

	"github.com/o-kit/micro-kit/misc/errcode"
)

func GetValue(key string) ([]byte, error) {
	return DefaultDatacenter.GetValue(key)
}

func PutValue(key string, value []byte) error {
	return DefaultDatacenter.PutValue(key, value)
}

func GetFromPrefix(key string) (map[string][]byte, error) {
	return DefaultDatacenter.GetPairsFromPrefix(key)
}

func GetValueByAllDatacenter(key string) (string, error) {
	dcs, err := GetDatacenters()
	if err != nil {
		return "", errors.WithMessagef(errcode.ErrNotFound, "key = %v", key)
	}
	for _, dc := range dcs {
		val, err := NewDatacenter(dc).GetValue(key)
		if err != nil || len(val) == 0 {
			continue
		}
		return string(val), nil
	}
	return "", errors.WithMessagef(errcode.ErrNotFound, "key = %v", key)
}

type Pair struct {
	Key   string
	Value []byte
}

// 往consul中放入key - value
func (dc *Datacenter) PutValue(key string, value []byte) error {
	_, err := dc.getConsul().KV().Put(&api.KVPair{
		Key:   key,
		Value: value,
	}, dc.getWriteOption())
	return errors.WithStack(err)
}

// 从consul中删除key对应的value
func (dc *Datacenter) DeleteKey(key string) error {
	_, err := dc.getConsul().KV().Delete(key, dc.getWriteOption())
	return err
}

// 根据key从数据中心中获取value
func (dc *Datacenter) GetValue(key string) ([]byte, error) {
	kvp, _, err := dc.getConsul().KV().Get(key, dc.getQueryOption())
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if kvp == nil {
		return nil, nil
	}
	return kvp.Value, nil
}

func (dc *Datacenter) GetPairs(key string) ([]Pair, error) {
	pairs, _, err := dc.getConsul().KV().List(key, dc.getQueryOption())
	if err != nil {
		return nil, errors.WithStack(err)
	}
	ret := make([]Pair, len(pairs))
	for idx := range pairs {
		ret[idx] = Pair{pairs[idx].Key, pairs[idx].Value}
	}
	return ret, nil
}

func (dc *Datacenter) GetPairsPrefix(prefix string) (map[string]string, error) {
	pairs, _, err := dc.getConsul().KV().List(prefix, dc.getQueryOption())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	ret := make(map[string]string, len(pairs))
	for _, item := range pairs {
		ret[item.Key] = string(item.Value)
	}
	return ret, nil
}

func (dc *Datacenter) GetPairsFromPrefix(prefix string) (map[string][]byte, error) {
	pairs, _, err := dc.getConsul().KV().List(prefix, dc.getQueryOption())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	ret := make(map[string][]byte, len(pairs))
	for _, item := range pairs {
		ret[item.Key] = item.Value
	}
	return ret, nil
}
