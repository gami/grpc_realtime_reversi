## About
「STARTING gRPC」第7章のためのサンプルです。

protoc \
	-Iproto \
	--go_out=plugins=grpc:. \
	proto/*