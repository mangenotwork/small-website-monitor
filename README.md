# small-website-monitor
small website monitor ( 小型站点监测平台 ) 


# 简介

small website monitor ( 小型站点监测平台 ) 是利用客户端视角对站点进行监测。
为什么叫 “小型站点监测平台” 小型在这里指简单(使用简单，零上手难度)，轻量(部署简单方便一键部署)，局限，
局限在这里表示从客户端视角监测服务端具有局限性，服务端对客户端是黑盒的，不过是从客户端视角监测的真实性是有保障的。

small website monitor拥有站点监测报警（通过邮件进行报警通知），死链检查，证书检查，TDK检查功能。

small website monitor还扩展了Mysql监测工具，SqlServer监测工具，Redis监测工具。

small website monitor也是一款基于gin+vue+Bootstrap5的标准实践项目，简洁直观的项目结构与编码细节，供大家参考学习；

# 架构选型

### v1
整体结构MVC
后端: gin 
JS: vue3, jq
UI: Bootstrap_v5

### V2
整体架构CS
Master(MVC gin+vue3+jq+Bootstrap_v5) 管理，监测汇总等
Slave 监测，采集等

### V3
整体架构CS+插件
CS: Master+Slave
插件: 分析插件用于分析采样数据提供解决方案实现智能分析问题
使用插件的目的是为了扩展，不同类型的问题分析各异，提升整体系统的可插拔性


### 进度计划
- 首页 - 设置
- 首页 - 图表
