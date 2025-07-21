# syntax=docker/dockerfile:1

FROM golang:1.21

WORKDIR /wotlk
COPY . .
COPY gitconfig /etc/gitconfig

RUN apt-get update
RUN apt-get install -y protobuf-compiler
RUN go get google.golang.org/protobuf@v1.31.0
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31.0

ENV NVM_DIR="/root/.nvm"
ENV NODE_VERSION=19.8.0

RUN curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.38.0/install.sh | bash && \
    . "$NVM_DIR/nvm.sh" && \
    nvm install ${NODE_VERSION} && \
    nvm use ${NODE_VERSION} && \
    nvm alias default ${NODE_VERSION}

ENV PATH="/root/.nvm/versions/node/v${NODE_VERSION}/bin:${PATH}"

EXPOSE 8080/tcp
