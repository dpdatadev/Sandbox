# anything existing as a library/API, we will call directly from the code to avoid using system/popen
# in this case -l uuid --> uuid_generate_random()
uuid | cut -f1 -d-