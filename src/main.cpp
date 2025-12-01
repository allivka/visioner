#include "visioner.hpp"

using namespace vislib;

motor::SpeedRange rpmSpeedRange(-600, 600);
motor::SpeedRange motorUseSpeedRange(-1500, 1500);
motor::SpeedRange motorInterfaceAngularSpeedRange(
    motorUseSpeedRange.mapValueToRange(-1500, rpmSpeedRange) * 2 * PI, motorUseSpeedRange.mapValueToRange(1500, rpmSpeedRange) * 2 * PI);

double wheelR = 0.1;
double motorDistance = 0.3;

platform::PlatformMotorConfig config({
    motor::MotorInfo(90, motorDistance, wheelR, motorUseSpeedRange, motorInterfaceAngularSpeedRange),
    motor::MotorInfo(180, motorDistance, wheelR, motorUseSpeedRange, motorInterfaceAngularSpeedRange),
    motor::MotorInfo(-90, motorDistance, wheelR, motorUseSpeedRange, motorInterfaceAngularSpeedRange),
    motor::MotorInfo(0, motorDistance, wheelR, motorUseSpeedRange, motorInterfaceAngularSpeedRange)
});

platform::Platform<V5::motor::V5MotorController> plat(config);

vislib_mpu6050::Gyroscope gyro;

void setup() {  
    Vex5.begin();
    Serial.begin(9600);

    ::gyro.initialize();

    auto e = plat.init(util::Array<VEX5_PORT_t>({(VEX5_PORT_t)1, (VEX5_PORT_t)2, (VEX5_PORT_t)3, (VEX5_PORT_t)4}));
    
    if(e) {
        Serial.print(e.msg.c_str());
    }
    
    Serial.println();
    Serial.println(motorInterfaceAngularSpeedRange.lowest);
    Serial.println(motorInterfaceAngularSpeedRange.highest);
    
    delay(5000);
    
}

void go(double angle, double speed) {
    auto speeds = platform::calculators::calculatePlatformLinearSpeeds(config, angle, speed);
    if(speeds) {
        Serial.println(("Ooops, something went wrong in calculating speeds for linear movement of the platform: " + speeds.Err().msg).c_str());
        return;
    }
    
    Serial.print("1: ");
    Serial.print(speeds()[0]);
    Serial.print("; 2: ");
    Serial.print(speeds()[1]);
    Serial.print("; 3: ");
    Serial.print(speeds()[2]);
    Serial.print("; 4: ");
    Serial.println(speeds()[3]);
    
    auto err = plat.setSpeeds(speeds());
    
    if(err.errcode != util::ErrorCode::success) {
        Serial.println(("Ooops, something went wrong in applying speeds to motors for linear movement of the platform: " + err.msg).c_str());
        Serial.println(err.msg.c_str());
    }
}

void move(double angle, double speed, ull_t delayMs) {
    go(angle, speed);
    delay(delayMs);
}

double speed = 200;
ull_t sectionTime = 1000;

void loop() {
    // move(0, speed, sectionTime);
    // move(90, speed, sectionTime);
    // move(180, speed, sectionTime);
    // move(-90, speed, sectionTime);
    move(0, speed, sectionTime);
    move(45, speed, sectionTime);
    move(90, speed, sectionTime);
    move(135, speed, sectionTime);
    move(180, speed, sectionTime);
    move(-135, speed, sectionTime);
    move(-90, speed, sectionTime);
    move(-45, speed, sectionTime);
}