FROM docker/for-desktop-kernel:5.10.47-0b705d955f5e283f62583c4e227d64a7924c138f AS ksrc
FROM ubuntu:latest

WORKDIR /
COPY --from=ksrc /kernel-dev.tar /
RUN tar xf kernel-dev.tar && rm kernel-dev.tar

RUN apt-get update && apt-get -y install wget vim gcc make
RUN wget https://golang.org/dl/go1.17.1.linux-amd64.tar.gz && \
    tar -C /usr/local -xvzf go1.17.1.linux-amd64.tar.gz && \
    rm go1.17.1.linux-amd64.tar.gz
RUN apt-get update && apt-get install -y software-properties-common
RUN add-apt-repository ppa:hadret/bpfcc
RUN add-apt-repository ppa:hadret/libbpf
RUN add-apt-repository ppa:hadret/bpftrace
RUN apt-get update && apt-get install -y libbpfcc-dev bpfcc-tools

ENV PATH=$PATH:/usr/local/go/bin
ENV GOPATH=/root/
COPY vimrc bashrc /root

CMD mount -t debugfs none /sys/kernel/debug && /bin/bash
