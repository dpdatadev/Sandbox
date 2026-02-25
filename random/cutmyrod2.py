


if __name__ == '__main__':
    rod = [1, 5, 8, 9, 10, 17, 17, 20]
    '''
    #Ignore all this - see ChatGPT output instead..
    The classic "Rod Cutting" Algorithm was the first algorithm exercize I was exposed to.
    I will admit, it took me a little time to fully understand what was being asked.
    Perhaps I'm a little slow in mathematics compared to others - but I really just need to feel that I 
    grasp what is being asked. After pouring over the various examples for this "problem" on several iterations,
    I find an opinion that the examples given are usually blown up to demonstrate harder to program problems with 
    a trivial/"easy" to understand example. If I was being paid to deliver the answer to this question at a real job,
    I would solve the problem in the most direct way possible that doesn't really flex any fancy memoization and tabulation muscles. 
    
    The problem:
    Given a rod of length N (inches = array position, so int arr[8] = 8 inch rod (yeah..)).
    I digress:
    Identify WHICH CUT(S) offer the most maximally optimized value.
   
    Instead of writing a function to recursively subtract elements using a for loop and range which leads to O(n^2) time at best, 
    Rather we should just measure the distance in value in proportion to the difference in length.
    With some assumptions being made about the rod, of course (more on that below).

    As we can see - the array is 8 inches in length. The full length Rod is worth $20:
    (prices[max(len(prices) - 1)]) = 20 (position [8] on the array gives value of 20)

    Following my algorithm, we should cut the rod at inch 2 (since we'll get $5) and inch 6 (since we'll get $17).
    Ergo, $22 > $20,
    But why?

    My Algorithm is arguably simpler than the traditional Memoization vs Tabulation/Dynamic Programming Recursion examples found on most websites and textbooks.
    
    Since ALL the academic examples assume a linear progression of cut positions to values (the array values are ALWAYS in ascending order of value).
    (This makes sense .. since more size = more better .. or does size really matter?)
    I digress again.

    Ergo again - measuring the distance in inches COMPARED to the increase in value should give us the rational positions
    to thus complete the optimization.

    Hands on Example with code and examples:
    prices = [1,  5,  8,  9,  10,  17,  17,  20]
              ^   ^   ^   ^   ^    ^    ^    ^
    inches = [1,  2,  3,  4,  5,   6,   7,   8]

    The distance in value between inch 1 ($1) and 2 ($5) is +$4
    Compared to only +$3 for inch 2 ($5) to 3 ($8) 
    The distance in value between inch 5 ($10) and 6 ($17) is a whopping +$7
    Compared to only a meager +$1 from inch 4 to 5 ($9 to $10)

    The Cherry on Top? My method is guaranteed to execute in O(n) (linear time and complexity) 
    Because there should only be ONE traversal of the array/list...period. (todo)
    '''

    
    
    #https://chatgpt.com/c/697aa600-f29c-832f-a8e4-13d9b6dee15f
    #chat gpt
    max_diff = float('-inf')
    indices = []

    for i in range(1, len(rod)):
        diff = rod[i] - rod[i - 1]
        if diff > max_diff:
            max_diff = diff
            indices = [i]
        elif diff == max_diff:
            indices.append(i)

    print(indices, max_diff)
    '''
    from dataclasses import dataclass
    
    @dataclass
    class RodPiece:
        cut_piece_inch: int
        cut_piece_value: int


    def slice_rod(prices: list[int], debug: bool = True) -> list[RodPiece]:
        rod_pieces = []
        for cut in range(1, (len(prices)+1)):
            cut_piece = cut
            cut_piece_value = (prices[cut - 1])
            r = RodPiece(cut_piece, cut_piece_value)
            if debug:
                print(f"{cut_piece}:{cut_piece_value}")
            rod_pieces.append(r)  
        return rod_pieces  
    '''
    # TODO
        

   
      