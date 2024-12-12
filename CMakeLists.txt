cmake_minimum_required(VERSION 3.10)

project(LLWC)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED True)

add_executable(LLWC src/main.cpp)

set(CSV_DIR ${CMAKE_CURRENT_SOURCE_DIR}/external/pokeapi/data/v2/csv)

message(STATUS "Resource directory: ${CSV_DIR}")


add_custom_target(
    copy_resources ALL
    COMMAND ${CMAKE_COMMAND} -E make_directory ${CMAKE_BINARY_DIR}/resources
    COMMAND ${CMAKE_COMMAND} -E copy_directory ${CSV_DIR} ${CMAKE_BINARY_DIR}/resources
    COMMENT "Copying resources to build directory"
)

add_dependencies(LLWC copy_resources)

install(TARGETS LLWC
    RUNTIME DESTINATION bin
)

install(DIRECTORY ${CSV_DIR}
    DESTINATION share/LLWC/resources
)

message(STATUS "Resources will be copied to ${CMAKE_BINARY_DIR}/resources during build")
message(STATUS "Resources will be installed to share/LLWC/resources during installation")