# 简介

这是一个用于将 s3 bucket 分割成多个子用户的 rclone auth proxy 插件

# 如何使用

### 安装

懒得配 Github Actions, 直接从源码构建 (等有空的时候配上)

```sh
go install github.com/shynome/s3-split@latest
```

### 启动

所有用户都使用同一个 `SECRET_ACCESS_KEY`(复杂一些就行了), 认证反而依靠于 `ACCESS_KEY_ID` (因为这是 rclone seve s3 里的逻辑, 不改了就按这样来工作量最小可维护性最好)

```sh
rclone serve s3 --auth-proxy s3-split --auth-key ACCESS_KEY_ID,SECRET_ACCESS_KEY
```

### 配置

在 `rclone.conf` 配置 `s3users` 块, 添加以下内容

```conf
[local]
type = local

[s3users]
# ACCESS_KEY_ID 就是密码, 记得设置的复杂点, 分用户是靠路径分割的
# [ACCESS_KEY_ID] = [后端:路径]
ACCESS_KEY_ID = local:/tmp/s3-split/ACCESS_KEY_ID

[s3test]
type = s3
provider = Rclone
endpoint = http://127.0.0.1:8080/
access_key_id = ACCESS_KEY_ID
secret_access_key = SECRET_ACCESS_KEY
use_multipart_uploads = false
```

### 测试

```sh
# 创建测试用目录
mkdir -p /tmp/s3-split/ACCESS_KEY_ID/bucket1
# 使用 webdav 测试管理文件
rclone serve webdav --addr 127.0.0.1:8000 s3test:/bucket1
```

在文件管理器打开 `webdav://127.0.0.1:8000/` , 看看上传是否正常
