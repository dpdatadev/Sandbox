
#pragma once

#include <stdio.h>
#include <stdlib.h>
#include <assert.h>
#include <sqlite3.h>
#include <uuid/uuid.h>
#include "include/httplib.h"

#include <vector>
#include <array>
#include <list>
#include <memory>
#include <shared_mutex>
#include <optional>
#include <map>
#include <cstring>
#include <iostream>
#include <fstream>
#include <algorithm>
#include <bit>

using namespace std;
using namespace httplib;

// Helpers and Macros
#define LEFT(a, b) Services::extract_left_chars(a, b)
#define NEWUUID Services::nuuid()
#define UUIDSIZE 7

#define HIGH 0x1
#define LOW 0x0

#define INPUT 0x0
#define OUTPUT 0x1
#define INPUT_PULLUP 0x2

// see ChatGPT notes
namespace Math
{
#define PI 3.1415926535897932384626433832795
#define HALF_PI 1.5707963267948966192313216916398
#define TWO_PI 6.283185307179586476925286766559
#define DEG_TO_RAD 0.017453292519943295769236907684886
#define RAD_TO_DEG 57.295779513082320876798154814105
#define EULER 2.718281828459045235360287471352

// m prefix to avoid clash with std
#define m_min(a, b) ((a) < (b) ? (a) : (b))
#define m_max(a, b) ((a) > (b) ? (a) : (b))
#define m_abs(x) ((x) > 0 ? (x) : -(x))
#define m_constrain(amt, low, high) ((amt) < (low) ? (low) : ((amt) > (high) ? (high) : (amt)))
#define m_round(x) ((x) >= 0 ? (long)((x) + 0.5) : (long)((x) - 0.5))
#define m_radians(deg) ((deg) * DEG_TO_RAD)
#define m_degrees(rad) ((rad) * RAD_TO_DEG)
#define m_sq(x) ((x) * (x))
};

// Helpers for files and bytes
namespace Bits
{
#define lowByte(w) ((uint8_t)((w) & 0xff))
#define highByte(w) ((uint8_t)((w) >> 8))
#define bitRead(value, bit) (((value) >> (bit)) & 0x01)
#define bitSet(value, bit) ((value) |= (1UL << (bit)))
#define bitClear(value, bit) ((value) &= ~(1UL << (bit)))
#define bitToggle(value, bit) ((value) ^= (1UL << (bit)))
#define bitWrite(value, bit, bitvalue) ((bitvalue) ? bitSet(value, bit) : bitClear(value, bit))
#define multipleBitsSet(i) (((i) & ((i) - 1)) != 0)
};

static inline bool is_hex(char c, int &v)
{
    if (0x20 <= c && isdigit(c))
    {
        v = c - '0';
        return true;
    }
    else if ('A' <= c && c <= 'F')
    {
        v = c - 'A' + 10;
        return true;
    }
    else if ('a' <= c && c <= 'f')
    {
        v = c - 'a' + 10;
        return true;
    }
    return false;
}

static inline void hexdump(const void *data, size_t len)
{
    const unsigned char *p = data;

    for (size_t i = 0; i < len; ++i)
    {
        printf("%02X ", p[i]);

        if ((i + 1) % 16 == 0)
            printf("\n");
    }

    printf("\n");
}

static inline bool from_hex_to_i(const std::string &s, size_t i, size_t cnt,
                                 int &val)
{
    if (i >= s.size())
    {
        return false;
    }

    val = 0;
    for (; cnt; i++, cnt--)
    {
        if (!s[i])
        {
            return false;
        }
        int v = 0;
        if (is_hex(s[i], v))
        {
            val = val * 16 + v;
        }
        else
        {
            return false;
        }
    }
    return true;
}

static inline std::string from_i_to_hex(size_t n)
{
    const char *charset = "0123456789abcdef";
    std::string ret;
    do
    {
        ret = charset[n & 15] + ret;
        n >>= 4;
    } while (n > 0);
    return ret;
}

static inline size_t to_utf8(int code, char *buff)
{
    if (code < 0x0080)
    {
        buff[0] = (code & 0x7F);
        return 1;
    }
    else if (code < 0x0800)
    {
        buff[0] = (0xC0 | ((code >> 6) & 0x1F));
        buff[1] = (0x80 | (code & 0x3F));
        return 2;
    }
    else if (code < 0xD800)
    {
        buff[0] = (0xE0 | ((code >> 12) & 0xF));
        buff[1] = (0x80 | ((code >> 6) & 0x3F));
        buff[2] = (0x80 | (code & 0x3F));
        return 3;
    }
    else if (code < 0xE000)
    { // D800 - DFFF is invalid...
        return 0;
    }
    else if (code < 0x10000)
    {
        buff[0] = (0xE0 | ((code >> 12) & 0xF));
        buff[1] = (0x80 | ((code >> 6) & 0x3F));
        buff[2] = (0x80 | (code & 0x3F));
        return 3;
    }
    else if (code < 0x110000)
    {
        buff[0] = (0xF0 | ((code >> 18) & 0x7));
        buff[1] = (0x80 | ((code >> 12) & 0x3F));
        buff[2] = (0x80 | ((code >> 6) & 0x3F));
        buff[3] = (0x80 | (code & 0x3F));
        return 4;
    }

    // NOTREACHED
    return 0;
}

