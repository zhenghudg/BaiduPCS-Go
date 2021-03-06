package pcscache

import (
	"github.com/iikira/BaiduPCS-Go/baidupcs"
	"time"
)

type dirCache struct {
	fdl      map[string]*baidupcs.FileDirectoryList
	lifeTime time.Duration // 生命周期
}

// Set 设置网盘目录缓存
func (dc dirCache) Set(path string, dirInfo *baidupcs.FileDirectoryList) {
	dc.fdl[path] = dirInfo
}

// Existed 检测缓存是否存在
func (dc dirCache) Existed(path string) bool {
	_, existed := dc.fdl[path]
	return existed
}

// Get 取出网盘目录缓存
func (dc dirCache) Get(path string) *baidupcs.FileDirectoryList {
	return dc.fdl[path]
}

// FindFileDirectory 网盘目录缓存内查找文件
func (dc dirCache) FindFileDirectory(path, name string) *baidupcs.FileDirectory {
	fdl := dc.Get(path)
	if fdl == nil {
		return nil
	}
	for _, fd := range *fdl {
		if fd.Filename == name {
			return fd
		}
	}
	return nil
}

func (dc dirCache) SetLifeTime(t time.Duration) {
	dc.lifeTime = t
}

// GC 缓存回收
func (dc dirCache) GC() {
	go func() {
		ticker := time.NewTicker(dc.lifeTime)
		for {
			select {
			case <-ticker.C:
				dc.DelAll()
			}
		}
	}()
}

// Del 删除网盘目录缓存
func (dc dirCache) Del(path string) {
	delete(dc.fdl, path)
}

// DelAll 清空网盘目录缓存
func (dc dirCache) DelAll() {
	for k := range dc.fdl {
		delete(dc.fdl, k)
	}
}
