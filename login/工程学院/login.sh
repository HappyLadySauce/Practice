#!/bin/bash

URL="http://172.17.10.100/eportal/InterFace.do"

USER_ID="2233229042"
PASSWORD="14f71ce05967a3c3011a7f304e725713d150d894e81d34c6a6675323b7ca76f1b7331c08131e0fe9ed1e101039111cd62c3f5328d967f247283b47e5f70da22b66ce755cdd58cd0cb670ccba5b6172753f12f99c40f55acd1b7177970cadbbd8d3a9615508d968c620ef6f4d6a734943bea39a301ac6ad2df2985a401a2a9434"
SERVICE="教学区免费上网"
QUERY_STRING="wlanuserip%3D172.18.121.5%26wlanacname%3DNAS%26ssid%3DRuijie%26nasip%3D172.17.10.10%26mac%3D000000000000%26t%3Dwireless-v2-plain%26url%3Dhttp%3A%252F%252Fwww.baidu.com%252F"

POST_DATA="method=login&userId=${USER_ID}&password=${PASSWORD}&service=${SERVICE}&queryString=${QUERY_STRING}&operatorPwd=&operatorUserId=&validcode=&passwordEncrypt=true"

curl --interface vmbr0 \
     -X POST \
     -H "Content-Type: application/x-www-form-urlencoded; charset=UTF-8" \
     -H "User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) Firefox/132.0" \
     -d "${POST_DATA}" \
     "${URL}"