#!/bin/bash
export PATH="/usr/sbin:/usr/bin:/sbin:/bin:$PATH"

# 颜色输出
print_color() {
    case $1 in
        "green") echo -e "\033[32m$2\033[0m" ;;
        "red") echo -e "\033[31m$2\033[0m" ;;
        "yellow") echo -e "\033[33m$2\033[0m" ;;
        "blue") echo -e "\033[34m$2\033[0m" ;;
    esac
}

# 基本配置（可按需覆盖）
HOST="${HOST:-auth.cqtbi.edu.cn}"
WLANACIP="${WLANACIP:-10.34.5.3}"
WLANACNAME="${WLANACNAME:-ZAX-BRAS-HC}"
VLAN="${VLAN:-0}"
URL_REDIRECT="${URL_REDIRECT:-http://www.msftconnecttest.com/redirect}"
CONNECT_TIMEOUT="${CONNECT_TIMEOUT:-5}"
MAX_TIME="${MAX_TIME:-15}"
# 可选：强制指定认证服务器IP（绕过DNS），例如：export AUTH_IP_OVERRIDE=1.2.3.4
AUTH_IP_OVERRIDE="10.34.2.9" 

# 账号信息（interface,name -> userId:passwd）
# 注意键使用 "接口,名称"（不带额外引号），例如 ["wlan0,张三"]
declare -A ACCOUNTS=(
    # ["veth0,SSJ"]="X250223201150:@Ssddffqxc547"
    # ["veth2,SSJ"]="X250223201150:@Ssddffqxc547"
    ["eth1,SSJ"]="X250223201150:@Ssddffqxc547"
)

# 计数器文件
COUNTER_FILE="login_counter.txt"
exec_count=0
if [[ -f "$COUNTER_FILE" ]]; then
    exec_count=$(cat "$COUNTER_FILE" 2>/dev/null || echo 0)
fi
[[ "$exec_count" =~ ^[0-9]+$ ]] || exec_count=0
exec_count=$((exec_count+1))
echo "$exec_count" > "$COUNTER_FILE"

# 日志文件（简化版）
LOG_FILE="response_log.txt"
echo "执行次数: $exec_count" > "$LOG_FILE"
echo "时间: $(date '+%Y-%m-%d %H:%M:%S')" >> "$LOG_FILE"

require_cmd() {
    command -v "$1" >/dev/null 2>&1 || {
        print_color red "缺少命令: $1；请安装后重试"
        exit 1
    }
}

has_cmd() { command -v "$1" >/dev/null 2>&1; }

require_cmd ip
require_cmd curl
has_cmd ping && PING_AVAILABLE=1 || PING_AVAILABLE=0
has_cmd getent && GETENT_AVAILABLE=1 || GETENT_AVAILABLE=0
has_cmd nslookup && NSLOOKUP_AVAILABLE=1 || NSLOOKUP_AVAILABLE=0
has_cmd host && HOST_AVAILABLE=1 || HOST_AVAILABLE=0

# 接口是否存在
has_iface() {
    local iface="$1"
    [[ -z "$iface" ]] && return 1
    ip link show "$iface" >/dev/null 2>&1
}

# 选择接口：优先使用账户指定接口，否则检测默认路由接口
pick_iface() {
    local prefer="$1"
    if has_iface "$prefer"; then
        echo "$prefer"
        return 0
    fi
    local def
    def=$(ip route show default | awk '/default/{for(i=1;i<=NF;i++){if($i=="dev"){print $(i+1); exit}}}')
    echo "$def"
}

# 获取IP、MAC、网关
get_ipv4() { ip -4 addr show dev "$1" | awk '/inet /{print $2}' | cut -d/ -f1 | head -n1; }
get_mac() { cat "/sys/class/net/$1/address" 2>/dev/null | tr '[:upper:]' '[:lower:]'; }
get_gw() { ip route show default | awk '/default/{print $3; exit}'; }

# 提取系统DNS服务器列表（仅展示简要信息）
get_dns_servers() {
    local servers=()
    if [[ -f /etc/resolv.conf ]]; then
        while read -r a b; do
            if [[ "$a" == "nameserver" && "$b" =~ ^[0-9] ]]; then
                servers+=("$b")
            fi
        done < /etc/resolv.conf
    fi
    servers+=("10.34.5.253" "114.114.114.114" "8.8.8.8" "1.1.1.1")
    echo "${servers[@]}"
}

# 解析域名为IP（优先系统解析，其次指定DNS服务器）
resolve_host() {
    local host="$1"; local ip=""; local dnslist; dnslist=$(get_dns_servers)

    if [[ -n "$AUTH_IP_OVERRIDE" ]]; then
        echo "$AUTH_IP_OVERRIDE"; return 0
    fi
    if [[ $GETENT_AVAILABLE -eq 1 ]]; then
        ip=$(getent hosts "$host" | awk '{print $1}' | head -n1)
    fi
    if [[ -z "$ip" && $NSLOOKUP_AVAILABLE -eq 1 ]]; then
        for dns in $dnslist; do
            out=$(nslookup "$host" "$dns" 2>/dev/null)
            ip=$(echo "$out" | awk '/Address [0-9]+:/{print $3; exit}')
            [[ -z "$ip" ]] && ip=$(echo "$out" | awk '/^Address: /{print $2; exit}')
            if [[ -n "$ip" ]]; then break; fi
        done
    fi
    if [[ -z "$ip" && $HOST_AVAILABLE -eq 1 ]]; then
        out=$(host "$host" 2>/dev/null)
        ip=$(echo "$out" | awk '/has address/{print $4; exit}')
    fi
    echo "$ip"
}

