FROM ubuntu:latest


# Update
RUN apt update

# GCC and other essentials
RUN apt install -y libpq-dev build-essential

# Networking tools
RUN apt install -y \
        traceroute \
        curl \
        iputils-ping \
        bridge-utils \
	apt-transport-https \
    	ca-certificates \
    	gnupg-agent \
    	software-properties-common \
        dnsutils \
        netcat-openbsd \
        jq \
        postgresql-client \
        nmap \
        net-tools \
        && rm -rf /var/lib/apt/lists/*

# Docker GPG key
RUN curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add -

# Docker PPA
RUN add-apt-repository \
    "https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" \
    && apt update

# Docker daemon
RUN apt install -y docker-ce docker-ce-cli containerd.io

# Golang
RUN apt install -y golang-1.10

ENV PATH="/usr/lib/go-1.10/bin:$PATH"

WORKDIR /usr/bin/theangrydev_runner

COPY . .

# Install dependencies
# TODO: Use go modules instead
RUN go get github.com/docker/docker/client


EXPOSE 8080

ENTRYPOINT ["./entrypoint.sh"]
