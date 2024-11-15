cmake_minimum_required(VERSION 3.10...3.31)
project(ecal-go LANGUAGES C CXX)

find_package(eCAL CONFIG REQUIRED)

add_library(utils INTERFACE)
target_include_directories(utils INTERFACE .)

add_library(ecal_go_core)
target_sources(ecal_go_core PRIVATE
  ./ecal/core.h
  ./ecal/core.cpp
)
target_link_libraries(ecal_go_core PRIVATE eCAL::core)

add_library(ecal_go_publisher)
target_sources(ecal_go_publisher PRIVATE
  ./ecal/publisher/publisher.h
  ./ecal/publisher/publisher.cpp
)
target_link_libraries(ecal_go_publisher PRIVATE eCAL::core utils)

add_library(ecal_go_subscriber)
target_sources(ecal_go_subscriber PRIVATE
  ./ecal/subscriber/subscriber.h
  ./ecal/subscriber/subscriber.cpp
)
target_link_libraries(ecal_go_subscriber PRIVATE eCAL::core utils)

# Subpackages that use cgo
set(subpackages
  "ecal"
  "ecal/publisher"
  "ecal/subscriber"
  "ecal/logging"
)

foreach(package ${subpackages})
  cmake_path(GET package STEM package_name)

  # Have content embedded here to evalue the package_name variable
  file(GENERATE
    OUTPUT "${CMAKE_SOURCE_DIR}/${package}/package_user.go"
    CONTENT "\
package ${package_name}

// #cgo CPPFLAGS: -I$<TARGET_PROPERTY:eCAL::core,INTERFACE_INCLUDE_DIRECTORIES>
// #cgo LDFLAGS: -L$<TARGET_LINKER_FILE_DIR:eCAL::core>
// #cgo LDFLAGS: -Wl,-rpath,$<TARGET_LINKER_FILE_DIR:eCAL::core>
import \"C\"
"
  )
endforeach()
