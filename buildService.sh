#!/bin/bash
##################################################################################
#This file contains the setup by step process for building QAPI app.
# Step 1 - Setup the GO Compiler
# Step 2 - run godeps
# Step 3 - Invoke Build statement
# Step 4 - Packging with tar or zip
#################################################################################
fun_usage()
{
cat << EOF
usage: $0 command

This Scripts builds & packages  qapi components in sevaral flavours and sevaral versions, including Macos, Linux, Arm for both 32 and 64 bit architecture
Dependencies : This Script is dependent on the following
1. go lang setup (PATH variable including the go binary path)


Commands:

    Syntax:
    buildService <args>                      Builds qFactory components based on the argument values passed.

Flags:
  -d, --debug                       Output in debug mode
  -p, --path                        APP_HOME ( APP_WORKSPACE ) path
  -n, --APP name                    Application name
  -v, --version                     Release Version number in semantic verisonformat ( ex 0.0.1 )
EOF
}

fun_Notify() {
  if [ "$DEBUG" = "1" ] ; then echo -e "\nINFO : $1" ;
fi
}

fun_Error() {
   echo -e "\n ERROR : $1"
   exit -1
}

export GOOS=""
export GOARCH=""
export APP_NAME="goelastic"

##########################################################################################
###
### Check 1 : check if APP_HOME has set or qapi-home-path has been passed as an argument
###
##########################################################################################
DEBUG=0 ## Set DEBUG to default value
i=1

if [ $# -lt 1 ]; then
    fun_usage
    exit 1
fi

while [[ $# -gt 1 ]]
do
key="$1"
case $key in
    -p|--path)
    SOURCE_PATH="$2"
    shift # past argument
    ;;
    -v|--version)
    BUILD_VERSION="$2"
    shift # past argument
    ;;
    -n|--name)
    APP_NAME="$2"
    shift # past argument
    ;;
    -h|--help)
    fun_usage
    exit
    #shift # past argument
    ;;
    -d|--debug)
    DEBUG=1
    #shift # past argument
    ;;
    *)
            # unknown option
    ;;
esac
shift # past argument or value
done

### Validate SOURCE_PATH (or)APP_HOME
if [ "$SOURCE_PATH" == "" -a "$APP_HOME" == "" ]
then
	fun_Error "Missing mandatory input -p Path to qapi home (or) APP_HOME enviornment varialble \n";
elif [ -d "$SOURCE_PATH" ]
then
	export APP_HOME=$SOURCE_PATH
elif [ -d "$APP_HOME" ]
then
    fun_Notify "Reading from envrionment variable APP_HOME=${APP_HOME}"
else
	fun_Error "Invalid path [ $SOURCE_PATH ]"
fi

#2).
### Validate BUILD_VERSION  param (or)VERSION_ID enviornment variable
if [ "$BUILD_VERSION" == "" -a "$VERSION_ID" == "" ]
then
	fun_Error "Missing mandatory input -v Version string (ex 0.0.1) (or) VERSION_ID enviornment varialble";
elif [ "$BUILD_VERSION" != "" ]
then
	export VERSION_ID=$BUILD_VERSION
elif [ "$VERSION_ID" != "" ]
then
	fun_Notify "Reading from envrionment variable VERSION_ID=${VERSION_ID}"
else
	fun_Error "Invalid Version $VERSION_ID "
fi



# Install runtime configuration
# Directory heirarchy where the pkg, src, and bin files are located for this app source
##########################################################################################
###
### Step 2  : Update enviornment variables for PATH & GOROOT
###
##########################################################################################

#source $APP_HOME/build/installBase.sh -p $APP_HOME

env | grep -i go

if [ ! -f  "`which go`" ]
then
	#fun_Error "go binary path not setup - Update your PATH varialble or invoke $APP_HOME/bin/installBase.sh"
fun_Error "The program 'go' is currently not installed. You can install it by typing:\n $APP_HOME/build/installBase.sh"

fi

fun_Notify "Current Location : `pwd`"
fun_Notify "Go Environment `go env`"


#for GOOS in darwin linux windows; do
for GOOS in darwin linux ; do
  for GOARCH in amd64 ; do
    echo "Building  goelastic_${VERSION_ID}_${GOOS}_${GOARCH}"
    export GOOS=$GOOS
    export GOARCH=$GOARCH
    ##Make the actual binary
  #go build -ldflags="-s -w"
  BUILD_VERSION="$2"
  shift # past argument
  ;;
  -n|--name)
  APP_NAME="$2"
    go build {$APP_NAME}_{$BUILD_VERSION}_{$GOOS}
  done
done

export GOOS=""
export GOARCH=""
