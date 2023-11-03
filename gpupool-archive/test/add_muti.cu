#include <iostream>
#include <math.h>
#include <chrono>
#include <omp.h>

// Kernel function to add the elements of two arrays
__global__
void add(int n, float *x, float *y)
{
  for (int i = 0; i < n; i++)
    y[i] = x[i] + y[i];
}

int main(int argc, char** argv) 
{
    if (argc != 3) 
    {
        fprintf(stderr, "must provide exactly 1 argument!\n");
        return 1;
    }

    unsigned long long num_of_gpu = atoll(argv[1]);
    unsigned long long sec = atoll(argv[2]);

    std::cout << "init " << num_of_gpu << " gpus" << std::endl;
    
    int N = 1<<10;
    float *x = (float*)malloc(N*sizeof(float));
    float *y = (float*)malloc(N*sizeof(float));
    float *x_device[8];
    float *y_device[8];

    // initialize x and y arrays on the host
    for (int i = 0; i < N; i++) {
        x[i] = 1.0f;
        y[i] = 2.0f;
    }
  
    #pragma omp parallel num_threads(num_of_gpu)
    {
        int dev_id = omp_get_thread_num();
        cudaSetDevice(dev_id);

        cudaMalloc(&x_device[dev_id], N * sizeof(float));
        cudaMemcpy(x_device[dev_id], x, N * sizeof(float), cudaMemcpyHostToDevice);
        cudaMalloc(&y_device[dev_id], N * sizeof(float));
        cudaMemcpy(y_device[dev_id], y, N * sizeof(float), cudaMemcpyHostToDevice);

        // Run kernel on the GPU
        if (dev_id == 0)
            std::cout << "Start loop for " << sec << "s" << std::endl;
        auto start = std::chrono::steady_clock::now();
        while (std::chrono::steady_clock::now() - start < std::chrono::seconds(sec)) {
            add<<<1, 1>>>(N, x_device[dev_id], y_device[dev_id]);
            cudaDeviceSynchronize();
        } 

        #pragma omp barrier
        if (dev_id == 0)
            std::cout << "End loop" << std::endl;

        // Free memory
        cudaFree(x_device[dev_id]);
        cudaFree(y_device[dev_id]);
    }
    
    return 0; 
}

// nvcc add_muti.cu -Xcompiler -fopenmp -o add_multi