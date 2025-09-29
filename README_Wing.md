

## 插件

又实现了一个插件，可以单独编译

```bash
cd plugins/examples
make qqmusic
```

添加配置
```toml
[PluginConfig.wing-qqmusic]
baseURL = "http://192.168.3.50:12000"
```

## 源码修改

### 修改


添加一个agent，直接读取qq音乐的信息 `core/agents/wingqq`

修改 `adapters/taglib/taglib.go` , 直接通过接口获取文件媒体信息


### 打包
```bash
# 先下载 https://github.com/navidrome/cross-taglib 放到项目 var 下 var/taglib-linux-amd64.tar.gz
# wsl
cd navidrome
dos2unix Dockerfile
# 本地测试
make docker-build PLATFORMS=linux/amd64
# 正式打包
make docker-image IMAGE_PLATFORMS=linux/amd64 DOCKER_TAG=ccr.ccs.tencentyun.com/wingao/navidrome:develop 
docker push ccr.ccs.tencentyun.com/wingao/navidrome:develop
```
