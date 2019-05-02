set -x
set -e
make
cd src

if [ ! -x gomobile ]
then
  go get golang.org/x/mobile/cmd/gomobile
  gomobile init
fi

GOPATH=`pwd`/.. go get gopkg.in/yaml.v1
GOPATH=`pwd`/.. go get github.com/alecthomas/log4go
GOPATH=`pwd`/.. go get github.com/rcrowley/go-metrics
GOPATH=`pwd`/.. go get github.com/inconshreveable/go-vhost
GOPATH=`pwd`/.. go get github.com/nsf/termbox-go

GOPATH=`pwd`/.. go get golang.org/x/mobile/cmd/gomobile

gomobile clean
 GOPATH=`pwd`/.. gomobile bind -v -target ios -tags debug -gcflags=" -E -K -j -l -r -v -w" ngrok/client
 GOPATH=`pwd`/.. gomobile bind -v -target ios -tags debug -gcflags=" -E -K -j -l -r -v -w" ngrok/server