# 连通性检测：通过指定接口访问测速站（默认 https://test.ustc.edu.cn/）
# 仅以 HTTP 2xx/3xx 判定“已联网”；ping 仅记录详情不作为依据
check_online() {
    local iface="$1"
    local ipaddr; ipaddr=$(get_ipv4 "$iface")
    local url="${CHECK_URL:-https://test.ustc.edu.cn/}"
    # 使用 USTC 测速站进行 HTTPS 连通性检测，绑定到接口或源IP，忽略响应体
    local code=""; local curl_extra=()
    [[ "${CURL_INSECURE}" == "1" ]] && curl_extra+=( -k )
    code=$(curl --interface "${ipaddr:-$iface}" -sS -o /dev/null --connect-timeout 3 --max-time 5 ${curl_extra[@]} "$url" -w "%{http_code}" 2>/dev/null)
    if [[ "$code" =~ ^[0-9]+$ ]] && [[ $code -ge 200 && $code -lt 400 ]]; then
        ONLINE_CHECK_DETAIL="curl:${url} code=${code} via=${ipaddr:-$iface}${CURL_INSECURE:+ (insecure)}"
        return 0
    fi
    # 记录 ping 结果，但不据此认定为“已联网”（避免门户前 ping 可达导致误判）
    if [[ $PING_AVAILABLE -eq 1 ]]; then
        if [[ -n "$ipaddr" ]]; then
            if ping -I "$ipaddr" -c 1 -W 1 "test.ustc.edu.cn" >/dev/null 2>&1; then
                ONLINE_CHECK_DETAIL="curl:${url} code=${code:-NA}; ping ok via ${ipaddr}"
            fi
        else
            if ping -I "$iface" -c 1 -W 1 "test.ustc.edu.cn" >/dev/null 2>&1; then
                ONLINE_CHECK_DETAIL="curl:${url} code=${code:-NA}; ping ok via ${iface}"
            fi
        fi
    fi
    [[ -z "$ONLINE_CHECK_DETAIL" ]] && ONLINE_CHECK_DETAIL="unreachable via ${iface} ip=${ipaddr} to ${url} (curl:${code:-NA})"
    return 1
}

# curl通用选项（不使用 --compressed）
CURL_COMMON_OPTS=(-sS --connect-timeout "$CONNECT_TIMEOUT" --max-time "$MAX_TIME")
[[ "${CURL_INSECURE}" == "1" ]] && CURL_COMMON_OPTS+=( -k )

