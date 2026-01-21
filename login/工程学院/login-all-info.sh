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

# 获取脚本所在目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
LOG_FILE="${SCRIPT_DIR}/login-all.log" # 更改日志文件名为 login-all.log
COUNTER_FILE="${SCRIPT_DIR}/counter-all.txt" # 更改计数器文件名为 counter-all.txt
RESPONSE_LOG_FILE="${SCRIPT_DIR}/response_log.txt" # 保持原有的响应日志文件

# 计数器函数
get_and_update_counter() {
    local count=1

    # 读取当前计数
    if [[ -f "${COUNTER_FILE}" ]]; then
        count=$(cat "${COUNTER_FILE}" 2>/dev/null || echo 1)
        # 确保count是数字
        if ! [[ "$count" =~ ^[0-9]+$ ]]; then
            count=1
        fi
    fi

    # 更新计数器文件
    echo $((count + 1)) > "${COUNTER_FILE}"

    # 返回当前执行次数
    echo "$count"
}

# 日志函数
log_message() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] $1" >> "${LOG_FILE}"
}

# 清理日志函数（每月第一天清理）
cleanup_logs() {
    local current_day=$(date '+%d')
    if [[ "$current_day" == "01" ]]; then
        if [[ -f "${LOG_FILE}" ]]; then
            # 保留最后100行，删除其余内容
            tail -n 100 "${LOG_FILE}" > "${LOG_FILE}.tmp" && mv "${LOG_FILE}.tmp" "${LOG_FILE}"
            log_message "INFO: 主日志文件已清理，保留最近100条记录"
        fi
        if [[ -f "${RESPONSE_LOG_FILE}" ]]; then
            # 清空响应日志文件
            > "${RESPONSE_LOG_FILE}"
            log_message "INFO: 响应日志文件已清空"
        fi
    fi
}

# 网络连接测试函数
test_internet_connection() {
    local interface=$1
    log_message "INFO: 开始通过网卡 ${interface} 进行网络连接测试"

    # 使用ping测试，超时5秒，发送3个包，指定网卡
    if ping -I "${interface}" -c 3 -W 5 8.8.8.8 > /dev/null 2>&1; then
        log_message "INFO: 网卡 ${interface} 网络连接正常"
        return 0
    else
        log_message "INFO: 网卡 ${interface} 网络连接异常"
        return 1
    fi
}

# 登录URL
# 配置参数
URL="http://172.17.10.100/eportal/InterFace.do"
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

