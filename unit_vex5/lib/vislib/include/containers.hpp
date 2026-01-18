#pragma once

#include "errordef.hpp"

namespace vislib::core {

template<typename T, size_t SIZE> class DefinedArray {
private:
    T data[SIZE];

public:
    DefinedArray() noexcept {
        for(size_t i = 0; i < SIZE; i++) data[i] = T();
    }

    template<size_t N> DefinedArray(const T (&p_data)[N]) noexcept {
        static_assert(N == SIZE, "Array size must match DefinedArray size");
        for(size_t i = 0; i < SIZE; i++) data[i] = p_data[i];
    }

    DefinedArray(const T (&p_data)[SIZE]) noexcept {
        for(size_t i = 0; i < SIZE; i++) data[i] = p_data[i];
    }

    DefinedArray(const DefinedArray<T, SIZE>& other) noexcept {
        for(size_t i = 0; i < SIZE; i++) data[i] = other.data[i];
    }

    DefinedArray(DefinedArray<T, SIZE>&& other) = default;
    DefinedArray<T, SIZE>& operator=(DefinedArray<T, SIZE>&& other) = default;

    DefinedArray<T, SIZE>& operator=(const DefinedArray<T, SIZE>& other) noexcept {
        if (this != &other) {
            for(size_t i = 0; i < SIZE; i++) data[i] = other.data[i];
        }
        return *this;
    }

    inline T& operator[](size_t index) noexcept {
        return data[index];
    }
    
    inline const T& operator[](size_t index) const noexcept {
        return data[index];
    }

    [[nodiscard]] Result<T&> at(size_t index) noexcept {
        if (index >= SIZE) {
            return {ErrorCode::indexOutOfRange, "Index out of range in 'defined array' element access"};
        }
        return data[index];
    }

    [[nodiscard]] Result<const T&> at(size_t index) const noexcept {
        if (index >= SIZE) {
            return {ErrorCode::indexOutOfRange, "Index out of range in 'defined array' element access"};
        }
        return data[index];
    }

    static constexpr size_t size() noexcept {
        return SIZE;
    }
    
    static inline constexpr bool empty() noexcept {
        return SIZE == 0;
    }
};

template<typename T> class Array {
protected:
    size_t size = 0;
    T *data = nullptr;

public:
    Array() = default;

    Array(const size_t p_size) noexcept : size(p_size), data(new T[size]) { }

    template<size_t N> explicit Array(const T (&p_data)[N]) noexcept : size(N), data(new T[size]) {
        for(size_t i = 0; i < size; i++) data[i] = p_data[i];
    }

    Array(const T p_data[], size_t p_size) noexcept : size(p_size), data(new T[size]) {
        for(size_t i = 0; i < size; i++) data[i] = p_data[i];
    }

    Array(const Array<T>& other) noexcept : size(0), data(nullptr) {
        if (other.size > 0 && other.data != nullptr) {
            data = new T[other.size];
            if (data) {
                size = other.size;
                for(size_t i = 0; i < size; i++) data[i] = other.data[i];
            }
        }
    }

    Array(Array<T>&& other) noexcept : size(other.size), data(other.data) {
        other.size = 0;
        other.data = nullptr;
    }

    ~Array() noexcept {
        delete [] data;
    }

    Array<T>& operator=(const Array<T>& other) noexcept {
        if (this == &other) return *this;
        T* newData = nullptr;
        size_t newSize = 0;
        if (other.size > 0 && other.data != nullptr) {
            newData = new T[other.size];
            if (newData) {
                newSize = other.size;
                for(size_t i = 0; i < newSize; i++) newData[i] = other.data[i];
            }
        }
        delete[] data;
        data = newData;
        size = newSize;
        return *this;
    }

    Array<T>& operator=(Array<T>&& other) noexcept {
        if (this != &other) {
            delete[] data;
            data = other.data;
            size = other.size;
            other.data = nullptr;
            other.size = 0;
        }
        return *this;
    }

    bool operator==(const Array<T>& other) const noexcept {
        if (size != other.size) return false;
        if(data == other.data && data == nullptr) return true;
        if((data == nullptr || other.data == nullptr) && data != other.data) return false;
        for(size_t i = 0; i < size; i++) {
            if (data[i] != other.data[i]) return false;
        }
        return true;
    }

    inline T& operator[](size_t index) noexcept {
        return data[index];
    }
    
    inline const T& operator[](size_t index) const noexcept {
        return data[index];
    }

    [[nodiscard]] Result<T&> at(size_t index) noexcept {
        if (data == nullptr) {
            return Error(ErrorCode::emptyArray, "could not access data of an empty array");
        }
        if (index >= size) {
            return Error(ErrorCode::indexOutOfRange, "index out of range in array element access");
        }
        return data[index];
    }

    [[nodiscard]] Result<const T&> at(size_t index) const noexcept {
        if (data == nullptr) {
            return Error(ErrorCode::emptyArray, "could not access data of an empty array");
        }
        if (index >= size) {
            return Error(ErrorCode::indexOutOfRange, "index out of range in array element access");
        }
        return data[index];
    }

    inline constexpr size_t Size() const noexcept {
        return size;
    }
    
    inline constexpr bool empty() const noexcept {
        return size == 0;
    }
    
    void clear() noexcept {
        delete [] data;
        size = 0;
        data = nullptr;
    }

    inline T* Data() noexcept {
        return data;
    }
    
    inline const T* Data() const noexcept {
        return data;
    }

    Array<T> operator+(const Array<T>& other) const noexcept {
        Array<T> newArr;
        newArr.size = size + other.size;
        newArr.data = new T[newArr.size];
        for(size_t i = 0; i < size; i++) newArr.data[i] = data[i];
        for(size_t i = 0; i < other.size; i++) newArr.data[size + i] = other.data[i];
        return newArr;
    }
};

class String : public Array<char> {
private:
    static size_t c_strlen(const char* str) noexcept {
        if (!str) return 0;
        const char* ptr = str;
        while (*ptr) ptr++;
        return static_cast<size_t>(ptr - str);
    }

