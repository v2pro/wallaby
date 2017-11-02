#!/bin/bash
# go默认使用go1.4
#
# 如果需要使用go1.5，请打开下面两行注释：
# export GOROOT=/usr/local/go1.5.1
# export PATH=$GOROOT/bin:$PATH
#
# 如果需要使用go1.6 ，请打开下面两行注释：
# export GOROOT=/usr/local/go1.6.2
# export PATH=$GOROOT/bin:$PATH
#
# 如果需要使用go1.7 ，请打开下面两行注释：
# export GOROOT=/usr/local/go1.7.5
# export PATH=$GOROOT/bin:$PATH
#
#

module=wallaby
output="output"

# change build timestamp 
sed -i -e '/ProxyBuildTimestamp = / s/= .*/= '`date +%s`'/' config/config.go


go build -o $module cmd/proxy/main.go    #编译目标文件
ret=$?
if [ $ret -ne 0 ];then
    echo "===== $module build failure ====="
    exit $ret
else
    echo -n "===== $module build successfully! ====="
fi

rm -rf $output
mkdir -p $output/{bin,conf,log,data}

# 填充output目录, output的内容即为待部署内容
(
    #cp -rf control.sh ${output}/ &&             # 拷贝部署脚本control.sh至output目录
    mv ${module} ${output}/bin/ &&              # 移动需要部署的文件到output目录下
    echo -e "===== Generate output ok ====="
) || { echo -e "===== Generate output failure ====="; exit 2; } # 填充output目录失败后, 退出码为 非0

