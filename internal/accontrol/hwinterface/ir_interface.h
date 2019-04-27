#include <stdio.h>
#include "irslinger.h"

int runIR(uint32_t pin, const char *code)
{
        int frequency = 38000;           // The frequency of the IR signal in Hz
        double dutyCycle = 0.5;          // The duty cycle of the IR signal. 0.5 means for every cycle,
                                         // the LED will turn on for half the cycle time, and off the other half
        int leadingPulseDuration = 9000; // The duration of the beginning pulse in microseconds
        int leadingGapDuration = 4500;   // The duration of the gap in microseconds after the leading pulse
        int onePulse = 562;              // The duration of a pulse in microseconds when sending a logical 1
        int zeroPulse = 562;             // The duration of a pulse in microseconds when sending a logical 0
        int oneGap = 1688;               // The duration of the gap in microseconds when sending a logical 1
        int zeroGap = 562;               // The duration of the gap in microseconds when sending a logical 0
        int sendTrailingPulse = 1;       // 1 = Send a trailing pulse with duration equal to "onePulse"
                                         // 0 = Don't send a trailing pulse

        int result = irSling(
                pin,
                frequency,
                dutyCycle,
                leadingPulseDuration,
                leadingGapDuration,
                onePulse,
                zeroPulse,
                oneGap,
                zeroGap,
                sendTrailingPulse,
                code);

        return result;
}