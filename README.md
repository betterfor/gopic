# gopic

[简体中文](https://github.com/owenthereal/goup/blob/master/README_ZH.md)

According to [PicGo](https://github.com/Molunerfinn/PicGo)

**GoPic is a terminal tool for quickly uploading images and getting URL links to images.**

Support Image Gallery:

- [x] github

- [x] gitee

- [x] smms

- [x] qiniu

- [x] imgur

- [x] aliyun

- [x] upyun

- [x] tcyun

All support image gallery is here, you can develop third-part plugins with [core](./core/plugins/uploader), welcome to
give me PR.

### Quick Start

Download and install

 ```shell
go get -u github.com/betterfor/gopic
```

Show help:

```shell
$ gopic -h
Gopic is a tool for uploading images.
It's easily, quickly, conveniently.
After your uploading images, you can get a link to save in your blog|markdown|article...'

Usage:
  gopic [command]

Available Commands:
  config      configuration
  help        Help about any command
  upload      upload local file to remote before completing setting

Flags:
      --config string   config file (default is $HOME/.gopic/config.yaml)
  -d, --debug           Help message for debug
  -h, --help            help for gopic
```

Before upload:

```shell
$ gopic config set [key]=[value]
```

> require is replace '*'

### github

- repoName(*): the name of warehouse, example: username/repo
- branch: project branch, default is master
- token(*): github personal access tokens
- path: storage path in github, default is images
- customUrl: convert url to custom url, example: https://xxx.com
- enableTimeFile: bool, after path, use date, like 2016/01/02

### gitee

- repoName(*): the name of warehouse, example: username/repo
- branch: project branch, default is master
- path: storage path in gitee, default is images
- email(*): user email in gitee
- password(*): user password in gitee
- clientId(*): api clientId for generate accessToken
- clientSecret(*): api clientSecret for generate accessToken
- enableTimeFile: bool, after path, use date, like 2016/01/02

### aliyun

- bucket(*): storage bucket, example: oss-cn-beijing
- area(*): storage area
- accessKeyId(*): access key id
- accessKeySecret(*): access key secret
- customUrl: custom domain name, example: https://xxx.com
- path: storage path, example: img/

### imgur

- clientId(*): client id
- proxy: proxy, example: http://127.0.0.1:1080

### qiniu

- accessKey(*): api access key
- secretKey(*): api secret key
- bucket(*): storage bucket
- area(*): storage area, example: z0
- url(*): access to the address, example: http://xxx.yyy.glb.clouddn.com
- option: url suffix, example: ?imageslim
- path: storage path, example: img/

### smms

- token(*): token

### tcyun

- secretId(*): api secret id
- secretKey(*): api secret key
- appId(*): app id, example: 1234567890
- bucket(*): storage bucket
- area(*): storage area, example: tj
- path: storage path, example: img/
- customUrl: custom domain, example: https://xxx.com

### ypyun

- bucket(*): storage bucket
- operator(*): setting operator
- password(*): operator password
- url(*): custom domain: example: http://xxx.test.upcdn.net
- options: url suffix, example: imgslim
- path: storage path, example: img/

## support typora upload

After configure gopic, use typora settings

File->Preferences->Image->Image Upload Setting->Image Uploader select **Custom Comman**,
**Command** file in `gopic upload`，test **Test Uploader**

![typora upload](https://cdn.jsdelivr.net/gh/betterfor/cloudImage/images/2021/07/20/typora_gopic_cn.png)

## LICENSE

[Apache 2.0](https://github.com/betterfor/gopic/blob/master/LICENSE)