var bigInt = require("big-integer");

// 64 round_constants
SEVEN = bigInt(7);
FORTYTWO = bigInt(42);
len_round_constants = 64;
round_constants = new Array();
for (i = 0; i < 64; i++) {
  round_constants[i] = bigInt(i).pow(SEVEN).xor(FORTYTWO);
}

modulus = bigInt(2).pow(256).subtract(bigInt(2).pow(32).multiply(351)).add(1);
nsteps = 8192;

// Forward MiMC
THREE = bigInt(3);
input = bigInt(3);
trace = bigInt(input);
console.time("forward-mimc");
for (i = 1; i < nsteps; i++) {
  trace = trace.pow(THREE).add(round_constants[i%len_round_constants]).mod(modulus);
}
console.timeEnd("forward-mimc");
output = bigInt(trace)

// Reverse MiMC
console.time("reverse-mimc");
rtrace = bigInt(output)
little_fermat_expt = bigInt(modulus).multiply(2).subtract(1).divide(3)
for (i = nsteps - 1; i > 0; i--) {
  rtrace = rtrace.subtract(round_constants[i%len_round_constants]).modPow(little_fermat_expt, modulus);
}
console.timeEnd("reverse-mimc");

// print if input matches
if ( rtrace.equals(input) ) {
  console.log("PASS");
} else {
  console.log("FAIL");
}
