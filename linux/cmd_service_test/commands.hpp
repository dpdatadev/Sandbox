#include "lib.hpp"

//todo command factory

class Command {
    public:
        string _commandText;
        Command() = delete;
        virtual ~Command() = default;
        virtual string getUuid();
        Command(string uuid);
        Command(string uuid, string commandText);
    private:
        string _uuid;
};