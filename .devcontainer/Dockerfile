ARG BASE_IMAGE=docker.io/ubuntu:oracular-20241120

###
# Build eCAL 6 from source
###
FROM ${BASE_IMAGE} AS builder

RUN apt-get update && apt-get install --no-install-recommends  -y\
    ca-certificates \
    cmake \
    g++ \
    git \
    libasio-dev \
    libprotobuf-dev \
    libprotoc-dev \
    libtclap-dev \
    libyaml-cpp-dev \
    ninja-build \
    protobuf-compiler \
    && rm -rf /var/lib/apt/lists/*

# filter=tree:0 will fetch the commits but their trees will be fetched on-demand
# So we end up with the full commit history but only download trees for HEAD
# This is needed over shallow clone for the tag-based auto-versioning
# Reset added to ensure the same commit is always used, even if master advances
RUN git clone --single-branch --filter=tree:0 --branch master https://github.com/eclipse-ecal/ecal.git /ecal \
    && git -C ./ecal reset --hard 7ba8e0a80328b0fbc2924eaca40e98669f7933bc
WORKDIR /ecal

# Only download the submodules for dependencies not available as system packages
RUN git submodule update --init --single-branch --depth 1 \
    thirdparty/recycle/recycle \
    thirdparty/tcp_pubsub/tcp_pubsub \
    thirdparty/ecaludp/ecaludp

# TODO: Blocker for installing to / or / usr
# - Fix usage of `INSTALL_INCLUDE_DIR` instead of `eCAL_install_include_dir`
# - Get CMakeFunctions to use GNUInstallDirs

# Configure a minimal build with less vendored dependencies
RUN cmake --preset core -G Ninja \
 -DCMAKE_BUILD_TYPE=Release \
 -DCMAKE_STAGING_PREFIX=staging \
 -DBUILD_TIME=OFF \
 -DECAL_THIRDPARTY_BUILD_ASIO=OFF \
 -DECAL_THIRDPARTY_BUILD_PROTOBUF=OFF \
 -DECAL_THIRDPARTY_BUILD_YAML-CPP=OFF \
 -DECAL_THIRDPARTY_BUILD_TCLAP=OFF

RUN cmake --build --preset core

# Component based install avoids the unneeded `app` component
RUN cmake --install ./out/core/build --strip --component configuration \
    && cmake --install ./out/core/build --strip --component sdk \
    && cmake --install ./out/core/build --strip --component Unspecified

###
# Dev base build-stage, for consuming eCAL
###
FROM ${BASE_IMAGE} AS dev-base

LABEL org.opencontainers.image.source https://github.com/downercase/ecal-go

RUN apt-get update && apt-get install --no-install-recommends  -y\
    g++ \
    libyaml-cpp0.8 \
    && rm -rf /var/lib/apt/lists/*

# copy the ca-certificate.crt from the build stage
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /ecal/staging /usr/local

# Update dynamic linker so it knows about our new libraries
RUN ldconfig

###
# Dev build-stage, for consuming eCAL and building/testing the C interface library
###
FROM dev-base AS dev-cpp

RUN apt-get update && apt-get install --no-install-recommends  -y\
    cmake \
    libprotobuf-dev \
    ninja-build \
    && rm -rf /var/lib/apt/lists/*

###
# Dev build-stage, for consuming eCAL and building/testing the Go code
###
FROM dev-base AS dev-go

RUN apt-get update && apt-get install --no-install-recommends  -y\
    golang \
    && rm -rf /var/lib/apt/lists/*

