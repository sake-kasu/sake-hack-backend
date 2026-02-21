FROM node:24.11.1-slim

ENV GO_VERSION=1.26.0

RUN apt-get update && \
  apt-get install -y \
  curl \
  unzip \
  git \
  ca-certificates \
  build-essential \
  less \
  vim \
  zsh \
  file \
  wget \
  procps

# go
RUN ARCH=$(uname -m) && \
  if [ "$ARCH" = "x86_64" ]; then \
  GOARCH="amd64"; \
  elif [ "$ARCH" = "aarch64" ]; then \
  GOARCH="arm64"; \
  else \
  echo "Unsupported architecture: $ARCH" && exit 1; \
  fi && \
  curl -LO https://go.dev/dl/go${GO_VERSION}.linux-${GOARCH}.tar.gz && \
  tar -C /usr/local -xzf go${GO_VERSION}.linux-${GOARCH}.tar.gz && \
  rm go${GO_VERSION}.linux-${GOARCH}.tar.gz

ENV PATH="/usr/local/go/bin:/bun/bin:$PATH"
ENV GOPATH="/home/node/go"
ENV PATH="$PATH:$GOPATH/bin"

# AWS SSO環境変数の設定
ENV AWS_CONFIG_FILE="/home/node/.aws/config"
ENV AWS_SHARED_CREDENTIALS_FILE="/home/node/.aws/credentials"

# golang-migrate
RUN ARCH=$(uname -m) && \
  if [ "$ARCH" = "x86_64" ]; then \
  GOARCH="amd64"; \
  elif [ "$ARCH" = "aarch64" ]; then \
  GOARCH="arm64"; \
  else \
  echo "Unsupported architecture: $ARCH" && exit 1; \
  fi && \
  curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.3/migrate.linux-${GOARCH}.tar.gz | tar xvz \
  && mv migrate /usr/local/bin/migrate \
  && chmod +x /usr/local/bin/migrate

# AWS CLI v2
RUN ARCH=$(uname -m) && \
  if [ "$ARCH" = "x86_64" ]; then \
  curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"; \
  elif [ "$ARCH" = "aarch64" ]; then \
  curl "https://awscli.amazonaws.com/awscli-exe-linux-aarch64.zip" -o "awscliv2.zip"; \
  else \
  echo "Unsupported architecture: $ARCH" && exit 1; \
  fi && \
  unzip awscliv2.zip && \
  ./aws/install && \
  rm -rf awscliv2.zip aws

# AWS Session Manager Plugin
RUN ARCH=$(uname -m) && \
  if [ "$ARCH" = "x86_64" ]; then \
  curl -fL "https://s3.amazonaws.com/session-manager-downloads/plugin/latest/ubuntu_64bit/session-manager-plugin.deb" -o "session-manager-plugin.deb"; \
  elif [ "$ARCH" = "aarch64" ]; then \
  curl -fL "https://s3.amazonaws.com/session-manager-downloads/plugin/latest/ubuntu_arm64/session-manager-plugin.deb" -o "session-manager-plugin.deb"; \
  else \
  echo "Unsupported architecture: $ARCH" && exit 1; \
  fi && \
  # ダウンロードファイルの検証 \
  if [ ! -f "session-manager-plugin.deb" ] || [ ! -s "session-manager-plugin.deb" ]; then \
  echo "Error: session-manager-plugin.deb download failed or empty file" && exit 1; \
  fi && \
  # ファイル形式の簡易チェック \
  file session-manager-plugin.deb && \
  # インストール実行 \
  dpkg -i session-manager-plugin.deb && \
  rm session-manager-plugin.deb

# zsh
RUN su - node -c "\
  sh -c \"\$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)\" --unattended && \
  git clone https://github.com/zsh-users/zsh-autosuggestions ~/.oh-my-zsh/custom/plugins/zsh-autosuggestions && \
  git clone https://github.com/zsh-users/zsh-syntax-highlighting ~/.oh-my-zsh/custom/plugins/zsh-syntax-highlighting \
  "

RUN mkdir -p /home/node/.zsh && \
  curl -fsSL https://raw.githubusercontent.com/git/git/master/contrib/completion/git-prompt.sh -o /home/node/.zsh/git-prompt.sh && \
  curl -fsSL https://raw.githubusercontent.com/git/git/master/contrib/completion/git-completion.bash -o /home/node/.zsh/git-completion.bash && \
  chown -R node:node /home/node/.zsh

# AWS設定ディレクトリの作成
RUN mkdir -p /home/node/.aws && \
  chown -R node:node /home/node/.aws

# TestContainers統合テスト用
# ホストのdockerソケット(GID=991)にアクセスできるようにする
RUN groupadd -g 991 docker || true && \
  usermod -aG docker node

RUN chsh -s /bin/zsh node

USER root

# golangci-lint
RUN ARCH=$(uname -m) && \
  if [ "$ARCH" = "x86_64" ]; then \
  GOARCH="amd64"; \
  elif [ "$ARCH" = "aarch64" ]; then \
  GOARCH="arm64"; \
  else \
  echo "Unsupported architecture: $ARCH" && exit 1; \
  fi && \
  curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b /usr/local/bin latest

# redoclyはoapi-codegen用
RUN npm config set prefix /home/node/.npm-global && \
  npm install -g @anthropic-ai/claude-code @redocly/cli
ENV PATH="/home/node/.npm-global/bin:$PATH"

WORKDIR /workspace
