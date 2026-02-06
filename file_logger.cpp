using namespace std;
#include <fstream>
#include <mutex>
#include <stdexcept>
#include <ostream>
#include "singleton.cpp"
#pragma once

using std::lock_guard;
using std::mutex;
using std::ofstream;
using std::runtime_error;
using std::string;

class FileLogger : public Singleton<FileLogger>
{
    private:
    // mutex is to protect access to file (which is shared across threads).
    static mutex m;
    // Lock mutex before accessing file.
    lock_guard<mutex> lock;
    FileLogger() = default;
    ~FileLogger() = default;
public:
    FileLogger() : lock(FileLogger::m) {}
    static void writeToFile(const string &fileName = "log.txt", const string &message)
    {

        // Try to open file.
        ofstream f{fileName};
        if (!f.is_open())
        {
            throw runtime_error("unable to open file");
        }

        // Write message to file.
        f.write(message.c_str(), message.size());

        // file will be closed first when leaving scope (regardless of exception)
        // mutex will be unlocked second (from lock destructor) when leaving scope
        // (regardless of exception).

        if (f.is_open())
        {
            f.close();
        }
    }
};