#pragma once

#include <vislib_robo.hpp>
#include <Arduino.h>

#define DEFAULT_INTERRUPT_MODE RISING

namespace vislib::binds::arduino {


struct MillisGetter_t {
    core::Result<size_t> operator()() const noexcept {
        return static_cast<size_t>(millis());
    }
};

auto millisGetter = core::TimeGetter<size_t>(MillisGetter_t{});

using port_t = uint8_t;

using InterruptTable = CallbackTable<port_t>;

CallbackPortInitializer<port_t> interruptPortInitializer = [](const port_t& port) -> core::Error {
    pinMode(port, INPUT);
    return {};
};

CallbackInitializer<port_t> interruptInitializer = [](const CallbackBase<port_t>&) -> core::Error {
    return {};
};

// CallbackAttacher<port_t> interruptAttacher = [](const CallbackBase<port_t>& callback) -> util::Error {
//     CallbackFunctor functor = callback;
//
//     attachInterrupt(digitalPinToInterrupt(callback.port), [=]() -> void {functor();}, DEFAULT_INTERRUPT_MODE);
//     return {};
// };

CallbackChecker<port_t> interruptChecker = [](const CallbackBase<port_t>& callback) -> bool {
    return digitalRead(callback.port) == HIGH;
};

};
