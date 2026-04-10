apt update -y
apt install -y libvips-dev libvpx9 libopus0 ffmpeg

go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
go install github.com/amacneil/dbmate@latest
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
