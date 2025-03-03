# JIK-面试刷题平台

## 项目介绍

​	基于 **Nunu** + **Redis** + **Elasticsearch** 的面试刷题平台；管理员可创建并批量管理题目；用户可 **分词检索** 题目，在线刷题并查看刷题日历；项目基于 **Sentinel Go** ， **Nacos** 等全面优化性能和安全性。基于 Nginx + Linux 管理面板将项目部署上线。

项目地址：jiaking.top

## 技术选型

- Nunu 框架
- MySQL 数据库 + GORM +Redis
- Elasticsearch 搜索引擎
- Sentinel Go 流量控制
- Nacos 配置中心
- 多角度项目优化：性能、安全性、可用性

## 模块介绍

### 用户模块

登录（同端互斥登录），注册，注销，获取用户登录状态，签到（ BitMap 实现），获取签到记录

（管理员）查看用户列表，添加用户，删除用户，修改用户

### 题库模块

查看题库列表 (限流，熔断)

（管理员）查看题库列表，添加题库，删除题库，修改题库

### 题目模块

查看题目列表 (限流，熔断，反爬虫)

（管理员）查看题目列表，添加题目，删除题目，修改题目，批量删除题目，批量添加删除题目题库关系

## 项目效果

### 主页

![image-20250228161228724](C:\Users\DELL\AppData\Roaming\Typora\typora-user-images\image-20250228161228724.png)

### 题库

![image-20250228161256025](C:\Users\DELL\AppData\Roaming\Typora\typora-user-images\image-20250228161256025.png)

### 用户管理

![image-20250303151807168](C:\Users\DELL\AppData\Roaming\Typora\typora-user-images\image-20250303151807168.png)

### 题库管理

![image-20250303151829238](C:\Users\DELL\AppData\Roaming\Typora\typora-user-images\image-20250303151829238.png)

### 题目管理

![image-20250303151914796](C:\Users\DELL\AppData\Roaming\Typora\typora-user-images\image-20250303151914796.png)
