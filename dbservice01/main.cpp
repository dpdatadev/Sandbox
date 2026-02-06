//#include "include/httplib.h"
#include <sqlite3.h>
#include <iostream>

int db_test() {
    sqlite3* db;
    int rc = sqlite3_open(":memory:", &db);
    if (rc) {
        std::cerr << "Can't open database: " << sqlite3_errmsg(db) << std::endl;
        return rc;
    } else {
        std::cout << "Opened database successfully" << std::endl;
    }
    sqlite3_close(db);
    return 0;
}

int main() {
   // TODO
   
}