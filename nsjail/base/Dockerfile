FROM ubuntu:jammy-20221003 as nsjail-build
WORKDIR /workdir
RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y -q \
    git

# nsjail dependencies from
# <https://github.com/google/nsjail/blob/master/Dockerfile>

RUN apt-get -y update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    autoconf \
    bison \
    flex \
    gcc \
    g++ \
    git \
    libprotobuf-dev \
    libtool \
    make \
    pkg-config \
    protobuf-compiler
# Optional dependencies
RUN apt-get -y update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    libnl-route-3-dev

RUN git clone https://github.com/google/nsjail.git .
RUN git checkout 3.1
RUN git submodule update --init --recursive

RUN make

FROM ubuntu:jammy-20221003
RUN apt-get -y update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    libprotobuf23 \
    libnl-route-3-200 \
    libcap2-bin \
    ca-certificates
COPY --from=nsjail-build /workdir/nsjail /bin/
