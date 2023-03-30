#include "Matrix.h"
#include <cuda_runtime.h>
#include <stdlib.h>
#include <chrono>
#include <iostream>

#define GET_ELEMENT(i, j, m, n) (i * n + j)
#define CHECK(call) \
{ \
    cudaError_t cudaStatus = (call); \
    if(cudaStatus != cudaSuccess) \
    { \
        fprintf(stderr, "CUDA error at %s:%d - %s\n", __FILE__, __LINE__, cudaGetErrorString(cudaStatus)); \
        exit(1); \
    } \
}

extern "C"
{

    __global__ void mat_mul(int *mat1, int *mat2, int *result) {
        const int bid = blockIdx.x;
        const int tid = threadIdx.x;
        // 每个线程计算一行
        const int row = bid * THREAD_NUM + tid;
        for (int c = 0; c < R_SIZE; c++) {
            for (int n = 0; n < R_SIZE; n++) {
            result[row*R_SIZE+c] += mat1[row*R_SIZE+n] * mat2[n*R_SIZE+c];
            }
        }
    }


    void generate_random_matrix(struct Matrix_int *mat, int m, int n)
    {
        assert(mat != NULL);
        mat->width = n;
        mat->height = m;
        if (mat->elements == NULL)
        {
            mat->elements = (int *)malloc(m * n * sizeof(int));
        }

        for(int i = 0; i < n; ++i)
        {
            for(int j = 0; j < m; ++j)
            {
                mat->elements[i * n + j] = rand() % 100;
            }
        }
    }

    void destroy_matrix(struct Matrix_int *mat)
    {
        assert(mat != NULL);
        free(mat->elements);
    }

    void MatrixMulCPU(int *a, int *b, int *c, int m, int n, int p)
    {
        for(int i = 0; i < m; ++i)
        {
            for(int j = 0; j < p; ++j)
            {
                for(int k = 1; k <= n; ++k)
                {
                    c[i*n + j] +=  a[ i*n +k - 1] * b[(k - 1)*p  + j ];
                }
            }
        }
    }


}

