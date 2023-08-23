//go:build k8s

/*
 * @Author: hugo lee hugo2lee@gmail.com
 * @Date: 2023-08-22 17:57
 * @LastEditors: hugo lee hugo2lee@gmail.com
 * @LastEditTime: 2023-08-22 21:29
 * @FilePath: /geektime-basic-go/webook/config/k8s.go
 * @Description:
 *
 * Copyright (c) 2023 by hugo, All Rights Reserved.
 */

// 使用 k8s 这个编译标签
package config

var Config = config{
	DB: DBConfig{
		// 本地连接
		DSN: "root:root@tcp(webook-live-mysql:3308)/webook",
	},
	Redis: RedisConfig{
		Addr: "webook-live-redis:6380",
	},
}
