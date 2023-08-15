package model

// Uri 站点的Uri存储
type Uri struct {
	// 外链监测站点主键
	HostID string

	// 各层的uri, map 层号:[]uri
	UriLayer1 map[int][]string

	// 死链接
	DeadLink []string

	// 监测点 Uri
	UriPoint []string
}
