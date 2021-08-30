# save-to-disk
[原例程](https://github.com/pion/webrtc/tree/master/examples/save-to-disk)

## 运行信令服务
进入到上级目录signal工程目录下,执行：
```
$ make build
$ cd build
//默认19801端口
$ ./client -port=19801
```

## 运行服务
进入工程目录，执行：
```
$ make build
$ cd build 
$ ./client -addr="127.0.0.1:19801"
```

## 打开网页客户端
将demo文件夹下的demo.html拖入浏览器中，`Browser base64 Session Description`和`Golang base64 Session Description`都出现内容后，点击`start session`

在文件执行目录下，会出现.ivf和.ogg两个文件，使用格式工厂混流后播放录制的视频。