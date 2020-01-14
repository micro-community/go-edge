# x-edge

Edge framework for device connections modified from go-micro.
It supports raw tcp/udp.

## 说明

+ 当前项目还处于draft状态，变动可能会很大.
+ master分支，服务设计并不规范，但是可以用run.
+ dev-experiment 是做规范化micro接口风格的规范设计.


## principal concept

+ config better than code 配置优于代码
+ conversion better than config 约定优于配置
+ protocol data unit of each type mappping to a handler 每一个类协议数据包mapping 到一个具体handle
+ 约定协议数据包的分割方式

## How to Use

```bash
go get -u -v github.com/micro-community/x-edge

```
