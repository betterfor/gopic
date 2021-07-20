# gopic

[简体中文](https://github.com/owenthereal/goup/blob/master/README_ZH.md)

According to [PicGo](https://github.com/Molunerfinn/PicGo)

**GoPic is a terminal tool for quickly uploading images and getting URL links to images.**

Support Image Gallery:

[x] github

[x] gitee

[x] smms

[x] qiniu

[x] imgur

[x] aliyun

[x] upyun

[x] tcyun

All support image gallery is here, you can develop third-part plugins with [core](./core/plugins/uploader), welcome to
give me PR.

### Quick Start

Download and install

 ```bash
go get -u github.com/betterfor/gopic
```

## support typora upload

with typora settings

File->Preferences->Image->Image Upload Setting->Image Uploader select **Custom Comman**,
**Command** file in `gopic upload`，test **Test Uploader**

## LICENSE

[Apache 2.0](https://github.com/owenthereal/goup/blob/master/LICENSE)