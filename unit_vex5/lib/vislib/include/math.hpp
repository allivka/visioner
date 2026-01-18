#pragma once

#include <math.h>
#include "types.hpp"
#include "containers.hpp"
#include "errordef.hpp"
#include "abstract.hpp"
#include "exception.hpp"

namespace vislib::core {

template <typename T> inline constexpr T absF(const T& x) noexcept(numberNoexcept<T>()) {
    return (x < static_cast<T>(0)) ? (x * static_cast<T>(-1)) : x;
}

template <typename T> inline constexpr T square(const T& x) noexcept(numberNoexcept<T>()) {
    return x * x;
}

template <typename T> inline constexpr char signF(const T& x) noexcept(numberNoexcept<T>()) {
    return (x > static_cast<T>(0)) ? 1 : ((x < static_cast<T>(0)) ? -1 : 0);
}

template <typename T> inline T simpleMul(const T& value, size_t count) noexcept(numberNoexcept<T>()) {
    T buffer = value;
    
    for(size_t i = 1; i < count; i++) {
        buffer += value;
    }
    
    return buffer;
}

template <typename T> inline T simplePow(const T& value, size_t count) noexcept(numberNoexcept<T>()) {
    T buffer = value;
    
    for(size_t i = 1; i < count; i++) {
        buffer *= value;
    }
    
    return buffer;
}

template <typename T> inline constexpr T minF(const T& x, const T& y) noexcept(numberNoexcept<T>()) {
    return (x < y) ? x : y;
}

template <typename T> inline constexpr T maxF(const T& x, const T& y) noexcept(numberNoexcept<T>()) {
    return (x > y) ? x : y;
}

template <typename T> inline constexpr T minEq(const T& x, const T& y) noexcept(numberNoexcept<T>()) {
    return (x <= y) ? x : y;
}

template <typename T> inline constexpr T maxEq(const T& x, const T& y) noexcept(numberNoexcept<T>()) {
    return (x >= y) ? x : y;
}

inline double cosDegrees(double angle) noexcept {
    return cos(angle * M_PI / 180.0);
}

inline double sinDegrees(double angle) noexcept {
    return sin(angle * M_PI / 180.0);
}

template <typename T> inline constexpr T deg2Rad(T angle) noexcept(numberNoexcept<T>()) {
    return angle * M_PI / 180.0;
}

template <typename T> inline constexpr T rad2Deg(T angle) noexcept(numberNoexcept<T>()) {
    return angle * 180.0 / M_PI;
}


template <typename T> class Range {
public:
    T lowest = 0;
    T highest = 0;
    
    bool isInfiniteLow;
    bool isInfiniteHigh;
    
    template<typename D> static D map(const D& x, const D& in_min, const D& in_max, const D& out_min, const D& out_max) noexcept(numberNoexcept<D>()) {
        if (in_max == in_min) {
            return out_min;
        }
        return (x - in_min) * (out_max - out_min) / (in_max - in_min) + out_min;
    }

    template<typename D> static D map(D x, const Range<D>& in, const Range<D>& out) noexcept(numberNoexcept<D>()) {
        return map(x, in.lowest, in.highest, out.lowest, out.highest);
    }

    Range() : isInfiniteLow(true), isInfiniteHigh(true) {};
    Range(const Range&) = default;
    Range(const T& p_lowest, const T& p_highest, bool isInfiniteLow = false, bool isInfiniteHigh = false)
        noexcept(noexcept(T(p_lowest)) && noexcept(T(p_highest)) && numberNoexcept<T>())
        : lowest(p_lowest), highest(p_highest), isInfiniteLow(isInfiniteLow), isInfiniteHigh(isInfiniteHigh) {}

    inline constexpr bool contains(const T& v) const noexcept(numberNoexcept<T>()) {
        return (v >= lowest || isInfiniteLow) && (v <= highest || isInfiniteHigh);
    }

    inline constexpr T restrict(const T& v) const noexcept(numberNoexcept<T>()) {
        
        return  ((v < lowest && isInfiniteLow) || (v > highest && isInfiniteHigh)) ? v :
                ((v < lowest) ? lowest : 
                ((v > highest) ? highest :
                    v));
    }

    inline constexpr T mapValueFromRange(const T& v, Range<T> r) const noexcept(numberNoexcept<T>()) {
        return map(v, r, *this);
    }

    inline constexpr T mapValueToRange(const T& v, Range<T> r) const noexcept(numberNoexcept<T>()) {
        return map(v, *this, r);
    }

};

template <typename T, typename TimeType = T> class Integrator {
protected:
    T integral{};
    TimeType lastTime{};
    
    Range<T> limit{};
    
public:
    
    Integrator() = default;

    Integrator(const TimeType& startTime, const Range<T> limit) noexcept(numberNoexcept<T>() && numberNoexcept<TimeType>())
        : lastTime(startTime), limit(limit) {}
    
    void start(const TimeType& startTime, const Range<T> limit) noexcept(numberNoexcept<T>() && numberNoexcept<TimeType>()) {
        lastTime = startTime;
        this->limit = limit;
    }
    
    Result<T> update(const TimeType& currentTime, const T& value) noexcept(numberNoexcept<T>() && numberNoexcept<TimeType>()) {
        if(currentTime < lastTime) return Error{ErrorCode::invalidArgument, "Current time is less than last operation time. Integrator state wasn't changed"};
        
        T result = integral + value * (currentTime - lastTime);
        
        if (result == integral) return integral;
        
        if(!limit.contains(result)) result = limit.restrict(result);
        
        integral = result;
        lastTime = currentTime;
        
        return integral;
        
    }
    
    T getIntegral() const noexcept {
        return integral;
    }
    
    void setIntegral(const T& value) noexcept(numberNoexcept<T>()) {
        integral = limit.restrict(value);
    }
    
    void resetIntegral() {
        integral = T();
    }
    
    void resetIntegral(const TimeType& currentTime) {
        integral = T();
        lastTime = currentTime;
    }
    
    Range<T> getLimit() const {
        return limit;
    }

};

template <typename T>
class Vector {
protected:
    Array<T> data;

public:
    Vector() = default;
    Vector(const Vector&) = default;
    Vector(Vector&&) noexcept = default;
    Vector& operator=(const Vector&) = default;
    Vector& operator=(Vector&&) noexcept = default;

