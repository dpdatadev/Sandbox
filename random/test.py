if __name__ == "__main__":
    import sys
    print("\n")
    print(sys.implementation)
    print("\n")
    print("Hello, Rasperry Pi Pico!")
    from picozero import pico_led
    
    pico_led.blink()
    
    import machine
    import time

    adcpin = 4
    sensor = machine.ADC(adcpin)
  
    def ReadTemperature():
        adc_value = sensor.read_u16()
        volt = (3.3/65535) * adc_value
        temperature = 27 - (volt - 0.706)/0.001721
        return round(temperature, 1)
  
    while True:
        temperature = ReadTemperature()
        print(f"Reading Temperature (Celcius): {temperature}")
        time.sleep(5)