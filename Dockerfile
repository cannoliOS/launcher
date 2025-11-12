FROM debian:bullseye

RUN apt-get update && apt-get install -y \
    build-essential \
    pkg-config \
    libsdl2-dev \
    libsdl2-ttf-dev \
    libsdl2-image-dev \
    libsdl2-gfx-dev \
    wget \
    ca-certificates \
    git

RUN wget -q https://dl.google.com/go/go1.25.3.linux-arm64.tar.gz && \
    tar -C /usr/local -xzf go1.25.3.linux-arm64.tar.gz && \
    rm go1.25.3.linux-arm64.tar.gz

ENV PATH=$PATH:/usr/local/go/bin

WORKDIR /build

COPY go.mod go.sum ./

RUN GOWORK=off go mod download

COPY . .
RUN GOWORK=off CGO_ENABLED=1 GOOS=linux GOARCH=arm64 go build -v -gcflags="all=-N -l" -o cannoliOS app/cannoli.go
RUN GOWORK=off CGO_ENABLED=1 GOOS=linux GOARCH=arm64 go build -v -gcflags="all=-N -l" -o igm app/igm.go

CMD ["/bin/bash"]