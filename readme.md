
## Workspace  

// protoc --go_out=. *.proto // uses official goprotobuf
// protoc --go_out=plugins=grpc:. *.proto // official with grpc support
// protoc --gofast_out=. *.proto // uses gogoprotobuf
protoc -I=$SRC_DIR --go_out=$DST_DIR $SRC_DIR/addressbook.proto

protoc --go_out=protobuf -I=proto proto/*.proto

protoc --go-grpc_out=protobuf -I=proto proto/*.proto

protoc --gogo_out=protobuf -I=proto proto/*.proto  

protoc --gogo_out=protobuf -I=${GOPATH}/pkg/mod -I=proto proto/*.proto

I can't remember why/how it works, but all of the imported .proto files are under ~/.local/include. I didn't think that was an automatic search path for protoc or linking and it's not in my ld_library_path or similar. but if you need additional .proto's to import, you can clone them there. and make sure the .proto you're importing is at ~/.local/include/import/path/name, then use the import "import/path/name.proto";


