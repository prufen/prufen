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

# Latest Busybox version as of 2022-10-16
# $ skopeo inspect --format "{{.Name}}@{{.Digest}}" docker://docker.io/library/busybox
FROM docker.io/library/busybox@sha256:9810966b5f712084ea05bf28fc8ba2c8fb110baa2531a10e2da52c1efc504698 AS busybox

# rutsky/cpp-test:17
# $ skopeo inspect --format "{{.Name}}@{{.Digest}}" docker://docker.io/rutsky/cpp-test:17
FROM docker.io/rutsky/cpp-test@sha256:0a0fde45421054d4dc05aae18c6bee5a7073135b0a1d93f09c89e477019e49b7 AS rutsky-cpp-test-17

FROM ubuntu:jammy-20221003
RUN apt-get -y update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    libprotobuf23 \
    libnl-route-3-200 \
    libcap2-bin \
    ca-certificates
COPY --from=nsjail-build /workdir/nsjail /bin/

# Copy images for running containers inside them.
COPY --from=busybox / /images/busybox/
COPY --from=rutsky-cpp-test-17 / /images/rutsky-cpp-test-17/

# Let's drop some potentially dangerous files and attributes. Should not be necessary, but let's do it just in case.

# Reset permissions on directories to be able to search all directories.
RUN find /images/ -type d -not -perm 775 -exec chmod -c 00755 {} \;

# Remove SetUID, SetGID and sticky bits from files.
# Executable bit shold stay, otherwise there will be no executable binaries.
RUN find /images/ -type f -perm /7000 -exec chmod -c -st {} \;

# Check is there any non-regular files.
RUN find /images/ ! -type f -and ! -type d -and ! -type l -print -exec false {} + || (echo "non-regular files in images"; exit 1)

# Check that all files and directories are of reasoable path length
RUN find /images/ -regextype posix-extended -regex ".{200,}" -print -exec false {} + || (echo "too long file/directory path in images"; exit 1)

RUN find /images/ -perm /7000 -printf "%M (%m) %p\n" -exec false {} + || (echo "files/directories with SetGID/SetUID/sticky bits found"; exit 1)

# TODO: Maybe reset ownership too?
# TODO: What about chattr?
# TODO: Add test for resetting permissions.
