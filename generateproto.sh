#!/bin/bash
protoc -I pb/ pb/capture.proto --go_out=plugins=grpc:pb
