# x-edge

Edge framework for device connections modified from go-micro.
It supports raw tcp/udp.

## Struct

<div align="center">
    <img src="https://github.com/micro-community/x-edge/blob/master/Struct.png">
</div>

+ Edge结构取消了go-micro中的Registry、Broker、Selector，保留了Client、Server、Transport和Coder

## Data Flow

<div align="center">
    <img src="https://github.com/micro-community/x-edge/blob/master/data%20flow.png">
</div>

+ 启动流程: go-micro  启动----> x-edge 启动--->x-edge监听---> go-micro 监听
+ transport从device（client端）收到tcp或者udp数据包（数据包格式可以自定义，默认是xml），调用edgeServer
+ edgeServer通过Coder解码，并通过router，找到相应的handler
+ 在service的handler中，处理相应的业务message proc
+ 在message proc可以调用其他Service的broker或者rpc，将数据send出去。同时，也可以通过edgeClient回复数据包给device（client端）

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
