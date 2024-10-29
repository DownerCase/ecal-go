ARG BASE_IMAGE=docker.io/ubuntu:oracular-20241009
FROM ${BASE_IMAGE} AS builder

RUN apt-get update && apt-get install --no-install-recommends  -y\
    ca-certificates \
    cmake \
    g++ \
    git \
    libasio-dev \
    libprotobuf-dev \
    libprotoc-dev \
    libsimpleini-dev \
    libtclap-dev \
    libyaml-cpp-dev \
    ninja-build \
    protobuf-compiler \
    && rm -rf /var/lib/apt/lists/*

RUN git clone --single-branch --filter=tree:0 --branch master https://github.com/eclipse-ecal/ecal.git /ecal
WORKDIR /ecal
RUN git submodule update --init --single-branch --depth 1 \
    thirdparty/recycle/recycle \
    thirdparty/tcp_pubsub/tcp_pubsub \
    thirdparty/ecaludp/ecaludp

# TODO: Blocker for using CMAKE_STAGING_PREFIX
# - Fix usage of `INSTALL_INCLUDE_DIR` instead of `eCAL_install_include_dir`
# - Get CMakeFunctions to use GNUInstallDirs

RUN cmake --preset core -G Ninja \
 -DCMAKE_BUILD_TYPE=Release \
 -DCMAKE_STAGING_PREFIX=staging \
 -DBUILD_TIME=OFF \
 -DECAL_THIRDPARTY_BUILD_ASIO=OFF \
 -DECAL_THIRDPARTY_BUILD_PROTOBUF=OFF \
 -DECAL_THIRDPARTY_BUILD_YAML-CPP=OFF \
 -DECAL_THIRDPARTY_BUILD_TCLAP=OFF

RUN cmake --build --preset core

RUN cmake --install ./out/core/build --strip --component configuration \
    && cmake --install ./out/core/build --strip --component sdk

FROM ${BASE_IMAGE} AS dev

LABEL org.opencontainers.image.source https://github.com/downercase/ecal-go

RUN apt-get update && apt-get install --no-install-recommends  -y\
    g++ \
    golang \
    libyaml-cpp0.8 \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /ecal/staging /usr/local

# Update dynamic linker so it knows about our new libraries
RUN ldconfig
