package model

import (
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/mangenotwork/common/log"
	"github.com/mangenotwork/common/utils"
)

// WebSite 监测站点
type WebSite struct {
	// 主键自增, 用于数据查询,分页的key
	ID string

	Host string

	// 监测频率单位 ms
	Rate int64

	// health 指定的生命监测uri
	HealthUri string

	// 探寻站点深度 默认 2
	UriDepth int64

	// 更新 每层uri的时间 单位小时
	UriUpdateRate int64

	// 设置超过这个响应时间报警 单位ms
	AlarmResTime int64

	HostIP string

	Created int64
}

func (m *WebSite) Add() error {
	DB.Open()
	defer func() {
		_ = DB.Conn.Close()
	}()
	return DB.Conn.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(WebSiteTable))
		if b == nil {
			return fmt.Errorf(WebSiteTable + "表不存在")
		}
		c := b.Cursor()
		k, _ := c.Last()
		var id int64 = 0
		if len(string(k)) > 0 {
			id = utils.AnyToInt64(string(k))
			log.Info("kInt = ", id)
		}
		m.ID = utils.AnyToString(id + 1)
		value, err := utils.AnyToJsonB(m)
		if err != nil {
			log.Error(err)
			return err
		}
		log.Info("m.ID = ", m.ID)
		err = b.Put([]byte(m.ID), value)
		if err != nil {
			log.Error(err)
			return err
		}
		return nil
	})
}

// List 分页获取
func (m *WebSite) List(pg, size int) ([]*WebSite, error) {
	DB.Open()
	defer func() {
		_ = DB.Conn.Close()
	}()
	start := (pg - 1) * size
	end := pg * size
	data := make([]*WebSite, 0)
	err := DB.Conn.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(WebSiteTable))
		if b == nil {
			return fmt.Errorf(WebSiteTable + "表不存在")
		}
		c := b.Cursor()
		i := 0
		for k, v := c.First(); k != nil; k, v = c.Next() {
			i++
			if i > start || i < end {
				fmt.Printf("key=%s, value=%s\n", k, v)
				value := &WebSite{}
				e := json.Unmarshal(v, value)
				if e != nil {
					log.Error("数据解析错误")
				}
				data = append(data, value)
			}
		}
		return nil
	})
	return data, err
}

// Get 指定获取
func (m *WebSite) Get(k string) (*WebSite, error) {
	value := &WebSite{}
	err := DB.Get(WebSiteTable, k, value)
	return value, err
}

// Update 更新数据
func (m *WebSite) Update(k string) error {
	return DB.Set(WebSiteTable, k, m)
}

// Delete 删除数据
func (m *WebSite) Delete(k string) error {
	return DB.Delete(WebSiteTable, k)
}
