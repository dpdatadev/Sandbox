
#pragma once

#include <stdio.h>
#include <stdlib.h>
#include <assert.h>
#include <vector>
#include <array>
#include <list>
#include <memory>
#include <map>
#include <cstring>
#include <iostream>
#include <bit>

#include "include/httplib.h"

using namespace std;
using namespace httplib;

// Helpers and Macros
#define HIGH 0x1
#define LOW 0x0

#define INPUT 0x0
#define OUTPUT 0x1
#define INPUT_PULLUP 0x2

namespace Math
{
#define PI 3.1415926535897932384626433832795
#define HALF_PI 1.5707963267948966192313216916398
#define TWO_PI 6.283185307179586476925286766559
#define DEG_TO_RAD 0.017453292519943295769236907684886
#define RAD_TO_DEG 57.295779513082320876798154814105
#define EULER 2.718281828459045235360287471352

#define min(a, b) ((a) < (b) ? (a) : (b))
#define max(a, b) ((a) > (b) ? (a) : (b))
#define abs(x) ((x) > 0 ? (x) : -(x))
#define constrain(amt, low, high) ((amt) < (low) ? (low) : ((amt) > (high) ? (high) : (amt)))
#define round(x) ((x) >= 0 ? (long)((x) + 0.5) : (long)((x) - 0.5))
#define radians(deg) ((deg) * DEG_TO_RAD)
#define degrees(rad) ((rad) * RAD_TO_DEG)
#define sq(x) ((x) * (x))
};

namespace Bits
{
#define lowByte(w) ((uint8_t)((w) & 0xff))
#define highByte(w) ((uint8_t)((w) >> 8))

#define bitRead(value, bit) (((value) >> (bit)) & 0x01)
#define bitSet(value, bit) ((value) |= (1UL << (bit)))
#define bitClear(value, bit) ((value) &= ~(1UL << (bit)))
#define bitToggle(value, bit) ((value) ^= (1UL << (bit)))
#define bitWrite(value, bit, bitvalue) ((bitvalue) ? bitSet(value, bit) : bitClear(value, bit))
};

inline bool is_hex(char c, int &v)
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

inline bool from_hex_to_i(const std::string &s, size_t i, size_t cnt,
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

inline std::string from_i_to_hex(size_t n)
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

inline size_t to_utf8(int code, char *buff)
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
inline std::string base64_encode(const std::string &in)
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

inline bool is_file(const std::string &path)
{
    struct stat st;
    return stat(path.c_str(), &st) >= 0 && S_ISREG(st.st_mode);
}

inline bool is_dir(const std::string &path)
{
    struct stat st;
    return stat(path.c_str(), &st) >= 0 && S_ISDIR(st.st_mode);
}

inline bool is_valid_path(const std::string &path)
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
};