// NOTE: This code came up with the following stackoverflow post:
// https://stackoverflow.com/questions/180947/base64-decode-snippet-in-c
static inline std::string base64_encode(const std::string &in)
{
    static const auto lookup =
        "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/";

    std::string out;
    out.reserve(in.size());

    int val = 0;
    int valb = -6;

    for (uint8_t c : in)
    {
        val = (val << 8) + c;
        valb += 8;
        while (valb >= 0)
        {
            out.push_back(lookup[(val >> valb) & 0x3F]);
            valb -= 6;
        }
    }

    if (valb > -6)
    {
        out.push_back(lookup[((val << 8) >> (valb + 8)) & 0x3F]);
    }

    while (out.size() % 4)
    {
        out.push_back('=');
    }

    return out;
}

static inline bool is_file(const std::string &path)
{
    struct stat st;
    return stat(path.c_str(), &st) >= 0 && S_ISREG(st.st_mode);
}

static inline bool is_dir(const std::string &path)
{
    struct stat st;
    return stat(path.c_str(), &st) >= 0 && S_ISDIR(st.st_mode);
}

static inline bool is_valid_path(const std::string &path)
{
    size_t level = 0;
    size_t i = 0;

    // Skip slash
    while (i < path.size() && path[i] == '/')
    {
        i++;
    }

    while (i < path.size())
    {
        // Read component
        auto beg = i;
        while (i < path.size() && path[i] != '/')
        {
            i++;
        }

        auto len = i - beg;
        assert(len > 0);

        if (!path.compare(beg, len, "."))
        {
            ;
        }
        else if (!path.compare(beg, len, ".."))
        {
            if (level == 0)
            {
                return false;
            }
            level--;
        }
        else
        {
            level++;
        }

        // Skip slash
        while (i < path.size() && path[i] == '/')
        {
            i++;
        }
    }

    return true;
}

// Manual memory management
namespace c_mem
{
    static void *safe_malloc(size_t size)
    {
        void *ptr = malloc(size);
        assert(ptr != NULL && "Memory allocation failed");
        return ptr;
    }

    static void safe_free(void *ptr)
    {
        if (ptr != NULL)
        {
            free(ptr);
        }
    }

    static bool is_aligned(void *ptr, size_t alignment)
    {
        return ((uintptr_t)ptr % alignment) == 0;
    }

    static void *aligned_malloc(size_t size, size_t alignment)
    {
        void *ptr = NULL;
        int result = posix_memalign(&ptr, alignment, size);
        assert(result == 0 && "Aligned memory allocation failed");
        return ptr;
    }
};

/* Utility/Helper functions */
namespace Services
{
#ifndef _t_factory_h_
#define _t_factory_h_
    // Add method to get singleton factory instance in implementation
    // https://chatgpt.com/c/6986b331-7e80-8325-9ffb-8c51562e1709
    // https://chatgpt.com/c/6986b331-7e80-8325-9ffb-8c51562e1709
    template <class T>
    class base_creator
    {
    public:
        virtual ~base_creator() {};
        virtual std::unique_ptr<T> create() = 0;
    };

    template <class derived_type, class base_type>
    class derived_creator : public base_creator<base_type>
    {
    public:
        std::unique_ptr<base_type> create()
        {
            return std::make_unique<derived_type>();
        }
    };

    // TODO, implement RAII changes - https://chatgpt.com/c/6986b331-7e80-8325-9ffb-8c51562e1709
    // https://chatgpt.com/c/6986b331-7e80-8325-9ffb-8c51562e1709

    template <class _key, class base_type>
    class thread_safe_factory
    {
    public:
        void register_type(_key &id, std::unique_ptr<base_creator<base_type>> _creator)
        {
            std::unique_lock lock(_mutex);
            _function_map[id].insert_or_assign(id, std::move(_creator)); // or emplace, right now accept overwrites
        }

        // Create an instance (READ LOCK)
        [[nodiscard]]
        std::unique_ptr<base_type> create(const _key &id)
        {
            std::shared_lock lock(_mutex);

            auto it = _function_map.find(id);

            if (it == _function_map.end())
                return nullptr;

            return it->second->create();
        }

        // Optional: query existence
        bool contains(const _key &id) const
        {
            std::shared_lock lock(_mutex);
            return _function_map.contains(id);
        }

    private:
        std::unordered_map<_key, std::unique_ptr<base_creator<base_type>>> _function_map;
        mutable std::shared_mutex _mutex;
    };

#endif /* defined _t_factory_h_ */

    // Based on ANSI SQL "LEFT() function"
    static inline std::string extract_left_chars(const std::string &original_str, size_t n)
    {
        // Ensure n does not exceed the actual string length to avoid exceptions
        size_t length_to_extract = std::min(n, original_str.length());
        // The second parameter to substr is the number of characters to include
        return original_str.substr(0, length_to_extract);
    }

    // Generate new UUID
    static std::string nuuid()
    {

        uuid_t binuuid;
        char uuid_str[37];

        uuid_generate_random(binuuid);
        uuid_unparse_lower(binuuid, uuid_str);

        return std::string(uuid_str);
    }

    // Use std::system
    static bool exec(const char *cmd,
                     char *result,
                     size_t result_size)
    {
        if (!cmd || !result || result_size == 0)
            return false;

        FILE *fp = popen(cmd, "r");
        if (!fp)
            return false;

        if (!fgets(result,
                   result_size,
                   fp))
        {
            pclose(fp);
            return false;
        }

        pclose(fp);
        return true;
    }
};
