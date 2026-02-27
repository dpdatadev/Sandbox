# PROTOTYPE - PicoBoo
# Client Script of NETEVENTS API
# Polls API server on same domain (centralized hub)
# Consumes "Events" network API which trivially determines "device_added"/"device_removed"


# Upload script to RP2040 (RENAME to main.py), enter wifi credentials
# Be notified of changes on the local network via ipconf/arp/nmap
# The Pico creates a physical alarm (adjust duration of blink and buzz) 
# and optionally sends email to admin.

import network
import urequests
import time
from machine import Pin
import ujson

# -----------------------------
# CONFIG
# -----------------------------
WIFI_SSID = "your_ssid"
WIFI_PASS = "your_password"
SERVER_URL = "http://your_server:8080/api/network/events"
POLL_INTERVAL = 30  # seconds

# -----------------------------
# HARDWARE SETUP
# -----------------------------
led = Pin("LED", Pin.OUT)
buzzer = Pin(15, Pin.OUT)  # adjust GPIO pin as needed

# -----------------------------
# WIFI CONNECT
# -----------------------------
def connect_wifi():
    wlan = network.WLAN(network.STA_IF)
    wlan.active(True)
    wlan.connect(WIFI_SSID, WIFI_PASS)

    while not wlan.isconnected():
        time.sleep(1)

    print("Connected:", wlan.ifconfig())
    return wlan

# -----------------------------
# ALERT FUNCTIONS
# -----------------------------
def alert_device_added(event):
    print("Device Added:", event["ip"])
    flash_led(3)
    buzz(3)

def alert_device_removed(event):
    print("Device Removed:", event["ip"])
    flash_led(1)
    buzz(1)

def flash_led(times):
    for _ in range(times):
        led.on()
        time.sleep(0.2)
        led.off()
        time.sleep(0.2)

def buzz(times):
    for _ in range(times):
        buzzer.on()
        time.sleep(0.2)
        buzzer.off()
        time.sleep(0.2)

# -----------------------------
# MAIN POLL LOOP
# -----------------------------
def poll_server():
    try:
        response = urequests.get(SERVER_URL)
        data = response.json()
        response.close()

        events = data.get("events", [])

        for event in events:
            if event["type"] == "device_added":
                alert_device_added(event)
            elif event["type"] == "device_removed":
                alert_device_removed(event)

    except Exception as e:
        print("Error polling server:", e)

# -----------------------------
# MAIN
# -----------------------------
def main():
    connect_wifi()

    while True:
        poll_server()
        time.sleep(POLL_INTERVAL)

main()
