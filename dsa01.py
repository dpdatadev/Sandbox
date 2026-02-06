# Selection sort array A
# O(n) loop over array
# O(1) initial index of max
# O(i) search for max in A[:i]
# O(1) check for larger value
# O(1) new max found
# O(1) swap
numbers = [85, 12, 2, 3, 400, 71]

def sel_sort(nums):
    # count index descending to 1
    for i in range(len(nums) - 1, 0, -1):
        m = i # largest index set
        for j in range(i): # now iterate over the same array in the opposite direction to compare the first with the last
            if nums[m] < nums[j]:
                m = j # new max found
        nums[m], nums[i] = nums[i], nums[m] # swap O(1)



if __name__ == '__main__':
    print(numbers)
    sel_sort(numbers)
    print(numbers)