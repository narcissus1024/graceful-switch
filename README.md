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

- [ ] 配置为空时，点击save会保存一个名为“.json”的文件
- [ ] 加载系统ssh,与应用配置进行合并（通过Host字段）
- [ ] 在系统配置选项中，点击保存，存在覆盖ssh配置文件为空的情况
- [ ] 保存功能存在多个问题
- [ ] 系统配置会保存在内部数据目录中（不应该存在）