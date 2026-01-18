#pragma once

namespace vislib::core {

template <typename T> constexpr T&& move(T& t) noexcept {
    return static_cast<T&&>(t);
}

template<class T> struct remove_reference {
    typedef T type;
};

template<class T> struct remove_reference<T&> {
    typedef T type;
};

template<class T> struct remove_reference<T&&> {
    typedef T type;
};

template<typename T> inline constexpr T&& forward(typename remove_reference<T>::type& t) noexcept {
    return static_cast<T&&>(t);
}

} //namespace vislib::core