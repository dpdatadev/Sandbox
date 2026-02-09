//#include "commands.hpp"
#include "lib.hpp"

using namespace Services;

std::string __unit_test_01(const char* testString) {
    std::string cppString(testString);
    return cppString;
}

//2/6 ongoing project learning C++

int main(int argc, char* argv[])
{

    cout << "Testing Code.." << endl;

    size_t uuid_size(7);
    cout << LEFT(NEWUUID, uuid_size) << endl;

    std::string _testString{"WouldYouLikeToGoOnADateWithMeNextSaturday\n\n"};
    
    cout << _testString << endl;

    size_t cutAllSize{0};

    std::string cut_testString = extract_left_chars(_testString, cutAllSize);

    cout << cut_testString << endl;

    // const char cString[15]{"Hello, World!\n"};

    // system("clear");
    // cout << __unit_test_01(cString) << endl;

    return 0;
}