#include <iostream>
#include "node.h"

void Node::displayNodeData() {
    cout << "Node data: " << data << endl;
}

void Node::displayNextNodeData() {
    if (next != nullptr) {
        cout << "Next node data: " << next->data << endl;
    } else {
        cout << "No next node." << endl;
    }
}

void Node::displayPrevNodeData() {
    if (prev != nullptr) {
        cout << "Previous node data: " << prev->data << endl;
    } else {
        cout << "No previous node." << endl;
    }
}

int main()
{
    cout << "Hello, Raspberry Pi DLL!" << endl;
    Node* headNode = new Node(DLL_HEAD);
    Node* nextNode = new Node(DLL_MIDDLE);
    Node* tailNode = new Node(DLL_TAIL);
    nextNode->prev = headNode;
    headNode->next = nextNode;
    tailNode->prev = nextNode;
    nextNode->next = tailNode;
    tailNode->next = nullptr;


    cout << "Node position and display:" << "\n" << endl;
    headNode->displayNodeData();
    headNode->displayNextNodeData();
    headNode->displayPrevNodeData();
    cout << "\n";
    nextNode->displayNodeData();
    nextNode->displayNextNodeData();
    nextNode->displayPrevNodeData();
    cout << "\n";
    tailNode->displayNodeData();
    tailNode->displayNextNodeData();
    tailNode->displayPrevNodeData();

    cout << "\nCleaning up memory..." << endl;

    delete headNode;
    delete nextNode;
    delete tailNode;

    cout << "\n::END::" << endl;

    return 0;
}