cd server
haxe -main Server -D static_link -cpp cpp
cd ..
go build
./hx-cpp-server

