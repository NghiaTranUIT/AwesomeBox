// This #include statement was automatically added by the Particle IDE.
#include "HttpClient/HttpClient.h"

// ------------
// Blink an LED
// ------------

/*-------------

We've heavily commented this code for you. If you're a pro, feel free to ignore it.

Comments start with two slashes or are blocked off by a slash and a star.
You can read them, but your device can't.
It's like a secret message just for you.

Every program based on Wiring (programming language used by Arduino, and Particle devices) has two essential parts:
setup - runs once at the beginning of your program
loop - runs continuously over and over

You'll see how we use these in a second.

This program will blink an led on and off every second.
It blinks the D7 LED on your Particle device. If you have an LED wired to D0, it will blink that LED as well.

-------------*/


// First, we're going to make some variables.
// This is our "shorthand" that we'll use throughout the program:

int led1 = D0; // Instead of writing D0 over and over again, we'll write led1
// You'll need to wire an LED to this one to see it blink.

int led2 = D7; // Instead of writing D7 over and over again, we'll write led2
// This one is the little blue LED on your board. On the Photon it is next to D7, and on the Core it is next to the USB jack.

int btn = D5;
int val = 0;

// Having declared these variables, let's move on to the setup function.
// The setup function is a standard part of any microcontroller program.
// It runs only once when the device boots up or is reset.

unsigned int nextTime = 0;    // Next time to contact the server
HttpClient http;

// Headers currently need to be set at init, useful for API keys etc.
http_header_t headers[] = {
    //  { "Content-Type", "application/json" },
    //  { "Accept" , "application/json" },
    { "Accept" , "*/*"},
    { NULL, NULL } // NOTE: Always terminate headers will NULL
};

http_request_t request;
http_response_t response;

void setup() {

  // We are going to tell our device that D0 and D7 (which we named led1 and led2 respectively) are going to be output
  // (That means that we will be sending voltage to them, rather than monitoring voltage that comes from them)

  // It's important you do this here, inside the setup() function rather than outside it or in the loop function.

Serial.begin(9600);

  pinMode(led1, OUTPUT);
  pinMode(led2, OUTPUT);
  pinMode(btn, INPUT);
}

// Next we have the loop function, the other essential part of a microcontroller program.
// This routine gets repeated over and over, as quickly as possible and as many times as possible, after the setup function is called.
// Note: Code that blocks for too long (like more than 5 seconds), can make weird things happen (like dropping the network connection).  The built-in delay function shown below safely interleaves required background activity, so arbitrarily long delays can safely be done if you need them.

void loop() {

        if (nextTime > millis()) {
        return;
    }

  val = digitalRead(btn);  // read input value
  if (val == HIGH)
  {         // check if the input is HIGH (button released)
    digitalWrite(led1, HIGH);  // turn LED OFF
    digitalWrite(led2, HIGH);

    // Request path and body can be set at runtime or at setup.
    request.hostname = "172.20.10.2";
    request.port = 10000;
    request.path = "/in-case-of-fire";

    // The library also supports sending a body with your request:
    //request.body = "{\"key\":\"value\"}";

    // Get request
    http.get(request, response, headers);
    Serial.print("Application>\tResponse status: ");
    Serial.println(response.status);

    Serial.print("Application>\tHTTP Response Body: ");
    Serial.println(response.body);

// Push IFTTT
   // Spark.publish("Fire", "Post Twitter");

    nextTime = millis() + 10000;
  }
  else
  {
        digitalWrite(led1, LOW);  // turn LED OFF
    digitalWrite(led2, LOW);
  }
}
