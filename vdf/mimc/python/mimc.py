import time

modulus = 2**256 - 2**32 * 351 + 1
little_fermat_expt = (modulus*2-1)//3
round_constants = [(i**7) ^ 42 for i in range(64)]
input = 3
steps = 8192

# Forward MiMC
start_time = time.time()
inp = input
for i in range(1,steps):
    inp = (inp**3 + round_constants[i % len(round_constants)]) % modulus
print("forward-mimc %.4f sec" % (time.time() - start_time))
output = inp

# Reverse MIMC
rtrace = output
start_time = time.time()
for i in range(steps - 1, 0, -1):
    rtrace = pow(rtrace-round_constants[i%len(round_constants)], little_fermat_expt, modulus)
print("reverse-mimc: %.4f sec" % (time.time() - start_time))

if rtrace == input:
    print("PASS")
else:
    print("FAIL")
