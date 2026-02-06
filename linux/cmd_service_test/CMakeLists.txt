# Specifies the minimum required version of CMake
cmake_minimum_required(VERSION 3.12)

# Defines the project name and language(s)
project(CommandServer VERSION 1.0 LANGUAGES CXX)

# Specifies the C++ standard to use
set(CMAKE_CXX_STANDARD 20)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

# Adds an executable target named "HelloWorldProject" from the source file
add_executable(CommandServer main.cpp)
