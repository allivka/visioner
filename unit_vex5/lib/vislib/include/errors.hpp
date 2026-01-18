#pragma once

#include "containers.hpp"

namespace vislib::core {

class Error {
public:
    ErrorCode errcode = ErrorCode::success;
    String msg = "";

    operator ErrorCode() const noexcept { return errcode; }

    operator String() const noexcept { return msg; }

    operator bool() const noexcept { return errcode != ErrorCode::success; }

    bool isError() const noexcept { return errcode != ErrorCode::success; }

    bool operator==(const Error& err) const noexcept {
        return errcode == err.errcode;
    }

    Error() noexcept = default;

    Error(ErrorCode code, const String& p_msg = "") noexcept : errcode(code), msg(p_msg) {}
};

template <typename T> class Result {
protected:

    T value{};
    Error err{};

public:

    Result(const T& v) noexcept(noexcept(T(v))) : value(v) { }

    Result(const Error& e) noexcept : err(e) { }

    constexpr bool isError() const noexcept { return err.isError(); }
    constexpr bool isOk() const noexcept { return !err.isError(); }

    bool getValue(T& out) const noexcept(noexcept(out = value)) {
        if (isOk()) {
            out = value;
            return true;
        }
        return false;
    }

    bool getError(Error& out) const noexcept(noexcept(out = err)) {
        if (isError()) {
            out = err;
            return true;
        }
        return false;
    }

    operator T() const noexcept(noexcept(T(value))) { return value; }
    operator Error() const noexcept { return err; }

    constexpr operator bool() const noexcept { return isError(); }

    T operator()() const noexcept(noexcept(T(value))) { return value; }

    T Value() const noexcept(noexcept(T(value))) { return value; }

    Error error() const noexcept { return err; }
};

template <typename T>
class Result<T&> : public Result<T*> {
public:
    Result(T& v) noexcept : Result<T*>(&v) {}
    Result(Error e) noexcept : Result<T*>(e) {}

    T& operator()() const noexcept { return *(this->value); }
    T& Value() const noexcept { return *(this->value); }
};


} //namespace vislib::core
