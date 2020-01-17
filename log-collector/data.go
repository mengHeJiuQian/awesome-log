package main

/**
 * Create Time : 2020/1/15 下午3:09
 * Update Time :
 * Author : sheldon
 * Decription :
 */

type logConfig struct {
	Topic    string `json:"topic"`
	LogPath  string `json:"log_path"`
	Service  string `json:"service"`
	SendRate int    `json:"send_rate"`
}
