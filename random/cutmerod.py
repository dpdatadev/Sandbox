# Python program to find maximum
# profit from rod of size n 

def cutRod(price):
    n = len(price)
    dp = [0] * (n + 1)

    # Find maximum value for all 
    # rod of length i.
    for i in range(1, n + 1):
        for j in range(1, i + 1):
            dp[i] = max(dp[i], price[j - 1] + dp[i - j])
            print(dp[i])
            print(price[j-1])
            print(dp[i-j])

    return dp[n]

if __name__ == "__main__":
    price = [1, 5, 8, 9, 10, 17, 17, 20]
    print(cutRod(price))