__global__ void mm_kernel(int *mat_1, int *mat_2, int *mat_3, int m, int n, int p)
{
    __shared__ int mat_1_tile[BLOCK_DIM][BLOCK_DIM];
    __shared__ int mat_2_tile[BLOCK_DIM][BLOCK_DIM];

    int acc_sum{0};

    for (size_t tile_idx{0};
         tile_idx < ceilf(static_cast<float>(n) / BLOCK_DIM); ++tile_idx)
    {
        size_t i{blockIdx.y * blockDim.y + threadIdx.y};
        size_t j{tile_idx * blockDim.x + threadIdx.x};
        if ((i < m) && (j < n))
        {
            mat_1_tile[threadIdx.y][threadIdx.x] = mat_1[i * n + j];
        }
        else
        {
            mat_1_tile[threadIdx.y][threadIdx.x] = 0;
        }
        i = tile_idx * blockDim.y + threadIdx.y;
        j = blockIdx.x * blockDim.x + threadIdx.x;
        if ((i < n) && (j < p))
        {
            mat_2_tile[threadIdx.y][threadIdx.x] = mat_2[i * p + j];
        }
        else
        {
            mat_2_tile[threadIdx.y][threadIdx.x] = 0;
        }
        __syncthreads();
        for (size_t k{0}; k < BLOCK_DIM; ++k)
        {
            acc_sum += mat_1_tile[threadIdx.y][k] * mat_2_tile[k][threadIdx.x];
        }
        __syncthreads();
    }

    // 2D block and 2D thread
    // Each thread computes one cell in mat_3.
    size_t i{blockIdx.y * blockDim.y + threadIdx.y};
    size_t j{blockIdx.x * blockDim.x + threadIdx.x};

    if ((i < m) && (j < p))
    {
        mat_3[i * p + j] = acc_sum;
    }
}


    void MatrixMul(Matrix_int *a, Matrix_int *b, Matrix_int *c, int m, int n, int k)
    {
        assert(a != NULL && b != NULL && c != NULL);
        assert(a->elements != NULL && b->elements != NULL && c->elements != NULL);
        dim3 blockNum(32, 32);
        dim3 threadsPerBlock(BLOCK_DIM, BLOCK_DIM);
        int *d_a = nullptr, *d_b = nullptr, *d_c = nullptr;
        int size_a = (a->width) * (a->height) * sizeof(int);
        int size_b = (b->width) * (b->height) * sizeof(int);
        int size_c = (c->width) * (c->height) * sizeof(int);

        CHECK(cudaMalloc((void**)&d_a, size_a));
        CHECK(cudaMalloc((void**)&d_b, size_b));
        CHECK(cudaMalloc((void**)&d_c, size_c));

        cudaMemcpy(d_a, a->elements, size_a, ::cudaMemcpyHostToDevice);
        cudaMemcpy(d_b, b->elements, size_b, ::cudaMemcpyHostToDevice);
        cudaDeviceSynchronize();

        mm_kernel<<<blockNum, threadsPerBlock>>>((int *)d_a, (int *)d_b, (int *)d_c, m, n, k);
        cudaMemcpy(c->elements, d_c, size_c, ::cudaMemcpyDeviceToHost);
        cudaDeviceSynchronize();
    }


    void test()
    {
        Matrix_int a {1000, 1000, NULL};
        Matrix_int b {1000, 1000, NULL};
        Matrix_int c {1000, 1000, NULL};

        generate_random_matrix(&a, a.height, a.width);
        generate_random_matrix(&b, b.height, b.width);
        generate_random_matrix(&c, c.height, c.width);

        // auto start_time_cpu = std::chrono::high_resolution_clock::now();
        // MatrixMulCPU(a.elements, b.elements, c.elements, 1024, 1024, 1024);
        // auto end_time_cpu = std::chrono::high_resolution_clock::now();
        // auto time_elasped = end_time_cpu - start_time_cpu;
        // std::cout << "CPU duration: " << std::chrono::duration_cast<std::chrono::microseconds>(time_elasped).count() << std::endl;

        auto start_time_gpu = std::chrono::high_resolution_clock::now();
        MatrixMul(&a, &b, &c, 8192, 8192, 8192);
        auto end_time_gpu = std::chrono::high_resolution_clock::now();
        auto time_elasped_gpu = end_time_gpu - start_time_gpu;
        std::cout << "GPU duration: " << std::chrono::duration_cast<std::chrono::microseconds>(time_elasped_gpu).count() << std::endl;

        int *mat1, *mat2, *result;
        int *g_mat1, *g_mat2, *g_mat_result;

        mat1 = (int*) malloc(M_SIZE * sizeof(int));
        mat2 = (int*) malloc(M_SIZE * sizeof(int));
        result = (int*) malloc(M_SIZE * sizeof(int));

        // initialize
        for (int i = 0; i < M_SIZE; i++) {
        mat1[i] = rand()/1000000;
        mat2[i] = rand()/1000000;
        result[i] = 0;
        }

        cudaMalloc((void **)&g_mat1, sizeof(int) * M_SIZE);
        cudaMalloc((void **)&g_mat2, sizeof(int) * M_SIZE);
        cudaMalloc((void **)&g_mat_result, sizeof(int) * M_SIZE);

        cudaMemcpy(g_mat1, mat1, sizeof(int) * M_SIZE, cudaMemcpyHostToDevice);
        cudaMemcpy(g_mat2, mat2, sizeof(int) * M_SIZE, cudaMemcpyHostToDevice);

        auto start_time_gpu2 = std::chrono::high_resolution_clock::now();
        mat_mul<<<BLOCK_NUM, THREAD_NUM>>>(g_mat1, g_mat2, g_mat_result);
        auto end_time_gpu2 = std::chrono::high_resolution_clock::now();
        auto time_elasped_gpu2 = end_time_gpu2 - start_time_gpu2;
        std::cout << "GPU duration_2: " << std::chrono::duration_cast<std::chrono::microseconds>(time_elasped_gpu2).count() << std::endl;



        cudaMemcpy(result, g_mat_result, sizeof(int) * M_SIZE, cudaMemcpyDeviceToHost);

    }

    void test2(){
        int *mat1, *mat2, *result;
        int *g_mat1, *g_mat2, *g_mat_result;
    
        mat1 = (int*) malloc(M_SIZE * sizeof(int));
        mat2 = (int*) malloc(M_SIZE * sizeof(int));
        result = (int*) malloc(M_SIZE * sizeof(int));

        // initialize
        for (int i = 0; i < M_SIZE; i++) {
        mat1[i] = rand()/1000000;
        mat2[i] = rand()/1000000;
        result[i] = 0;
        }

        cudaMalloc((void **)&g_mat1, sizeof(int) * M_SIZE);
        cudaMalloc((void **)&g_mat2, sizeof(int) * M_SIZE);
        cudaMalloc((void **)&g_mat_result, sizeof(int) * M_SIZE);

        cudaMemcpy(g_mat1, mat1, sizeof(int) * M_SIZE, cudaMemcpyHostToDevice);
        cudaMemcpy(g_mat2, mat2, sizeof(int) * M_SIZE, cudaMemcpyHostToDevice);

        mat_mul<<<BLOCK_NUM, THREAD_NUM>>>(g_mat1, g_mat2, g_mat_result);

        cudaMemcpy(result, g_mat_result, sizeof(int) * M_SIZE, cudaMemcpyDeviceToHost);
    }

int main()
{
    test();
}
