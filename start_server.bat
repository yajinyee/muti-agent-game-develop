@echo off
title 吉伊卡哇 Game Server
echo ================================
echo  吉伊卡哇：像素大討伐 Server
echo  Port: 7777
echo ================================
cd /d d:\Kiro\server
go run cmd/gameserver/main.go
pause
