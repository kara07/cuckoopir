#ifndef MATRIX_H
#include <stdio.h>
#include <stdlib.h>
#include <assert.h>

#define BLOCK_DIM 32

#define BLOCK_NUM 32   //块数量
#define THREAD_NUM 256 // 每个块中的线程数
#define R_SIZE BLOCK_NUM * THREAD_NUM
#define M_SIZE R_SIZE * R_SIZE

struct Matrix_int
{
    int width, height;
    int stride;
    int *elements = NULL;
};

extern "C" void generate_random_matrix(struct Matrix_int *mat, int m, int n);


extern "C" void destroy_matrix(struct Matrix_int *mat);


#endif