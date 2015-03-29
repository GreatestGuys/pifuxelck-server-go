#!/bin/bash

Color_Off='\e[0m'
BGreen='\e[1;32m'
BYellow='\e[1;33m'
BWhite='\e[1;37m'

if [[ $# -ne 1 ]] ; then
  echo 'usage: ./deploy sauce-user-name'
  exit 1
fi

USER=$1
SSHCMD="ssh $USER@everythingissauce.com"

echo -e "${BWhite}[+] ${BYellow}Removing existing consoles...${Color_Off}"
$SSHCMD "sudo su -c 'rm -Rf /srv/prometheus/consoles/*' prometheus"
$SSHCMD "sudo su -c 'rm -Rf /srv/prometheus/console_libraries/*' prometheus"

echo -e "${BWhite}[+] ${BYellow}Deploying new consoles...${Color_Off}"
scp \
  prometheus/consoles/* \
  $USER@everythingissauce.com:/srv/prometheus/consoles/
scp \
  prometheus/console_libraries/* \
  $USER@everythingissauce.com:/srv/prometheus/console_libraries/

echo -e "${BWhite}[+] ${BGreen}All done!${Color_Off}"

exit 0
