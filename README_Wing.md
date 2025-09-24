
添加一个agent，直接读取qq音乐的信息

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
