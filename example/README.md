# Instruction

extractor 定义了，如何从UDP数据包/TCP数据流中如何抽取一个完成的协议数据数据包

## example说明

example设计了一个简单的xml交互协议,如下：

```xml
<?xml version="1.0" encoding="gb2312"?>
<PROTOCOL>
<VER>1.0</VER>
<NAME>danny-event-type</NAME>
<TYPE>Event</TYPE>
</PROTOCOL>
```

协议数据包以</PROTOCOL>进行分割。

## 关系说明

| 服务关系|||
| -------------| :-------------:|  -------------|
|服务|<--raw tcp/udp -->|设备|

|交互说明|||
|--------|:--------:|--------|
|设备|--- event通知---->|服务|
|服务|--- control请求 ---->|设备|

下述数据包标表示：网关/控制器/盒子/向当前服务发送的一个事件(Event)通知

```xml
<?xml version="1.0" encoding="gb2312"?>
<PROTOCOL>
<VER>1.0</VER>
<NAME>danny-event-type</NAME>
<GENDER>MALE</GENDER>
<TYPE>Event</TYPE>
<ADDR>Road.1</ADDR>
<PHONE>400-800-5555</PHONE>
<COMPANY>xxx</COMPANY>
<TIME>2019.12.1-11:11:11</TIME>
</PROTOCOL>
```
