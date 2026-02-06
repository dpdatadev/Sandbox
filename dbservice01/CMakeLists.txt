cmake_minimum_required(VERSION 3.19)

project(DBSERVICE01)

find_package(SQLite3 REQUIRED)

file(GLOB_RECURSE SOURCES CONFIGURE_DEPENDS *.cc *.h)

add_executable(DBSERVICE01 main.cpp)

target_include_directories(DBSERVICE01 PRIVATE /include)

target_link_libraries(DBSERVICE01
        SQLite::SQLite3
        )