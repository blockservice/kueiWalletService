#!/usr/bin/env bash

# Absolute path to this script. /home/user/bin/foo.sh
SCRIPT=$(readlink -f $0)
# Absolute path this script is in. /home/user/bin
SCRIPTPATH=`dirname ${SCRIPT}`

if [ ! -d ./bin ];then
  mkdir ./bin
fi

go build -o ${SCRIPTPATH}/bin/ews github.com/ChungkueiBlock/kueiWalletService/cmd/ews