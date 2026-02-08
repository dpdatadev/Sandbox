#include "lib.hpp"

// TODO, learning project 2/26
//Not quite sure what we're building(?)
// Command Framework
// Input → Scrubber → Factory → Command → Output
// https://chatgpt.com/c/6986b331-7e80-8325-9ffb-8c51562e1709

namespace Commands
{
    enum class CommandCategory
    {
        NIL = -1,
        TEXT = 0,
        WEB = 1,
        DATA = 2,
        OTHER = 3
    };

    struct CommandOutput
    {
        std::string uuid;
        std::string text;
        CommandCategory category;
    };

    class Command
    {
    public:
        Command() = delete;

        Command(std::string uuid,
                std::string commandText)
            : _uuid(std::move(uuid)),
              _commandText(std::move(commandText))
        {
        }

        virtual ~Command() = default;

        [[nodiscard]]
        virtual std::unique_ptr<CommandOutput>
        process() const = 0;

        const std::string &
        getUuid() const noexcept
        {
            return _uuid;
        }

        const std::string &
        getCommandText() const noexcept
        {
            return _commandText;
        }

    protected:
        std::string _uuid;
        std::string _commandText;
    };

    class CommandScrubber
    {
    public:
        static bool validate(
            const std::string &id,
            const std::string &text)
        {
            if (text.size() > 1024)
                return false;

            // forbid shell metacharacters
            if (text.contains(";") ||
                text.contains("&&") ||
                text.contains("sudo") ||
                text.contains("rm"))
                return false;

            return true;
        }
    };

    class TextCommand : public Command
    {
    public:
        using Command::Command; // inherit ctor

        std::unique_ptr<CommandOutput>
        process() const override
        {
            auto output =
                std::make_unique<CommandOutput>();

            output->uuid = _uuid;
            output->category =
                CommandCategory::TEXT;
            // TODO CommandExecute service
            //  Test processing logic:
            output->text =
                "Processed TEXT command: " +
                _commandText;

            return output;
        }
    };
};
