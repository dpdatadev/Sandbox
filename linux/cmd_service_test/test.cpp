#include <iostream>

using namespace std;

int main() {

    
    cout << "Testing..." << endl;
    static_assert((sizeof(int) < 5), "int is 4");
    return 0;
}