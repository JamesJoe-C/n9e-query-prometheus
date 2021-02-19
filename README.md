# 简介


Nightingale接入prometheus query接口做数据源，替代Alertmanager作为CNCF监控体系中告警模块的补充。


# 使用方式

正常安装夜莺后，安装插件。

直接从prometheus的API接口拿取prometheus的已存储数据，方便在已有prometheus+gragana监控体系的情况下，进一步提供告警机制。



## 快速构建 

    $ mkdir -p $GOPATH/src/github.com/n9e
    $ cd $GOPATH/src/github.com/n9e
    $ git clone https://github.com/JamesJoe-C/n9e-query-prometheus.git
    $ cd n9e-query-prometheus
    $ export GO111MODULE=on
    $ export GOPROXY=https://goproxy.cn
    $ go build
    $ cat plugin.test.json | ./n9e-query-prometheus

 
## 插件使用
夜莺配置stdin如下：
```
{
  "exporter_urls": [
    "http://prometheus.url"
  ],
  "query": "sum by (instance) (service_record:pod_qps)",
  "timeout": 5000,
  "metric_prefix":"sum_pod_qps",
  "nid":"dept",
  "endpoint":"192.168.1.1"
}
```
配置文件说明：
exporter_urls中对于线上需要用户名密码验证的prometheus，书写格式如下：https://user:password@prometheus.url

query可以书写任意prometheus的pql语句，但使用聚合函数时因metrice不存在，请书写metric_prefix字段。

endpoint对应夜莺资源树中机器。

nid对应夜莺组织，用于无设备相关监控选项。如果配置nid，endpoint请保持为空或不传递。

## 注意

目前仅测试到Gauge数据格式，后续更新其他数据格式DEMO。

因夜莺包缺少Nid字段，正在沟通中。本项目暂时缺少依赖，后续更新。



