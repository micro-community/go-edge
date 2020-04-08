# x-edge notes

record something for x-edge

## 使用web/service的框架结构，改造x-edge

问题：
+ 按照web/service的结构，micro.Service的位置：type service struct -> opts Options -> Service  micro.Service
    而edgeService的结构：
    type edgeService struct {
	    opts    service.Options
	    service micro.Service
    }

    考虑给edgeService重新定义 Options

+ web/service中微服务Service  micro.Service，在newOptions时创建了，在web/service的Init函数中调用了micro.Service的Init，但是却并未run。
    原因：主要是web/service使用micro.Service的client，不需要run。

    但是对于x-edge来说micro.Service的server需要run起来，因为x-edge需要被其他服务调用，向通过tcp/udp向下发送命令，或者x-edge需要用broker向订阅段推送消息。
    这里需要考虑为x-edge启动micro.Service的server

+ edege 的opinion结构，里面有server和tranport，newService函数里面，先new opinion函数，
    然后就用server的name，所以必须在new opinion函数中new server，而这个server是我们自己edge server
    所以在new opinion函数，就定义了defaultServer，且直接赋值
    同时，也定义了defaultTransport，且直接赋值
