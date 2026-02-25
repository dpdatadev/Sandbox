"""
PSEUDOCODE

class   NODE(item)
                data =  item
                left =  None
                right = None

func    SEARCH_KEY(root(Node), item) <return BOOL>

        is the node empty?
            return False

        is the current node(root) data == the search term?
            return True

        is the search term/key greater than the root key?
            return SEARCH(root.right, item) #recurse
        else
            return SEARCH(root.left, item) #recurse

"""


# Optimize the above algorithm from O(h)/O(h) to O(h)/O(1) <remove auxilliary space and recursion stack via iterative approach>
# Node structure
class Node:
    def __init__(self, item):
        self.data = item
        self.left = None
        self.right = None


def search(root, key):
    present = False

    # iterative traversal
    while root is not None:
        if root.data == key:
            present = True
            break
        elif key > root.data:
            root = root.right
        else:
            root = root.left

    return present


if __name__ == "__main__":
    # Creating BST
    #     6
    #   / \
    #   2   8
    #      / \
    #     7   9
    root = Node(6)
    root.left = Node(2)
    root.right = Node(8)
    root.right.left = Node(7)
    root.right.right = Node(9)

    key = 7
    # Searching for key in the BST
    print(search(root, key))
