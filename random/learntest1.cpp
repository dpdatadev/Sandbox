using namespace std;
#include <iostream>
#include <memory>
#include <vector>

class Data{
    public:
        Data(int val) : value(val) {}
        int getValue() const { return value; }
    private:
        int value;
};

shared_ptr<Data> data_ptr = make_shared<Data>(Data(50));

int change_value(int x) {
    x = x * x;
    return x;
}

// Reference that does not change the original value
int no_change_reference(const int& x) {
    return x * x;
}

int change_reference(int& x) {
    x = x * x;
    return x;
}

int* change_pointer(int* x) {
    *x = (*x) * (*x);
    return x;
}


int main() {
    int x = 100;
    cout << x << endl;
    x = change_value(x);
    cout << x << "\n" << endl;
    int* y = change_pointer(&x);
    cout << *y << endl;

    int z = 4;
    int* zz = &z;
    cout << *zz << endl;
    int& w = z;
    cout << w << endl;
    /*
    x = 100;
    cout << "\n" << x << endl;
    int h = change_value(x);
    cout << h << "\n" << endl;
    cout << x << "\n" << endl;
    */
    x = 100;
    cout << "\n" << x << endl;
    //int h1 = no_change_reference(x);
    //cout << h1 << "\n" << endl;
    //cout << x << "\n" << endl;
    int h = no_change_reference(x);
    cout << h << "\n" << endl;
    cout << x << "\n" << endl;

    vector<unique_ptr<Data>> data_vec = {make_unique<Data>(Data(1)), make_unique<Data>(Data(2)), make_unique<Data>(Data(3))};
    for(const auto& item : data_vec) {
        cout << "<<" << item->getValue() << ">>" << "\n" << endl;
    }
}