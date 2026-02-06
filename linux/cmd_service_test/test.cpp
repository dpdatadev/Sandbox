#include "lib.hpp"

using namespace std;

using namespace Services;

int main(int argc, char* argv[])
{

    cout << "Testing Code.." << endl;
    size_t uuid_size(7);
    cout << LEFT(NEWUUID, uuid_size) << endl;    

    return 0;
}