# 主执行逻辑
main() {
    # 获取执行次数
    local execution_count=$(get_and_update_counter)
    log_message "INFO: ===== 脚本开始执行（第 ${execution_count} 次） ======"

    # 初始化计数器
    local SUCCESS_COUNT=0
    local FAILURE_COUNT=0

    # 每月清理日志
    cleanup_logs

    # 创建或清空响应数据文件
    echo "登录响应数据记录" > "${RESPONSE_LOG_FILE}"
    echo "===========================================" >> "${RESPONSE_LOG_FILE}"
    echo "运行时间: $(date)" >> "${RESPONSE_LOG_FILE}"
    echo "===========================================" >> "${RESPONSE_LOG_FILE}"

    # 遍历账号进行登录
    for account in "${!ACCOUNTS[@]}"; do
        # 解析网卡名和用户名
        IFS=',' read -r interface name <<< "${account}"

        log_message "INFO: 处理账号 ${name} (网卡: ${interface})"
        echo "正在检查 ${name} (${interface}) 的网络连接... "

        # 测试当前网卡网络连接
        if test_internet_connection "${interface}"; then
            log_message "INFO: 网卡 ${interface} 网络连接正常，无需认证。"
            print_color "green" "✓ ${name} (${interface}) 网络连接正常，跳过认证。"
        else
            log_message "INFO: 网卡 ${interface} 网络连接异常，尝试进行认证。"
            # 解析用户ID和密码
            IFS=':' read -r userid password <<< "${ACCOUNTS[$account]}"

            # 构建POST数据
            POST_DATA="method=login&userId=${userid}&password=${password}&service=${SERVICE}&queryString=${QUERY_STRING}&operatorPwd=&operatorUserId=&validcode=&passwordEncrypt=true"

            log_message "INFO: 正在登录 ${name} (${interface})..."
            echo "正在登录 ${name} (${interface})..." # 保持控制台输出

            # 执行curl命令并捕获响应
            response=$(curl --interface "${interface}" \
                 -X POST \
                 -H "Content-Type: application/x-www-form-urlencoded; charset=UTF-8" \
                 -H "User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) Firefox/132.0" \
                 -d "${POST_DATA}" \
                 "${URL}" 2>/dev/null)

            local curl_exit_code=$?

            if [[ $curl_exit_code -eq 0 ]]; then
                log_message "INFO: ${name} (${interface}) 认证请求发送成功"
                log_message "DEBUG: ${name} (${interface}) 服务器响应: ${response}"

                # 检查响应中是否包含success
                if echo "${response}" | grep -q "success\|成功"; then
                    print_color "green" "✓ ${name} (${interface}) 登录成功"
                    log_message "SUCCESS: ${name} (${interface}) 登录认证成功"
                    status="成功"
                    SUCCESS_COUNT=$((SUCCESS_COUNT + 1))
                else
                    # 解析错误信息
                    error_msg="未知错误"
                    if echo "${response}" | grep -q "message"; then
                        error_msg=$(echo "${response}" | grep -o '"message":"[^"]*"' | cut -d':' -f2- | tr -d '"')
                    elif echo "${response}" | grep -q "error"; then
                        error_msg=$(echo "${response}" | grep -o '"error":"[^"]*"' | cut -d':' -f2- | tr -d '"')
                    elif echo "${response}" | grep -q "用户被踢下线"; then
                        error_msg="用户被踢下线"
                    elif echo "${response}" | grep -q "用户已在线"; then
                        error_msg="用户已在线"
                    elif echo "${response}" | grep -q "密码错误"; then
                        error_msg="密码错误"
                    elif echo "${response}" | grep -q "用户不存在"; then
                        error_msg="用户不存在"
                    elif [[ -z "${response}" ]]; then
                        error_msg="网络连接失败或服务器无响应"
                    fi
                    print_color "red" "✗ ${name} (${interface}) 登录失败: $error_msg"
                    log_message "WARNING: ${name} (${interface}) 登录认证状态未知或失败: ${error_msg}"
                    status="失败: $error_msg"
                    FAILURE_COUNT=$((FAILURE_COUNT + 1))
                fi
            else
                print_color "red" "✗ ${name} (${interface}) 认证请求失败，curl退出代码: ${curl_exit_code}"
                log_message "ERROR: ${name} (${interface}) 认证请求失败，curl退出代码: ${curl_exit_code}"
                log_message "ERROR: ${name} (${interface}) 错误信息: ${response}"
                status="失败: curl退出代码 ${curl_exit_code}"
                FAILURE_COUNT=$((FAILURE_COUNT + 1))
            fi

            # 显示完整响应数据
            print_color "yellow" "响应数据: $response"

            # 记录到响应数据文件
            echo -e "\n[$name ($interface)] - 状态: $status" >> "${RESPONSE_LOG_FILE}"
            echo "-------------------------------------------" >> "${RESPONSE_LOG_FILE}"
            echo "时间: $(date '+%Y-%m-%d %H:%M:%S')" >> "${RESPONSE_LOG_FILE}"
            echo "完整响应: $response" >> "${RESPONSE_LOG_FILE}"
            echo "-------------------------------------------" >> "${RESPONSE_LOG_FILE}"
        fi # End of test_internet_connection for current interface
    done
    log_message "INFO: 所有账号登录尝试完成。"
    print_color "blue" "响应数据已保存到 ${RESPONSE_LOG_FILE} 文件"
    log_message "INFO: 所有账号登录尝试完成。"
    print_color "blue" "响应数据已保存到 ${RESPONSE_LOG_FILE} 文件"

    log_message "INFO: 登录尝试总结: 成功 ${SUCCESS_COUNT} 个, 失败 ${FAILURE_COUNT} 个。"
    print_color "blue" "登录尝试总结: 成功 ${SUCCESS_COUNT} 个, 失败 ${FAILURE_COUNT} 个。"

    log_message "INFO: ========== 脚本执行完毕 =========="
}

# 执行主函数
main
