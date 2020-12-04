# gopic

** Gopic:一个用于快速上传图片并获取图片URL链接的命令行工具。

现支持的图床：
- github

这是当前已经支持的图床，你可以参考[core](./core)自行开发图床插件,欢迎提交PR。

### 快速开始
下载和安装
 ```bash
go get -u github.com/betterfor/gopic
```

设置配置
- github
 ```bash
gopic config set --file ./example.yaml
```