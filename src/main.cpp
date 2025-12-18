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


constexpr binds::arduino::port_t mpuInterruptPort = 2;
binds::mpu6050::GyroscopeDMP mpu;
util::IncrementTimer<int64_t> timer([]() -> util::Result<int64_t> {
    return static_cast<int64_t>(millis());
});

char *str;


void setup() {
    
    pinMode(LED_BUILTIN, OUTPUT);
    
    str = static_cast<char*>(malloc(500));
    
    Serial.begin(9600);
    delay(100);

    Serial.println("Initialized serial\n\nInitializing timer");
    
    auto e = timer.start();
    
    
    if (e) while (true) Serial.println(String("Failed to initialize timer: ") + String(e.error().msg.c_str()));

    delay(100);

    Serial.println("Initialized timer\n\nInitializing I2C protocol");

    Wire.begin();
    delay(100);

    Serial.println("Initialized I2C protocol\n\nInitializing MPU6050");

    mpu.initialize();

    Serial.println("Initialized MPU6050\n\nTesting MPU6050 connection");
    delay(100);

    if(!mpu.testConnection()) while(true) Serial.println("MPU6050 connection failed");

    Serial.println("Initialized MPU6050\n\nInitializing MPU6050 DMP driver interrupt table");
    delay(100);

    auto er = binds::mpu6050::GyroscopeDMP::initInterruptTable(
        util::Array<binds::arduino::port_t>(util::Array<binds::arduino::port_t>({mpuInterruptPort}))
    );

    if (er) while (true) Serial.println(er.msg.c_str());

    Serial.println("Initialized MPU6050 DMP driver interrupt table\n\nInitializing MPU6050 DMP");
    delay(100);

    // er = mpu.initDMP(mpuInterruptPort);
    // if (er) while (true) Serial.println(er.msg.c_str());

    Serial.println("Initialized MPU6050 DMP\n\n");
    delay(100);
    
    // Serial.println("Initializing platform controller");
    // Vex5.begin();
    
    // auto e = plat.init(util::Array<VEX5_PORT_t>({(VEX5_PORT_t)1, (VEX5_PORT_t)2, (VEX5_PORT_t)3, (VEX5_PORT_t)4}));
    
    // if(e) {
    //     Serial.print(e.msg.c_str());
    // }
    
    Serial.println("\nDone initialization");
    
    
    digitalWrite(LED_BUILTIN, HIGH);
    delay(500);
    
    
}

// void go(double angle, double speed) {
//     auto speeds = platform::calculators::calculatePlatformLinearSpeeds(config, angle, speed);
//     if(speeds) {
//         Serial.println(("Ooops, something went wrong in calculating speeds for linear movement of the platform: " + speeds.error().msg).c_str());
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
    
    const auto info = mpu.getGyroData();

    if (info) {
        sprintf(str, "Shit happened: %s\n", info.error().msg.c_str());
        goto usual;
    }

    sprintf(str, "[%d %d ms]: speedX = %s;\t speedY = %s;\t speedZ = %s;\t yaw = %s;\t pitch = %s;\t roll = %s\n",
        static_cast<size_t>(timer.getTime()() / 1000),
        static_cast<size_t>(timer.getTime()() % 1000),
        String(int(info().speed[0])).c_str(),
        String(int(info().speed[1])).c_str(),
        String(int(info().speed[2])).c_str(),
        String(int(info().ypr.yaw)).c_str(),
        String(int(info().ypr.pitch)).c_str(),
        String(int(info().ypr.roll)).c_str()
    );


usual:
    ++timer;
    Serial.println(static_cast<size_t>(timer.getTime()() / 1000));
    mpu.update(nullptr);

    Serial.print(str);

}