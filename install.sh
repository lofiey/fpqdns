#!/usr/bin/env bash
set -e

APP_NAME="fpqdns"
INSTALL_DIR="/opt/${APP_NAME}"
BIN_PATH="/usr/local/bin/${APP_NAME}"

echo "==> FPQDNS 一键安装开始"

# 1. 基础环境
echo "==> 检查系统依赖"
apt update
apt install -y curl wget ca-certificates git

# 2. 安装 Go（如果没有）
if ! command -v go >/dev/null 2>&1; then
  echo "==> 未检测到 Go，开始安装 Go 1.22"
  wget -q https://go.dev/dl/go1.22.6.linux-amd64.tar.gz
  rm -rf /usr/local/go
  tar -C /usr/local -xzf go1.22.6.linux-amd64.tar.gz
  echo 'export PATH=$PATH:/usr/local/go/bin' >/etc/profile.d/go.sh
  source /etc/profile.d/go.sh
fi

go version

# 3. 拉取源码
echo "==> 拉取 FPQDNS 源码"
rm -rf "${INSTALL_DIR}"
git clone https://github.com/lofiey/fpqdns.git "${INSTALL_DIR}"
cd "${INSTALL_DIR}"

# 4. 修复 go.mod（兜底）
echo "==> 初始化 go.mod"
cat > go.mod <<'EOF'
module fpqdns

go 1.22
EOF

# 5. 依赖整理
echo "==> 拉取 Go 依赖"
go mod tidy

# 6. 编译
echo "==> 编译 FPQDNS"
go build -o ${APP_NAME} ./cmd/dns-core

# 7. 安装二进制
echo "==> 安装到 ${BIN_PATH}"
install -m 755 ${APP_NAME} ${BIN_PATH}

# 8. systemd 服务
echo "==> 写入 systemd 服务"
cat > /etc/systemd/system/fpqdns.service <<EOF
[Unit]
Description=FPQDNS Server
After=network.target

[Service]
ExecStart=${BIN_PATH} -c ${INSTALL_DIR}/config.yaml
Restart=on-failure
LimitNOFILE=1048576

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable fpqdns

echo
echo "==> 安装完成"
echo "👉 启动服务: systemctl start fpqdns"
echo "👉 查看状态: systemctl status fpqdns"
echo "👉 Web 面板: http://服务器IP:8080"
