#pragma once

namespace vislib::core {

template <typename T> inline void swap(T& x, T& y) noexcept(noexcept(T(move(x))) && noexcept(x = move(y)) && noexcept(y = move(x))) {
    T temp = move(x);
    x = move(y);
    y = move(temp);
}

template <typename T> inline void deleter(T *ptr) noexcept {
    delete ptr;
}

template<typename T, void (*DELETER)(T*) = deleter> class UniquePtr {
protected:
    T *ptr = nullptr;

public:
    UniquePtr() = default;

    explicit UniquePtr(const T& v) {
        ptr = new T(v);
    }

    explicit UniquePtr(T *p_ptr) : ptr(p_ptr) {}

    UniquePtr(const UniquePtr&) = delete;

    UniquePtr(nullptr_t) : ptr(nullptr) {}

    UniquePtr(UniquePtr&& other) noexcept {
        ptr = other.ptr;
        other.ptr = nullptr;
    }

    ~UniquePtr() noexcept {
        if (ptr != nullptr) DELETER(ptr);
    }

    T& operator*() const noexcept {
        return *ptr;
    }

    T* operator->() const noexcept {
        return ptr;
    }

    inline explicit operator bool() const noexcept {
        return ptr != nullptr;
    }

    UniquePtr& operator=(const UniquePtr&) = delete;

    UniquePtr& operator=(UniquePtr&& other) noexcept {
        if (ptr != nullptr) DELETER(ptr);
        ptr = other.ptr;
        other.ptr = nullptr;
        return *this;
    }

    UniquePtr& operator=(const T& v) {
        if(ptr == nullptr) {
            ptr = new T(v);
        } else {
            *ptr = v;
        }
        return *this;
    }

    bool operator==(const UniquePtr<T>& other) const noexcept {
        return *ptr == *other.ptr;
    }

    bool operator!=(const UniquePtr<T>& other) const noexcept {
        return *ptr != *other.ptr;
    }

    bool operator==(nullptr_t) const noexcept {
        return ptr == nullptr;
    }

    bool operator!=(nullptr_t) const noexcept {
        return ptr != nullptr;
    }

    bool is(const UniquePtr<T>& other) const noexcept {
        return ptr == other.ptr;
    }

    inline T* get() const noexcept {
        return ptr;
    }

    T* release() noexcept {
        T* temp = ptr;
        ptr = nullptr;
        return temp;
    }

    void reset(T *new_ptr = nullptr) noexcept {
        if (ptr != nullptr) DELETER(ptr);
        ptr = new_ptr;
    }

};

} //namespace vislib::core
