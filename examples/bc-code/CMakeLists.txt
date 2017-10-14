# Set the minimum version of CMake that can be used
# To find the cmake version run
# $ cmake --version
cmake_minimum_required(VERSION 3.0)

# Set the project name
project (bc_demo)

set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -std=c++11")

# Create a sources variable with a link to all cpp files to compile
file(GLOB_RECURSE SOURCES src/*.cpp)


# Add an executable with the above sources
add_executable(bc_demo ${SOURCES} main.cpp)

# Set the direcoties that should be included in the build command for this target
# when running g++ these will be included as -I/directory/path/
target_include_directories(bc_demo
    PRIVATE ${PROJECT_SOURCE_DIR}/include
)

include_directories(${PROJECT_SOURCE_DIR}/include)

enable_testing()
find_package(GTest REQUIRED)
include_directories(${GTEST_INCLUDE_DIRS})

file(GLOB TEST_SRC_FILES test/*.cpp)
add_executable(bc_demo_test ${SOURCES} ${TEST_SRC_FILES})
target_link_libraries(bc_demo_test ${GTEST_LIBRARIES})


add_test(AllTests bc_demo_test)