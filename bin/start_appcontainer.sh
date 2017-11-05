#!/usr/bin/env bash
APPNAME=appcontainer
if [ -d $APPNAME ];then
    rm -rf $APPNAME
fi
mkdir $APPNAME
cd $APPNAME
mkdir {bin,lib,lib64,usr,var,run,proc,data,tmp,etc}
ENV=`uname`
if [[ "$ENV" == "Linux" ]];then
    echo "environment: Linux"
#    ldd /bin/bash
    cp /lib64/ld-linux-x86-64.so.2 lib64/
    cp /lib64/libc-2.17.so lib64/libc.so.6
    cp /lib64/libdl-2.17.so lib64/libdl.so.2
    cp /lib64/libtinfo.so.5.9 lib64/libtinfo.so.5
else
    echo "environment: Darwin"
#    otool -L /bin/bash
    mkdir -p usr/lib
    mkdir usr/lib/system
    cp /usr/lib/libncurses.5.4.dylib usr/lib
    cp /usr/lib/libSystem.B.dylib usr/lib
    cp /usr/lib/dyld usr/lib
    cp /usr/lib/libobjc.A.dylib usr/lib
    cp /usr/lib/libc++abi.dylib usr/lib
    cp /usr/lib/libc++.1.dylib usr/lib

    cp /usr/lib/system/libmathCommon.A.dylib usr/lib/system
    cp /usr/lib/system/libcache.dylib usr/lib/system
    cp /usr/lib/system/libcommonCrypto.dylib usr/lib/system
    cp /usr/lib/system/libcompiler_rt.dylib usr/lib/system
    cp /usr/lib/system/libcopyfile.dylib usr/lib/system
    cp /usr/lib/system/libcorecrypto.dylib usr/lib/system
    cp /usr/lib/system/libdispatch.dylib usr/lib/system
    cp /usr/lib/system/libdyld.dylib usr/lib/system
    cp /usr/lib/system/liblaunch.dylib usr/lib/system
    cp /usr/lib/system/libmacho.dylib usr/lib/system
    cp /usr/lib/system/libquarantine.dylib usr/lib/system
    cp /usr/lib/system/libremovefile.dylib usr/lib/system
    cp /usr/lib/system/libsystem_asl.dylib usr/lib/system
    cp /usr/lib/system/libsystem_blocks.dylib usr/lib/system
    cp /usr/lib/system/libsystem_c.dylib usr/lib/system
    cp /usr/lib/system/libsystem_configuration.dylib usr/lib/system
    cp /usr/lib/system/libsystem_coreservices.dylib usr/lib/system
    cp /usr/lib/system/libsystem_coretls.dylib usr/lib/system
    cp /usr/lib/system/libsystem_dnssd.dylib usr/lib/system
    cp /usr/lib/system/libsystem_info.dylib usr/lib/system
    cp /usr/lib/system/libsystem_kernel.dylib usr/lib/system
    cp /usr/lib/system/libsystem_m.dylib usr/lib/system
    cp /usr/lib/system/libsystem_malloc.dylib usr/lib/system
    cp /usr/lib/system/libsystem_network.dylib usr/lib/system
    cp /usr/lib/system/libsystem_networkextension.dylib usr/lib/system
    cp /usr/lib/system/libsystem_notify.dylib usr/lib/system
    cp /usr/lib/system/libsystem_platform.dylib usr/lib/system
    cp /usr/lib/system/libsystem_pthread.dylib usr/lib/system
    cp /usr/lib/system/libsystem_sandbox.dylib usr/lib/system
    cp /usr/lib/system/libsystem_secinit.dylib usr/lib/system
    cp /usr/lib/system/libsystem_symptoms.dylib usr/lib/system
    cp /usr/lib/system/libsystem_trace.dylib usr/lib/system
    cp /usr/lib/system/libunwind.dylib usr/lib/system
    cp /usr/lib/system/libxpc.dylib usr/lib/system
    cp /usr/lib/system/libkeymgr.dylib usr/lib/system
fi
cp /bin/bash bin
cp /bin/ls bin
echo "$APPNAME is ready."
chroot . /bin/bash



