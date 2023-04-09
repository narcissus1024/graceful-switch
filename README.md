# Graceful Switch

## Material

ssh config说明

1. 重复key, 以第一个为主

2. \#开头或空格被认为是注释

3. 从host开始到下一个host结束，这期间的配置被认为是同一个ssh配置；在这之外的被认为是全局的配置

4. kv可以使用 “=”或“ ”（空格）连接

5. 每行配置开头是否有空格或制表符无关紧要

参考：

[ssh官方文档](https://man.openbsd.org/ssh_config.5#SSH_CONFIG_FILE_FORMAT)

[ssh学院文档](https://www.ssh.com/academy/ssh/config)

## TODO 

### feature

### function

- [x] 加载系统ssh,与应用配置进行合并（通过Host字段）
- [ ] save时，如果内容没变，不存储
- [ ] 添加删除配置项功能
- [ ] 内容为空，save时，给出提示
- [ ] save后，刷新编辑页面text

### bugfix

- [x] 配置为空时，点击save会保存一个名为“.json”的文件
- [x] save时，将sys，id为0的配置存储了
- [x] inner为空，sys不为空。初次启动，点击系统配置，点击save，config变为空
- [x] 重复点击save，内容就没了
- [ ] 正则表达式有问题（ps:验证ssh是如何判断k==v格式的，k是key，=v是value吗？）
- [x] 在系统配置编辑，删除某段配置，save后，并没有删除【隐藏save按钮，不允许编辑】
- [ ] 解决注释问题、第一行非host问题
- [ ] 对inner索引和数据进行排序