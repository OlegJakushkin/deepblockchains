/*
 gcc mimc.c -o mimc -lgmp
 ./mimc
*/
/* Program to demonstrate time taken by function fun() */
#include <stdio.h>
#include <gmp.h>
#include <sys/time.h>
#define nsteps 8192
#define len_round_constants 64

// 2^256 - 2^32 * 351 +1 = 115792089237316195423570985008687907853269984665640564039457584006405596119041
void default_modulus(mpz_t modulus)
{
  mpz_t t0;
  mpz_t t1;
  mpz_t two_256;
  mpz_t two_32;

  mpz_init(t0);
  mpz_init(t1);
  mpz_init(two_256);
  mpz_init(two_32);
  mpz_init(modulus);

  // 2^256 - 351*2^32
  mpz_ui_pow_ui(two_256, 2, 256);
  mpz_ui_pow_ui(two_32, 2, 32);
  mpz_mul_ui(t0, two_32, 351);
  mpz_sub(t1, two_256, t0);
  mpz_add_ui(modulus, t1, 1);
  // gmp_printf("%Zd\n",  modulus);
}


int main()
{
  mpz_t modulus;
  default_modulus(modulus);

  // some temporary vars
  mpz_t t0, t1;
  mpz_init(t0);
  mpz_init(t1);

  // compute 64 round_constants
  mpz_t FORTYTWO;
  mpz_init(FORTYTWO);
  mpz_set_ui(FORTYTWO, 42);
  mpz_t round_constants[len_round_constants];
  for (unsigned long i = 0; i < len_round_constants; i++) {
    mpz_init(round_constants[i]);
    mpz_ui_pow_ui(t0, i, 7);
    mpz_xor(round_constants[i], t0, FORTYTWO);
  }

  // Forward MiMC
  mpz_t THREE;
  mpz_t input;
  mpz_t trace;
  mpz_init(THREE);
  mpz_set_ui(THREE, 3);
  mpz_init(input);
  mpz_set_ui(input, 3);
  mpz_init(trace);
  mpz_set_ui(trace, 3);
  struct timeval begin, end;
  gettimeofday(&begin,NULL);
  int i;
  for (i = 1; i < nsteps; i++) {
    mpz_pow_ui(t0, trace, 3);
    mpz_add(t1, t0, round_constants[i%len_round_constants]);
    mpz_mod(trace, t1, modulus);
  }
  gettimeofday(&end,NULL);
  mpz_t output;
  mpz_init(output);
  mpz_set(output, trace);
  int elapsed = ((end.tv_sec - begin.tv_sec) * 1000000) + (end.tv_usec - begin.tv_usec);
  printf("forward-mimc: %d microseconds\n", elapsed);

  // Reverse MiMC
  gettimeofday(&begin,NULL);
  mpz_t rtrace;
  mpz_init(rtrace);
  mpz_set(rtrace, output);
  mpz_t little_fermat_expt;
  mpz_init(little_fermat_expt);

  mpz_mul_ui(t1, modulus, 2);
  mpz_sub_ui(t0, t1, 1);
  mpz_div_ui(little_fermat_expt, t0, 3);

  for (i = nsteps - 1; i > 0; i--) {
    mpz_sub(t0, rtrace, round_constants[i%len_round_constants]);
    mpz_powm(rtrace, t0, little_fermat_expt, modulus);
  }
  gettimeofday(&end,NULL);
  //	fmt.Printf("reverse-mimc: %s\n", time.Since(start))
  elapsed = ((end.tv_sec - begin.tv_sec) * 1000000) + (end.tv_usec - begin.tv_usec);
  printf("reverse-mimc: %d microseconds\n", elapsed);

  if ( mpz_cmp(rtrace, input) != 0 ) {
    printf("FAIL\n");
  } else {
    printf("PASS\n");
  }
}
