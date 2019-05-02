set -x
set -e

export GOPATH="`pwd`"
export PATH="$GOPATH/bin:$PATH"

make
cd src

if [ ! -x gomobile ]
then
  go get golang.org/x/mobile/cmd/gomobile
  gomobile init
fi

go get gopkg.in/yaml.v1
go get github.com/alecthomas/log4go
go get github.com/rcrowley/go-metrics
go get github.com/inconshreveable/go-vhost
go get github.com/nsf/termbox-go
go get golang.org/x/mobile/cmd/gomobile
gomobile clean
gomobile bind -v -target ios -tags debug -gcflags=" -E -K -j -l -r -v -w" ngrok/client
gomobile bind -v -target ios -tags debug -gcflags=" -E -K -j -l -r -v -w" ngrok/server
