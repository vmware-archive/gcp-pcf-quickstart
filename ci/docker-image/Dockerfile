FROM ubuntu:16.04

# Packages
RUN DEBIAN_FRONTEND=noninteractive apt-get -y -qq update && apt-get -y -qq install \
  gcc \
  git-core \
  make \
  python-software-properties \
  software-properties-common \
  wget \
  curl \
  dnsutils \
  unzip \
  jq

WORKDIR /tmp/docker-build

# Golang
ENV GO_VERSION=1.12.1
ENV GO_SHA256SUM=2a3fdabf665496a0db5f41ec6af7a9b15a49fbe71a85a50ca38b1f13a103aeec

RUN curl -LO https://storage.googleapis.com/golang/go${GO_VERSION}.linux-amd64.tar.gz && \
    echo "${GO_SHA256SUM}  go${GO_VERSION}.linux-amd64.tar.gz" > go_${GO_VERSION}_SHA256SUM && \
    sha256sum -cw --status go_${GO_VERSION}_SHA256SUM
RUN tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz
ENV GOPATH /root/go
RUN mkdir -p /root/go/bin
ENV PATH $PATH:/usr/local/go/bin:$GOPATH/bin
RUN go get github.com/onsi/ginkgo
RUN go install github.com/onsi/ginkgo/...
RUN go get golang.org/x/lint/golint

# Google SDK
ENV GCLOUD_VERSION=157.0.0
ENV GCLOUD_SHA1SUM=383522491db5feb9f03053f29aaf6a1cf778e070

RUN wget https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-sdk-${GCLOUD_VERSION}-linux-x86_64.tar.gz \
    -O gcloud_${GCLOUD_VERSION}_linux_amd64.tar.gz && \
    echo "${GCLOUD_SHA1SUM}  gcloud_${GCLOUD_VERSION}_linux_amd64.tar.gz" > gcloud_${GCLOUD_VERSION}_SHA1SUM && \
    sha1sum -cw --status gcloud_${GCLOUD_VERSION}_SHA1SUM && \
    tar xvf gcloud_${GCLOUD_VERSION}_linux_amd64.tar.gz && \
    mv google-cloud-sdk / && cd /google-cloud-sdk  && ./install.sh

ENV PATH=$PATH:/google-cloud-sdk/bin

ENV TERRAFORM_VERSION=0.11.10
ENV TERRAFORM_SHA256SUM=43543a0e56e31b0952ea3623521917e060f2718ab06fe2b2d506cfaa14d54527

RUN wget https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip \
    -O terraform.zip && \
    echo "${TERRAFORM_SHA256SUM}  terraform.zip" > terraform_SHA256SUM && \
    sha256sum -cw --status terraform_SHA256SUM && \
    unzip terraform.zip && \
    mv terraform /usr/local/bin && \
    chmod a+x /usr/local/bin/terraform

# Cleanup
RUN rm -rf /tmp/docker-build
