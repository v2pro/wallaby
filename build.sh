#!/bin/bash
#
# export GOROOT=/usr/local/go1.8.3
# export PATH=$GOROOT/bin:$PATH
#
#
WORKSPACE=$(cd $(dirname $0) && pwd -P)
echo "wallaby workspace: $WORKSPACE"
cd $WORKSPACE
module=wallaby
output="output"

# change build timestamp 
sed -i -e '/ProxyBuildTimestamp[[:space:]]*=/ s/=.*/= '`date +%s`'/' config/config.go


go build -o $module cmd/proxy/main.go
ret=$?
if [ $ret -ne 0 ];then
    echo "===== $module build failure ====="
    exit $ret
else
    echo -n "===== $module build successfully! ====="
fi

rm -rf $output
mkdir -p $output/{bin,conf,log,data}

# fullfil output dir, output is the app to be deployed
(
    #cp -rf control.sh ${output}/ &&
    mv ${module} ${output}/bin/ &&
    cp -rf bin/servicectl.py ${output}/bin/ &&
    echo "===== Generate output ok ====="
) || { echo "===== Generate output failure ====="; exit 2; } 

