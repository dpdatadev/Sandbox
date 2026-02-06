#include <iostream>
using namespace std;

#define DLL_HEAD 1
#define DLL_TAIL 2
#define DLL_MIDDLE 0

class Node {
    public:
        int data;
        Node* prev;
        Node* next;
        Node(int val) {
            data = val;
            prev = nullptr;
            next = nullptr;
        }
        void displayNodeData();
        void displayNextNodeData();
        void displayPrevNodeData();
};