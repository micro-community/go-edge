# docs

Instructions

## Struct

<div align="center">
    <img src="images/Struct.png">
</div>

+ Edge结构取消了go-micro中的Registry、Broker、Selector，保留了Client、Server、Transport和Codec

## Data Flow

<div align="center">
    <img src="images/data%20flow.png">
</div>

+ 启动流程: go-micro  启动----> x-edge 启动--->x-edge监听---> go-micro 监听
+ transport从device（client端）收到tcp或者udp数据包（数据包格式可以自定义，默认是xml），调用edgeServer
+ edgeServer通过Codec解码，并通过router，找到相应的handler
+ 在service的handler中，处理相应的业务message proc
+ 在message proc可以调用其他Service的broker或者rpc，将数据send出去。同时，也可以通过edgeClient回复数据包给device（client端）