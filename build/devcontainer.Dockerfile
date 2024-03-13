ARG ALPINE_VERSION=3.19.1

FROM alpine:$ALPINE_VERSION

ARG USER_ID=1000
ARG USER_NAME=default
ARG GROUP_ID=1001
ARG GROUP_NAME=default

RUN apk update && apk upgrade && apk add \
    alpine-sdk \
    linux-headers \
    sudo

RUN mkdir -p /home/$USER_NAME && \
    mkdir -p /etc/sudoers.d && \
    echo "$USER_NAME ALL=(ALL) NOPASSWD:ALL" > /etc/sudoers.d/nopasswd && \
    addgroup -g $GROUP_ID $GROUP_NAME && \
    adduser -u $USER_ID -G $GROUP_NAME -s /bin/sh -D $USER_NAME

USER $USER_ID
RUN chown $USER_NAME:$GROUP_NAME /home/$USER_NAME

RUN sudo apk add \
    git \
    wget \
    make \
    cmake \
    tar \
    bash \
    bash-completion

# Install Golang and another stuff
ARG GOLANG_VERSION=1.22.1
RUN wget https://go.dev/dl/go$GOLANG_VERSION.linux-amd64.tar.gz -O ~/go-$GOLANG_VERSION.tar.gz && \
    sudo rm -rf /usr/local/go && \
    sudo tar -C /usr/local/ -xzf /home/$USER_NAME/go-$GOLANG_VERSION.tar.gz && \
    sudo ln -s -f /usr/local/go/bin/go /usr/bin/go && \
    rm -rf /home/$USER_NAME/go-$GOLANG_VERSION.tar.gz

RUN go install github.com/go-delve/delve/cmd/dlv@latest && \
    sudo cp /home/$USER_NAME/go/bin/dlv /usr/bin/dlv

ENTRYPOINT ["sleep", "infinity"]