for key in "${!ACCOUNTS[@]}"; do
    IFS=',' read -r interface name <<< "${key}"
    IFS=':' read -r userid passwd <<< "${ACCOUNTS[$key]}"
    [[ -z "$name" ]] && name="默认"

    # 选择接口并拉取网络信息
    iface=$(pick_iface "$interface")
    if [[ -z "$iface" ]]; then
        print_color red "✗ ${name} 未检测到接口"
        echo "账号: ${name} | 接口: - | 状态: 失败(无接口)" >> "$LOG_FILE"
        continue
    fi

    ipv4=$(get_ipv4 "$iface")
    mac=$(get_mac "$iface")
    gw=$(get_gw)

    if [[ -z "$ipv4" || -z "$mac" ]]; then
        print_color red "✗ ${name} (${iface}) 无法获取IP或MAC"
        echo "账号: ${name} | 接口: ${iface} | 网络: IP=${ipv4} MAC=${mac} GW=${gw} | 状态: 失败(缺少IP/MAC)" >> "$LOG_FILE"
        continue
    fi

    # 如已联网，则跳过登录
    if check_online "$iface"; then
        print_color green "✓ ${name} (${iface}) 已联网，无需登录 | 次数: $exec_count"
        echo "账号: ${name} | 接口: ${iface} | 次数: ${exec_count}" >> "$LOG_FILE"
        echo "时间: $(date '+%Y-%m-%d %H:%M:%S')" >> "$LOG_FILE"
        echo "网络: IP=${ipv4} MAC=${mac} GW=${gw}" >> "$LOG_FILE"
        echo "状态: 已联网(无需登录) | 连通性: ${ONLINE_CHECK_DETAIL}" >> "$LOG_FILE"
        echo "-------------------------------------------" >> "$LOG_FILE"
        continue
    fi

    # 记录未联网的探测细节，便于排查
    echo "连通性: ${ONLINE_CHECK_DETAIL}" >> "$LOG_FILE"

    print_color blue "登录 ${name} (${iface})"
    dns_list=$(get_dns_servers)
    AUTH_IP=$(resolve_host "$HOST")
    if [[ -n "$AUTH_IP" ]]; then
        RESOLVE_OPT=(--resolve "${HOST}:443:${AUTH_IP}")
    else
        RESOLVE_OPT=()
    fi

    FULL_URL="https://${HOST}/webauth.do?wlanacip=${WLANACIP}&wlanacname=${WLANACNAME}&wlanuserip=${ipv4}&mac=${mac}&vlan=${VLAN}&url=${URL_REDIRECT}"

    cookie_file="./cookie_${name}.txt"
    header_file=$(mktemp)
    body_file=$(mktemp)

    # 预先GET建立会话Cookie（不输出详尽信息）
    curl --interface "$iface" "${CURL_COMMON_OPTS[@]}" "${RESOLVE_OPT[@]}" -c "$cookie_file" -b "$cookie_file" -A "Mozilla/5.0" -D "$header_file" -o "$body_file" "$FULL_URL" >/dev/null 2>&1

    # 构建POST（字段与抓包一致）
    common_headers=(
        -H "Content-Type: application/x-www-form-urlencoded; charset=UTF-8"
        -H "Origin: https://${HOST}"
        -H "Referer: ${FULL_URL}"
        -A "Mozilla/5.0 (X11; Linux x86_64) Firefox/132.0"
        -c "$cookie_file" -b "$cookie_file"
    )

    data_fields=(
        --data-urlencode "scheme=https"
        --data-urlencode "serverIp=${HOST}:443"
        --data-urlencode "hostIp=http://127.0.0.1:8445/"
        --data-urlencode "loginType="
        --data-urlencode "auth_type=0"
        --data-urlencode "isBindMac1=0"
        --data-urlencode "pageid=282"
        --data-urlencode "templatetype=1"
        --data-urlencode "listbindmac=0"
        --data-urlencode "recordmac=1"
        --data-urlencode "isRemind=1"
        --data-urlencode "loginTimes="
        --data-urlencode "groupId="
        --data-urlencode "distoken="
        --data-urlencode "echostr="
        --data-urlencode "url=${URL_REDIRECT}"
        --data-urlencode "isautoauth="
        --data-urlencode "mobile="
        --data-urlencode "notice_pic_loop2=/portal/uploads/pc/demo2/images/bj.png"
        --data-urlencode "notice_pic_loop1=/portal/uploads/pc/demo2/images/logo.png"
        --data-urlencode "userId=${userid}"
        --data-urlencode "passwd=${passwd}"
        --data-urlencode "remInfo=on"
        --data-urlencode "isBindMac=bindmac"
    )

    http_code=$(curl --interface "$iface" "${CURL_COMMON_OPTS[@]}" "${RESOLVE_OPT[@]}" -X POST "${common_headers[@]}" "${data_fields[@]}" -D "$header_file" -o "$body_file" "$FULL_URL" -w "%{http_code}")
    response=$(cat "$body_file" 2>/dev/null)

    # 结果判定（只输出关键信息）
    status="失败"
    error_msg=""
    if [[ "$http_code" == "000" ]]; then
        error_msg="网络/TLS失败"
    elif [[ "$http_code" != "200" ]]; then
        error_msg="HTTP $http_code"
    elif [[ $response == *"Welcome to Drcom System"* ]] || [[ $response == *"LOGINSUCC"* ]]; then
        status="成功"
    else
        if [[ $response == *"密码错误"* ]]; then
            error_msg="密码错误"
        elif [[ $response == *"用户不存在"* ]]; then
            error_msg="用户不存在"
        elif [[ $response == *"用户已在线"* ]]; then
            error_msg="用户已在线"
        elif [[ -z "$response" ]]; then
            error_msg="服务器无响应"
        else
            error_msg="认证失败"
        fi
    fi

    if [[ "$status" == "成功" ]]; then
        print_color green "✓ ${name} 登录成功 | 次数: $exec_count"
    else
        print_color red "✗ ${name} 登录失败: ${error_msg} | 次数: $exec_count"
    fi

    # 简化日志：时间、计数器、账号、接口、网络信息、状态、错误
    echo "账号: ${name} | 接口: ${iface} | 次数: ${exec_count}" >> "$LOG_FILE"
    echo "时间: $(date '+%Y-%m-%d %H:%M:%S')" >> "$LOG_FILE"
    echo "网络: IP=${ipv4} MAC=${mac} GW=${gw}" >> "$LOG_FILE"
    echo "域名解析: ${HOST} -> ${AUTH_IP:-解析失败}" >> "$LOG_FILE"
    echo "状态: ${status}${error_msg:+ | 说明: ${error_msg}}" >> "$LOG_FILE"
    echo "-------------------------------------------" >> "$LOG_FILE"

    rm -f "$cookie_file" "$header_file" "$body_file" 2>/dev/null

done

print_color blue "完成 | 次数: $exec_count"
print_color blue "日志: $LOG_FILE"