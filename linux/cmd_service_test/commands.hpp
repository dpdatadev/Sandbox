#include "lib.hpp"

// TODO, learning project 2/26

namespace Commands
{
    // TODO, implement a Command Scrubber that can sanitize and prevent unwanted input from being executed
    // or only registered Commands can be executed - hijacked control can never happen (in theory?)
    // DATA Access Layer (Sqlite)
    // TODO

    enum CommandCategory
    {
        NIL = -1,
        TEXT = 0,
        WEB = 1,
        DATA = 2,
        OTHER = 3
    };

    struct CommandOutput
    {
        string uuid;
        string text;
        CommandCategory category;
    };

    // TODO, make sure only the CommandFactory can create Commands:
    // read the discussion .. https://stackoverflow.com/questions/56994499/force-class-construction-exclusively-inside-factory
    class Command
    {
    public:
        Command() = delete;
        string getUuid() { return this->_uuid; }
        Command(string uuid) { this->_uuid = uuid; }
        Command(string uuid, string commandText)
        {
            this->_uuid = uuid;
            this->_commandText = commandText;
        }
        virtual ~Command() = default;
        // How certain commands are processed are determined in derived class, which are instantiated by the CommandFactory
        virtual unique_ptr<CommandOutput> process() const = 0;

    private:
        string _uuid;
        string _commandText;
    };

    class TextCommand : public Command
    {
    public:
        TextCommand() = delete;
        ~TextCommand() = default;
        TextCommand(string uuid) : Command(uuid) {}
        TextCommand(string uuid, string commandText) : Command(uuid, commandText) {}
        unique_ptr<CommandOutput> process() const override
        {
            // TODO
            // call uuid service:
            size_t uuid_size(7);
            string new_uuid = LEFT(NEWUUID, uuid_size);
            // CommandCategory::TEXT;
            // string output = system();//todo

            return make_unique<CommandOutput>();
        }
    };
};
