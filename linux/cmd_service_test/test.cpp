#include "commands.hpp"

using namespace Services;

std::string __unit_test_01(const char* testString) {
    std::string cppString(testString);
    return cppString;
}

//2/6 ongoing project learning C++

int main(int argc, char* argv[])
{

    cout << "Testing Code.." << endl;
    //size_t uuid_size(7);
    //cout << LEFT(NEWUUID, uuid_size) << endl;

    const char cString[15]{"Hello, World!\n"};

    //system("clear");
    cout << __unit_test_01(cString) << endl;


    return 0;
}