    size_t char_capacity() const noexcept {
        const size_t s = this->Size();
        return s == 0 ? 0 : s - 1;
    }

    void ensure_null_terminated() noexcept {
        if (this->Size() == 0) return;
        this->Data()[ this->Size() - 1 ] = '\0';
    }

public:
    String() noexcept : Array<char>(1) {
        this->Data()[0] = '\0';
    }

    String(const char* cstr) noexcept : Array<char>(c_strlen(cstr) + 1)  {
        size_t len = c_strlen(cstr);
        size_t cap = char_capacity();
        if (len > cap) len = cap;
        for (size_t i = 0; i < len; i++) this->Data()[i] = cstr[i];
        ensure_null_terminated();
    }

    String(const char* cstr, size_t length) noexcept : Array<char>(length + 1) {
        size_t cap = char_capacity();
        size_t copy_len = length <= cap ? length : cap;
        for (size_t i = 0; i < copy_len; i++) this->Data()[i] = cstr[i];
        ensure_null_terminated();
    }

    String(size_t count, char ch) noexcept : Array<char>(count + 1) {
        size_t cap = char_capacity();
        size_t n = count <= cap ? count : cap;
        for (size_t i = 0; i < n; i++) this->Data()[i] = ch;
        ensure_null_terminated();
    }

    using Array<char>::Array;
    using Array<char>::operator=;

    const char* c_str() const noexcept {
        if (this->Size() == 0) return "";
        return this->Data();
    }

    size_t length() const noexcept {
        if (this->Size() == 0) return 0;
        const char* d = this->Data();
        for (size_t i = 0; i < char_capacity(); i++) {
            if (d[i] == '\0') return i;
        }
        return char_capacity();
    }

    String operator+(const String& other) const noexcept {
        size_t l1 = this->length();
        size_t l2 = other.length();
        size_t out_len = l1 + l2;
        String result(out_len + 1, '\0');
        for (size_t i = 0; i < l1 && i < result.char_capacity(); i++)
            result.Data()[i] = this->Data()[i];
        size_t pos = l1;
        for (size_t j = 0; j < l2 && pos < result.char_capacity(); j++, pos++)
            result.Data()[pos] = other.Data()[j];
        result.ensure_null_terminated();
        return result;
    }

    String operator+(const char* other) const {
        String rhs(other);
        return *this + rhs;
    }

    String& operator+=(const String& other) noexcept {
        size_t l1 = this->length();
        size_t l2 = other.length();
        size_t cap = char_capacity();
        size_t copy_len = (l1 + l2 <= cap) ? l2 : (cap > l1 ? cap - l1 : 0);
        for (size_t i = 0; i < copy_len; i++) this->Data()[l1 + i] = other.Data()[i];
        ensure_null_terminated();
        return *this;
    }

    String& operator+=(const char* other) noexcept {
        String rhs(other);
        *this += rhs;
        return *this;
    }

    bool operator==(const String& other) const noexcept {
        size_t l1 = length();
        size_t l2 = other.length();
        if (l1 != l2) return false;
        for (size_t i = 0; i < l1; i++) if (this->Data()[i] != other.Data()[i]) return false;
        return true;
    }

    bool operator==(const char* other) const noexcept {
        size_t other_len = c_strlen(other);
        size_t l = length();
        if (l != other_len) return false;
        for (size_t i = 0; i < l; i++) if (this->Data()[i] != other[i]) return false;
        return true;
    }

    inline bool operator!=(const String& other) const noexcept {
        return !(*this == other);
    }
    
    inline bool operator!=(const char* other) const noexcept {
        return !(*this == other);
    }
};

inline String operator+(const char* str1, const String& str2) noexcept {
    String temp(str1);
    return temp + str2;
}

inline bool operator==(const char* lhs, const String& rhs) noexcept {
    return rhs == lhs;
}

inline bool operator!=(const char* lhs, const String& rhs) noexcept {
    return rhs != lhs;
}

inline String to_string(const size_t value) noexcept {
    if (value == 0) return {"0"};
    char buf[32];
    size_t idx = 0;
    unsigned long long v = value;
    while (v > 0 && idx < sizeof(buf)-1) {
        buf[idx++] = static_cast<char>('0' + (v % 10));
        v /= 10;
    }
    String s(idx + 1, '\0');
    for (size_t i = 0; i < idx; i++) s[i] = buf[idx - 1 - i];
    s[idx] = '\0';
    return s;
}

inline String to_string(const long long value) noexcept {
    if (value < 0) {
        const auto uv = static_cast<size_t>(-value);
        String s = String("-") + to_string(uv);
        return s;
    }
    return to_string(static_cast<size_t>(value));
}

} //namespace vislib::core
