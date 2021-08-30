# save-to-disk
[原例程](https://github.com/pion/webrtc/tree/master/examples/save-to-disk)

功能：将网页捕获的音视频数据通过udp端口传输

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
$ export SIGNALADDR="127.0.0.1:19801" 
$ ./client
```

## 打开网页客户端
将demo文件夹下的demo.html拖入浏览器中，`Browser base64 Session Description`和`Golang base64 Session Description`都出现内容后，点击`start session`

使用vlc打开目录下的`rtp-fowwarder.sdp`文件，观看直播