#pragma once

#include <assert.h>

namespace vislib::core {

template <typename ReturnType = void, typename... ArgumentTypes>
class CallableInterface {
public:

    virtual ReturnType operator()(ArgumentTypes...) = 0;
    virtual CallableInterface* clone() const = 0;
    virtual ~CallableInterface() = default;
};

template <typename CallableType, typename ReturnType, typename... ArgumentTypes>
class CallableWrapper : public CallableInterface<ReturnType, ArgumentTypes...> {
    CallableType callable;

public:
    CallableWrapper(const CallableType& c) : callable(c) {}
    CallableWrapper(CallableType&& c) : callable(static_cast<CallableType&&>(c)) {}

    ReturnType operator()(ArgumentTypes... args) override {
        return callable(args...);
    }

    CallableInterface<ReturnType, ArgumentTypes...>* clone() const override {
        return new CallableWrapper<CallableType, ReturnType, ArgumentTypes...>(callable);
    }
};

template <typename ReturnType = void, typename... ArgumentTypes>
class Callable {
    CallableInterface<ReturnType, ArgumentTypes...>* impl = nullptr;

public:
    Callable() = default;

    Callable(const Callable& other) {
        if (other.impl == nullptr) {
            impl = nullptr;
        } else {
            impl = other.impl->clone();
        }
    }

    Callable& operator=(const Callable& other) {
        if (this == &other) return *this;
        if (impl) delete impl;
        impl = other.impl ? other.impl->clone() : nullptr;
        return *this;
    }

    Callable(Callable&& other) noexcept : impl(other.impl) {
        other.impl = nullptr;
    }

    Callable& operator=(Callable&& other) noexcept {
        if (this == &other) return *this;
        if (impl) delete impl;
        impl = other.impl;
        other.impl = nullptr;
        return *this;
    }

    template<typename C> Callable(C c) {
        impl = new CallableWrapper<C, ReturnType, ArgumentTypes...>(static_cast<C&&>(c));
    }

    ReturnType operator()(ArgumentTypes... args) {
        // static_assert(impl != nullptr, "Calling a null Callable object");
        assert(impl != nullptr);
        return (*impl)(args...);
    }

    ReturnType execute(ArgumentTypes... args) {
        // static_assert(impl != nullptr, "Calling a null Callable object");
        assert(impl != nullptr);
        return (*impl)(args...);
    }

    bool isValid() const noexcept {
        return impl != nullptr;
    }

    void reset() {
        delete impl;
        impl = nullptr;
    }

    ~Callable() {
        if (impl == nullptr) return;
        delete impl;
    }
};

template <typename ReturnType> class Callable<ReturnType> {
    CallableInterface<ReturnType>* impl = nullptr;

public:
    Callable() = default;

    Callable(const Callable& other) {
        if (other.impl == nullptr) {
            impl = nullptr;
        } else {
            impl = other.impl->clone();
        }
    }

    Callable& operator=(const Callable& other) {
        if (this == &other) return *this;
        if (impl) delete impl;
        impl = other.impl ? other.impl->clone() : nullptr;
        return *this;
    }

    Callable(Callable&& other) noexcept : impl(other.impl) {
        other.impl = nullptr;
    }

    Callable& operator=(Callable&& other) noexcept {
        if (this == &other) return *this;
        if (impl) delete impl;
        impl = other.impl;
        other.impl = nullptr;
        return *this;
    }

    template<typename C> Callable(C c) {
        impl = new CallableWrapper<C, ReturnType>(static_cast<C&&>(c));
    }

    ReturnType operator()() {
        // static_assert(impl != nullptr, "Calling a null Callable object");
        assert(impl != nullptr);
        return (*impl)();
    }

    ReturnType execute() {
        // static_assert(impl != nullptr, "Calling a null Callable object");
        assert(impl != nullptr);
        return (*impl)();
    }

    bool isValid() const noexcept {
        return impl != nullptr;
    }

    void reset() {
        delete impl;
        impl = nullptr;
    }

    ~Callable() {
        if (impl == nullptr) return;
        delete impl;
    }

    operator ReturnType() {
        return *(this->impl);
    }
};

} // namespace vislib::core
