#!/bin/bash

# 定义颜色输出函数
print_color() {
    case $1 in
        "green") echo -e "\033[32m$2\033[0m" ;;
        "red") echo -e "\033[31m$2\033[0m" ;;
        "yellow") echo -e "\033[33m$2\033[0m" ;;
        "blue") echo -e "\033[34m$2\033[0m" ;;
    esac
}

# 登录URL
URL="http://172.17.10.100/eportal/InterFace.do"

# 通用参数
SERVICE="教学区免费上网"
QUERY_STRING="wlanuserip%3D172.18.121.5%26wlanacname%3DNAS%26ssid%3DRuijie%26nasip%3D172.17.10.10%26mac%3D000000000000%26t%3Dwireless-v2-plain%26url%3Dhttp%3A%252F%252Fwww.baidu.com%252F"

# 定义账号信息数组
declare -A ACCOUNTS=(
    # ["macvlan2,YD"]="2233229224:2254680660a5d0d0b6f61e77a8ade8d26da49a34553048106e34b3daa01a3ca95b0c11e54a7df702631564c64d7a06771f9c1955afe93875f70c8ce39f600f78f3b33254bec7dbb3dc9874e14eac83dfcbdb5a3ed66b68c7405f1923ab69f8601f61a1fe1f695561f67d2b02fe04844e40e7f9041c0fd0208d001f29ee33f5f7"
    ["macvlan2,ZH"]="2233229098:8530f1b7533a342f21502b64bc7b01aaa5b2e8b4a04ec004ba67b02106976ee83bd197b678311268907c2b6bbe3e42ed6a5ed5b7eac68c87de1a50ba78901a90bac69e5383cb80bfa639a303bd02ebd1a7be488d2946271796f414fa7ddb0995f4bd3364d02bebe62efff3ad6151110ec219325a04dbd9ef36756ce10aad06f4"
    # ["macvlan4,ZZR"]="2333212031:1243bbf41ef91e2190fa42199c7ff623a3ea0bdf2d789d122dc2582215b0dc8cae2173f1ff1d8c5e493325efb633500d01b4343dd05f8ecd21a9afac63e1a57373bbcc809548a220d153fa093d7211bb56f9aad16cc98350019908303fca2be2a03ae76988d669b7717991a75ac514e1f59bdead4749e398a3a46232b9124926"
    ["macvlan3,YXS"]="2233229042:14f71ce05967a3c3011a7f304e725713d150d894e81d34c6a6675323b7ca76f1b7331c08131e0fe9ed1e101039111cd62c3f5328d967f247283b47e5f70da22b66ce755cdd58cd0cb670ccba5b6172753f12f99c40f55acd1b7177970cadbbd8d3a9615508d968c620ef6f4d6a734943bea39a301ac6ad2df2985a401a2a9434"
    ["macvlan1,QX"]="2433229383:8b8f5465625bbf238b848ae98290bf46b76d3af92571845d244a67c006c25fd3a2dd68fd8c97d7ad934d2af2a000bc95eb21297503ff5b7a1130162dfae5481995ae15575a83798717a62407838d79235b74068bd98a44aec778b2c0887dbd8a78676dcede165b1da81183a767daef3e60d2d6a99487937e996a6b4cc7065d69"
)

# 创建或清空响应数据文件
echo "登录响应数据记录" > response_log.txt
echo "===========================================" >> response_log.txt
echo "运行时间: $(date)" >> response_log.txt
echo "===========================================" >> response_log.txt

# 遍历账号进行登录
for account in "${!ACCOUNTS[@]}"; do
    # 解析网卡名和用户名
    IFS=',' read -r interface name <<< "${account}"

    # 解析用户ID和密码
    IFS=':' read -r userid password <<< "${ACCOUNTS[$account]}"

    # 构建POST数据
    POST_DATA="method=login&userId=${userid}&password=${password}&service=${SERVICE}&queryString=${QUERY_STRING}&operatorPwd=&operatorUserId=&validcode=&passwordEncrypt=true"

    echo "正在登录 ${name} (${interface})..."

    # 执行curl命令并捕获响应
    response=$(curl --interface "${interface}" \
         -X POST \
         -H "Content-Type: application/x-www-form-urlencoded; charset=UTF-8" \
         -H "User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) Firefox/132.0" \
         -d "${POST_DATA}" \
         "${URL}" 2>/dev/null)

    # 检查响应中是否包含success
    if [[ $response == *"success"* ]]; then
        print_color "green" "✓ ${name} (${interface}) 登录成功"
        status="成功"
    else
        # 解析错误信息
        error_msg="未知错误"
        if [[ $response == *"message"* ]]; then
            # 提取 message 字段内容 - 基本格式为 "message":"错误信息"
            error_msg=$(echo "$response" | grep -o '"message":"[^"]*"' | cut -d':' -f2- | tr -d '"')
        elif [[ $response == *"error"* ]]; then
            # 如果没有message字段，尝试提取error字段
            error_msg=$(echo "$response" | grep -o '"error":"[^"]*"' | cut -d':' -f2- | tr -d '"')
        elif [[ $response == *"用户被踢下线"* ]]; then
            error_msg="用户被踢下线"
        elif [[ $response == *"用户已在线"* ]]; then
            error_msg="用户已在线"
        elif [[ $response == *"密码错误"* ]]; then
            error_msg="密码错误"
        elif [[ $response == *"用户不存在"* ]]; then
            error_msg="用户不存在"
        elif [[ -z "$response" ]]; then
            error_msg="网络连接失败或服务器无响应"
        fi

        print_color "red" "✗ ${name} (${interface}) 登录失败: $error_msg"
        status="失败: $error_msg"
    fi
    
    # 显示完整响应数据
    print_color "yellow" "响应数据: $response"
    
    # 记录到响应数据文件
    echo -e "\n[$name ($interface)] - 状态: $status" >> response_log.txt
    echo "-------------------------------------------" >> response_log.txt
    echo "时间: $(date '+%Y-%m-%d %H:%M:%S')" >> response_log.txt
    echo "完整响应: $response" >> response_log.txt
    echo "-------------------------------------------" >> response_log.txt
done

echo -e "\n所有账号登录尝试完成"
print_color "blue" "响应数据已保存到 response_log.txt 文件"