    explicit Vector(const Array<T>& arr) noexcept(noexcept(T(arr[0]))) : data(arr) {}
    explicit Vector(Array<T>&& arr) noexcept(noexcept(T(arr[0]))) : data(core::move(arr)) {}

    operator Array<T>&() noexcept {
        return data;
    }

    operator const Array<T>&() const noexcept {
        return data;
    }

    Array<T>& raw() noexcept {
        return data;
    }

    const Array<T>& raw() const noexcept {
        return data;
    }

    size_t Size() const noexcept {
        return data.Size();
    }

    T& at(size_t i) noexcept {
        return data[i];
    }

    const T& at(size_t i) const noexcept {
        return data[i];
    }

    T& operator[](size_t i) noexcept {
        return data[i];
    }

    const T& operator[](size_t i) const noexcept {
        return data[i];
    }

    Vector operator+(const Vector& other) const noexcept(numberNoexcept<T>()) {
        Vector temp = *this;
        for (size_t i = 0; i < minF(temp.Size(), other.Size()); i++) {
            temp.at(i) += other.at(i);
        }
        return temp;
    }

    Vector& operator+=(const Vector& other) noexcept(numberNoexcept<T>()) {
        for (size_t i = 0; i < minF(Size(), other.Size()); i++) {
            at(i) += other.at(i);
        }
        return *this;
    }

    Vector operator-(const Vector& other) const noexcept(numberNoexcept<T>()) {
        Vector temp = *this;
        for (size_t i = 0; i < minF(temp.Size(), other.Size()); i++) {
            temp.at(i) -= other.at(i);
        }
        return temp;
    }

    Vector& operator-=(const Vector& other) noexcept(numberNoexcept<T>()) {
        for (size_t i = 0; i < minF(Size(), other.Size()); i++) {
            at(i) -= other.at(i);
        }
        return *this;
    }

    Vector operator*(const T& value) const noexcept(numberNoexcept<T>()) {
        Vector temp = *this;
        for (size_t i = 0; i < temp.Size(); i++) {
            temp.at(i) *= value;
        }
        return temp;
    }

    Vector& operator*=(const T& value) noexcept(numberNoexcept<T>()) {
        for (size_t i = 0; i < Size(); i++) {
            at(i) *= value;
        }
        return *this;
    }

    Vector operator/(const T& value) const noexcept(numberNoexcept<T>()) {
        if (value == T(0)) return *this;
        Vector temp = *this;
        for (size_t i = 0; i < temp.Size(); i++) {
            temp.at(i) /= value;
        }
        return temp;
    }

    Vector& operator/=(const T& value) noexcept(numberNoexcept<T>()) {
        if (value == T(0)) return *this;
        for (size_t i = 0; i < Size(); i++) {
            at(i) /= value;
        }
        return *this;
    }

    Vector operator-() const noexcept(numberNoexcept<T>()) {
        Vector temp = *this;
        for (size_t i = 0; i < temp.Size(); i++) {
            temp.at(i) = -temp.at(i);
        }
        return temp;
    }

    double module() const noexcept(numberNoexcept<T>()) {
        double buffer = 0;
        for (size_t i = 0; i < Size(); i++) {
            buffer += static_cast<double>(at(i)) * static_cast<double>(at(i));
        }
        return sqrt(buffer);
    }

    double dot(const Vector& other) const noexcept(numberNoexcept<T>()) {
        double buffer = 0;
        for (size_t i = 0; i < minF(Size(), other.Size()); i++) {
            buffer += static_cast<double>(at(i)) * static_cast<double>(other.at(i));
        }
        return buffer;
    }

    Vector normal() const noexcept(numberNoexcept<T>()) {
        auto m = module();
        return (m != 0) ? (*this / static_cast<T>(m)) : *this;
    }

    void normalize() noexcept(numberNoexcept<T>()) {
        auto m = module();
        if (m != 0) *this /= static_cast<T>(m);
    }
};

template <typename T>
Vector<T> operator*(const T& val, const Vector<T>& vec) noexcept(numberNoexcept<T>()) {
    return vec * val;
}


} //namespace vislib::core