# Set the minimum version of CMake that can be used
# To find the cmake version run
# $ cmake --version
cmake_minimum_required(VERSION 3.0)

# Set the project name
project (bc_demo)

set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -std=c++11")

# Create a sources variable with a link to all cpp files to compile
set(SOURCES
    src/domain/model.cpp
    src/repositories/repository.cpp
    src/services/service.cpp
    main.cpp
)


# Add an executable with the above sources
add_executable(bc_demo ${SOURCES})

# Set the direcoties that should be included in the build command for this target
# when running g++ these will be included as -I/directory/path/
target_include_directories(bc_demo
    PRIVATE ${PROJECT_SOURCE_DIR}/include
)