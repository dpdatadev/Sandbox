#include "engine.hpp"

// TODO, command builder
// TODO, command server
// TODO, file/data logger

// 2/5

// ideas:
//  TODO, service that runs commands and returns the output
//  Pico being used to monitor certain network activity
//  But the Pico runs MicroPython and has limited storage and capability
//  So the Pico can receive UDP messages with pre-programmed commands to periodically peform actions
//  The server/UDP messages serve as a scheduler and the actual data is managed on the server (Pi 5)

#define UPTIME_SIZE 64
#define NUUID_SIZE 64
#define SYSINFO_SIZE 256
#define IP_LIST_SIZE 512

// Testing
int main(void)
{

        Server svr;

        svr.Get("/info", [](const Request &req, Response &res)
                { 
                // Buffers to hold command output
                char nuuid[NUUID_SIZE];
                char sysinfo[SYSINFO_SIZE];
                char uptime[UPTIME_SIZE];
                

                const char sysinfostring[] = "uname -a";
                const char nuuidstring[] = "uuid | cut -f1 -d-";
                const char uptimestring[] = "uptime";
                

                Services::cmd(sysinfostring, sysinfo);
                Services::cmd(nuuidstring, nuuid);
                Services::cmd(uptimestring, uptime);

                strcat(nuuid, sysinfo);
                const char* outputReport = strcat(nuuid, uptime);
                static std::string formattedOutputReport = outputReport;

                puts("Operating System:");
                printf("%s\n", sysinfo);
                puts("New UUID:");
                printf("%s\n", nuuid);
                sleep(1);
                res.set_header("Access-Control-Allow-Origin", "*");         
                res.set_content(formattedOutputReport, "text/plain"); });

        svr.Get("/api/ips", [](const Request &req, Response &res)
                {
                char ipList[IP_LIST_SIZE];
                const char refreshIPList[] = "source ~/Documents/Code/Scripts/linux/getips.sh";

                Services::cmd(refreshIPList, ipList);

                res.set_header("Access-Control-Allow-Origin", "*");         
                res.set_content(ipList, "text/plain"); });

        svr.listen("0.0.0.0", 5000);
}