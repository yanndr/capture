# capture
Grpc service to capture an image from a video

## About
I needed to extract an image from a video from a personal website. I didn't want this to be part of the other website, so I created a service that exposed the function by GRPC. 
Why GRPC? Because I wanted to learn about that, [Go-kit](https://github.com/go-kit/kit) and [Promotheus](https://github.com/prometheus/client_golang)

## Installation
This service depends on [screengen](github.com/opennota/screengen), this require to have [ffmpeg](https://ffmpeg.org/) installed

```
go get github.com/yanndr/capture
```
