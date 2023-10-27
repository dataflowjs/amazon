package main

import (
	"gorm.io/gorm"
)

type Settings struct {
	gorm.Model
	TelegramName string `form:"telegramName"`
	TelegramKey  string `form:"telegramKey"`
	StartTime    int    `form:"startTime"`
	EndTime      int    `form:"endTime"`
	Pause        int    `form:"pause"`
	Proxy        string `form:"proxy"`
}

var Stg Settings

var Start bool
