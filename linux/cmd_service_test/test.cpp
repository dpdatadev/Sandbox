#include "commands.hpp"

using namespace Services;

//2/6 ongoing project learning C++

int main(int argc, char* argv[])
{

    cout << "Testing Code.." << endl;
    size_t uuid_size(7);
    cout << LEFT(NEWUUID, uuid_size) << endl;

    return 0;
}