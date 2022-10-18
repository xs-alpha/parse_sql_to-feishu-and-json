# parse_sql_to-feishu-and-json

Quickly convert sql to feishu form and json file

# use
put sql into parse1.sql and then turn off the file

run software and you will get two file

> you can change filename and filepath by clone source code and change filename located in tools/const.go



# 使用
把sql建表语句放到parse.sql文件中，关闭文件, 没有的话需要新建

运行软件，将会得到两个文件，

> 可以下载源码来更改文件名和路径，相对应的代码在tools./const.go

自动生成配置文件

支持解析
  - sql生成飞书表格 && 并生成json
  - 飞书表格字段到json
