cmake_minimum_required(VERSION 3.10...3.31)
project(ecal-go LANGUAGES C CXX)

find_package(eCAL CONFIG REQUIRED)

add_library(ecal_go)
target_sources(ecal_go PRIVATE
  ./ecal/ecal_go.h
  ./ecal/ecal_go.cpp
)
target_link_libraries(ecal_go PRIVATE eCAL::core)

file(GENERATE
  OUTPUT "${CMAKE_SOURCE_DIR}/ecal/package_user.go"
  INPUT "./cmake/package_user.go.in"
)
