FROM ubuntu:18.04 as nsjail-build
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
RUN git checkout 2.8
RUN git submodule update --init --recursive

RUN make

FROM ubuntu:18.04
RUN apt-get -y update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    libprotobuf10 \
    libnl-route-3-200 \
    libcap2-bin \
    ca-certificates
COPY --from=nsjail-build /workdir/nsjail /bin/
