using namespace std;
#include <iostream>
#include <vector>
#include "file_logger.cpp"
#include "singleton.cpp"

#define LOG_FILE_NAME "math_game_log.txt"
#define WELCOME_MESSAGE ":::Welcome to the Elementary MATH GAME!!!:::\n"

// TODO, interactive addition and subtraction game for kids

enum OperationType
{
    ADDITION,
    SUBTRACTION
};


class MathGame : public Singleton<MathGame>
{

private:
    vector<int> nums_to_add;
    vector<int> nums_to_subtract;

    void loadNumsToAdd()
    {
        nums_to_add = {1, 2, 3, 4, 5, 6, 7, 8, 9, 10};
    }
    void loadNumsToSubtract()
    {
        nums_to_subtract = {1, 2, 3, 4, 5, 6, 7, 8, 9, 10};
    }

public:
    void start()
    {
        this->loadNumsToAdd();
        this->loadNumsToSubtract();
        std::cout << "Welcome to the Math Game!" << std::endl;
        // Game logic goes here
        explainRules();
        //displayAdditionQuestion();
        //collectAnswer(ADDITION);
    }
    void explainRules()
    {
        std::cout << "In this game, you will practice addition and subtraction.\n"
                  << std::endl;
        std::cout << "You will be given numbers to add and subtract.\n"
                  << std::endl;
        std::cout << "Try to get as many correct answers as possible!\n"
                  << std::endl;
        std::cout << "Good luck!\n"
                  << std::endl;
        std::cout << "----------------------------------------\n"
                  << std::endl;
    }

    int getRandomInteger(vector<int> &numbers)
    {
        if (numbers.empty())
        {
            throw std::runtime_error("The input vector is empty.");
        }
        int randomIndex = rand() % numbers.size();
        return numbers[randomIndex];
    }

    bool checkAnswer(OperationType opType, int num1, int num2, int userAnswer)
    {
        int correctAnswer;
        if (opType == ADDITION)
        {
            correctAnswer = num1 + num2;
        }
        else if (opType == SUBTRACTION)
        {
            correctAnswer = num1 - num2;
        }
        else
        {
            throw std::invalid_argument("Invalid operation type.");
        }

        return userAnswer == correctAnswer;
    }

    void displayAdditionQuestion()
    {
        if (nums_to_add.size() < 2)
        {
            std::cout << "Not enough numbers to add." << std::endl;
            return;
        }
        int num1 = getRandomInteger(this->nums_to_add);
        int num2 = getRandomInteger(this->nums_to_add);
        std::cout << "What is " << num1 << " + " << num2 << "?" << std::endl;
    }
    void collectAnswer(OperationType opType)
    {
        int answer;
        std::cout << "Please enter your answer: ";
        std::cin >> answer;
        //todo
    }
};


// TODO
int main()
{
    // Setup singleton instance and Logger
    FileLogger &logger = Singleton<FileLogger>::getInstance();
    MathGame &game = Singleton<MathGame>::getInstance();
    game.start();
    cout << WELCOME_MESSAGE << endl;
    logger.writeToFile(LOG_FILE_NAME, WELCOME_MESSAGE);
    return 0;
}