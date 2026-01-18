#pragma once

namespace vislib::core {

enum class ErrorCode {
    success,
    failure,
    reachedTheTarget,
    initFailed,
    invalidArgument,
    invalidConfiguration,
    invalidResource,
    failedConnection,
    outOfRange,
    surpassedLimit,
    indexOutOfRange,
    emptyArray,
    zeroDivision
};

class Error;
template<typename T> class Result;

} //namespace vislib::core
