cd src
set -x
gomobile clean
 GOPATH=`pwd`/.. gomobile bind -v -target ios -tags debug -gcflags=" -E -K -j -l -r -v -w" ngrok/client
 GOPATH=`pwd`/.. gomobile bind -v -target ios -tags debug -gcflags=" -E -K -j -l -r -v -w" ngrok/server
