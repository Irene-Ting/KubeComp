#include <iostream>
#include <math.h>
#include <chrono>

// Kernel function to add the elements of two arrays
__global__
void add(int n, float *x, float *y)
{
  for (int i = 0; i < n; i++)
    y[i] = x[i] + y[i];
}

int main(int argc, char** argv) 
{
  if (argc != 2) 
  {
    fprintf(stderr, "must provide exactly 1 argument!\n");
    return 1;
  }

  unsigned long long sec = atoll(argv[1]);

  std::cout << "init" << std::endl;
  
  int N = 1<<10;
  float *x, *y;

  // Allocate Unified Memory – accessible from CPU or GPU
  cudaMallocManaged(&x, N*sizeof(float));
  cudaMallocManaged(&y, N*sizeof(float));

  // initialize x and y arrays on the host
  for (int i = 0; i < N; i++) {
    x[i] = 1.0f;
    y[i] = 2.0f;
  }

  // Run kernel on the GPU
  std::cout << "Start loop for " << sec << "s" << std::endl;
  auto start = std::chrono::steady_clock::now();
  while (std::chrono::steady_clock::now() - start < std::chrono::seconds(sec)) {
    add<<<1, 1>>>(N, x, y);
    cudaDeviceSynchronize();
  } 
  std::cout << "End loop" << std::endl;

  // Wait for GPU to finish before accessing on host
  // cudaDeviceSynchronize();

  // Free memory
  cudaFree(x);
  cudaFree(y);
  
  return 0;
}