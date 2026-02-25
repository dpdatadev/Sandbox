#include <iostream>
#include <array>
#include <vector>
#include <bits/stdc++.h>
#include <limits>
#include <math.h> // for c example

using namespace std;
// using namespace vector_ops;

typedef u_int32_t bit32;
typedef u_int8_t bit8;

// depc
int find_max_diff_nums(int nums[], size_t nums_size)
{

    float neg_inf = INFINITY;
    int num_tracker[nums_size];

    for (int i = 0; i <= (nums_size); i++)
    {
        printf(" %d -- ", nums[i]);
    }

    return 0;
}

// todo - over engineer example ...
/*
namespace vector_ops
{
    typedef struct vector_max
    {
        vector<int> indices;
        int maximum_difference;
    } vector_max;

    std::ostream &operator<<(std::ostream &out, const vector_max &data)
    {
        out << "Index value(s): " << data.indices
            << ", Maximum Distance Value: " << data.maximum_difference;
        return out;
    }

    vector_max* cpp_findMaxVectorIntDifference(const vector<int> &vec)
    {
        int length = vec.size();
        vector<int> tracker = {};

        vector_max diff_report = new vector_max();

        float max_diff = -std::numeric_limits<float>::infinity();

        if (length < 2)
        {
            exit(EXIT_FAILURE);
        }

        for (int i = 1; i <= length; i++)
        {
            int curr_val = vec[i];
            int prev_val = vec[i - 1];
            int diff = (curr_val - prev_val);

            if (diff >= max_diff)
            {
                max_diff = diff;
                tracker.clear();
                tracker.push_back(i);
            }
            else if (diff == max_diff)
            {
                tracker.push_back(i);
            }

            diffReport.indices = tracker;
            diffReport.maximum_difference = max_diff;
        }

        return &diffReport;
    }

}
*/

//incorrect version 2
void cpp_findMaxVectorIntDifference(const vector<int>& vec)
{
    int length = vec.size();

    float max_diff = -std::numeric_limits<float>::infinity();

    if (length < 2)
    {
        exit(EXIT_FAILURE);
    }

    for (int i = 0; i <= (length); i++)
    {
        int curr_val = vec[i];
        int prev_val = vec[i - 1];
        int diff = (curr_val - prev_val);

        if (diff > max_diff && diff > 1) // dont handle negative diffs or 0
        {
            max_diff = diff;
            cout << "Maximum difference " << max_diff << " found at index: " << i << endl;
        }
    }
}

// correct version
#include <vector>
#include <stdexcept>
#include <algorithm>

int findMaxDifference(const std::vector<int>& vec)
{
    if (vec.size() < 2) {
        throw std::invalid_argument("Vector must contain at least two elements");
    }

    int max_diff = vec[1] - vec[0];

    for (std::size_t i = 1; i < vec.size(); ++i) {
        int diff = vec[i] - vec[i - 1];
        max_diff = std::max(max_diff, diff);
    }

    return max_diff;
}


typedef struct Point
{
    std::string label;
    int n1;
    int n2;
} Point;

float getDistance(const Point &x, const Point &y)
{
    return sqrt(pow((x.n2 - x.n1), 2) + pow((y.n2 - y.n1), 2) * 1.0);
}

void distanceTest()
{
    Point p1, p2, p3, p4;
    p1.label = "x";
    p1.n1 = 1;
    p1.n2 = 13;

    p2.label = "x";
    p2.n1 = 2;
    p2.n2 = 10;

    p3.label = "y";
    p3.n1 = 4;
    p3.n2 = 6;

    p4.label = "y";
    p4.n1 = 8;
    p4.n2 = 4;

    array<Point, 2> firstCoordinates{p1, p3};

    array<Point, 2> secondCoordinates{p2, p4};

    cout << getDistance(firstCoordinates[0], firstCoordinates[1]) << endl;
    cout << getDistance(secondCoordinates[0], secondCoordinates[1]) << endl;
}

// masking
#define BIT(n) (1U << (n))

void test1()
{
    bit32 testbyte = 0b00000000;
    cout << testbyte << endl;
    testbyte |= (1U << 1);
    cout << testbyte << endl;

    if (0b00000010 & testbyte)
    {
        cout << "success.." << endl;
    }
    else
    {
        cerr << "FAIL" << endl;
    }
    bit32 tb2 = testbyte;

    tb2 |= BIT(0);

    if (0b00000001 & tb2)
    {
        cout << "success.." << endl;
    }
    else
    {
        cerr << "FAIL" << endl;
    }
}

short test2()
{
    short i = 0b0101;
    return i <<= 2;
}

// multiplies by 2^n (4)
int test3()
{
    int i = 5;
    return i <<= 2;
}

// divides by 2^n (4)
int test4()
{
    int i = 48;
    return i >>= 2;
}

int main()
{
    // cout << "Bitmasking Etc.,\n"
    //<< endl;
    bit32 bits = 0b00000000; // set to 0

    // cout << bits << endl;

    bits |= 0b00000001; // manually set bit

    // cout << bits << endl;

    bits ^= bits; // toggle bits (1 back to 0)

    // cout << bits << endl;

    bit32 y = 0b00000001;

    y = y << 3;

    bits |= BIT(3); // Set on (1 x 2^3 = 8, 0b00001000)

    // cout << bits << endl;

    bits ^= BIT(3); // toggle off

    // cout << bits << endl;

    bits ^= BIT(3); // toggle back

    if (bits & BIT(3))
    {
        // cout << "The bit is set.." << endl;
    }
    else
    {
        // cout << "The bit is not set.." << endl;
    }

    // packing multiple values into one integer

    u_int8_t status = 0;

    /*
    bit 0: power on
    bit 1: error
    bit 2: data ready
    */

    status |= BIT(0); // power on

    // cout << status << endl;

    status |= BIT(2); // data ready

    // cout << status << endl;

    // bit slicing
    // bits 7 - 5 = mode
    // bits 4 - 0 = value

    bit32 reg = 0b11010110;

    bit32 mode = (reg >> 5) & 0b111;

    // cout << mode << endl;

    // bitwise assembly of multi-byte values

    // test1();
    // cout << test2() << endl;
    // cout << test3() << endl;
    // cout << test4() << endl;

    vector<int> vec{27, 2, 3, 4, 50};
    // int nums[] = {1, 2, 3, 4, 5, 44, 1093, 4};
    // size_t nums_size = (sizeof(nums) / sizeof(nums[0]));

    // cout << find_max_diff_nums(nums, nums_size);

    cout << findMaxDifference(vec) << endl;

    // distanceTest();
    return 0;
}