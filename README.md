# small-website-monitor
small website monitor ( 小型站点监测平台 ) 


# 简介

small website monitor ( 小型站点监测平台 ) 是利用客户端视角对站点进行监测。
为什么叫 “小型站点监测平台” 小型在这里指简单(使用简单，零上手难度)，轻量(部署简单方便一键部署)，局限，
局限在这里表示从客户端视角监测服务端具有局限性，服务端对客户端是黑盒的，不过是从客户端视角监测的真实性是有保障的。

small website monitor拥有站点监测报警（通过邮件进行报警通知），死链检查，证书检查，TDK检查功能。

small website monitor还扩展了Mysql监测工具，SqlServer监测工具，Redis监测工具。

small website monitor也是一款基于gin+vue+Bootstrap5的标准实践项目，简洁直观的项目结构与编码细节，供大家参考学习；

# v1 基础架构
整体结构MVC
后端: gin 
UI: Bootstrap v5

# V2 基础架构
整体架构CS
Master(MVC) 管理，监测汇总等
Slave 监测，采集等

# V3 基础架构
整体架构CS+插件
CS: Master+Slave
插件: 分析插件用于分析采样数据提供解决方案实现智能分析问题
使用插件的目的是为了扩展，不同类型的问题分析各异，提升整体系统的可插拔性

# 期望
v1 : 基础公共
V2 : 分布式监测，地域监测点
v3 : 智能问题分析

### 进度计划
- 站点信息详情 链接复制，链接加入监测点
- 首页 - 查看日志
- 首页 - 删除
- 缓存记录请求信息
- 首页显示请求信息
- 首页 - 设置
- 首页 - 图表
