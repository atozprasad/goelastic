sudo yum install -y wget
sudo yum install -y git
wget https://dl.google.com/go/go1.10.1.linux-amd64.tar.gz
mkdir -p ~/Install
mkdir -p ~/code/go/
tar -C ~/Install -xzf go1.10.1.linux-amd64.tar.gz

mkdir -p ~/code/go/bin
mkdir -p ~/code/go/str
mkdir -p ~/code/go/pkg



cat >> ~/.bash_profile << IAM_DONE_HERE
export GOPATH=~/code/go/
export GOROOT=~/Install/go
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin:.
export BASEPATH=/Users/atozprasad/code/go
export GOBIN=$GOROOT/bin
export GO15VENDOREXPERIMENT=1
export GO_VERSION=1.10.1
export PATH=$PATH:$GOBIN:/usr/local/bin:/usr/bin:/usr/local/sbin:/usr/sbin:/bin
export GOARCH="amd64"
export GOCHAR="6"
export CC="gcc"
export GOGCCFLAGS="-fPIC -m64 -pthread -fmessage-length=0"
export CXX="g++"
export PATH=$PATH:$GOROOT/bin:$GOROOT:$GOPATH
IAM_DONE_HERE

source ~/.bash_profile

cd  ~/code/go/str
git clone https://github.com/atozprasad/goelastic.git
cd goelastic
go get github.com/gin-gonic/gin
go get github.com/olivere/elastic
go get github.com/teris-io/shortid
go build
