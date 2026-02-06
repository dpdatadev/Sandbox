#include "lib.hpp"

// see ChatGPT/notes

namespace Services
{
    // TODO - prototype in python https://chatgpt.com/c/697bf30b-02c4-8328-99ae-c144bd271709

    // TODO - build database to store commands

    // TODO - cmd function was used to generate UUID for various processes

    static int cmd(const char *cmd, char *result) // run a linux command and return the initial result
    {
        FILE *fp;
        char buf[64];
        int res;

        result[0] = 0;
        fp = popen(cmd, "r");
        if (fp == NULL)
            return 0;
        if (fgets(buf, sizeof(buf) - 1, fp) == NULL)
        {
            pclose(fp);
            return 0;
        }
        else
        {
            pclose(fp);
            res = sscanf(buf, "%s", result);
            if (res != 1)
                return 0;
            return 1;
        }
    }
};
