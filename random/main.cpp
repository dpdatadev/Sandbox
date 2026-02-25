using namespace std;
#include <iostream>
#include <vector>

int main() {

    cout << "Hello, World!\n" << endl;

    auto blowUp = [](int val) { return val * val; };

    std::size_t size = 10;

    vector<int> vec[size + 1];

    for(int i = size; i > 0; --i) {
        vec[i].push_back(blowUp(i));
        cout << "vec[" << i << "] = " << vec[i][0] << endl;
        cout << "Address of vec[" << i << "] is " << &vec[i] << endl;
    }


    return 0;
}