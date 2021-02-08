# 简介

**Nightingale接入prometheus query接口做数据源，替代Alertmanager作为CNCF监控体系中告警模块的补充。**


# 使用方式

正常安装夜莺后，安装插件。

作为Nightingale的插件，用于收集prometheus的指标

prometheus作为优秀的开源监控产品，本身不仅完整的指标体系，还拥有丰富的指标采集解决方案。

Prometheus-exporter-collector以插件形式集成在collector中，通过Nightingale插件采集，collector采集目标exporter指标并上报

## 快速构建 

    $ mkdir -p $GOPATH/src/github.com/n9e
    $ cd $GOPATH/src/github.com/n9e
    $ git clone https://github.com/JamesJoe-C/n9e-query-prometheus.git
    $ cd prometheus-exporter-collector
    $ export GO111MODULE=on
    $ export GOPROXY=https://goproxy.cn
    $ go build
    $ cat plugin.test.json | ./n9e-query-prometheus


 ### 配置参数
 Name                             |  type     | Description
 ---------------------------------|-----------|--------------------------------------------------------------------------------------------------
 exporter_urls                    | array     | Address to collect metric for prometheus exporter.
 append_tags                      | array     | Append tags for n9e metric default empty
 endpoint                         | string    | Field endpoint for n9e metric default empty
 ignore_metrics_prefix            | array     | Ignore metric prefix default empty
 timeout                          | int       | Timeout for access a exporter url default 500ms
 metric_prefix                    | string    | append metric prefix when push to n9e. e.g. 'xx_exporter.'
 metric_type                      | map       | specify metric type
 default_mapping_metric_type      | string    | Default conversion rule for Prometheus cumulative metrics. support COUNTER and SUBTRACT. default SUBTRACT
 ###
 
 ###
 
