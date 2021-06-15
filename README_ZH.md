# gopic

** Gopic **:一个用于快速上传图片并获取图片URL链接的命令行工具。

现支持的图床：

- github
- gitee
- smms

这是当前已经支持的图床，你可以参考[core](./core/cloud)自行开发图床插件,欢迎提交PR。

## 快速开始

下载和安装

 ```bash
go get -u github.com/betterfor/gopic
```

### github配置

通过设置

- repoName 仓库名称
- branch 仓库分支
- token 个人私钥
- path 存储路径

来配置github存储配置

**Example**

 ```bash
# set github config
$ gopic config --set github.branch=master
# use github config
$ gopic config --use github
# list config 
$ gopic config --list
current: github
github:
  repoName: betterfor/cloudImage
  branch: master
  token: xxx
  path: images/{2006/01/02}
```

[接口信息](https://docs.github.com/en/free-pro-team@latest/rest/reference/repos#create-or-update-file-contents)

**cdn**

该CDN方法针对于github上传设置，可以使用CDN的方式加快图片等静态文件的查看。

可用：

- [jsdelivr](https://www.jsdelivr.com/)

### gitee设置

预先申请OAuth的秘钥，然后可用 [api](https://gitee.com/api/v5/swagger#/postV5ReposOwnerRepoContentsPath) 上传文件

```shell
$ gopic config --list
gitee:
  repoName: zongl/cloudImage
  branch: master
  path: images/{2006/01/02}
  email: xxx
  password: xxx
  clientId: xxx
  clientSecret: xxx
```

### SMMS设置

用户申请完[API Secret](https://sm.ms/home/apitoken)后，填写配置后即可

```shell
$ gopic config --list
smms:
  token: xxx
```

## 支持typora上传

在typora的文件->偏好设置->图像->上传服务设定->上传服务选择**Custom Comman**,命令填写`gopic upload`，测试**验证图片上传选项**

## LICENSE

[Apache 2.0](https://github.com/owenthereal/goup/blob/master/LICENSE)