#pragma once

namespace vislib::core {

template <typename T> T&& declval() noexcept;

template <typename T> inline constexpr bool basicArithmeticNoexcept() noexcept {
    return noexcept(declval<T>() + declval<T>()) &&
        noexcept(declval<T>() - declval<T>()) &&
        noexcept(declval<T>() * declval<T>()) &&
        noexcept(declval<T>() / declval<T>()) &&
        noexcept(-declval<T>());
}

template <typename T> inline constexpr bool assignableNoexcept() noexcept {
    return noexcept(declval<T&>() = declval<T>());
}

template <typename T> inline constexpr bool arithmeticAssignableNoexcept() noexcept {
    return assignableNoexcept<T>() &&
        noexcept(declval<T&>() += declval<T>()) &&
        noexcept(declval<T&>() -= declval<T>()) &&
        noexcept(declval<T&>() *= declval<T>()) &&
        noexcept(declval<T&>() /= declval<T>());
}

template <typename T> inline constexpr bool comparableNoexcept() noexcept {
    return noexcept(declval<T>() < declval<T>()) &&
        noexcept(declval<T>() > declval<T>()) &&
        noexcept(declval<T>() >= declval<T>()) &&
        noexcept(declval<T>() <= declval<T>()) &&
        noexcept(declval<T>() == declval<T>()) &&
        noexcept(declval<T>() != declval<T>());
}

template <typename T> inline constexpr bool arithmeticNoexcept() noexcept {
    return basicArithmeticNoexcept<T>() && arithmeticAssignableNoexcept<T>();
}

template <typename T> inline constexpr bool numberNoexcept() noexcept {
    return arithmeticNoexcept<T>() && comparableNoexcept<T>();
}

template <typename T, typename U, typename... Rest> inline constexpr bool basicArithmeticNoexcept() noexcept {
    return basicArithmeticNoexcept<T>() && basicArithmeticNoexcept<U, Rest...>();
}

template <typename T, typename U, typename... Rest> inline constexpr bool arithmeticAssignableNoexcept() noexcept {
    return arithmeticAssignableNoexcept<T>() && arithmeticAssignableNoexcept<U, Rest...>();
}

template <typename T, typename U, typename... Rest> inline constexpr bool arithmeticNoexcept() noexcept {
    return arithmeticNoexcept<T>() && arithmeticNoexcept<U, Rest...>();
}

template <typename T, typename U, typename... Rest> inline constexpr bool numberNoexcept() noexcept {
    return numberNoexcept<T>() && numberNoexcept<U, Rest...>();
}

} // namespace vislib::core
