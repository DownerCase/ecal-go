ARG BASE_IMAGE=docker.io/alpine:3.21.2

###
# Build eCAL 6 from source
###
FROM ${BASE_IMAGE} AS builder

RUN apk add --no-cache \
    asio-dev \
    cmake \
    g++ \
    git \
    ninja-is-really-ninja \
    protobuf-dev \
    tclap-dev \
    yaml-cpp-dev \
    # fts functions which are missing musl libc
    musl-fts-dev

# filter=tree:0 will fetch the commits but their trees will be fetched on-demand
# So we end up with the full commit history but only download trees for HEAD
# This is needed over shallow clone for the tag-based auto-versioning
# Reset added to ensure the same commit is always used, even if master advances
RUN git clone --single-branch --filter=tree:0 --branch master https://github.com/DownerCase/ecal.git /ecal \
    && git -C /ecal reset --hard dcd73635dfcaec79634ec009060b19bc77e11f86 \
    # Only download the submodules for dependencies not available as system packages
    && git -C /ecal submodule update --init --single-branch --depth 1 \
    thirdparty/recycle/recycle \
    thirdparty/tcp_pubsub/tcp_pubsub \
    thirdparty/ecaludp/ecaludp

WORKDIR /ecal

# TODO: Blocker for installing to / or / usr
# - Fix usage of `INSTALL_INCLUDE_DIR` instead of `eCAL_install_include_dir`
# - Get CMakeFunctions to use GNUInstallDirs

RUN \
    # Configure a minimal build with less vendored dependencies
    cmake --preset core -G Ninja \
    -DCMAKE_BUILD_TYPE=Release \
    -DCMAKE_STAGING_PREFIX=staging \
    -DBUILD_TIME=OFF \
    -DECAL_THIRDPARTY_BUILD_ASIO=OFF \
    -DECAL_THIRDPARTY_BUILD_PROTOBUF=OFF \
    -DECAL_THIRDPARTY_BUILD_YAML-CPP=OFF \
    -DECAL_THIRDPARTY_BUILD_TCLAP=OFF \
    # Build core preset
    &&  cmake --build --preset core \
    # Component based install avoids the unneeded `app` component
    && cmake --install ./out/core/build --strip --component configuration \
    && cmake --install ./out/core/build --strip --component sdk \
    && cmake --install ./out/core/build --strip --component Unspecified

###
# Dev base build-stage, for consuming eCAL
###
FROM ${BASE_IMAGE} AS dev-base

LABEL org.opencontainers.image.source https://github.com/downercase/ecal-go

COPY --from=builder /ecal/staging /usr/local

###
# Dev build-stage, for consuming eCAL and building/testing the C interface library
###
FROM dev-base AS dev-cpp

RUN apk add --no-cache \
    cmake \
    g++ \
    protobuf-dev \
    # Ninja build compatible
    samurai

###
# Dev build-stage, for consuming eCAL and building/testing the Go code
###
FROM dev-base AS dev-go

RUN apk add --no-cache \
    g++ \
    go \
    musl-fts \
    yaml-cpp

