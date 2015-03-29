#!/bin/bash

Color_Off='\e[0m'
BGreen='\e[1;32m'
BYellow='\e[1;33m'
BWhite='\e[1;37m'

if [[ $# -ne 2 ]] ; then
  echo 'usage: ./deploy sauce-user-name my-sql-password'
  exit 1
fi

USER=$1
PASS=$2

TMP_FILE=`mktemp /tmp/pifuxelck.XXXXXXX.sh`
cat > $TMP_FILE << EOF
echo ""
echo -e "${BYellow}Launching new instance...${Color_Off}"
cd /srv/pifuxelck/
/usr/bin/nohup ./pifuxelck-server-go \
  --port 3002 \
  --mysql-host db.everythingissauce.com \
  --mysql-port 3306 \
  --mysql-user pifuxelck \
  --mysql-password $PASS \
  --mysql-db pifuxelck 2>stderr > pifuxelck-server-go.log < /dev/null &
exit
EOF

echo ""
echo -e "${BWhite}[+] ${BYellow}Building the executable...${Color_Off}"
rm pifuxelck-server-go
go build github.com/GreatestGuys/pifuxelck-server-go

echo ""
echo -e "${BWhite}[+] ${BYellow}Killing existing server instances...${Color_Off}"
ssh \
  $USER@everythingissauce.com \
  "sudo su -c 'killall pifuxelck-server-go' pifuxelck"

echo ""
echo -e "${BWhite}[+] ${BYellow}Deploying the executable...${Color_Off}"
scp \
  pifuxelck-server-go \
  $USER@everythingissauce.com:/srv/pifuxelck/pifuxelck-server-go

echo ""
echo -e "${BWhite}[+] ${BYellow}Deploying start up script...${Color_Off}"
scp $TMP_FILE $USER@everythingissauce.com:${TMP_FILE}

echo ""
echo -e "${BWhite}[+] ${BYellow}Running start up script...${Color_Off}"
ssh \
  $USER@everythingissauce.com \
  "chmod a+rxw $TMP_FILE && sudo su -c $TMP_FILE pifuxelck && rm $TMP_FILE"

echo ""
echo -e "${BWhite}[+] ${BGreen}All done!${Color_Off}"

rm $TMP_FILE
exit 0
