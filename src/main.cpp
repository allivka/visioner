#include "visioner.hpp"
#include <stdio.h>
#include <Wire.h>

using namespace vislib;

// motor::SpeedRange rpmSpeedRange(-600, 600);
// motor::SpeedRange motorUseSpeedRange(-1500, 1500);
// motor::SpeedRange motorInterfaceAngularSpeedRange(
//     motorUseSpeedRange.mapValueToRange(-1500, rpmSpeedRange) * 2 * PI, motorUseSpeedRange.mapValueToRange(1500, rpmSpeedRange) * 2 * PI);
//
// double wheelR = 0.1;
// double motorDistance = 0.3;

// platform::PlatformMotorConfig config({
//     motor::MotorInfo(90, motorDistance, wheelR, motorUseSpeedRange, motorInterfaceAngularSpeedRange),
//     motor::MotorInfo(180, motorDistance, wheelR, motorUseSpeedRange, motorInterfaceAngularSpeedRange),
//     motor::MotorInfo(-90, motorDistance, wheelR, motorUseSpeedRange, motorInterfaceAngularSpeedRange),
//     motor::MotorInfo(0, motorDistance, wheelR, motorUseSpeedRange, motorInterfaceAngularSpeedRange)
// });
//
// platform::Platform<V5::motor::V5MotorController> plat(config);

vislib_mpu6050::GyroscopeCalculator mpu;

struct MillisGetter : public util::TimeGetter<size_t, MillisGetter> {
    util::Result<size_t> getTimeImplementation() const {
        return static_cast<size_t>(millis());
    }
};

util::IncrementTimer<size_t, MillisGetter> timer{MillisGetter()};

char *str;

void setup() {
    
    str = static_cast<char*>(malloc(500));
    
    Serial.begin(9600);
    delay(100);
    
    Wire.begin();
    delay(100);
    
    Serial.println("Initializing MPU...");
    
    mpu.initialize();
    delay(100);
    
    Serial.println("Testing MPU6050 connection...");
    if(!mpu.testConnection()) {
        Serial.println("MPU6050 connection failed");
        while(true);
    }
    
    Serial.println("MPU6050 connection successful");
    
    Serial.println("Calibrating gyroscope...\n");
    
    mpu.calibrate();
    
    timer.start();
    
    Serial.println(timer.getTime()());
    
    Serial.println("initializing YPR calculator");
    
    auto temp = gyro::YPRElementCalculatorConfig<double, size_t, double>();
    temp.integralWeight = 0.98;
    
    mpu.initCalculator(
        temp,
        temp,
        temp
    );
    
    // Serial.println("Initializing platform controller");
    // Vex5.begin();
    
    // auto e = plat.init(util::Array<VEX5_PORT_t>({(VEX5_PORT_t)1, (VEX5_PORT_t)2, (VEX5_PORT_t)3, (VEX5_PORT_t)4}));
    
    // if(e) {
    //     Serial.print(e.msg.c_str());
    // }
    
    Serial.println("Done startup");
    
    // delay(5000);
    
}

// void go(double angle, double speed) {
//     auto speeds = platform::calculators::calculatePlatformLinearSpeeds(config, angle, speed);
//     if(speeds) {
//         Serial.println(("Ooops, something went wrong in calculating speeds for linear movement of the platform: " + speeds.Err().msg).c_str());
//         return;
//     }
    
//     Serial.print("1: ");
//     Serial.print(speeds()[0]);
//     Serial.print("; 2: ");
//     Serial.print(speeds()[1]);
//     Serial.print("; 3: ");
//     Serial.print(speeds()[2]);
//     Serial.print("; 4: ");
//     Serial.println(speeds()[3]);
    
//     auto err = plat.setSpeeds(speeds());
    
//     if(err.errcode != util::ErrorCode::success) {
//         Serial.println(("Ooops, something went wrong in applying speeds to motors for linear movement of the platform: " + err.msg).c_str());
//         Serial.println(err.msg.c_str());
//     }
// }

// void move(double angle, double speed, ull_t delayMs) {
//     go(angle, speed);
//     delay(delayMs);
// }

// double speed = 200;
// ull_t sectionTime = 1000;

void loop() {
    // move(0, speed, sectionTime);
    // move(90, speed, sectionTime);
    // move(180, speed, sectionTime);
    // move(-90, speed, sectionTime);
    // move(0, speed, sectionTime);
    // move(45, speed, sectionTime);
    // move(90, speed, sectionTime);
    // move(135, speed, sectionTime);
    // move(180, speed, sectionTime);
    // move(-135, speed, sectionTime);
    // move(-90, speed, sectionTime);
    // move(-45, speed, sectionTime);

    ++timer;

    auto info = mpu.calculateGyroData(timer.getTime()());
    
    if(info) {
        Serial.println("SHIT");
        Serial.println(info.Err().msg.c_str());
        return;
    }
    
    sprintf(str, "[%d %d ms]: speedX = %s;\tspeedY = %s;\tspeedZ = %s;\tyaw = %s;\tpitch = %s;\troll = %s\n",
        timer.getTime()() / 1000,
        timer.getTime()() % 1000,
        String(info().speed[0]).c_str(),
        String(info().speed[1]).c_str(),
        String(info().speed[2]).c_str(),
        String(info().ypr.yaw).c_str(),
        String(info().ypr.pitch).c_str(),
        String(info().ypr.roll).c_str()
    );
    
    Serial.print(str);
    
}