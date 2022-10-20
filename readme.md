
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


##

SQL:

- create table
- create temp table
- insert into temp table, ignore constraints
- insert from temp table into table, follow constraints

- identifiers in order of preference: 
-  unique row id's if src has them for that dataset
-  complete row hash of ordered src data, if all fields are meaningful for that dataset
-  incomplete row has of ordered src data, omitting unmeaningful fields

eg: schedule entry has schedulekey as a unique row id - use schedulekey
    staffmember omits any type of versioning id and it includes a last login timestamp - order data, including child slices, omit last login timestamp, hash resulting json (which is the src format)

